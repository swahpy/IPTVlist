package main

import (
	"fmt"

	"github.com/Hywfred/IPTVlist/checker"
)

func main() {
	//checker.Request("http://dbiptv.sn.chinamobile.com/PLTV/88888890/224/3221220000/index.m3u8", 16)
	// _ = checker.ParseFromM3uAddr("https://raw.githubusercontent.com/BurningC4/Chinese-IPTV/master/TV-IPV4.m3u")

	// _ = checker.ParseFromM3uAddr("https://cdn.jsdelivr.net/gh/hywfred/IPTVlist@latest/docs/tvlist.m3u")
	// checker.CheckAllSync(src_list, 16)
	res := checker.ParseRange("http://113.200.58.252:9901/tsfile/live/[1000-2000]_1.m3u8")
	for _, v := range res {
		fmt.Println(v)
	}
}
