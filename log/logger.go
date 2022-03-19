package log

import (
	"fmt"
	"log"
	"os"
)

// TODO log 代码的重构  看看能不能把每次error日志打印两遍的问题解决一下

var (
	FeatureLogger, ProfileLogger,StorageLogger *log.Logger
)

func Init(fLogFile, pLogFile,sLogFile string) error {
	var err error

	if FeatureLogger, err = RegisterLogger(fLogFile); err != nil {
		return fmt.Errorf("Init feature logger err(%v)", err)
	}
	if ProfileLogger, err = RegisterLogger(pLogFile); err != nil {
		return fmt.Errorf("Init profile logger err(%v)", err)
	}
	if StorageLogger, err = RegisterLogger(sLogFile); err != nil {
		return fmt.Errorf("Init logger err(%v)", err)
	}

	return nil
}

func RegisterLogger(logFile string) (*log.Logger, error) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	logger := log.New(file, "", log.LstdFlags|log.Lshortfile)
	return logger, nil
}
