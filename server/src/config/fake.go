//cashier

//@description cashier system
//@link https://github.com/gzw/cashier
//@authors gzw

package config

import (
	"errors"
	"strconv"
	"strings"
)

type fakeConfigContainer struct {
	data map[string]string
}

func (c *fakeConfigContainer) getdata(key string) string {
	key = strings.ToLower(key)
	return c.data[key]
}

func (c *fakeConfigContainer) Set(key, val string) error {
	key = strings.ToLower(key)
	c.data[key] = val
	return nil
}

func (c *fakeConfigContainer) String(key string) string {
	return c.getdata(key)
}

func (c *fakeConfigContainer) Strings(key string) []string {
	return strings.Split(c.getdata(key), ";")
}

func (c *fakeConfigContainer) Int(key string) (int, error) {
	return strconv.Atoi(c.getdata(key))
}

func (c *fakeConfigContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.getdata(key), 10, 64)
}

func (c *fakeConfigContainer) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.getdata(key))
}

func (c *fakeConfigContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(key, 64)
}

func (c *fakeConfigContainer) DIY(key string) (interface{}, error) {
	key = strings.ToLower(key)
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return nil, errors.New("key not found")
}

var _ ConfigContainer = new(fakeConfigContainer)

func NewFakeConfig() ConfigContainer {
	return &fakeConfigContainer{
		data: make(map[string]string),
	}
}
