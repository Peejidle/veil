// Package ipinfo looks up this machine's public IP address.
package ipinfo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strings"
)

const endpoint = "https://ipinfo.io/ip"

// PublicIP asks ipinfo.io which address our traffic appears to come from.
func PublicIP(ctx context.Context) (netip.Addr, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return netip.Addr{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("fetching public IP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return netip.Addr{}, fmt.Errorf("fetching public IP: unexpected status %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64))
	if err != nil {
		return netip.Addr{}, fmt.Errorf("reading public IP: %w", err)
	}
	addr, err := netip.ParseAddr(strings.TrimSpace(string(body)))
	if err != nil {
		return netip.Addr{}, fmt.Errorf("response was not an IP address: %w", err)
	}
	return addr, nil
}
