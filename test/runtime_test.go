package test

import (
	"fmt"
	"runtime"
	"testing"
)

func TestRuntime(t *testing.T) {
	fmt.Printf("numCPU: %v\n", runtime.NumCPU())
	fmt.Printf("numProcs: %v\n", runtime.GOMAXPROCS(0))
}

func TestName(t *testing.T) {
	fmt.Printf("%f\n", 5.28*3000*10000)
}
