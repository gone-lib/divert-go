//go:build windows
// +build windows

package main

import (
	"fmt"

	"github.com/gone-lib/divert-go"
)

func main() {
	ver, err := divert.GetVersionInfo()
	if err != nil {
		panic(err)
	}
	fmt.Println(ver)
}
