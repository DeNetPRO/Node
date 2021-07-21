package upnp

import (
	"fmt"

	"github.com/alex-gubin/fastupnp"
)

var InternetDevice *fastupnp.Device

func Init() {
	fmt.Println("Checking UPnP devices...")

	device, err := fastupnp.InitDevice()
	if err != nil {
		fmt.Println("Warn: manual port forwarding may be needed")
	}

	InternetDevice = device
}
