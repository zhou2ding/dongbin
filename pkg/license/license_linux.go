//go:build linux
// +build linux

package license

/*
#cgo CFLAGS: -I../../lib
#cgo LDFLAGS: -L../../lib -lc_linux_x86_64_37124
#include "c_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"blog/pkg/l"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"unsafe"
)

var vendorCode = "zhoudgonbin"

const feature = 51

type LicenseMgr struct {
}

var licenseMgr LicenseMgr

func InitLicenseMgr(v *viper.Viper, logger *zap.Logger) error {
	logger.Info("InitLicenseMgr linux")
	return nil
}

func UnInitLicenseMgr(logger *zap.Logger) {
	logger.Info("UnInitLicenseMgr linux")
}

func GetLicenseMgrInstance() *LicenseMgr {
	return &licenseMgr
}

func (mgr *LicenseMgr) CheckValid() (bool, error) {
	l.Logger().Info("Check license linux...")
	vcStr := C.CString(vendorCode)
	defer C.free(unsafe.Pointer(vcStr))

	var cHandle uint32
	status := C.c_login(
		feature,
		C.c_vendor_code_t(unsafe.Pointer(vcStr)),
		(*C.uint)(unsafe.Pointer(&cHandle)),
	)
	if status != 0 {
		l.Logger().Info("c_login failed")
		return false, nil
	}

	status = C.c_logout(
		C.uint(cHandle),
	)
	if status != 0 {
		l.Logger().Info("c_logout failed")
	}

	l.Logger().Info("License is valid")

	return true, nil
}
