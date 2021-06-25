package upnp

import (
	"context"
	"dfile-secondary-node/shared"
	"fmt"

	"gitlab.com/NebulousLabs/go-upnp"
)

var InternetDevice = &upnp.IGD{}

func InitIGD() error {
	const logInfo = "shared.InitIGD->"
	device, err := upnp.DiscoverCtx(context.Background())
	if err != nil {
		return fmt.Errorf("%s %w", logInfo, shared.GetDetailedError(err))
	}

	InternetDevice = device

	return nil
}

// ====================================================================================
