package cli

import (
	"fmt"
	"net/netip"
)

// validateNetworkProxy requires explicit opt-in and a private IPv4 listener.
// Without any allow-cidr the LAN proxy is reachable by any host able to access
// the selected interface. A host firewall is recommended in that mode.
func validateNetworkProxy(mode string, o connectOpts) error {
	listen := o.listen
	if listen == "" {
		listen = "127.0.0.1"
	}
	if mode != "proxy" {
		if o.allowLAN || len(o.allowCIDRs) != 0 || listen != "127.0.0.1" {
			return fmt.Errorf("--listen, --allow-lan and --allow-cidr are only supported in proxy mode")
		}
		return nil
	}
	if !o.allowLAN {
		if listen != "127.0.0.1" || len(o.allowCIDRs) != 0 {
			return fmt.Errorf("--listen and --allow-cidr require --allow-lan for non-loopback access")
		}
		return nil
	}
	ip, err := netip.ParseAddr(listen)
	if err != nil || !ip.Is4() || !ip.IsPrivate() {
		return fmt.Errorf("--listen must be a specific private IPv4 address assigned to this host (not 0.0.0.0): %q", listen)
	}

	for _, cidr := range o.allowCIDRs {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil || !prefix.Addr().Is4() || prefix.Bits() < 16 || prefix.Masked() != prefix || !prefix.Addr().IsPrivate() {
			return fmt.Errorf("--allow-cidr must be a canonical private IPv4 network (/16 or narrower): %q", cidr)
		}
	}
	return nil
}
