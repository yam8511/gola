package helper

import (
	"gola/internal/bootstrap"
	"net/url"
	"strconv"
)

// 取服務網址
func ServiceURL(conf *bootstrap.ServiceConf) *url.URL {
	addr := conf.IP
	if conf.Port > 0 {
		addr += ":" + strconv.Itoa(conf.Port)
	}

	link := &url.URL{Host: addr}

	if conf.Secure {
		link.Scheme = "https"
	} else {
		link.Scheme = "http"
	}

	return link
}

// 解析網址或IP
func ParseUrlOrIP(addr string) (*url.URL, error) {
	link, err := url.Parse(addr)
	if err != nil {
		var err2 error
		link, err2 = url.Parse("http://" + addr)
		if err2 == nil {
			return link, nil
		}
	}
	return link, nil
}
