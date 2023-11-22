//go:build linux || darwin

package pe

import (
	"fmt"
)

func Hi(p func([]byte) ([]byte, error)) error {
	fmt.Println("Hello World!")
	return nil
}
