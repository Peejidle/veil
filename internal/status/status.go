// Package status gathers everything veil reports about the tunnel.
package status

import (
	"context"
	"net/netip"
	"time"

	"github.com/Peejidle/veil/internal/ipinfo"
	"github.com/Peejidle/veil/internal/wg"
)

// Report is a snapshot of the tunnel's state at one moment.
type Report struct {
	Up        bool
	PublicIP  netip.Addr
	Handshake time.Time // zero if the tunnel is down or hasn't handshaked yet
}

// Check collects a Report for the named interface.
func Check(ctx context.Context, iface string) (Report, error) {
	var r Report
	var err error

	r.Up, err = wg.IsUp(iface)
	if err != nil {
		return Report{}, err
	}

	r.PublicIP, err = ipinfo.PublicIP(ctx)
	if err != nil {
		return Report{}, err
	}

	if r.Up {
		r.Handshake, err = wg.LatestHandshake(ctx, iface)
		if err != nil {
			return Report{}, err
		}
	}
	return r, nil
}
