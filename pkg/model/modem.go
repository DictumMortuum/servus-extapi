package model

import "time"

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
