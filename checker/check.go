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

// HttpClient provides an http client with a timeout for users
func HttpClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout * time.Second,
	}
}

// FetchSource performs a get request for the given url and returns the body content
// @params url: target url
// @params timeout: timeout for request
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

// Do is method of struct urlChecker.
// It performs GET request for the url and saves the result to log file.
func (c *urlChecker) Do() {
	fmt.Println(c.url)
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
		strChan <- c.url
	} else {
		logrus.WithFields(logrus.Fields{
			"url":    c.url,
			"result": fmt.Sprint(resp.StatusCode),
		}).Errorln("Failed!")
	}
}

// CheckRangeSycn checks all the urls in the range given by url
// @params url: a url containing a range, through which we could get a sequence of urls within this range.
// The range must be a number range.
func CheckRangeSycn(url string, timeout time.Duration) {
	// generate all the urls according to the url range
	urls := ParseRange(url)
	// make use of CheckAllSync to complete this check
	CheckAllSync(urls, timeout)
}

// CheckAllSync checks all the urls simultaneously
// @params urls: all the url that need to check
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
