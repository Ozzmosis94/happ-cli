package cli

import "testing"

func TestValidateNetworkProxy(t *testing.T) {
	cases := []struct {
		name   string
		mode   string
		opts   connectOpts
		reject bool
	}{
		{"default proxy", "proxy", connectOpts{}, false},
		{"default TUN", "tun", connectOpts{}, false},
		{"LAN without CIDR", "proxy", connectOpts{listen: "192.168.1.210", allowLAN: true}, false},
		{"LAN single CIDR", "proxy", connectOpts{listen: "192.168.1.210", allowLAN: true, allowCIDRs: []string{"192.168.1.0/24"}}, false},
		{"LAN with separate Docker network", "proxy", connectOpts{listen: "192.168.1.210", allowLAN: true, allowCIDRs: []string{"192.168.1.0/24", "172.18.0.0/16"}}, false},
		{"LAN Docker only", "proxy", connectOpts{listen: "192.168.1.210", allowLAN: true, allowCIDRs: []string{"172.18.0.0/16"}}, false},
		{"no opt-in", "proxy", connectOpts{listen: "192.168.1.210"}, true},
		{"wildcard", "proxy", connectOpts{listen: "0.0.0.0", allowLAN: true}, true},
		{"public IP", "proxy", connectOpts{listen: "1.2.3.4", allowLAN: true}, true},
		{"too wide", "proxy", connectOpts{listen: "192.168.1.210", allowLAN: true, allowCIDRs: []string{"192.0.0.0/8"}}, true},
		{"noncanonical CIDR", "proxy", connectOpts{listen: "192.168.1.210", allowLAN: true, allowCIDRs: []string{"192.168.1.10/24"}}, true},
		{"mixed valid and invalid", "proxy", connectOpts{listen: "192.168.1.210", allowLAN: true, allowCIDRs: []string{"172.18.0.0/16", "8.8.8.0/24"}}, true},
		{"LAN with TUN", "tun", connectOpts{listen: "192.168.1.210", allowLAN: true}, true},
		{"CIDR without LAN", "proxy", connectOpts{allowCIDRs: []string{"192.168.1.0/24"}}, true},
		{"IPv6 rejected", "proxy", connectOpts{listen: "fc00::1", allowLAN: true}, true},
		{"loopback LAN rejected", "proxy", connectOpts{listen: "127.0.0.1", allowLAN: true}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateNetworkProxy(tc.mode, tc.opts)
			if (err != nil) != tc.reject {
				t.Fatalf("validateNetworkProxy() error = %v, reject = %v", err, tc.reject)
			}
		})
	}
}
