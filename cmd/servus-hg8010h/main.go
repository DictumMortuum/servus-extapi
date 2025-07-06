package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus-extapi/pkg/prometheus"
	"github.com/DictumMortuum/servus-extapi/pkg/telnet"
	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/gin-gonic/gin"
	tl "github.com/ziutek/telnet"
)

// https://hack-gpon.org/ont-huawei-hg8010h/
// https://nova.gr/upload/editor/pdf-documents/diepafes/byod-final-_13-2.pdf

var (
	re_rx_power         = regexp.MustCompile(`RxPower\s+:\s+([\-\.\d]+) \(dBm\)`)
	re_tx_power         = regexp.MustCompile(`TxPower\s+:\s+([\-\.\d]+) \(dBm\)`)
	re_dropped          = regexp.MustCompile(`Dropped packets\s+:\s+(\d+)`)
	re_bip              = regexp.MustCompile(`Bip err\s+:\s+(\d+)`)
	re_rx_omci_overflow = regexp.MustCompile(`Rx omci overflow\s+:\s+(\d+)`)
	re_rx_oversize      = regexp.MustCompile(`Rx oversize\s+:\s+(\d+)`)
)

func parseStats(host, user, password string) (*model.Ont, error) {
	var stats model.Ont

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

	err = telnet.Sendln(t, "display onu info")
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
	stats.Status = strings.Contains(raw, "status:O5")

	err = telnet.Sendln(t, "display optic")
	if err != nil {
		return nil, err
	}

	data, err = t.ReadBytes('>')
	if err != nil {
		return nil, err
	}

	raw = string(data)

	refs := re_rx_power.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.RxPower = util.Atof(match[1])
	}

	refs = re_tx_power.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.TxPower = util.Atof(match[1])
	}

	err = telnet.Sendln(t, "display pon statistics")
	if err != nil {
		return nil, err
	}

	data, err = t.ReadBytes('>')
	if err != nil {
		return nil, err
	}

	raw = string(data)

	refs = re_dropped.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.DroppedPackets = util.Atoi(match[1])
	}

	refs = re_bip.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.BIPErrors = util.Atoi(match[1])
	}

	refs = re_rx_omci_overflow.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.RxOMCIOverflow = util.Atoi(match[1])
	}

	refs = re_rx_oversize.FindAllStringSubmatch(raw, -1)
	if len(refs) > 0 {
		match := refs[0]
		stats.RxOversize = util.Atoi(match[1])
	}

	return &stats, nil
}

func getStats() ([]string, error) {
	s, err := parseStats(config.Cfg.ServusHost+":23", config.Cfg.ServusUser, config.Cfg.ServusPass)
	if err != nil {
		return nil, err
	}
	s.Host = config.Cfg.ServusHost
	return s.Strings(), nil
}

func Version(c *gin.Context) {
	rs := map[string]any{
		"version": "v0.0.1",
	}
	c.AbortWithStatusJSON(200, rs)
}

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/version", Version)
	r.GET("/metrics", gin.WrapF(prometheus.Metrics(getStats)))
	r.GET("/readiness", gin.WrapF(prometheus.ReadinessHandler()))
	r.GET("/liveness", gin.WrapF(prometheus.LivenessHandler()))
	log.Fatal(r.Run(":8080"))
}
