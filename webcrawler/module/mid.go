package module

import (
	"fmt"
	"github.com/summerKK/go-code-snippet-library/webcrawler/errors"
	"github.com/summerKK/go-code-snippet-library/webcrawler/toolkit/sn"
	"net"
	"strconv"
	"strings"
)

var DefaultSNGen = sn.NewGenerator(1, 0)

// 包名+Sn[|host:port]
var midTemplate = "%s%d|%s"

func GenMid(mtype MType, sn uint64, addr net.Addr) (mid MID, err error) {
	if !Legalletter(mtype) {
		err = fmt.Errorf("illegal module letter:%v\n", err)
		return
	}
	letter := legalletterMap[mtype]
	var midStr string
	if addr == nil {
		midStr = fmt.Sprintf(midTemplate, letter, sn, "")
		midStr = midStr[:len(midStr)-1]
	} else {
		midStr = fmt.Sprintf(midTemplate, letter, sn, addr.String())
	}

	mid = MID(midStr)
	return
}

func SplitMid(mid MID) (s []string, err error) {
	midStr := string(mid)
	if len(midStr) <= 1 {
		err = errors.NewIllegalParamsError("insufficient MID")
		return
	}

	letter := MType(midStr[:1])
	if _, ok := legalletterMap[letter]; !ok {
		err = errors.NewIllegalParamsError(
			fmt.Sprintf("illegal module type letter: %s", letter))
		return
	}

	var snStr string
	var addr string
	snAndAddr := midStr[1:]
	index := strings.LastIndex(snAndAddr, "|")
	if index < 0 {
		snStr = snAndAddr
		// 如果index < 0,表示不存在addr(D12345).这里要判断midStr[1:]后面是否是个合法数字
		if !legalSN(snStr) {
			return nil, errors.NewIllegalParamsError(
				fmt.Sprintf("illegal module SN: %s", snStr))
		}
	} else {
		snStr = snAndAddr[:index]
		// 如果index < 0,表示不存在addr(D12345).这里要判断midStr[1:]后面是否是个合法数字
		if !legalSN(snStr) {
			return nil, errors.NewIllegalParamsError(
				fmt.Sprintf("illegal module SN: %s", snStr))
		}
		// 去掉 `|`
		addr = snAndAddr[index+1:]
		index = strings.LastIndex(addr, ":")
		if index <= 0 {
			return nil, errors.NewIllegalParamsError(
				fmt.Sprintf("illegal module address: %s", addr))
		}
		ipStr := addr[:index]
		if ip := net.ParseIP(ipStr); ip == nil {
			return nil, errors.NewIllegalParamsError(
				fmt.Sprintf("illegal module IP: %s", ipStr))
		}
		portStr := addr[index+1:]
		if _, err := strconv.ParseUint(portStr, 10, 64); err != nil {
			return nil, errors.NewIllegalParamsError(
				fmt.Sprintf("illegal module port: %s", portStr))
		}
	}
	s = []string{string(letter), snStr, addr}

	return
}

// legalSN 用于判断序列号的合法性。
func legalSN(snStr string) bool {
	_, err := strconv.ParseUint(snStr, 10, 64)
	if err != nil {
		return false
	}
	return true
}
