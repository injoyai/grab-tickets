package main

import (
	"github.com/injoyai/goutil/net/stun"
)

/*
GetNetIP
通过stun服务器获取公网IP信息
*/
func GetNetIP() (string, error) {
	addr, err := stun.GetNetAddr()
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}
