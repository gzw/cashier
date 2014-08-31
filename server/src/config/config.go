//cashier

//@description cashier system
//@link https://github.com/gzw/cashier
//@authors gzw
package config

import (
	"fmt"
)

type ConfigContainer interface {
	Set(key, val string) error // suport section::key type in given key when using ini type
	String(key string) string
	Strings(key string) []string
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	DIY(key string) (interface{}, error)
}

type Config interface {
	Parse(key string) (ConfigContainer, error)
}

var adapters = make(map[string]Config)

func Register(name string, adapter Config) {
	if adapter == nil {
		panic("Config: Register adapter is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("Config: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

func NewConfig(adapterName, filename string) (ConfigContainer, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("config: unkownn adaptername %q (forgtten import?)", adapterName)
	}
	return adapter.Parse(filename)
}
