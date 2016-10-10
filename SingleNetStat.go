package main

/*
 Copy paste: https://github.com/abrander/agento/blob/master/plugins/agents/netstat/SingleNetStats.go
*/

import (
	"strconv"
)

type SingleNetStats struct {
	RxBytes   uint64
	RxPackets uint64
	TxBytes   uint64
	TxPackets uint64
}

func (s *SingleNetStats) ReadArray(data []string) {
	l := len(data)

	if l > 1 {
		s.RxBytes, _ = strconv.ParseUint(data[1], 10, 64)
	}

	if l > 2 {
		s.RxPackets, _ = strconv.ParseUint(data[2], 10, 64)
	}

	if l > 9 {
		s.TxBytes, _ = strconv.ParseUint(data[9], 10, 64)
	}

	if l > 10 {
		s.TxPackets, _ = strconv.ParseUint(data[10], 10, 64)
	}
}
