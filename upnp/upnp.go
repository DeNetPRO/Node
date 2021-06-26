package upnp

import (
	"context"
	"fmt"

	"gitlab.com/NebulousLabs/go-upnp"
)

var InternetDevice *upnp.IGD

func InitIGD() {
	fmt.Println("Searching UPnP devices...")

	device, err := upnp.DiscoverCtx(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	InternetDevice = device
}

// ====================================================================================
