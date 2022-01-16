package plugin

import (
	"os/exec"
	"runtime"
)

type Runtime string

const (
	RuntimeDefault Runtime = RuntimeExec
	RuntimePython  Runtime = "python"
	RuntimeDotnet  Runtime = "dotnet"
	RuntimeNode    Runtime = "node"
	RuntimeExec    Runtime = "exec"
	RuntimeJava    Runtime = "java"
)

type RuntimeContext interface {
	Name() Runtime
	Extensions() []string
	Executable() string
	Args() []string
	Command(path string) *exec.Cmd
}

type runtimeContext struct {
	name       Runtime
	extensions []string
	executable string
	args       []string
}

func (r *runtimeContext) Name() Runtime {
	return r.name
}

func (r *runtimeContext) Extensions() []string {
	return r.extensions
}

func (r *runtimeContext) Executable() string {
	return r.executable
}

func (r *runtimeContext) Args() []string {
	return r.args
}

func (r *runtimeContext) Command(path string) *exec.Cmd {
	if r.executable == "" {
		return exec.Command(path, r.args...)
	}
	args := append(r.args, path)
	return exec.Command(r.executable, args...)
}

func NewDotnet() RuntimeContext {
	return &runtimeContext{
		name:       RuntimeDotnet,
		extensions: []string{".dll"},
		executable: "dotnet",
		args:       []string{"run"},
	}
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}
func NewExec() RuntimeContext {
	extensions := []string{}
	if isWindows() {
		extensions = append(extensions, ".exe")
	}
	return &runtimeContext{
		name:       RuntimeExec,
		extensions: extensions,
		executable: "",
		args:       []string{},
	}
}

func NewPython() RuntimeContext {
	return &runtimeContext{
		name:       RuntimePython,
		extensions: []string{".py"},
		executable: "python3",
		args:       []string{},
	}
}

func NewNode() RuntimeContext {
	return &runtimeContext{
		name:       RuntimeNode,
		extensions: []string{".js"},
		executable: "node",
		args:       []string{},
	}
}

func NewJava() RuntimeContext {
	return &runtimeContext{
		name:       RuntimeJava,
		extensions: []string{".jar"},
		executable: "java",
		args:       []string{"-jar"},
	}
}

var runtimeContextMap = map[Runtime]RuntimeContext{
	RuntimeDotnet: NewDotnet(),
	RuntimeExec:   NewExec(),
	RuntimeJava:   NewJava(),
	RuntimeNode:   NewNode(),
	RuntimePython: NewPython(),
}

func GetRuntimeContext(name Runtime) RuntimeContext {
	ctx, ok := runtimeContextMap[name]
	if ok {
		return ctx
	}
	return runtimeContextMap[RuntimeDefault]
}

func GetRuntimeContextFromString(name string) RuntimeContext {
	runtime := GetRuntime(name)
	return GetRuntimeContext(runtime)
}

func GetRuntime(name string) Runtime {
	runtime := Runtime(name)
	switch runtime {
	case RuntimeDotnet, RuntimeExec, RuntimeJava, RuntimeNode, RuntimePython:
		return runtime
	default:
		return RuntimeDefault
	}
}
