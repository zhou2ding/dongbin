package filemanager

import (
	"blog/pkg/v"
	"blog/pkg/l"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var (
	once     sync.Once
	instance *FileManager
)

type FileManager struct {
	file *file
}

func GetFileManager() *FileManager {
	once.Do(func() {
		instance = &FileManager{}
	})
	return instance
}

func (c *FileManager) Init() error {
	l.GetLogger().Info("Init FileManager")
	dir := v.GetViper().GetString("storage.homedir")
	if dir == "" {
		return errors.New("no home directory in configuration")
	}

	c.file = &file{
		homeDir: strings.Trim(dir, "\\") + "\\" + "fssdir" + "\\",
	}
	if err := pathExistsOrCreate(c.file.homeDir); err != nil {
		return err
	}
	l.GetLogger().Debug("init home dir", zap.String("homeDir", c.file.homeDir))
	return nil
}

func (c *FileManager) Read(path string) ([]byte, error) {
	l.GetLogger().Info("Read file start", zap.String("file name", path))
	content, err := ioutil.ReadFile(c.file.homeDir + autoDir + strings.TrimPrefix(path, "/"))
	if err != nil {
		if os.IsNotExist(err) {
			l.GetLogger().Warn("ReadFile find no result")
			return nil, nil
		}
		l.GetLogger().Error("ReadFile error", zap.Error(err))
		return nil, err
	}

	return content, nil
}

func (c *FileManager) Write(mode int, fileName string, createAt int64, fileData []byte) error {
	l.GetLogger().Info("Write file start")
	write := writeValue{
		Mode:     mode,
		Dir:      c.file.homeDir + autoDir,
		FileName: fileName,
		CreateAt: createAt,
	}
	switch write.Mode {
	case AutoPath:
		if err := c.file.writeWithoutPath(&write, fileData); err != nil {
			l.GetLogger().Error("writeWithoutPath failed", zap.Error(err))
			return err
		}
	}
	return nil
}
