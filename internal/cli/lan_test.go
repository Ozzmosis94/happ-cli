package cli

import "testing"

func TestValidateNetworkProxy(t *testing.T) {
	cases := []struct {
		name string
		mode string
		opts connectOpts
		reject bool
	}{
		{"default proxy", "proxy", connectOpts{}, false},
		{"default TUN", "tun", connectOpts{}, false},
		{"LAN", "proxy", connectOpts{listen: "192.168.1.210", allowLAN: true, allowCIDR: "192.168.1.0/24"}, false},
		{"LAN 10.x", "proxy", connectOpts{listen: "10.20.0.3", allowLAN: true, allowCIDR: "10.20.0.0/16"}, false},
		{"no opt-in", "proxy", connectOpts{listen: "192.168.1.210"}, true},
		{"no source CIDR", "proxy", connectOpts{listen: "192.168.1.210", allowLAN: true}, true},
		{"wildcard", "proxy", connectOpts{listen: "0.0.0.0", allowLAN: true, allowCIDR: "192.168.1.0/24"}, true},
		{"public IP", "proxy", connectOpts{listen: "1.2.3.4", allowLAN: true, allowCIDR: "1.2.3.0/24"}, true},
		{"too wide", "proxy", connectOpts{listen: "192.168.1.210", allowLAN: true, allowCIDR: "192.0.0.0/8"}, true},
		{"different subnet", "proxy", connectOpts{listen: "192.168.1.210", allowLAN: true, allowCIDR: "192.168.2.0/24"}, true},
		{"noncanonical CIDR", "proxy", connectOpts{listen: "192.168.1.210", allowLAN: true, allowCIDR: "192.168.1.10/24"}, true},
		{"LAN with TUN", "tun", connectOpts{listen: "192.168.1.210", allowLAN: true, allowCIDR: "192.168.1.0/24"}, true},
		{"CIDR without LAN", "proxy", connectOpts{allowCIDR: "192.168.1.0/24"}, true},
		{"IPv6 rejected", "proxy", connectOpts{listen: "fc00::1", allowLAN: true, allowCIDR: "fc00::/64"}, true},
		{"loopback LAN rejected", "proxy", connectOpts{listen: "127.0.0.1", allowLAN: true, allowCIDR: "127.0.0.0/24"}, true},
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
