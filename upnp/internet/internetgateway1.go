package internet

import (
	"errors"
	"net/url"

	"git.denetwork.xyz/dfile/dfile-secondary-node/upnp/soap"
)

const (
	URN_LANDevice_1           = "urn:schemas-upnp-org:device:LANDevice:1"
	URN_WANConnectionDevice_1 = "urn:schemas-upnp-org:device:WANConnectionDevice:1"
	URN_WANDevice_1           = "urn:schemas-upnp-org:device:WANDevice:1"

	URN_WANIPConnection_1  = "urn:schemas-upnp-org:service:WANIPConnection:1"
	URN_WANPPPConnection_1 = "urn:schemas-upnp-org:service:WANPPPConnection:1"
)

func AddPortMapping(NewRemoteHost, NewInternalClient, NewProtocol, NewPortMappingDescription string, NewExternalPort, NewInternalPort uint16, NewEnabled bool, NewLeaseDuration uint32, clientURL *url.URL) error {
	if clientURL == nil {
		return errors.New("url is empty")
	}
	type addRequest struct {
		NewRemoteHost             string
		NewExternalPort           string
		NewProtocol               string
		NewInternalPort           string
		NewInternalClient         string
		NewEnabled                string
		NewPortMappingDescription string
		NewLeaseDuration          string
	}

	exPort, err := soap.MarshalU16(NewExternalPort)
	if err != nil {
		return err
	}

	inPort, err := soap.MarshalU16(NewInternalPort)
	if err != nil {
		return err
	}

	duration, err := soap.MarshalU32(NewLeaseDuration)
	if err != nil {
		return err
	}

	request := addRequest{
		NewRemoteHost:             NewRemoteHost,
		NewExternalPort:           exPort,
		NewProtocol:               NewProtocol,
		NewInternalPort:           inPort,
		NewInternalClient:         NewInternalClient,
		NewEnabled:                soap.MarshalBoolean(NewEnabled),
		NewPortMappingDescription: NewPortMappingDescription,
		NewLeaseDuration:          duration,
	}

	response := interface{}(nil)

	if err = soap.PerformAction(URN_WANPPPConnection_1, "AddPortMapping", clientURL, request, response); err != nil {
		return err
	}

	return nil
}

func DeletePortMapping(NewRemoteHost string, NewExternalPort uint16, NewProtocol string, clientURL *url.URL) error {
	if clientURL == nil {
		return errors.New("url is empty")
	}

	type deleteRequest struct {
		NewRemoteHost   string
		NewExternalPort string
		NewProtocol     string
	}

	port, err := soap.MarshalU16(NewExternalPort)
	if err != nil {
		return err
	}

	request := deleteRequest{
		NewRemoteHost:   NewRemoteHost,
		NewProtocol:     NewProtocol,
		NewExternalPort: port,
	}

	response := interface{}(nil)

	if err := soap.PerformAction(URN_WANPPPConnection_1, "DeletePortMapping", clientURL, request, response); err != nil {
		return err
	}

	return nil
}

func GetExternalIPAddress(clientURL *url.URL) (string, error) {
	if clientURL == nil {
		return "", errors.New("url is empty")
	}

	request := interface{}(nil)

	response := &struct {
		NewExternalIPAddress string
	}{}

	if err := soap.PerformAction(URN_WANPPPConnection_1, "GetExternalIPAddress", clientURL, request, response); err != nil {
		return "", err
	}

	return response.NewExternalIPAddress, nil
}
