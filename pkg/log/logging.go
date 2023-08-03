package log

import "fmt"

func Infof(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}

func Debugf(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}

func Errorf(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}

func Panicf(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}
