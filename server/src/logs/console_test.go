package logs

import (
	"fmt"
	"testing"
)

func TestConsole(t *testing.T) {

	log := NewLogger(100)
	log.EnableFuncCallDepth(true)
	err := log.SetLogger("console", "")
	if err != nil {
		fmt.Println("log.SetLogger: %q", err.Error())
	}
	log.Info("Info")
	log.Warn("Warn")
	log.Debug("Debug")
	log.Error("Error")
	log.Critical("Critical")

	log2 := NewLogger(10000)
	log2.SetLogger("console", `{"level":1}`)
	log2.EnableFuncCallDepth(true)
	log2.Trace("trace")
	log2.Info("info")
	log2.Warn("warning")
	log2.Debug("debug")
	log2.Critical("critical")
}

func BenchmarkConsole(b *testing.B) {
	log := NewLogger(10000)
	log.EnableFuncCallDepth(true)
	log.SetLogger("console", "")
	for i := 0; i < b.N; i++ {
		log.Trace("trace")
	}
}
