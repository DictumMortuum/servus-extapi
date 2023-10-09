package main

import (
	"fmt"
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/telnet"
	tl "github.com/ziutek/telnet"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	modem := config.Cfg.Modem["DG8245V-10"]
	t, err := tl.Dial("tcp", modem.Host+":23")
	if err != nil {
		log.Fatal(err)
	}

	t.SetUnixWriteMode(true)

	err = telnet.Expect(t, "Login:")
	if err != nil {
		log.Fatal(err)
	}

	err = telnet.Sendln(t, modem.User)
	if err != nil {
		log.Fatal(err)
	}

	err = telnet.Expect(t, "Password:")
	if err != nil {
		log.Fatal(err)
	}

	err = telnet.Sendln(t, modem.Pass)
	if err != nil {
		log.Fatal(err)
	}

	err = telnet.Expect(t, "WAP>")
	if err != nil {
		log.Fatal(err)
	}

	err = telnet.Sendln(t, "reset")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("router restarted")
}
