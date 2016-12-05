// +build debug

package dcc

import "fmt"

func debug(f string, args ...interface{}) {
	fmt.Printf(f, args...)
}
