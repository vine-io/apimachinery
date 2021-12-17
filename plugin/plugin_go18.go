// +build go1.8,!windows,amd64,!static_build,!gccgo

// MIT License
//
// Copyright (c) 2021 Lack
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.


package plugin

import (
	"fmt"
	"path/filepath"
	"plugin"
	"runtime"
)

// loadPlugins loads all plugins for the OS and Arch
// that vine service is built for inside the provided path
func loadPlugins(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	pattern := filepath.Join(abs, fmt.Sprintf(
		"*-%s-%s.%s",
		runtime.GOOS,
		runtime.GOARCH,
		getLibExt(),
	))
	libs, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, lib := range libs {
		if _, err := plugin.Open(lib); err != nil {
			return err
		}
	}
	return nil
}

// getLibExt returns a platform specific lib extension for
// the platform that vine service is running on
func getLibExt() string {
	switch runtime.GOOS {
	case "windows":
		return "dll"
	default:
		return "so"
	}
}
