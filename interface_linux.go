//go:build linux
// +build linux

package main

import (
	"log"
	"net"
	"os/exec"

	"github.com/milosgajdos/tenus"
	"github.com/songgao/water"
)

func ifaceSetup(localCIDR string) *water.Interface {

	lIP, lNet, err := net.ParseCIDR(localCIDR)
	if nil != err {
		log.Fatalln("\nlocal ip is not in ip/cidr format")
	}

	iface, err := water.New(water.Config{DeviceType: water.TUN})

	if nil != err {
		log.Println("Unable to allocate TUN interface:", err)
		panic(err)
	}

	log.Println("Interface allocated:", iface.Name())
	link, err := tenus.NewLinkFrom(iface.Name())
	if nil != err {
		log.Fatalln("Unable to get interface info", err)
	}

	err = link.SetLinkMTU(MTU)
	if nil != err {
		log.Fatalln("Unable to set MTU on interface")
	}

	err = link.SetLinkIp(lIP, lNet)
	if nil != err {
		log.Fatalln("Unable to set IP to ", lIP, "/", lNet, " on interface")
	}

	//ip -6 addr flush tun0
	err = exec.Command("ip", "-6", "addr", "flush", iface.Name()).Run()
	if err != nil {
		log.Fatalln("Unable to setup interface:", err)
	}

	err = link.SetLinkUp()
	if nil != err {
		log.Fatalln("Unable to UP interface")
	}

	return iface
}
