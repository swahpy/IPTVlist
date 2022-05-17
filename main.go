package main

import (
	"github.com/Hywfred/IPTVlist/checker"
)

func main() {
	//checker.Request("http://dbiptv.sn.chinamobile.com/PLTV/88888890/224/3221220000/index.m3u8", 16)
	// _ = checker.ParseFromM3uAddr("https://raw.githubusercontent.com/BurningC4/Chinese-IPTV/master/TV-IPV4.m3u")

	// _ = checker.ParseFromM3uAddr("https://cdn.jsdelivr.net/gh/hywfred/IPTVlist@latest/docs/tvlist.m3u")
	// checker.CheckAllSync(src_list, 16)
	res := checker.ParseRange("http://zteres.sn.chinamobile.com:6060/yinhe/2/ch0000009099000000[0000-9999]/index.m3u8?virtualDomain=yinhe.live_hls.zte.com")
	checker.CheckAllSync(res, 16)
}
