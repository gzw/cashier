package logs

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestFile(t *testing.T) {
	fmt.Println("start TestFile--------------!")
	log := NewLogger(1000)
	log.SetLogger("file", `{"filename":"test.log"}`)
	log.Trace("trace")
	fmt.Println("1")
	log.Info("info")
	log.Debug("debug")
	log.Warn("warn")
	log.Critical("Critical")
	log.Error("error")

	time.Sleep(time.Second * 4)
	f, err := os.Open("test.log")
	if err != nil {
		t.Fatal(err)
	}
	b := bufio.NewReader(f)
	linenum := 0

	for {
		line, _, err := b.ReadLine()
		if err != nil {
			break
		}
		if len(line) > 0 {
			linenum++
		}
		fmt.Println(string(line))
	}

	// if linenum != 5 {
	// 	t.Fatal(linenum, "not line 6")
	// }
	os.Remove("test.go")
}

func TestFileRotate(t *testing.T) {
	log := NewLogger(10000)
	log.SetLogger("file", `{"filename":"test3.log","maxlines":4}`)
	log.Trace("test")
	log.Info("info")
	log.Debug("debug")
	log.Warn("warning")
	log.Error("error")
	log.Critical("critical")
	time.Sleep(time.Second * 4)
}

func BenchmarkFile(b *testing.B) {
	log := NewLogger(100000)
	log.SetLogger("file", `{"filename":"test4.log"}`)
	for i := 0; i < b.N; i++ {
		log.Trace("trace")
	}
	os.Remove("test4.log")
}
