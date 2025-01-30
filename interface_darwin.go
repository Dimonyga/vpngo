//go:build darwin
// +build darwin

package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"

	"github.com/songgao/water"
)

// ifaceSetup returns new interface OR PANIC!
func ifaceSetup(localCIDR string) *water.Interface {
	lIP, _, err := net.ParseCIDR(localCIDR)
	if err != nil {
		log.Println("Unable to parse CIDR:", err)
		return nil
	}
	iface, err := water.New(water.Config{DeviceType: water.TUN})

	if nil != err {
		log.Println("Unable to allocate TUN interface:", err)
		panic(err)
	}

	log.Println("Interface allocated:", iface.Name())

	if err := exec.Command("ifconfig", iface.Name(), "inet", fmt.Sprint(lIP), lIP.String(), "mtu", strconv.FormatInt(MTU, 10), "up").Run(); err != nil {
		log.Fatalln("Unable to setup interface:", err)
	}

	return iface
}
