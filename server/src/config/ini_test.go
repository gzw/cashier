//cashier

//@description cashier system
//@link https://github.com/gzw/cashier
//@authors gzw

package config

import (
	"os"
	"testing"
)

var inicontext = `
;comment one
#comment two
appname = cashier
httpport = 8080
mysqlport = 3600
PI = 3.1415926
runmode = "dev"
autorender = false
copyrequestbody = true
[demo]
key1 = "gzw"
key2 = "gzw2"
CaseInsensitive = true
peers = one;two;three
`

func TestIni(t *testing.T) {
	f, err := os.Create("testini.conf")
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.WriteString(inicontext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testini.conf")
	iniconf, err := NewConfig("ini", "testini.conf")
	if err != nil {
		t.Fatal(err)
	}

	if iniconf.String("appname") != "cashier" {
		t.Fatal("appname not equal to cashier")
	}
	if port, err := iniconf.Int("httpport"); err != nil || port != 8080 {
		t.Error(port)
		t.Fatal(err)
	}
	if mysqlport, err := iniconf.Int("mysqlport"); err != nil || mysqlport != 3600 {
		t.Error(mysqlport)
		t.Fatal(err)
	}
	if pi, err := iniconf.Float("PI"); err != nil || pi != 3.1415926 {
		t.Error(pi)
		t.Fatal(err)
	}
	if runmode := iniconf.String("runmode"); runmode != "dev" {
		println(runmode)
		t.Error(runmode)
		t.Fatal("run mode is not equal to dev")
	}

	if v, err := iniconf.Bool("autorender"); err != nil || v != false {
		t.Error(v)
		t.Fatal(err)
	}
	if v, err := iniconf.Bool("copyrequestbody"); err != nil || v != true {
		t.Error(v)
		t.Fatal(err)
	}

	if err := iniconf.Set("name", "guozw"); err != nil {
		t.Fatal(err)
	}

	if iniconf.String("name") != "guozw" {
		t.Fatal("get name error")
	}

	if iniconf.String("demo::key1") != "gzw" {
		t.Fatal("get demo.key1 error")
	}
	if iniconf.String("demo::key2") != "gzw2" {
		t.Fatal("get demo.key2 error")
	}

	if v, err := iniconf.Bool("demo::CaseInsensitive"); err != nil || v != true {
		t.Error(v)
		t.Fatal(err)
	}
	data := iniconf.Strings("demo::peers")
	if len(data) != 3 {
		t.Fatal("get strings error")
	}
	if data[0] != "one" {
		t.Fatal("get first params error not equal to one")
	}
	if data[1] != "two" {
		t.Fatal("get second params error not equal to two")
	}
	if data[2] != "three" {
		t.Fatal("get third params error not equal to three")
	}

}
