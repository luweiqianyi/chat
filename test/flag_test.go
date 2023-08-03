package test

import (
	"flag"
	"fmt"
	"testing"
	"time"
)

var fs flag.FlagSet

type Config struct {
	property1 int
	property2 string
	property3 time.Duration
}

// TestFlag 使用flag包为变量赋默认值
func TestFlag(t *testing.T) {
	var cfg Config

	fs.IntVar(&cfg.property1, "property1", 666, "usage: show")
	fs.StringVar(&cfg.property2, "property2", "666", "usage: show")
	fs.DurationVar(&cfg.property3, "property3", time.Duration(666), "usage: show")
	fmt.Printf("%#v\n", cfg)
}

func TestTimeFormat(t *testing.T) {
	fmt.Printf("time format: %v\n", time.Now().String())
}
