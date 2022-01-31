package runtime_test

import (
	"testing"

	"github.com/dapr/dapr/pkg/apis/components/v1alpha1"
	"github.com/dapr/dapr/pkg/runtime"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createComponentDefinition(name, componentType string) *v1alpha1.Component {
	return &v1alpha1.Component{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Spec: v1alpha1.ComponentSpec{
			Type: componentType,
		},
	}
}

func TestComponentDefinitionQueueCanEnqueue(t *testing.T) {
	q := runtime.NewComponentDefinitionQueue()
	result := q.Enqueue(createComponentDefinition("name", "type"))
	require.True(t, result)
}

func TestComponentDefinitionQueueEnqueueFailsWhenSpecTypeOrNameMissing(t *testing.T) {
	q := runtime.NewComponentDefinitionQueue()
	result := q.Enqueue(createComponentDefinition("", ""))
	require.False(t, result)
}

func TestComponentDefinitionQueueCanDequeue(t *testing.T) {
	q := runtime.NewComponentDefinitionQueue()
	result := q.Enqueue(createComponentDefinition("name", "type"))
	require.True(t, result)
	c, result := q.Dequeue()
	require.True(t, result)
	require.NotNil(t, c)
}

func TestDependencyGraphCanInsert(t *testing.T) {
	g := runtime.NewGraph()
	require.NotNil(t, g)

	g.Upsert("test", "test")
	n, exists := g.Lookup("test")
	require.True(t, exists)
	require.NotNil(t, n)
	require.Equal(t, n, "test")
}

func TestDependencyGraphCanUpdate(t *testing.T) {
	g := runtime.NewGraph()

	g.Upsert("test", "test")
	n, exists := g.Lookup("test")
	require.True(t, exists)
	require.NotNil(t, n)
	require.Equal(t, n, "test")

	g.Upsert("test", "new")
	n, exists = g.Lookup("test")
	require.True(t, exists)
	require.NotNil(t, n)
	require.Equal(t, n, "new")

}

func TestDependencyGraphCanLookupData(t *testing.T) {
	g := runtime.NewGraph()
	g.Upsert("test", "data")
	data, exists := g.Lookup("test")
	require.True(t, exists)
	require.Equal(t, "data", data)
}

func TestDependencyGraphLookupReturnsFalseWhenMissing(t *testing.T) {
	g := runtime.NewGraph()
	g.Upsert("test", "data")
	_, exists := g.Lookup("other")
	require.False(t, exists)
}

func TestDependencyGraphNextReturnsFalseWhenEmpty(t *testing.T) {
	g := runtime.NewGraph()
	_, exists := g.Next()
	require.False(t, exists)
}

func TestDependencyGraphNextReturnsFalseWhenCircularDependency(t *testing.T) {
	g := runtime.NewGraph()
	g.Upsert("first", "first", "middle")
	g.Upsert("middle", "middle", "last")
	g.Upsert("last", "last", "first")
	_, exists := g.Next()
	require.False(t, exists)
}

func TestDependencyGraphNextReturnsTrueWhenHierarchy(t *testing.T) {
	g := runtime.NewGraph()
	g.Upsert("parent", "parent")
	g.Upsert("child", "child", "parent")
	id, exists := g.Next()
	require.True(t, exists)
	require.Equal(t, "parent", id)
	id, exists = g.Next()
	require.True(t, exists)
	require.Equal(t, "child", id)
}
