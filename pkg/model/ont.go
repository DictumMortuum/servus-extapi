package model

import (
	"fmt"
	"io"
)

type Ont struct {
	Host           string  `json:"host"`
	Uptime         int64   `json:"uptime"`
	Status         bool    `json:"status"`
	RxPower        float64 `json:"rx_power"`
	TxPower        float64 `json:"tx_power"`
	DroppedPackets int     `json:"dropped_packets"`
	BIPErrors      int     `json:"bip_errors"`
	RxOMCIOverflow int     `json:"rx_omci_overflow"`
	RxOversize     int     `json:"rx_oversize"`
}

func (m Ont) Write(w io.Writer) {
	var isEnabled int = 0
	if m.Status {
		isEnabled = 1
	}

	fmt.Fprintf(w, "uptime,ont,hostname,%s=%d\n", m.Host, m.Uptime)
	fmt.Fprintf(w, "status,ont,hostname,%s=%d\n", m.Host, isEnabled)
	fmt.Fprintf(w, "rx_power,ont,hostname,%s=%f\n", m.Host, m.RxPower)
	fmt.Fprintf(w, "tx_power,ont,hostname,%s=%f\n", m.Host, m.TxPower)
	fmt.Fprintf(w, "dropped_packets,ont,hostname,%s=%d\n", m.Host, m.DroppedPackets)
	fmt.Fprintf(w, "bip_errors,ont,hostname,%s=%d\n", m.Host, m.BIPErrors)
	fmt.Fprintf(w, "rx_omci_overflow,ont,hostname,%s=%d\n", m.Host, m.RxOMCIOverflow)
	fmt.Fprintf(w, "rx_oversize,ont,hostname,%s=%d\n", m.Host, m.RxOversize)
}

func (m Ont) Strings() []string {
	var isEnabled int = 0
	if m.Status {
		isEnabled = 1
	}

	return []string{
		fmt.Sprintf("uptime,ont,hostname,%s=%d", m.Host, m.Uptime),
		fmt.Sprintf("status,ont,hostname,%s=%d", m.Host, isEnabled),
		fmt.Sprintf("rx_power,ont,hostname,%s=%f", m.Host, m.RxPower),
		fmt.Sprintf("tx_power,ont,hostname,%s=%f", m.Host, m.TxPower),
		fmt.Sprintf("dropped_packets,ont,hostname,%s=%d", m.Host, m.DroppedPackets),
		fmt.Sprintf("bip_errors,ont,hostname,%s=%d", m.Host, m.BIPErrors),
		fmt.Sprintf("rx_omci_overflow,ont,hostname,%s=%d", m.Host, m.RxOMCIOverflow),
		fmt.Sprintf("rx_oversize,ont,hostname,%s=%d", m.Host, m.RxOversize),
	}
}
