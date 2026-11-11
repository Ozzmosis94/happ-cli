package xray

import (
	"encoding/json"
	"testing"

	"github.com/aimuzov/happ-cli/internal/link"
)

func TestLANInboundsAndSourceRouting(t *testing.T) {
	srv := &link.Server{Protocol: "vless", Address: "example.com", Port: 443, UUID: "test", Network: "tcp", Security: "none"}
	for _, tc := range []struct {
		name string
		lan bool
		wantInbounds int
	}{
		{"loopback-only", false, 2},
		{"LAN-and-loopback", true, 4},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := Options{SocksPort: 10808, HTTPPort: 10809}
			if tc.lan {
				opts.Listen = "192.168.1.210"
				opts.LANCIDR = "192.168.1.0/24"
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
				if cfg.Routing == nil || len(cfg.Routing.Rules) != 2 {
					t.Fatal("missing LAN access restrictions")
				}
				rules := cfg.Routing.Rules
				if rules[0].OutboundTag != "proxy" || len(rules[0].Source) != 1 || rules[0].Source[0] != "192.168.1.0/24" {
					t.Fatal("allowed-source rule incorrect")
				}
				if rules[1].OutboundTag != "block" {
					t.Fatal("deny rule incorrect")
				}
			} else if cfg.Routing != nil {
				t.Fatal("unexpected routing in local-only mode")
			}
			if _, err := json.Marshal(cfg); err != nil {
				t.Fatal(err)
			}
		})
	}
}
