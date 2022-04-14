package standalone_test

import (
	"fmt"
	"io/fs"
	"path"
	"runtime"
	"testing"
	"testing/fstest"

	"github.com/dapr/dapr/pkg/plugin/standalone"
	"github.com/stretchr/testify/require"
)

func TestCanMatchRuntimeContext(t *testing.T) {
	extensions := []string{".jar", ".exe", "", ".dll", ".py", ".js"}
	name := "dapr-gomemory-v0.0.1%s"
	mockFS := fstest.MapFS{}

	for _, ext := range extensions {
		fileNameWithExtension := fmt.Sprintf(name, ext)
		file := &fstest.MapFile{
			Data: []byte(""),
		}
		mockFS[fileNameWithExtension] = file
	}

	t.Run("can match java", func(t *testing.T) {
		CanMatchRuntimeContext(t, mockFS, ".jar", standalone.RuntimeJava)
	})

	t.Run("can match dotnet", func(t *testing.T) {
		CanMatchRuntimeContext(t, mockFS, ".dll", standalone.RuntimeDotnet)
	})

	t.Run("can match python", func(t *testing.T) {
		CanMatchRuntimeContext(t, mockFS, ".py", standalone.RuntimePython)
	})

	t.Run("can match javascript", func(t *testing.T) {
		CanMatchRuntimeContext(t, mockFS, ".js", standalone.RuntimeNode)
	})

	t.Run("can match exec", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			CanMatchRuntimeContext(t, mockFS, ".exe", standalone.RuntimeExec)
		} else {
			CanMatchRuntimeContext(t, mockFS, "", standalone.RuntimeExec)
		}
	})
}

func CanMatchRuntimeContext(t *testing.T, filesystem fs.FS, ext string, name standalone.Runtime) {
	dirList, err := fs.ReadDir(filesystem, ".")
	files := []string{}
	// filter out directories and files that don't have the same extension
	for _, file := range dirList {
		if file.IsDir() {
			continue
		}

		extension := path.Ext(file.Name())
		if extension == ext {
			files = append(files, file.Name())
		} else if ext == "" && extension == ".1" {
			files = append(files, file.Name())
		}
	}
	require.Nil(t, err)
	require.Equal(t, 1, len(files))

	file := files[0]
	contexts := standalone.MatchRuntimeContext(file)
	require.Equal(t, 1, len(contexts))

	context := contexts[0]
	require.Equal(t, name, context.Name())
}
