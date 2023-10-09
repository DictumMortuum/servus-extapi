package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/DictumMortuum/servus/pkg/models"
)

type result struct {
	Id     int            `json:"id"`
	Result map[string]any `json:"result"`
}

func getStats(ip, csrf, login_uid, session_id string) ([]result, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	jsonStr := []byte(`[{"jsonrpc":"2.0","method":"GET","params":"Device.Services.VoiceService.1.CallControl.Line@","id":1},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Line.1.LinkStatus","id":1},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Line.1.CurrentProfile","id":2},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Line.1.LineEncoding","id":3},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Line.1.StandardUsed","id":4},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Channel.1.DownstreamCurrRate","id":5},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Channel.1.UpstreamCurrRate","id":6},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Line.1.DownstreamMaxBitRate","id":7},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Line.1.UpstreamMaxBitRate","id":8},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Line.1.DownstreamNoiseMargin","id":9},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Line.1.UpstreamNoiseMargin","id":10},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Line.1.DownstreamAttenuation","id":11},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Line.1.UpstreamAttenuation","id":12},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Line.1.DownstreamPower","id":13},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Line.1.UpstreamPower","id":14},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Line.1.LastChange","id":15},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Channel.1.Downdelay","id":16},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Channel.1.Updelay","id":17},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Channel.1.DownINP","id":18},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Channel.1.UpINP","id":19},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Channel.1.Stats.Total.DownCRCErrors","id":20},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Channel.1.Stats.Total.UpCRCErrors","id":21},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Channel.1.Stats.Total.DownFECErrors","id":22},{"jsonrpc":"2.0","method":"GET","params":"Device.DSL.Channel.1.Stats.Total.UpFECErrors","id":23}]`)
	req, err := http.NewRequest("POST", "https://"+ip+"/data/data.cgi?csrf_token="+csrf, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/109.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Referer", "https://192.168.2.254/status-and-support.html")
	req.Header.Set("Cookie", "ID=dfuser; login_uid="+login_uid+"; session_id="+session_id+"; username=undefined")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rs []result
	err = json.Unmarshal(raw, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func parseStats(rs []result) *models.Modem {
	var stats models.Modem

	for _, item := range rs {
		for key, val := range item.Result {
			if key == "Device.DSL.Line.1.LinkStatus" {
				if val == "Up" {
					stats.Status = true
				} else {
					stats.Status = false
				}
			} else if key == "Device.DSL.Channel.1.DownstreamCurrRate" {
				stats.CurrentDown = util.Atoi(val.(string))
			} else if key == "Device.DSL.Channel.1.UpstreamCurrRate" {
				stats.CurrentUp = util.Atoi(val.(string))
			} else if key == "Device.DSL.Line.1.DownstreamMaxBitRate" {
				stats.MaxDown = util.Atoi(val.(string))
			} else if key == "Device.DSL.Line.1.UpstreamMaxBitRate" {
				stats.MaxUp = util.Atoi(val.(string))
			} else if key == "Device.DSL.Channel.1.Stats.Total.UpFECErrors" {
				stats.FECUp = util.Atoi(val.(string))
			} else if key == "Device.DSL.Channel.1.Stats.Total.DownFECErrors" {
				stats.FECDown = util.Atoi(val.(string))
			} else if key == "Device.DSL.Channel.1.Stats.Total.UpCRCErrors" {
				stats.CRCUp = util.Atoi(val.(string))
			} else if key == "Device.DSL.Channel.1.Stats.Total.DownCRCErrors" {
				stats.CRCDown = util.Atoi(val.(string))
			} else if key == "Device.DSL.Line.1.UpstreamNoiseMargin" {
				stats.SNRUp = util.Atof(val.(string)) / 10
			} else if key == "Device.DSL.Line.1.DownstreamNoiseMargin" {
				stats.SNRDown = util.Atof(val.(string)) / 10
			} else if key == "Device.Services.VoiceService.1.CallControl.Line@" {
				if voip, ok := val.([]any); ok {
					for _, telephone_number := range voip {
						tel, ok := telephone_number.(map[string]any)

						if tel["DirectoryNumber"] != "" && ok {
							if v, ok := tel["Status"]; ok {
								stats.VoipStatus = v == "Up"
							}
						}
					}
				}
			}
		}
	}

	return &stats
}

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	modem := config.Cfg.Modem["SpeedportPlus2"]
	rs, err := getStats(modem.Host, modem.Extra["csrf"], modem.Extra["login_uid"], modem.Extra["session_id"])
	if err != nil {
		log.Fatal(err)
	}

	s := parseStats(rs)
	s.Host = modem.Host
	err = saveStats(s, modem.Modem)
	if err != nil {
		log.Fatal(err)
	}
}
