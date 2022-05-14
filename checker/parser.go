package checker

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
)

type channel struct {
	name        string
	tvg_id      string
	tvg_name    string
	tvg_logo    string
	group_title string
	url         string
}

var channels []channel

// parse m3u file and return IPTV source list
func ParseFromM3uAddr(m3u string) error {
	// get all the sources info
	err, sources := FetchSource(m3u, 16)
	if err != nil {
		return err
	}
	// parse sources to list
	reg, err := regexp.Compile(`#EXTINF:-1\s*(tvg-id="(?P<tvg_id>.*?)")?\s*(tvg-name="(?P<tvg_name>.*?)")?\s*(tvg-logo="(?P<tvg_logo>.*?)")?\s*(group-title="(?P<group_title>.*?)")?.*,(?P<name>.*)\s+(?P<url>.*)`)
	if err != nil {
		return err
	}
	res := reg.FindAllStringSubmatch(sources, -1)
	channels = make([]channel, len(res))
	for i, v := range res {
		channels[i].tvg_id = v[reg.SubexpIndex("tvg_id")]
		channels[i].tvg_name = v[reg.SubexpIndex("tvg_name")]
		channels[i].tvg_logo = v[reg.SubexpIndex("tvg_logo")]
		channels[i].group_title = v[reg.SubexpIndex("group_title")]
		channels[i].name = v[reg.SubexpIndex("name")]
		channels[i].url = v[reg.SubexpIndex("url")]
	}
	fmt.Printf("%#v\n", channels[0])
	return nil
}

// ParseRange parses the given url which has a range mark to get the url series
func ParseRange(url string) []string {
	reg, err := regexp.Compile(`(?P<prefix>.*)\[(?P<low>\d+)\-(?P<high>\d+)\](?P<suffix>.*)`)
	if err != nil {
		fmt.Println("regexp compile error: ", err)
		logrus.Fatalln(err)
	}
	res := reg.FindAllStringSubmatch(url, -1)
	if res == nil {
		fmt.Println("No matches found: ", url)
		logrus.Fatalln("No matches found: ", url)
	}
	prefix := res[0][reg.SubexpIndex("prefix")]
	suffix := res[0][reg.SubexpIndex("suffix")]
	lowstr := res[0][reg.SubexpIndex("low")]
	highstr := res[0][reg.SubexpIndex("high")]
	low, err := strconv.Atoi(lowstr)
	if err != nil {
		fmt.Printf("%s unable to be converted to int.\n", lowstr)
		logrus.WithFields(logrus.Fields{
			"object": lowstr,
		}).Fatalln(err)
	}
	high, err := strconv.Atoi(highstr)
	if err != nil {
		fmt.Printf("%s unable to be converted to int.\n", lowstr)
		logrus.WithFields(logrus.Fields{
			"object": lowstr,
		}).Fatalln(err)
	}
	result := make([]string, high-low+1)
	format := "%s%0" + fmt.Sprint(len(highstr)) + "d%s"
	for i := low; i <= high; i++ {
		result[i-low] = fmt.Sprintf(format, prefix, i, suffix)
	}
	return result
}
