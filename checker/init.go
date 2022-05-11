package checker

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func init() {
	initLog()
}

func initLog() {
	logpath, err := os.Getwd()
	if err != nil {
		logrus.Fatalln(err)
	}
	logname := "iptv.log"
	path := filepath.Join(logpath, logname)
	logfile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logrus.Fatalln(err)
	}
	logrus.RegisterExitHandler(func() {
		if logfile == nil {
			return
		}
		logfile.Close()
	})
	logrus.SetOutput(logfile)
}
