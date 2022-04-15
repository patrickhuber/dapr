package env

import (
	"os"
	"strings"
)

type Env interface {
	Lookup(key string) (string, bool)
	Get(key string) string
	List() map[string]string
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

func (e *env) List() map[string]string {
	result := map[string]string{}
	list := os.Environ()
	for _, item := range list {
		splits := strings.SplitN(item, "=", 2)
		if len(splits) == 1 {
			result[splits[0]] = ""
		} else if len(splits) == 2 {
			result[splits[0]] = splits[1]
		} // skip if no match
	}
	return result
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
