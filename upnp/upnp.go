package upnp

import (
	"context"
	"dfile-secondary-node/shared"

	"gitlab.com/NebulousLabs/go-upnp"
)

var InternetDevice *upnp.IGD

func InitIGD() {
	const logInfo = "shared.InitIGD->"
	device, err := upnp.DiscoverCtx(context.Background())
	if err != nil {
		shared.LogError(logInfo, shared.GetDetailedError(err))
	}

	InternetDevice = device

}

// ====================================================================================
