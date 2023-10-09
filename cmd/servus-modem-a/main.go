package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/telnet"
	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/DictumMortuum/servus/pkg/models"
	tl "github.com/ziutek/telnet"
)

var (
	re_max      = regexp.MustCompile(`Max:\s+Upstream rate = (\d+) Kbps, Downstream rate = (\d+) Kbps`)
	re_cur      = regexp.MustCompile(`Path:\s+\d+, Upstream rate = (\d+) Kbps, Downstream rate = (\d+) Kbps`)
	re_fec_down = regexp.MustCompile(`\nFECErrors:\s+(\d+)`)
	re_fec_up   = regexp.MustCompile(`ATUCFECErrors:\s+(\d+)`)
	re_crc_down = regexp.MustCompile(`\nCRCErrors:\s+(\d+)`)
	re_crc_up   = regexp.MustCompile(`ATUCCRCErrors:\s+(\d+)`)
	re_bytes    = regexp.MustCompile(`bytessent\s+= (\d+)\s+,bytesreceived\s+= (\d+)`)
	re_snr      = regexp.MustCompile(`display dsl snr up=([\d\.]+) down=([\d\.]+) success`)
	re_voip     = regexp.MustCompile(`Status\s+:Enable`)
)

func parseStats(host, user, password, voip string) (*models.Modem, error) {
	var stats models.Modem

	t, err := tl.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	defer t.Close()

	t.SetUnixWriteMode(true)
	var data []byte

	err = telnet.Expect(t, "Login:")
	if err != nil {
		return nil, err
	}

	err = telnet.Sendln(t, user)
	if err != nil {
		return nil, err
	}

	err = telnet.Expect(t, "Password:")
	if err != nil {
		return nil, err
	}

	err = telnet.Sendln(t, password)
	if err != nil {
		return nil, err
	}

	err = telnet.Expect(t, "WAP>")
	if err != nil {
		return nil, err
	}

	err = telnet.Sendln(t, "display xdsl connection status")
	if err != nil {
		return nil, err
	}

	data, err = t.ReadBytes('>')
	if err != nil {
		return nil, err
	}

	raw := string(data)

	// TODO: need to parse On Line: 0 Days 3 Hour 17 Min 24 Sec to unix timestamp
	stats.Uptime = 0
	stats.Status = strings.Contains(raw, "Status: Up")

	refs := re_max.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.MaxUp = util.Atoi(match[1])
		stats.MaxDown = util.Atoi(match[2])
	}

	refs = re_cur.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.CurrentUp = util.Atoi(match[1])
		stats.CurrentDown = util.Atoi(match[2])
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

	err = telnet.Sendln(t, "display xdsl statistics")
	if err != nil {
		return nil, err
	}

	data, err = t.ReadBytes('>')
	if err != nil {
		return nil, err
	}

	raw = string(data)

	refs = re_bytes.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.DataUp = util.Atoi64(match[1])
		stats.DataDown = util.Atoi64(match[2])
	}

	err = telnet.Sendln(t, "display dsl snr")
	if err != nil {
		return nil, err
	}

	data, err = t.ReadBytes('>')
	if err != nil {
		return nil, err
	}

	raw = string(data)

	refs = re_snr.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.SNRUp = util.Atof(match[1])
		stats.SNRDown = util.Atof(match[2])
	}

	err = telnet.Sendln(t, "display waninfo interface "+voip)
	if err != nil {
		return nil, err
	}

	data, err = t.ReadBytes('>')
	if err != nil {
		return nil, err
	}

	raw = string(data)

	refs = re_voip.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		stats.VoipStatus = true
	}

	return &stats, nil
}

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	modem := config.Cfg.Modem["DG8245V-10"]

	s, err := parseStats(modem.Host+":23", modem.User, modem.Pass, modem.Voip)
	if err != nil {
		log.Fatal(err)
	}

	s.Host = modem.Host
	err = saveStats(s, modem.Modem)
	if err != nil {
		log.Fatal(err)
	}
}
