package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Stats struct {
	Device string

	RxBytes      uint64
	rxBytesSlice []uint64

	TxBytes      uint64
	txBytesSlice []uint64

	RxPackets      uint64
	rxPacketsSlice []uint64

	TxPackets      uint64
	txPacketsSlice []uint64

	Flows uint64

	lastNetStat    SingleNetStats
	currentNetStat SingleNetStats
	sampleRate     uint
	slicePosition  uint
}

func NewStats(device string, sampleRate uint) *Stats {
	return &Stats{
		Device:         device,
		rxBytesSlice:   make([]uint64, sampleRate, sampleRate),
		txBytesSlice:   make([]uint64, sampleRate, sampleRate),
		rxPacketsSlice: make([]uint64, sampleRate, sampleRate),
		txPacketsSlice: make([]uint64, sampleRate, sampleRate),
		sampleRate:     sampleRate,
	}
}

func (s *Stats) findFlowStats() error {

	file, err := os.Open("/proc/sys/net/netfilter/nf_conntrack_count")
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)

	s.Flows, err = strconv.ParseUint(strings.Trim(string(data), "\n"), 10, 64)

	if err != nil {
		return fmt.Errorf("Cannot parse nf_conntrack_count: %s", err.Error())

	}

	return nil
}

func (s *Stats) findNetStats() error {

	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		data := strings.Fields(strings.Trim(text, " "))
		if len(data) != 17 {
			continue
		}

		if strings.Trim(data[0], ":") != s.Device {
			continue
		}

		if strings.HasSuffix(data[0], ":") {

			s.currentNetStat.ReadArray(data)
		}
	}

	s.addNetStats()

	return nil
}

func (s *Stats) addNetStats() {

	s.rxBytesSlice[s.slicePosition] = s.currentNetStat.RxBytes - s.lastNetStat.RxBytes
	s.txBytesSlice[s.slicePosition] = s.currentNetStat.TxBytes - s.lastNetStat.TxBytes
	s.rxPacketsSlice[s.slicePosition] = s.currentNetStat.RxPackets - s.lastNetStat.RxPackets
	s.txPacketsSlice[s.slicePosition] = s.currentNetStat.TxPackets - s.lastNetStat.TxPackets

	s.slicePosition++

	if s.slicePosition >= s.sampleRate {
		s.slicePosition = 0
	}

	s.lastNetStat = s.currentNetStat

	sumSlice(&s.rxBytesSlice, &s.RxBytes)
	sumSlice(&s.txBytesSlice, &s.TxBytes)
	sumSlice(&s.rxPacketsSlice, &s.RxPackets)
	sumSlice(&s.txPacketsSlice, &s.TxPackets)

}

func sumSlice(s *[]uint64, to *uint64) {

	*to = 0
	for _, val := range *s {
		*to += val
	}
}
