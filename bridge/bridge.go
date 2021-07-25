package bridge

// #cgo CFLAGS: -I. -I../WFA
// #cgo LDFLAGS: -L. -lwfabridge -L../WFA/build -lwfa -ljson-c
// #include <stdlib.h>
// #include <wfa_bridge.h>
import "C"
import (
	"encoding/json"
	"unsafe"
)

type Alignment struct {
	Reference string `json:"pattern_alg"`
	Sequence  string `json:"text_alg"`
	Ops       string `json:"ops_alg"`
	Score     int    `json:"score"`
	Opsn      []int  `json:"opsn"`
	Opsc      []byte `json:"opsc"`
}

func AffineWaveformAlign(sequenceA string, sequenceB string) (Alignment, error) {
	a := C.CString(sequenceA)
	defer C.free(unsafe.Pointer(a))

	b := C.CString(sequenceB)
	defer C.free(unsafe.Pointer(b))

	var alignment Alignment
	err := json.Unmarshal([]byte(C.GoString(C.align(a, b))), &alignment)
	return alignment, err
}
