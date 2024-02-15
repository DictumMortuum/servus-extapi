package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/go-rod/rod"
)

func parseVoice(page *rod.Page, stats *model.Modem) {
	status := page.MustElement(`div.table-row:nth-child(2) > div:nth-child(3)`).MustText()
	stats.VoipStatus = status == "Up"
}

func parseDSL(page *rod.Page, stats *model.Modem) {
	status := page.MustElement(`#link_status`).MustText()
	stats.Status = status == "Up"

	ds_current_rate := page.MustElement(`#ds_current_rate`).MustText()
	stats.CurrentDown = util.Atoi(strings.TrimSuffix(ds_current_rate, " Kbps"))

	us_current_rate := page.MustElement(`#us_current_rate`).MustText()
	stats.CurrentUp = util.Atoi(strings.TrimSuffix(us_current_rate, " Kbps"))

	ds_maximum_rate := page.MustElement(`#ds_maximum_rate`).MustText()
	stats.MaxDown = util.Atoi(strings.TrimSuffix(ds_maximum_rate, " Kbps"))

	us_maximum_rate := page.MustElement(`#us_maximum_rate`).MustText()
	stats.MaxUp = util.Atoi(strings.TrimSuffix(us_maximum_rate, " Kbps"))

	ds_snr := page.MustElement(`#ds_noise_margin`).MustText()
	stats.SNRDown = util.Atof(strings.TrimSuffix(ds_snr, " dB"))

	us_snr := page.MustElement(`#us_noise_margin`).MustText()
	stats.SNRUp = util.Atof(strings.TrimSuffix(us_snr, " dB"))

	ds_crc := page.MustElement(`#ds_crc`).MustText()
	stats.CRCDown = util.Atoi(ds_crc)

	us_crc := page.MustElement(`#us_crc`).MustText()
	stats.CRCUp = util.Atoi(us_crc)

	ds_fec := page.MustElement(`#ds_fec`).MustText()
	stats.FECDown = util.Atoi(ds_fec)

	us_fec := page.MustElement(`#us_fec`).MustText()
	stats.FECUp = util.Atoi(us_fec)

	stats.DataDown = 0
	stats.DataUp = 0
}

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	modem := config.Cfg.Modem["SpeedportPlus2"]
	var s model.Modem
	browser := rod.New().MustConnect().Trace(false).Timeout(30 * time.Second)
	defer browser.MustClose()

	page := browser.MustPage("http://192.168.2.254")
	defer page.Close()
	page.MustElement(`#userName`).MustWaitVisible()
	page.MustElement(`#userName`).MustInput(modem.User)
	page.MustElement(`div.row:nth-child(4) > div:nth-child(1) > input:nth-child(1)`).MustInput(modem.Pass)
	page.MustElement(`.button`).MustClick()
	page.MustElement(`li.main-menu:nth-child(6) > a:nth-child(1) > span:nth-child(1)`).MustWaitVisible()

	dsl := page.MustNavigate(`http://192.168.2.254/status-and-support.html#sub=1&subSub=66`)
	defer dsl.Close()
	dsl.MustWaitStable()
	parseDSL(dsl, &s)
	// screenshot(dsl, "scr.png")
	dsl.MustElement(`#\33  > a:nth-child(1)`).MustClick()
	dsl.MustWaitStable()
	// screenshot(dsl, "scr2.png")
	parseVoice(dsl, &s)

	s.Host = modem.Host
	err = saveStats(&s, modem.Modem)
	if err != nil {
		log.Fatal(err)
	}

	err = os.RemoveAll("/tmp/rod")
	if err != nil {
		log.Fatal(err)
	}
}
