package checker

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func init() {
	initLog()
	go writeToFile()
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

var strChan chan string

func writeToFile() {
	strChan = make(chan string)
	// create m3u file
	m3upath, err := os.Getwd()
	if err != nil {
		logrus.Fatalln(err)
	}
	m3uname := "iptv.m3u"
	path := filepath.Join(m3upath, m3uname)
	m3ufile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logrus.Fatalln(err)
	}
	for str := range strChan {
		_, err = m3ufile.WriteString(str + "\n")
		if err != nil {
			logrus.Fatalln(err)
		}
	}
	err = m3ufile.Close()
	if err != nil {
		logrus.Fatalln(err)
	}
}
