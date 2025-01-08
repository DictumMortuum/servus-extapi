package model

import (
	"fmt"
	"io"
	"time"
)

type Modem struct {
	Host           string    `json:"host"`
	Uptime         int64     `json:"uptime"`
	Status         bool      `json:"status"`
	VoipStatus     bool      `json:"voip_status"`
	VoipCallStatus string    `json:"voip_call_status"`
	VoipLastCall   time.Time `json:"voip_last_call"`
	CurrentUp      int       `json:"current_up"`
	CurrentDown    int       `json:"current_down"`
	MaxUp          int       `json:"max_up"`
	MaxDown        int       `json:"max_down"`
	DataUp         int64     `json:"data_up"`
	DataDown       int64     `json:"data_down"`
	FECUp          int       `json:"fec_up"`
	FECDown        int       `json:"fec_down"`
	CRCUp          int       `json:"crc_up"`
	CRCDown        int       `json:"crc_down"`
	SNRUp          float64   `json:"snr_up"`
	SNRDown        float64   `json:"snr_down"`
}

func (m Modem) Write(w io.Writer) {
	var isEnabled int = 0
	if m.Status {
		isEnabled = 1
	}

	var isVoipEnabled int = 0
	if m.VoipStatus {
		isVoipEnabled = 1
	}

	var isInVoipCall int = 0
	if m.VoipCallStatus == "incall" {
		isInVoipCall = 1
	} else if m.VoipCallStatus == "disconnecting" {
		isInVoipCall = 0
	} else if m.VoipCallStatus == "idle" {
		isInVoipCall = 0
	} else if m.VoipCallStatus == "connecting" {
		isInVoipCall = 1
	} else if m.VoipCallStatus == "calling" {
		isInVoipCall = 1
	}

	fmt.Fprintf(w, "uptime,modem,hostname,%s=%d\n", m.Host, m.Uptime)
	fmt.Fprintf(w, "status,modem,hostname,%s=%d\n", m.Host, isEnabled)
	fmt.Fprintf(w, "voip_status,modem,hostname,%s=%d\n", m.Host, isVoipEnabled)
	fmt.Fprintf(w, "voip_call_status,modem,hostname,%s=%d\n", m.Host, isInVoipCall)
	fmt.Fprintf(w, "snr_up,modem,hostname,%s=%f\n", m.Host, m.SNRUp)
	fmt.Fprintf(w, "snr_down,modem,hostname,%s=%f\n", m.Host, m.SNRDown)
	fmt.Fprintf(w, "max_up,modem,hostname,%s=%d\n", m.Host, m.MaxUp)
	fmt.Fprintf(w, "max_down,modem,hostname,%s=%d\n", m.Host, m.MaxDown)
	fmt.Fprintf(w, "fec_up,modem,hostname,%s=%d\n", m.Host, m.FECUp)
	fmt.Fprintf(w, "fec_down,modem,hostname,%s=%d\n", m.Host, m.FECDown)
	fmt.Fprintf(w, "data_up,modem,hostname,%s=%d\n", m.Host, m.DataUp)
	fmt.Fprintf(w, "data_down,modem,hostname,%s=%d\n", m.Host, m.DataDown)
	fmt.Fprintf(w, "current_up,modem,hostname,%s=%d\n", m.Host, m.CurrentUp)
	fmt.Fprintf(w, "current_down,modem,hostname,%s=%d\n", m.Host, m.CurrentDown)
	fmt.Fprintf(w, "crc_up,modem,hostname,%s=%d\n", m.Host, m.CRCUp)
	fmt.Fprintf(w, "crc_down,modem,hostname,%s=%d\n", m.Host, m.CRCDown)
}
