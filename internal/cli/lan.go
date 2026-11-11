package cli

import (
	"fmt"
	"net/netip"
)

// validateNetworkProxy requires explicit opt-in and a private IPv4 listener.
// The default proxy remains loopback-only; TUN modes never expose LAN listeners.
func validateNetworkProxy(mode string, o connectOpts) error {
	listen := o.listen
	if listen == "" {
		listen = "127.0.0.1"
	}
	if mode != "proxy" {
		if o.allowLAN || o.allowCIDR != "" || listen != "127.0.0.1" {
			return fmt.Errorf("--listen, --allow-lan and --allow-cidr are only supported in proxy mode")
		}
		return nil
	}
	if !o.allowLAN {
		if listen != "127.0.0.1" || o.allowCIDR != "" {
			return fmt.Errorf("--listen and --allow-cidr require --allow-lan for non-loopback access")
		}
		return nil
	}
	ip, err := netip.ParseAddr(listen)
	if err != nil || !ip.Is4() || !ip.IsPrivate() {
		return fmt.Errorf("--listen must be a specific private IPv4 address assigned to this host (not 0.0.0.0): %q", listen)
	}
	prefix, err := netip.ParsePrefix(o.allowCIDR)
	if err != nil || !prefix.Addr().Is4() || prefix.Bits() < 16 || prefix.Masked() != prefix || !prefix.Addr().IsPrivate() || !prefix.Contains(ip) {
		return fmt.Errorf("--allow-cidr must be a private IPv4 network (/16 or narrower) containing %s", listen)
	}
	return nil
}
