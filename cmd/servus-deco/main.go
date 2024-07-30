package main

import (
	"fmt"
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/deco"
)

func printDecos(c *deco.Client) error {
	result, err := c.DeviceList()
	if err != nil {
		return err
	}

	count := 0

	for _, device := range result.Result.DeviceList {
		var status int
		if device.InetStatus == "online" {
			status = 1
			count++
		} else {
			status = 0
		}

		fmt.Printf(
			"device,deco,ip,%s,role,%s,inet_error,%s,nickname,%s,group_status,%s=%d\n",
			device.DeviceIP,
			device.Role,
			device.InetErrorMsg,
			device.Nickname,
			device.GroupStatus,
			status)

		fmt.Printf(
			"signal24,deco,ip,%s,role,%s,inet_error,%s,nickname,%s,group_status,%s=%v\n",
			device.DeviceIP,
			device.Role,
			device.InetErrorMsg,
			device.Nickname,
			device.GroupStatus,
			device.SignalLevel.Band24)

		fmt.Printf(
			"signal5,deco,ip,%s,role,%s,inet_error,%s,nickname,%s,group_status,%s=%v\n",
			device.DeviceIP,
			device.Role,
			device.InetErrorMsg,
			device.Nickname,
			device.GroupStatus,
			device.SignalLevel.Band5)
	}

	fmt.Printf("total,deco=%d\n", count)
	return nil
}

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	c := deco.New(config.Cfg.Deco.Host)
	err = c.Authenticate(config.Cfg.Deco.Pass)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = printDecos(c)
	if err != nil {
		log.Fatal(err.Error())
	}
}
