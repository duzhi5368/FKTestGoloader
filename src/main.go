package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"unsafe"

	"github.com/dearplain/goloader"
)

// LoadPlugins Load all plugin in path
func loadPlugins(path string) (functions map[string]*func(interface{}, map[string]interface{}) (interface{}, error)) {

	functions = make(map[string]*func(interface{}, map[string]interface{}) (interface{}, error))
	// Fix Path
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	// Read path
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Load Plugin
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".o") {
			var runFuncName = f.Name()
			if strings.LastIndexAny(runFuncName, ".") > -1 {
				runFuncName = runFuncName[0:strings.LastIndexAny(runFuncName, ".")]
			}

			file, err := os.Open(path + f.Name())
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			reloc, err := goloader.ReadObj(file)

			if err != nil {
				fmt.Println(err)
				return
			}
			symPtr := loaderRegister()
			codeModule, err := goloader.Load(reloc, *symPtr)
			if err != nil {
				fmt.Println("Load error:", err)
			}

			runFuncPtr := codeModule.Syms["main."+runFuncName]
			funcPtrContainer := (uintptr)(unsafe.Pointer(&runFuncPtr))
			runFunc := (*func(interface{}, map[string]interface{}) (interface{}, error))(unsafe.Pointer(&funcPtrContainer))
			functions[runFuncName] = runFunc
		}
	}

	return
}

func loaderRegister() *map[string]uintptr {
	symPtr := make(map[string]uintptr)
	goloader.RegSymbol(symPtr)

	// Register
	w := sync.WaitGroup{}
	rw := sync.RWMutex{}
	goloader.RegTypes(symPtr, &w, w.Wait, &rw)

	// For request
	httpClient := new(http.Client)
	goloader.RegTypes(symPtr, strings.NewReader, http.NewRequest, httpClient, ioutil.ReadAll, httpClient.Do)
	return &symPtr
}

func main() {
	functions := loadPlugins("./plugin/")

	var vars = make(map[string]interface{})
	vars["method"] = "GET"
	vars["url"] = "https://www.google.com/"
	vars["data"] = ""
	b, _ := (*functions["request"])(nil, vars)
	if b != nil {
		fmt.Println(string(b.([]byte)))
	} else {
		fmt.Println("plugin call is successed, but request failed.")
	}
}
