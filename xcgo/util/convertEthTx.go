package util

/*
#include<cgo.h>
*/
import "C"

////convert eth tx to Top tx
func ConvertEthTx(rawTx string) bool {
	if C.convertEthTx(C.CString(rawTx)) != 0 {
		return false
	}
	return true
}
