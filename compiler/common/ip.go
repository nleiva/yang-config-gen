package common

import (
	"fmt"
	"net/netip"
)

var Masks map[string]string = map[string]string{
	"8":  "255.0.0.0",
	"9":  "255.128.0.0",
	"10": "255.192.0.0",
	"11": "255.224.0.0",
	"12": "255.240.0.0",
	"13": "255.248.0.0",
	"14": "255.252.0.0",
	"15": "255.254.0.0",
	"16": "255.255.0.0",
	"17": "255.255.128.0",
	"18": "255.255.192.0",
	"19": "255.255.224.0",
	"20": "255.255.240.0",
	"21": "255.255.248.0",
	"22": "255.255.252.0",
	"23": "255.255.254.0",
	"24": "255.255.255.0",
	"25": "255.255.255.128",
	"26": "255.255.255.192",
	"27": "255.255.255.224",
	"28": "255.255.255.240",
	"29": "255.255.255.248",
	"30": "255.255.255.252",
	"31": "255.255.255.254",
	"32": "255.255.255.255",
}

var IntMasks map[int]string = map[int]string{
	8:  "255.0.0.0",
	9:  "255.128.0.0",
	10: "255.192.0.0",
	11: "255.224.0.0",
	12: "255.240.0.0",
	13: "255.248.0.0",
	14: "255.252.0.0",
	15: "255.254.0.0",
	16: "255.255.0.0",
	17: "255.255.128.0",
	18: "255.255.192.0",
	19: "255.255.224.0",
	20: "255.255.240.0",
	21: "255.255.248.0",
	22: "255.255.252.0",
	23: "255.255.254.0",
	24: "255.255.255.0",
	25: "255.255.255.128",
	26: "255.255.255.192",
	27: "255.255.255.224",
	28: "255.255.255.240",
	29: "255.255.255.248",
	30: "255.255.255.252",
	31: "255.255.255.254",
	32: "255.255.255.255",
}

var ACLMasks map[int]string = map[int]string{
	0:  "255.255.255.255",
	1:  "127.255.255.255",
	2:  "63.255.255.255",
	3:  "31.255.255.255",
	4:  "15.255.255.255",
	5:  "7.255.255.255",
	6:  "3.255.255.255",
	7:  "1.255.255.255",
	8:  "0.255.255.255",
	9:  "0.127.255.255",
	10: "0.63.255.255",
	11: "0.31.255.255",
	12: "0.15.255.255",
	13: "0.7.255.255",
	14: "0.3.255.255",
	15: "0.1.255.255",
	16: "0.0.255.255",
	17: "0.0.127.255",
	18: "0.0.63.255",
	19: "0.0.31.255",
	20: "0.0.15.255",
	21: "0.0.7.255",
	22: "0.0.3.255",
	23: "0.0.1.255",
	24: "0.0.0.255",
	25: "0.0.0.127",
	26: "0.0.0.63",
	27: "0.0.0.31",
	28: "0.0.0.15",
	29: "0.0.0.7",
	30: "0.0.0.3",
	31: "0.0.0.1",
	32: "0.0.0.0",
}

func NextIP(ip string, times int) (string, error) {
	a, err := netip.ParseAddr(ip)
	if err != nil {
		return "", fmt.Errorf("can't parse ip %s: %w", ip, err)
	}
	n := a.Next()
	for times > 1 {
		n = n.Next()
		times--
	}
	return n.String(), nil
}

func PrevIP(ip string, times int) (string, error) {
	a, err := netip.ParseAddr(ip)
	if err != nil {
		return "", fmt.Errorf("can't parse ip %s: %w", ip, err)
	}
	n := a.Prev()
	for times > 1 {
		n = n.Prev()
		times--
	}
	return n.String(), nil
}

func SplitPrefix(pfx string) (string, string, error) {
	prefix, err := netip.ParsePrefix(pfx)
	if err != nil {
		return "", "", fmt.Errorf("can't parse prefix %s: %w", pfx, err)
	}

	addr := prefix.Addr().String()
	mask := fmt.Sprintf("%v", prefix.Bits())

	return addr, mask, nil
}

func Range(pfx string) (string, string, error) {
	prefix, err := netip.ParsePrefix(pfx)
	if err != nil {
		return "", "", fmt.Errorf("can't parse prefix %s: %w", pfx, err)
	}

	if !prefix.Addr().Is4() {
		return "", "", fmt.Errorf("only IPv4 is supported")
	}

	ones := prefix.Bits()
	if ones > 30 {
		return "", "", fmt.Errorf("no usable IPs in this prefix")
	}

	firstIP := prefix.Addr().Next()

	// Calculate number of hosts in the range
	numHosts := uint32(1) << (32 - ones)

	lastIP, err := NextIP(firstIP.String(), int(numHosts)-2)
	if err != nil {
		return "", "", fmt.Errorf("can't calculate last IP for %s: %w", pfx, err)
	}

	return firstIP.String(), lastIP, nil
}
