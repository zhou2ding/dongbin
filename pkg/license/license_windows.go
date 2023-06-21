//go:build windows
// +build windows

package license

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
	"blog/pkg/l"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
	"syscall"
	"unsafe"
)

const vendorCode string = "zhoudongbin"

const feature = 51

type LicenseMgr struct {
	libPath string
}

var licenseMgr LicenseMgr
var once sync.Once

func InitLicenseMgr(v *viper.Viper, logger *zap.Logger) error {
	logger.Info("InitLicenseMgr windows")
	return GetLicenseMgrInstance().init(v, logger)
}

func UnInitLicenseMgr(logger *zap.Logger) {
	logger.Info("UnInitLicenseMgr windows")
}

func GetLicenseMgrInstance() *LicenseMgr {
	return &licenseMgr
}

func (mgr *LicenseMgr) init(v *viper.Viper, logger *zap.Logger) error {
	var err error = nil
	once.Do(func() {
		libPath := v.GetString("license.lib_path")
		if len(libPath) == 0 {
			logger.Info("lib path not set")
			err = errors.New("lib path not set")
		}
		mgr.libPath = libPath
	})

	return err
}

func (mgr *LicenseMgr) CheckValid() (bool, error) {
	l.Logger().Info("Check license windows...")
	handle, err := syscall.LoadDLL(mgr.libPath)
	if err != nil {
		return false, err
	}

	login, err := handle.FindProc("c_login")
	if err != nil {
		return false, err
	}

	logout, err := handle.FindProc("c_logout")
	if err != nil {
		return false, err
	}

	vcstr := C.CString(vendorCode)
	defer C.free(unsafe.Pointer(vcstr))
	var cHandle uint32
	ret, _, lastErr := login.Call(feature, uintptr(unsafe.Pointer(vcstr)), uintptr(unsafe.Pointer(&cHandle)))
	if lastErr != nil {
		l.Logger().Info("login.Call lastErr", zap.String("lastErr", lastErr.Error()))
	}
	if ret != 0 {
		l.Logger().Warn("login.Call failed")
		return false, nil
	}

	ret, _, lastErr = logout.Call(uintptr(cHandle))
	if lastErr != nil {
		l.Logger().Info("logout.Call lastErr", zap.String("lastErr", lastErr.Error()))
	}
	if ret != 0 {
		l.Logger().Warn("logout.Call failed")
	}

	l.Logger().Info("License is valid")

	return true, nil
}
