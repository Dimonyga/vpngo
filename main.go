package main

import (
	"context"
	"crypto/rc4"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/songgao/water"
)

const (
	BUFFERSIZE = 1518 //(MTU + Ethernet header)
	MTU        = BUFFERSIZE - 18 - 20 - 8
)

func rcvrThread(ctx context.Context, wg *sync.WaitGroup, conn net.PacketConn, iface *water.Interface, cfg *Config) {
	payload := make([]byte, BUFFERSIZE)
	dst := make([]byte, BUFFERSIZE)
	c, err := rc4.NewCipher([]byte(cfg.Main.Secret))
	if nil != err {
		log.Fatalln("Unable to create RC4 cipher:", err)
	}
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("rcvrThread: Exit")
			return
		default:
			err := conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
			if nil != err {
				log.Fatalln("Unable to set read deadline:", err)
			}
			n, client, err := conn.ReadFrom(payload)
			if err != nil {
				continue
			}
			// ReadFromUDP can return 0 bytes on timeout and skip broken packets < 22 bytes
			if n < 22 {
				log.Println("Small packet")
				continue
			}
			if cfg.Neighbor.Port == 0 {
				cfg.SetNeighborPort(client.(*net.UDPAddr).Port)
				log.Println("Neighbor port set to: ", cfg.Neighbor.Port)
			}
			if cfg.Neighbor.EAddress == "0.0.0.0" {
				cfg.SetNeighborEaddress(client.(*net.UDPAddr).IP.String())
				log.Println("Neighbor address set to: ", cfg.Neighbor.EAddress)
			}
			c.Reset()
			c.XORKeyStream(dst, payload[:n])
			_, err = iface.ReadWriteCloser.Write(dst[:n])
			if nil != err {
				log.Println("Error writing to local interface: ", err)
			}
		}
	}
}

func sndrThread(ctx context.Context, wg *sync.WaitGroup, conn net.PacketConn, iface *water.Interface, cfg *Config) {
	var packet IPPacket = make([]byte, BUFFERSIZE)
	ds := make([]byte, BUFFERSIZE)
	c, err := rc4.NewCipher([]byte(cfg.Main.Secret))
	if nil != err {
		log.Fatalln("Unable to create RC4 cipher:", err)
	}
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("sndrThread: Exit")
			return
		default:
			plen, err := iface.ReadWriteCloser.Read(packet[:MTU])
			if err != nil {
				log.Println("Error reading from TUN device: ", err)
				break
			}
			if !packet.IsIPv4() {
				log.Println("Not an IPv4 packet")
				continue
			}
			if cfg.Neighbor.Port == 0 {
				log.Println(cfg.Neighbor)
				log.Println("Neighbor port is not set")
				continue
			}
			c.Reset()
			c.XORKeyStream(ds, packet[:plen])
			neighborAddr := &net.UDPAddr{IP: net.ParseIP(cfg.Neighbor.EAddress), Port: cfg.Neighbor.Port}
			_, err = conn.WriteTo(ds[:plen], neighborAddr)
			if nil != err {
				log.Println("Error writing to socket: ", err)
			}
		}
	}

}

func main() {
	whoami := flag.String("whoami", "", "whoami")
	flag.Parse()
	cfg := getConfig(*whoami)
	iface := ifaceSetup(fmt.Sprint(cfg.Me.IAddress, "/", cfg.Me.IMask))
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: cfg.Me.Port})
	if err != nil {
		log.Fatalln("Unable to get UDP socket:", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go rcvrThread(ctx, &wg, udpConn, iface, &cfg)
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go sndrThread(ctx, &wg, udpConn, iface, &cfg)
	}

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGTERM)
	signal.Notify(exitChan, syscall.SIGINT)

	<-exitChan
	log.Println("Main: signal received")
	cancel()
	log.Println("Main: Waiting threads")
	iface.Close()
	wg.Wait()
	udpConn.Close()
	log.Println("Main: Exit")
	time.Sleep(1 * time.Second)
	os.Exit(0)
}
