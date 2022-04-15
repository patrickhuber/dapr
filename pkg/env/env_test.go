package env_test

import (
	"os"
	"testing"

	"github.com/dapr/dapr/pkg/env"
	"github.com/stretchr/testify/require"
)

func TestLookup(t *testing.T) {
	const testKey = "IS_SET"
	const testValue = "TEST"
	const testNotKey = "NOT_SET"

	os.Setenv(testKey, testValue)
	e := env.NewOS()
	t.Run("when is set", func(t *testing.T) {
		v, ok := e.Lookup(testKey)
		require.True(t, ok)
		require.Equal(t, testValue, v)
	})
	t.Run("when is not set", func(t *testing.T) {
		_, ok := e.Lookup(testNotKey)
		require.False(t, ok)
	})
}

func TestGet(t *testing.T) {
	const testKey = "IS_SET"
	const testValue = "TEST"
	const testNotKey = "NOT_SET"

	os.Setenv(testKey, testValue)
	e := env.NewOS()
	t.Run("when is set", func(t *testing.T) {
		v := e.Get(testKey)
		require.Equal(t, testValue, v)
	})
	t.Run("when not set", func(t *testing.T) {
		v := e.Get(testNotKey)
		require.Equal(t, "", v)
	})
}

func TestList(t *testing.T) {
	const testKey = "IS_SET"
	const testValue = "TEST"
	os.Setenv(testKey, testValue)
	e := env.NewOS()
	t.Run("when is set", func(t *testing.T) {
		match := false
		for k, v := range e.List() {
			if k == testKey {
				require.Equal(t, testValue, v)
				match = true
			}
		}
		require.True(t, match)
	})
}

func TestSet(t *testing.T) {
	const testKey = "IS_SET"
	const testValue = "TEST"
	e := env.NewOS()

	err := e.Set(testKey, testValue)
	require.Nil(t, err)

	v := e.Get(testKey)
	require.Equal(t, testValue, v)
}

func TestUnset(t *testing.T) {
	const testKey = "IS_SET"
	const testValue = "TEST"
	e := env.NewOS()

	err := e.Set(testKey, testValue)
	require.Nil(t, err)

	v := e.Get(testKey)
	require.Equal(t, testValue, v)

	err = e.Unset(testKey)
	require.Nil(t, err)

	v = e.Get(testKey)
	require.Equal(t, "", v)

}

func TestClear(t *testing.T) {

	const testKey = "IS_SET"
	e := env.NewOS()
	// reset env to avoid bleed over into other tests
	defer func(original map[string]string) {
		for k, v := range original {
			e.Set(k, v)
		}
	}(e.List())
	e.Clear()
	v := e.Get(testKey)
	require.Equal(t, "", v)
}
