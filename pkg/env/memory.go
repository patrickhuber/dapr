package env

type memory struct {
	variables map[string]string
}

func NewMemory() Env {
	return &memory{
		variables: make(map[string]string),
	}
}

func (m *memory) Get(key string) string {
	return m.variables[key]
}

func (m *memory) List() map[string]string {
	return m.variables
}

func (m *memory) Lookup(key string) (string, bool) {
	value, ok := m.variables[key]
	return value, ok
}

func (m *memory) Set(key string, value string) error {
	m.variables[key] = value
	return nil
}

func (m *memory) Unset(key string) error {
	delete(m.variables, key)
	return nil
}

func (m *memory) Clear() {
	for key := range m.variables {
		delete(m.variables, key)
	}
}
