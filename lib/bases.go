package lib

import "math/big"

func Base10ToBase62(v int) string {
	bi := big.NewInt(int64(v))
	return bi.Text(62)
}
