package runtime

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	components_v1alpha1 "github.com/dapr/dapr/pkg/apis/components/v1alpha1"
	"github.com/dapr/dapr/pkg/components"
	diag "github.com/dapr/dapr/pkg/diagnostics"
	"github.com/dapr/dapr/pkg/modes"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	secretstores "github.com/dapr/components-contrib/secretstores"
	"github.com/dapr/kit/logger"
)

type ComponentRegistryV2 struct {
	ComponentRegistry
}

type ComponentInitializer interface {
	Category() ComponentCategory
	Initialize(*components_v1alpha1.Component) error
}

type ComponentFactory struct {
	logger       logger.Logger
	initializers map[ComponentCategory]ComponentInitializer
}

func NewComponentFactory(logger logger.Logger, initializers ...ComponentInitializer) *ComponentFactory {
	initializerMap := map[ComponentCategory]ComponentInitializer{}
	for _, i := range initializers {
		initializerMap[i.Category()] = i
	}
	return &ComponentFactory{
		logger:       logger,
		initializers: initializerMap,
	}
}

func (f *ComponentFactory) Initialize(component *components_v1alpha1.Component) error {
	category, err := f.getComponentCategory(component)
	if err != nil {
		return err
	}

	initializer, ok := f.initializers[category]
	if !ok {
		return fmt.Errorf("unable to load initializer for component category %s", category)
	}

	errChan := make(chan error)
	defer close(errChan)
	go func() {
		errChan <- initializer.Initialize(component)
	}()

	timeout, err := time.ParseDuration(component.Spec.InitTimeout)
	if err != nil {
		timeout = defaultComponentInitTimeout
	}

	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	case <-time.After(timeout):
		return fmt.Errorf("init timeout for component %s exceeded after %s", component.Name, timeout.String())
	}

	f.logger.Infof("init timeout for component %s exceeded after %s", component.Name, timeout.String())
	diag.DefaultMonitoring.ComponentLoaded()

	return nil
}

func (f *ComponentFactory) getComponentCategory(component *components_v1alpha1.Component) (ComponentCategory, error) {
	for _, category := range componentCategoriesNeedProcess {
		if strings.HasPrefix(component.Spec.Type, fmt.Sprintf("%s.", category)) {
			return category, nil
		}
	}
	return "", fmt.Errorf("incorrect type %s", component.Spec.Type)
}

type runtimeV2 struct {
	componentDefinitionManager *ComponentDefinitionManager
}

type ComponentAuthorizationPolicy func(comp *components_v1alpha1.Component) bool

// ComponentDefinitionManager manages the component definition lifecycle
type ComponentDefinitionManager struct {
	loader                    components.ComponentLoader
	authorizationPolicy       ComponentAuthorizationPolicy
	componentLifecycleManager *ComponentLifecycleManager
	logger                    logger.Logger
}

func (m *ComponentDefinitionManager) Load() error {
	comps, err := m.loader.LoadComponents()
	if err != nil {
		return err
	}
	m.logFound(comps)

	authorized := m.filterAuthorized(comps)
	m.process(authorized)

	return nil
}

func (m *ComponentDefinitionManager) process(comps []components_v1alpha1.Component) {
	for _, comp := range comps {
		m.componentLifecycleManager.Process(&comp)
	}
}

func (m *ComponentDefinitionManager) logFound(comps []components_v1alpha1.Component) {
	for _, comp := range comps {
		m.logger.Debugf("found component. name: %s, type: %s/%s", comp.ObjectMeta.Name, comp.Spec.Type, comp.Spec.Version)
	}
}

func (m *ComponentDefinitionManager) filterAuthorized(comps []components_v1alpha1.Component) []components_v1alpha1.Component {
	authorized := []components_v1alpha1.Component{}
	for _, c := range comps {
		if m.authorizationPolicy(&c) {
			authorized = append(authorized, c)
		}
	}
	return authorized
}

// ComponentLifecycleStatus contains the list of component lifecycles
type ComponentLifecycleStatus string

const (
	ComponentLifecycleStatusReady   ComponentLifecycleStatus = "ready"
	ComponentLifecycleStatusPending ComponentLifecycleStatus = "pending"
)

// ComponentLifecycleManager manages the component lifecycle
type ComponentLifecycleManager struct {
	logger    logger.Logger
	factory   *ComponentFactory
	ch        chan struct{}
	registry  *ComponentRegistryV2
	workQueue *ComponentDefinitionQueue
}

func (m *ComponentLifecycleManager) NewComponentLifecycleManager(
	logger logger.Logger,
	registry *ComponentRegistryV2) *ComponentLifecycleManager {

	return &ComponentLifecycleManager{
		logger:   logger,
		registry: registry,
		ch:       make(chan struct{}),
	}
}

func (m *ComponentLifecycleManager) Process(comp *components_v1alpha1.Component) {
	// enqueue the component for processing
	// the work queue takes care of de-duplication and will return false if the item was not queued
	if !m.workQueue.Enqueue(comp) {
		return
	}

	// async signal to the processor there is a new message to process
	go func() {
		m.ch <- struct{}{}
	}()
}

func (m *ComponentLifecycleManager) RunAsync() (doneChan chan struct{}, errChan chan error) {
	doneChan = make(chan struct{})
	errChan = make(chan error)
	go func() {
		defer close(doneChan)
		defer close(errChan)

		for {
			// receive message and check for closed channel
			_, open := <-m.ch
			if !open {
				// post that we are done to the done channel
				doneChan <- struct{}{}
				return
			}

			// process all components that we can
			// each new message send above may unblock multiple components in the queue
			for {
				component, exists := m.workQueue.Dequeue()
				if !exists {
					break
				}

				err := m.factory.Initialize(component)
				if err == nil {
					continue
				}

			}
		}
	}()
	return doneChan, errChan
}

func (m *ComponentLifecycleManager) Shutdown() {
	m.logger.Info("Shutting down all components")

}

type SecretStoreResolver interface {
	Get(name string) secretstores.SecretStore
}

// ComponentInterpolator takes a component definition and interpolates secret references
type ComponentInterpolator struct {
	resolver SecretStoreResolver
	mode     modes.DaprMode
	logger   logger.Logger
}

func NewComponentInterpolator(resolver SecretStoreResolver, mode modes.DaprMode, logger logger.Logger) *ComponentInterpolator {
	return &ComponentInterpolator{
		resolver: resolver,
		mode:     mode,
		logger:   logger,
	}
}

func (i *ComponentInterpolator) Interpolate(comp *components_v1alpha1.Component) (*components_v1alpha1.Component, error) {
	for index, m := range comp.Spec.Metadata {
		if m.SecretKeyRef.Name == "" {
			continue
		}

		secretStoreName := i.getSecretStoreName(m.SecretKeyRef.Name)

		if i.shouldDecodeSecrets(secretStoreName) {
			dynamicValue, err := i.decodeSecret(m.Value.Raw)
			if err != nil {
				i.logger.Errorf("error decoding secret: %s", err)
				continue
			}
			m.Value = *dynamicValue
			comp.Spec.Metadata[index] = m
		}

		secretStore := i.getSecretStore(secretStoreName)
		if secretStore == nil {
			i.logger.Errorf("missing secret store %s", secretStoreName)
			continue
		}

		response, err := secretStore.GetSecret(secretstores.GetSecretRequest{
			Name: m.SecretKeyRef.Name,
			Metadata: map[string]string{
				"namespace": comp.ObjectMeta.Namespace,
			},
		})
		if err != nil {
			i.logger.Errorf("error getting secret %s", err)
			continue
		}

		// Use the SecretKeyRef.Name key if SecretKeyRef.Key is not given
		secretKeyName := m.SecretKeyRef.Key
		if secretKeyName == "" {
			secretKeyName = m.SecretKeyRef.Name
		}

		val, ok := response.Data[secretKeyName]
		if !ok {
			continue
		}

		comp.Spec.Metadata[index].Value = components_v1alpha1.DynamicValue{
			JSON: v1.JSON{
				Raw: []byte(val),
			},
		}
	}
	return comp, nil
}

func (i *ComponentInterpolator) getSecretStore(name string) secretstores.SecretStore {
	if name == "" {
		return nil
	}
	return i.resolver.Get(name)
}

func (i *ComponentInterpolator) getSecretStoreName(name string) string {
	if name == "" && i.mode == modes.KubernetesMode {
		// there should be some constants somewhere to do these lookups
		return "kubernetes"
	}
	return name
}

func (i *ComponentInterpolator) shouldDecodeSecrets(name string) bool {
	return i.mode == modes.KubernetesMode && name == kubernetesSecretStore
}

func (i *ComponentInterpolator) decodeSecret(val []byte) (*components_v1alpha1.DynamicValue, error) {
	var jsonVal string
	err := json.Unmarshal(val, &jsonVal)
	if err != nil {
		return nil, err
	}

	dec, err := base64.StdEncoding.DecodeString(jsonVal)
	if err != nil {
		return nil, err
	}

	return &components_v1alpha1.DynamicValue{
		JSON: v1.JSON{
			Raw: dec,
		},
	}, nil
}

func NewSecretStoreCacheFacade(secretStore secretstores.SecretStore) secretstores.SecretStore {
	return &secretStoreCacheFacade{
		cache:       map[string]secretstores.GetSecretResponse{},
		secretStore: secretStore,
	}
}

type secretStoreCacheFacade struct {
	cache       map[string]secretstores.GetSecretResponse
	secretStore secretstores.SecretStore
}

func (s *secretStoreCacheFacade) Init(secretstores.Metadata) error {
	return nil
}

func (s *secretStoreCacheFacade) BulkGetSecret(request secretstores.BulkGetSecretRequest) (secretstores.BulkGetSecretResponse, error) {
	return s.secretStore.BulkGetSecret(request)
}

func (s *secretStoreCacheFacade) GetSecret(request secretstores.GetSecretRequest) (secretstores.GetSecretResponse, error) {
	resp, ok := s.cache[request.Name]
	if ok {
		return resp, nil
	}
	return s.secretStore.GetSecret(request)
}

// ComponentDefinitionQueue contains the queue of pending components to be processed by the component definition manager
type ComponentDefinitionQueue struct {
	lock  sync.RWMutex
	graph *Graph
}

func NewComponentDefinitionQueue() *ComponentDefinitionQueue {
	return &ComponentDefinitionQueue{
		graph: NewGraph(),
	}
}

// Enqueue adds teh component spec to the list work list. If the component is unchanged it returns without modification.
func (q *ComponentDefinitionQueue) Enqueue(component *components_v1alpha1.Component) bool {
	if component == nil || component.Spec.Type == "" || component.Name == "" {
		return false
	}

	// calculate the ID, this will be used to lookup items in the dependency graph
	ID := fmt.Sprintf("%s:%s", component.Spec.Type, component.Name)

	// if the item exists and is unchanged, exit
	oldComponent, exists := q.lookup(ID)
	if exists && reflect.DeepEqual(oldComponent.Spec.Metadata, component.Spec.Metadata) {
		return false
	}

	// create dependencies
	dependencies := []string{}
	if component.SecretStore != "" {
		dependencies = append(dependencies, component.SecretStore)
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	// update the component graph. This marks the component definition as unprocessed
	q.graph.Upsert(ID, component, dependencies...)
	return true
}

func (q *ComponentDefinitionQueue) lookup(ID string) (*components_v1alpha1.Component, bool) {
	q.lock.RLock()
	defer q.lock.RUnlock()

	instance, exists := q.graph.Lookup(ID)
	if !exists {
		return nil, false
	}
	return instance.(*components_v1alpha1.Component), true
}

// Dequeue removes the next item from the list that can be processed
func (q *ComponentDefinitionQueue) Dequeue() (*components_v1alpha1.Component, bool) {

	// locks access to the graph
	q.lock.Lock()

	// get the next node in the graph
	ID, exists := q.graph.Next()
	if !exists {
		q.lock.Unlock()
		return nil, false
	}
	q.lock.Unlock()

	component, exists := q.lookup(ID)
	if !exists {
		return nil, false
	}
	return component, true
}

type Node struct {
	Processed bool
	ID        string
	Data      interface{}
	Edges     map[string]*Node
}

type Graph struct {
	Nodes map[string]*Node
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: map[string]*Node{},
	}
}

// newNode creates a node with the given ID and registers it in the node list
func (g *Graph) newNode(ID string) *Node {
	n := &Node{
		ID:        ID,
		Edges:     map[string]*Node{},
		Processed: false,
	}
	g.Nodes[ID] = n
	return n
}

// Upsert inserts or updates the node with ID and adds dependencies. Upsert also marks the node as dirty
func (g *Graph) Upsert(ID string, data interface{}, dependencies ...string) {
	n, exists := g.Nodes[ID]
	if !exists {
		n = g.newNode(ID)
	}

	n.Data = data
	n.Processed = false

	for _, dep := range dependencies {
		d, nodeExists := g.Nodes[dep]
		if !nodeExists {
			d = g.newNode(dep)
		}
		n.Edges[dep] = d
	}
}

// Next returns the next unprocessed node without dependencies
// If the next node is found, Next() also returns true, false otherwise
func (g *Graph) Next() (string, bool) {
	found := false
	ID := ""
	for _, n := range g.Nodes {
		unprocessedEdges := false
		for _, e := range n.Edges {
			if !e.Processed {
				unprocessedEdges = true
			}
		}
		// skip if:
		// * there are unprocessed edges
		// * the current node has been processed
		// * the data is nil
		if unprocessedEdges || n.Processed || n.Data == nil {
			continue
		}
		ID = n.ID
		found = true
		n.Processed = true
		break
	}
	if !found {
		return "", false
	}
	for _, m := range g.Nodes {
		delete(m.Edges, ID)
	}
	return ID, found
}

// Lookup gets the data for the node. If the node exists it also returns true
func (g *Graph) Lookup(ID string) (data interface{}, exists bool) {
	n, exists := g.Nodes[ID]
	if !exists {
		return nil, false
	}
	return n.Data, true
}
