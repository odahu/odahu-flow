package feedback

import (
	"go.uber.org/zap"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type LoggerEngine interface {
	Post(tag string, message interface{}) error
	Close() error
}

type Logger struct {
	logEngine LoggerEngine
}

type Wrapper func(content string) (wrappedContent interface{})

func NewLogger(engine LoggerEngine) Logger {
	return Logger{
		logEngine: engine,
	}
}


// LogDir logs content of each .json file in `path` but wrap it content using
// Wrapper before.
// `tag` is used to add extra information to LoggerEngine
func (l Logger) LogDir(path string, tag string, wrap Wrapper) error {
	items, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, item := range items {
		file := item.Name()
		if item.IsDir() {
			zap.S().Infof("%s is directory (not file). Skip logging", file)
			continue
		}
		if strings.HasSuffix(file, ".json") {
			fp := filepath.Join(path, file)
			data, err := ioutil.ReadFile(fp)
			if err != nil {
				return err
			}

			logContent := wrap(string(data))

			if err := l.logEngine.Post(tag, logContent); err != nil {
				zap.S().Errorw("Error during logging", zap.Error(err))
			} else {
				zap.S().Infof("%s successfully logged", file)
			}
		} else {
			zap.S().Infof("%s has not .json extension. Skip logging", file)
		}
	}
	return nil
}
