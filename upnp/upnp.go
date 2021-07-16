package upnp

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	"git.denetwork.xyz/dfile/dfile-secondary-node/upnp/internet"
	"git.denetwork.xyz/dfile/dfile-secondary-node/upnp/ssdp"
)

const (
	clientUrlEnd = "/ctl/IPConn"
)

type Device struct {
	location  string
	clientUrl *url.URL
	port      int
}

var InternetDevice *Device

func Init() {
	fmt.Println("Checking UPnP devices...")

	device, err := initDevice()
	if err != nil {
		fmt.Println("Warn: manual port forwarding may be needed")
	}

	InternetDevice = device
}

func initDevice() (*Device, error) {
	device := &Device{}
	times := 3
	timeSleep := time.Millisecond * 500
	for i := 0; i < times; i++ {
		list, err := ssdp.Search(ssdp.All, 1, "")
		if err != nil {
			return nil, err
		}

		var address string
		for _, srv := range list {
			if srv.Type == internet.URN_WANIPConnection_1 {
				address = srv.Location
			}
		}

		if address == "" {
			time.Sleep(timeSleep)
			continue
		}

		addressSplit := strings.Split(address, "http://")

		if len(addressSplit) != 2 {
			return nil, errors.New("invalid address")
		}

		ipSplit := strings.Split(addressSplit[1], "/")

		if len(ipSplit) != 2 {
			return nil, errors.New("invalid address")
		}

		device.location = ipSplit[0]
		clientAddress := "http://" + device.location + clientUrlEnd
		clientUrl, err := url.Parse(clientAddress)
		if err != nil {
			return device, err
		}
		device.clientUrl = clientUrl

		return device, nil
	}

	return nil, errors.New("something wrong")
}

func (d *Device) Location() string {
	return d.location
}

func (d *Device) PublicIP() (string, error) {
	return internet.GetExternalIPAddress(d.clientUrl)
}

func (d *Device) Forward(port int) error {
	d.port = port
	return internet.AddPortMapping("", getInternalIP(), "TCP", "Test port mapping", uint16(port), uint16(port), true, 0, d.clientUrl)
}

func (d *Device) Close() error {
	return internet.DeletePortMapping("", uint16(d.port), "TCP", d.clientUrl)
}

func getInternalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
