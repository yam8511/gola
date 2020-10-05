package helper

import (
	"gola/internal/logger"
	"net"
)

// 本機可用IP
func LocalAvailableIPs() (availableIPs []string) {
	availableIPs = []string{}

	cfgs, err := net.Interfaces()
	if err != nil {
		logger.Warn("取本機取網卡資訊失敗: " + err.Error())
		return
	}

	for _, cfg := range cfgs {
		ips, err := cfg.Addrs()
		if err != nil {
			logger.Warn("取本機取網卡資訊失敗: " + err.Error())
			continue
		}

		for _, ip := range ips {
			iip, _, err := net.ParseCIDR(ip.String())
			if err != nil {
				logger.Warn("取本機取網卡資訊失敗: " + err.Error())
				continue
			}
			if iip.String() != "" {
				availableIPs = append(availableIPs, iip.String())
			}
		}
	}

	return availableIPs
}
