package tool

import (
	"github.com/valyala/gozstd"
)

func Compress(rowData []byte) []byte {
	res := make([]byte, len(rowData))
	return gozstd.Compress(res, rowData)
}

func Decompress(rowData []byte) (res []byte, err error) {
	res = make([]byte, len(rowData)*2)
	return gozstd.Decompress(res, rowData)
}
