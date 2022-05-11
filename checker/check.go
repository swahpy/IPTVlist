package checker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func HttpClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout * time.Second,
	}
}

func FetchSource(url string, timeout time.Duration) (error, string) {
	client := HttpClient(timeout)
	resp, err := client.Get(url)
	if err != nil {
		return err, ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, ""
	}
	return nil, string(body)
}

func Check(url string, c *http.Client) (error, int) {
	resp, err := c.Get(url)
	if err != nil {
		return err, -1
	}
	return nil, resp.StatusCode
}

type urlChecker struct {
	url    string
	client *http.Client
}

func (c *urlChecker) Do() {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithFields(logrus.Fields{
				"url": c.url,
			}).Errorln(err)
		}
	}()
	resp, err := c.client.Get(c.url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url": c.url,
		}).Errorln(err)
	}
	if resp.StatusCode == http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"url": c.url,
		}).Infoln("Passed!")
	} else {
		logrus.WithFields(logrus.Fields{
			"url":    c.url,
			"result": fmt.Sprint(resp.StatusCode),
		}).Errorln("Failed!")
	}
}

func CheckRangeSycn(url string, timeout time.Duration) {

}

func CheckAllSync(urls []string, timeout time.Duration) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	client := HttpClient(timeout)
	var wg sync.WaitGroup
	wg.Add(len(urls))
	p := New(runtime.NumCPU())
	for _, url := range urls {
		uc := urlChecker{
			url:    url,
			client: client,
		}
		go func() {
			p.Run(&uc)
			wg.Done()
		}()
	}
	wg.Wait()
	p.Shutdown()
}
