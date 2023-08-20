package basefile

import (
	"github.com/pkg/errors"
	"os"
)

const (
	AutoPath = iota
	ManualPath
)

const (
	timeTemplateYMD = "2006-01-02"
	autoDir         = "BFQueryAuto\\"
	DisplayPrefix   = "display"
	LoopBackIP      = "127.0.0.1"
)

type writeValue struct {
	Mode     int
	Dir      string
	FileName string
	CreateAt int64
}

func pathExistsOrCreate(path string) error {
	fInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) { //create if not exist
			err = os.MkdirAll(path, 0666)
			if err != nil {
				return err
			} else {
				return nil
			}
		}
	}

	if !fInfo.IsDir() {
		return errors.Errorf("%s exists but not a directory", path)
	}

	return nil
}
