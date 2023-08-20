package basefile

import (
	"blog/pkg/l"
	"go.uber.org/zap"
	"os"
	"path"
)

type file struct {
	homeDir string
}

func (c *file) writeWithoutPath(val *writeValue, data []byte) error {
	if err := pathExistsOrCreate(val.Dir); err != nil {
		l.Logger().Error("pathExistsOrCreate error", zap.Error(err))
		return err
	}

	fullName := path.Base(val.FileName)
	fullPath := val.Dir + fullName
	l.Logger().Info("writeWithoutPath", zap.String("full path", fullPath))
	if err := os.WriteFile(fullPath, data, 0666); err != nil {
		l.Logger().Error("WriteFile error", zap.Error(err))
		return err
	}
	return nil
}
