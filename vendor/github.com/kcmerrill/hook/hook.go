package hook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"sync"
)

// Contains our triggers AND priority mapping
var triggers map[string]map[int][]interface{}

// Trigger lock, for those pesky race conditions
var tl *sync.Mutex

// Filter alias to trigger, for optics and code readability
var Filter func(string, ...interface{})

// Plugin alias to trigger, for optics and code readability
var Plugin func(string, ...interface{})

// Register a trigger with priority 0
func Register(trigger string, args ...interface{}) {
	tl.Lock()
	defer tl.Unlock()

	var f interface{}
	var priority int

	if len(args) >= 2 {
		priority, _ = args[0].(int)
		f = args[1]
	} else {
		priority = 0
		f = args[0]
	}

	str := args[len(args)-1]

	// if the last element you pass in is a string, and not a function, then lets default to exec
	if _, isString := str.(string); isString {
		f = func(i interface{}) {
			payload, err := json.Marshal(i)
			if err != nil {
				return
			}

			in := bytes.NewReader(payload)
			cmd := exec.Command("bash", "-c", str.(string))
			cmd.Stdin = in

			if results, err := cmd.Output(); err == nil {
				fmt.Println("results", string(results))
				json.Unmarshal(results, &i)
			}
		}
	}

	_, exists := triggers[trigger]
	if !exists {
		triggers[trigger] = make(map[int][]interface{})
		for p := 0; p <= 100; p++ {
			triggers[trigger][p] = make([]interface{}, 0)
		}
	}
	triggers[trigger][priority] = append(triggers[trigger][priority], f)
}

// Trigger a, uh, trigger
func Trigger(trigger string, args ...interface{}) {
	tl.Lock()
	priorities, exists := triggers[trigger]
	tl.Unlock()

	if exists {
		for p := 0; p <= 100; p++ {
			for _, f := range priorities[p] {
				params := make([]reflect.Value, len(args))
				for idx := range args {
					params[idx] = reflect.ValueOf(args[idx])
				}
				reflect.ValueOf(f).Call(params)
			}
		}
	}
}

// Giddy Up
func init() {
	Filter = Trigger
	Plugin = Trigger
	triggers = make(map[string]map[int][]interface{})
	tl = &sync.Mutex{}
}
