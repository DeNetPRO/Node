package upnp

import (
	"context"
	"dfile-secondary-node/shared"
	"fmt"

	"gitlab.com/NebulousLabs/go-upnp"
)

var InternetDevice *upnp.IGD

func InitIGD() {
	const logInfo = "shared.InitIGD->"

	fmt.Println("Enabling UPnP...")

	device, err := upnp.DiscoverCtx(context.Background())
	if err != nil {
		shared.LogError(logInfo, shared.GetDetailedError(err))
	}

	InternetDevice = device

}

// ====================================================================================
