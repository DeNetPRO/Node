package upnp

import (
	"context"
	"fmt"

	"gitlab.com/NebulousLabs/go-upnp"
)

var InternetDevice *upnp.IGD

func InitIGD() {
	fmt.Println("Checking UPnP devices...")

	device, err := upnp.DiscoverCtx(context.Background())
	if err != nil {
		fmt.Println(err)
		fmt.Println("Warn: manual port forwarding may be needed")
	}

	InternetDevice = device
}

// ====================================================================================
