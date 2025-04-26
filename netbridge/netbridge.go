package netbridge

import (
	"net"

	port "rounds.com.ar/watcher/netbridge/port"
	upnp "rounds.com.ar/watcher/netbridge/upnp"
)

type Netbridge struct {
	PublicIP   	net.IP
	Ports 		map[int32]port.Port

	upnpManager *upnp.UPnPManager
}

func NewNetbridge() (*Netbridge, error) {
	var netbridge *Netbridge = &Netbridge{
		PublicIP:    net.IPv4(0, 0, 0, 0),
		Ports:       make(map[int32]port.Port),
		upnpManager: upnp.NewUPnPManager(),
	}

	publicIP, err := netbridge.upnpManager.GetPublicIP()
	if err != nil {
		return nil, err
	}

	netbridge.PublicIP = publicIP

	return netbridge, nil
}
