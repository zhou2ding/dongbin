package filemanager

import (
	"blog/pkg/l"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
)

type file struct {
	homeDir string
}

func (c *file) writeWithoutPath(val *writeValue, data []byte) error {
	if err := pathExistsOrCreate(val.Dir); err != nil {
		l.GetLogger().Error("pathExistsOrCreate error", zap.Error(err))
		return err
	}

	fullName := path.Base(val.FileName)
	fullPath := val.Dir + fullName
	l.GetLogger().Info("writeWithoutPath", zap.String("full path", fullPath))
	if err := ioutil.WriteFile(fullPath, data, 0666); err != nil {
		l.GetLogger().Error("WriteFile error", zap.Error(err))
		return err
	}
	return nil
}
