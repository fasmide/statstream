package main

/*
 Copy paste: https://github.com/abrander/agento/blob/master/plugins/agents/netstat/SingleNetStats.go
*/

import (
	"strconv"
)

type SingleNetStats struct {
	RxBytes      uint64
	RxPackets    uint64
	RxErrors     uint64
	RxDropped    uint64
	RxFifo       uint64
	RxFrame      uint64
	RxCompressed uint64
	RxMulticast  uint64
	TxBytes      uint64
	TxPackets    uint64
	TxErrors     uint64
	TxDropped    uint64
	TxFifo       uint64
	TxCollisions uint64
	TxCarrier    uint64
	TxCompressed uint64
}

func (s *SingleNetStats) ReadArray(data []string) {
	l := len(data)

	if l > 1 {
		s.RxBytes, _ = strconv.ParseUint(data[1], 10, 64)
	}

	if l > 2 {
		s.RxPackets, _ = strconv.ParseUint(data[2], 10, 64)
	}

	if l > 3 {
		s.RxErrors, _ = strconv.ParseUint(data[3], 10, 64)
	}

	if l > 4 {
		s.RxDropped, _ = strconv.ParseUint(data[4], 10, 64)
	}

	if l > 5 {
		s.RxFifo, _ = strconv.ParseUint(data[5], 10, 64)
	}

	if l > 6 {
		s.RxFrame, _ = strconv.ParseUint(data[6], 10, 64)
	}

	if l > 7 {
		s.RxCompressed, _ = strconv.ParseUint(data[7], 10, 64)
	}

	if l > 8 {
		s.RxMulticast, _ = strconv.ParseUint(data[8], 10, 64)
	}

	if l > 9 {
		s.TxBytes, _ = strconv.ParseUint(data[9], 10, 64)
	}

	if l > 10 {
		s.TxPackets, _ = strconv.ParseUint(data[10], 10, 64)
	}

	if l > 11 {
		s.TxErrors, _ = strconv.ParseUint(data[11], 10, 64)
	}

	if l > 12 {
		s.TxDropped, _ = strconv.ParseUint(data[12], 10, 64)
	}

	if l > 13 {
		s.TxFifo, _ = strconv.ParseUint(data[13], 10, 64)
	}

	if l > 14 {
		s.TxCollisions, _ = strconv.ParseUint(data[14], 10, 64)
	}

	if l > 15 {
		s.TxCarrier, _ = strconv.ParseUint(data[15], 10, 64)
	}

	if l > 16 {
		s.TxCompressed, _ = strconv.ParseUint(data[16], 10, 64)
	}
}
