package upnp

import (
	"fmt"

	fastupnp "github.com/DeNetPRO/srcPRO/turbo-upnp"
)

var InternetDevice *fastupnp.Device

// Initializes internet gateway device.
// Remember to enable UPnP on your router.
func Init() {
	fmt.Println("Checking UPnP devices...")

	device, err := fastupnp.InitDevice()
	if err != nil {
		fmt.Println("Warn: manual port forwarding may be needed")
	}

	InternetDevice = device
}
