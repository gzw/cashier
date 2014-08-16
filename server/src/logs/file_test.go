package logs

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

func Testfile(t *testing.T) {
	fmt.Println("start--------------!")
	log := NewLogger(1000)
	log.SetLogger("file", "test.log")
	log.Trace("trace")
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

	// b := bufio.NewReader(f)
	// linenum := 0
	// for {
	// 	line, _, err := b.ReadLine()
	// 	if err != nil {
	// 		break
	// 	}
	// 	if len(line) > 0 {
	// 		linenum++
	// 	}
	// }
	// if linenum != 6 {
	// 	t.Fatal(linenum, "not line 6")
	// }
	// os.Remove("test.log")

}
