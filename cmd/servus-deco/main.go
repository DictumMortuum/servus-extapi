package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/deco"
)

func printDevices(c *deco.Client) error {
	result, err := c.ClientList()
	if err != nil {
		return err
	}

	mappings, err := getMappings()
	if err != nil {
		return err
	}

	for _, device := range result.Result.ClientList {
		var status int
		if device.Online {
			status = 1
		} else {
			status = 0
		}

		nickname := strings.ReplaceAll(device.Name, ",", " ")
		for _, mapping := range mappings {
			if mapping.Mac == device.MAC {
				nickname = mapping.Alias
				break
			}
		}

		fmt.Printf(
			"client,deco,nickname,%s,ip,%s,mac,%s,type,%s,interface,%s=%d\n",
			nickname,
			device.IP,
			device.MAC,
			device.ClientType,
			device.Interface,
			status)
	}

	fmt.Printf("client_total,deco=%d\n", len(result.Result.ClientList))

	return nil
}

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

	err = printDevices(c)
	if err != nil {
		log.Fatal(err.Error())
	}

}
