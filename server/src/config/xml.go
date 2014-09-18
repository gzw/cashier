//cashier

//@description cashier system
//@link https://github.com/gzw/cashier
//@authors gzw

package config

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"x2j"
)

// XmlConfig is a xml config parser and implements Config interface
type XMLConfig struct {
}

func (c *XMLConfig) test() {
	x2j.
}