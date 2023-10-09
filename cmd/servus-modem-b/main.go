package main

import (
	"log"
	"regexp"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/telnet"
	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/DictumMortuum/servus/pkg/models"
	tl "github.com/ziutek/telnet"
)

var (
	re_status     = regexp.MustCompile(`Status\s+: Up`)
	re_max        = regexp.MustCompile(`Max Rate\(Kbps\)    : (\d+)\s+(\d+)`)
	re_cur        = regexp.MustCompile(`Current Rate\(Kbps\): (\d+)\s+(\d+)`)
	re_fec_down   = regexp.MustCompile(`FEC Errors[ :]+(\d+)`)
	re_fec_up     = regexp.MustCompile(`ATU CFEC Errors[ :]+(\d+)`)
	re_crc_down   = regexp.MustCompile(`CRC Errors[ :]+(\d+)`)
	re_crc_up     = regexp.MustCompile(`ATU CCRC Errors[ :]+(\d+)`)
	re_bytes_down = regexp.MustCompile(`Receive Blocks[ :]+(\d+)`)
	re_bytes_up   = regexp.MustCompile(`Transmit Blocks[ :]+(\d+)`)
	re_snr        = regexp.MustCompile(`Noise Margin\(dB\)[ :]+([\d\.]+)\s+([\d\.]+)`)
)

func getStats(host, user, password string) (string, error) {
	t, err := tl.Dial("tcp", host)
	if err != nil {
		return "", err
	}
	defer t.Close()

	t.SetUnixWriteMode(true)
	var data []byte

	err = telnet.Expect(t, "ADSL2PlusRouter login:")
	if err != nil {
		return "", err
	}

	err = telnet.Sendln(t, user)
	if err != nil {
		return "", err
	}

	err = telnet.Expect(t, "Password:")
	if err != nil {
		return "", err
	}

	err = telnet.Sendln(t, password)
	if err != nil {
		return "", err
	}

	err = telnet.Expect(t, "> ")
	if err != nil {
		return "", err
	}

	err = telnet.SendSlowly(t, "adsl stats\n")
	if err != nil {
		return "", err
	}

	data, err = t.ReadBytes('>')
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func parseStats(raw string) *models.Modem {
	var stats models.Modem

	refs := re_status.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		stats.Status = true
	}

	refs = re_max.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.MaxDown = util.Atoi(match[1])
		stats.MaxUp = util.Atoi(match[2])
	}

	refs = re_cur.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.CurrentDown = util.Atoi(match[1])
		stats.CurrentUp = util.Atoi(match[2])
	}

	refs = re_snr.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.SNRDown = util.Atof(match[1])
		stats.SNRUp = util.Atof(match[2])
	}

	refs = re_crc_down.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.CRCDown = util.Atoi(match[1])
	}

	refs = re_crc_up.FindAllStringSubmatch(raw, 1)
	if len(refs) > 0 {
		match := refs[0]
		stats.CRCUp = util.Atoi(match[1])
	}

	refs = re_fec_down.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.FECDown = util.Atoi(match[1])
	}

	refs = re_fec_up.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.FECUp = util.Atoi(match[1])
	}

	refs = re_bytes_down.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.DataDown = util.Atoi64(match[1])
	}

	refs = re_bytes_up.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.DataUp = util.Atoi64(match[1])
	}

	return &stats
}

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	modem := config.Cfg.Modem["TD5130"]

	raw, err := getStats(modem.Host+":23", modem.User, modem.Pass)
	if err != nil {
		log.Fatal(err)
	}

	s := parseStats(raw)
	s.Host = modem.Host
	err = saveStats(s, modem.Modem)
	if err != nil {
		log.Fatal(err)
	}
}
