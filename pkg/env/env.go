package env

import "os"

type Env interface {
	Lookup(key string) (string, bool)
	Get(key string) string
	List() []string
	Set(key string, value string) error
	Unset(key string) error
	Clear()
}

type env struct {
}

func NewOS() Env {
	return &env{}
}

func (e *env) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (e *env) Get(key string) string {
	return os.Getenv(key)
}

func (e *env) List() []string {
	return os.Environ()
}

func (e *env) Set(key, value string) error {
	return os.Setenv(key, value)
}

func (e *env) Unset(key string) error {
	return os.Unsetenv(key)
}

func (e *env) Clear() {
	os.Clearenv()
}
