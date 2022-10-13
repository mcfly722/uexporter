package main

import (
	"io/fs"
	"io/ioutil"

	"github.com/dop251/goja"

	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/jsEngine"
)

// IOUtil ...
type IOUtil struct {
	context   context.Context
	eventLoop jsEngine.EventLoop
	runtime   *goja.Runtime
}

// Constructor ...
func (ioUtil IOUtil) Constructor(context context.Context, eventLoop jsEngine.EventLoop, runtime *goja.Runtime) {
	runtime.Set("IOUtil", &IOUtil{
		context:   context,
		eventLoop: eventLoop,
		runtime:   runtime,
	})
}

// ReadAll ...
func (ioUtil *IOUtil) ReadAll(name string) (string, error) {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ReadDir ...
func (ioUtil *IOUtil) ReadDir(name string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(name)
}
