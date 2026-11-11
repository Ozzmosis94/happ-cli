package xray

import (
	"encoding/json"
	"testing"

	"github.com/aimuzov/happ-cli/internal/link"
)

func TestLANInboundsAndSourceRouting(t *testing.T) {
	srv := &link.Server{Protocol: "vless", Address: "example.com", Port: 443, UUID: "test", Network: "tcp", Security: "none"}
	for _, tc := range []struct {
		name         string
		lan          bool
		cidrs        []string
		wantInbounds int
		wantRouting  bool
	}{
		{"loopback-only", false, nil, 2, false},
		{"LAN without CIDR", true, nil, 4, false},
		{"LAN with single CIDR", true, []string{"192.168.1.0/24"}, 4, true},
		{"LAN and Docker CIDRs", true, []string{"192.168.1.0/24", "172.18.0.0/16"}, 4, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := Options{SocksPort: 10808, HTTPPort: 10809, LANEnabled: tc.lan, LANCIDRs: tc.cidrs}
			if tc.lan {
				opts.Listen = "192.168.1.210"
			}
			cfg, err := BuildConfig(srv, opts)
			if err != nil {
				t.Fatal(err)
			}
			if len(cfg.Inbounds) != tc.wantInbounds {
				t.Fatalf("got %d inbounds; want %d", len(cfg.Inbounds), tc.wantInbounds)
			}
			if cfg.Inbounds[0].Listen != "127.0.0.1" {
				t.Fatal("loopback inbound missing")
			}
			if tc.lan {
				if cfg.Inbounds[2].Listen != "192.168.1.210" || cfg.Inbounds[2].Tag != "socks-lan" {
					t.Fatal("incorrect LAN SOCKS inbound")
				}
				if cfg.Inbounds[3].Listen != "192.168.1.210" || cfg.Inbounds[3].Tag != "http-lan" {
					t.Fatal("incorrect LAN HTTP inbound")
				}
			}
			if (cfg.Routing != nil) != tc.wantRouting {
				t.Fatalf("routing presence = %v, expected %v", cfg.Routing != nil, tc.wantRouting)
			}
			if tc.wantRouting {
				if len(cfg.Routing.Rules) != 2 {
					t.Fatal("expected allow and deny rules")
				}
				rules := cfg.Routing.Rules
				if rules[0].OutboundTag != "proxy" || len(rules[0].Source) != len(tc.cidrs) {
					t.Fatal("allowed-source rule incorrect")
				}
				for i, cidr := range tc.cidrs {
					if rules[0].Source[i] != cidr {
						t.Fatalf("source %d: %s != %s", i, rules[0].Source[i], cidr)
					}
				}
				if rules[1].OutboundTag != "block" {
					t.Fatal("missing deny rule")
				}
			}
			if _, err := json.Marshal(cfg); err != nil {
				t.Fatal(err)
			}
		})
	}
}
