//cashier

//@description cashier system
//@link https://github.com/gzw/cashier
//@authors gzw

package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// JosonConfig is a json config parser and implements Config interface
type JosonConfig struct {
}

// Parse returns a ConfigContainer with parsed json config map
func (js *JosonConfig) Parse(filename string) (ConfigContainer, errors) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	x := &JsonConfigContainer{
		data:make(map(string)interface{})
	}

	content , err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &x.data)
	if err != nil {
		var wrappingArray []interface{}
		err2 := json.Unmarshal(content, &wrappingArray)
		if err2 != nil {
			return nil, err
		}
		x.data["rootArray"] = wrappingArray
	}
	return x, nil
}

type JsonConfigContainer struct {
	data map[string]interface{}
	sync.RWMutex
}


// section.key or key
func (c *JsonConfigContainer) getdata(key string) interface{} {
	c.RLock()
	defer c.RUnlock()
	if len(key) == 0 {
		return nil
	}
	sectionkey := strings.Split(key, "::")
	if len(sectionkey) >= 2 {
		cruval, ok := c.data[sectionkey[0]]
		if !ok {
			return nil
		}
		for _, key := range sectionkey[1:] {
			if v, ok := cruval.(map[string]interface{}); !ok {
				return nil
			} else if cruval, ok = v[key]; !ok {
				return nil
			}
		}
		return cruval
	} else {
		if v, ok := c.data[key]; ok {
			return v
		}
	}
	return nil
}

// section.key or key
// func (c *JsonConfigContainer) getdata(key string) interface{} {
// 	c.RLock()
// 	defer c.RUnlock()
// 	if len(key) == 0 {
// 		return nil
// 	}
// 	sectionkey := strings.Split(key, "::")
// 	if len(sectionkey) >= 2 {
// 		cruval, ok := c.data[sectionkey[0]]
// 		if !ok {
// 			return nil
// 		}
// 		for _, key := range sectionkey[1:] {
// 			if v, ok := cruval.(map[string]interface{}); !ok {
// 				return nil
// 			} else if cruval, ok = v[key]; !ok {
// 				return nil
// 			}
// 		}
// 		return cruval
// 	} else {
// 		if v, ok := c.data[key]; ok {
// 			return v
// 		}
// 	}
// 	return nil
// }

