package plugin

// Middleware is a HTTP middleware component definition.
type Plugin struct {
	Name          string
	FactoryMethod FactoryMethod
}

// Registry is the interface for callers to get registered HTTP middleware.
type Registry interface {
	Register(plugins ...Plugin)
	Create(name, version string) (Client, error)
}

type registry struct {
	clients map[string]FactoryMethod
}

// FactoryMethod is the method creating middleware from metadata.
type FactoryMethod func() (Client, error)

func NewLoader(name string, factory FactoryMethod) Plugin {
	return Plugin{
		Name:          name,
		FactoryMethod: factory,
	}
}

func NewRegistry() Registry {
	return &registry{
		clients: make(map[string]FactoryMethod),
	}
}

func (r *registry) Register(plugins ...Plugin) {
	for _, p := range plugins {
		r.clients[p.Name] = p.FactoryMethod
	}
}

func (r *registry) Create(name, version string) (Client, error) {
	return nil, nil
}
