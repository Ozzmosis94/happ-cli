package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDefault(t *testing.T) {
	c := Default()
	if c.Connect.Mode != "proxy" {
		t.Errorf("mode = %q", c.Connect.Mode)
	}
	if c.Connect.SocksPort != 10808 {
		t.Errorf("socks-port = %d", c.Connect.SocksPort)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()

	c := Default()
	c.Connect.Mode = "tun"
	c.Connect.HealthCheck = true
	c.Connect.AllowCIDRs = CIDRs{"192.168.1.0/24", "172.18.0.0/16"}

	if err := c.Save(dir); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Connect.Mode != "tun" {
		t.Errorf("mode = %q", loaded.Connect.Mode)
	}
	if !loaded.Connect.HealthCheck {
		t.Error("health-check not preserved")
	}
	if !reflect.DeepEqual(loaded.Connect.AllowCIDRs, c.Connect.AllowCIDRs) {
		t.Fatalf("CIDRs changed: got %v want %v", loaded.Connect.AllowCIDRs, c.Connect.AllowCIDRs)
	}
}

func TestAllowCIDRLegacyAndList(t *testing.T) {
	for _, tc := range []struct {
		name string
		raw  string
		want CIDRs
		errorExpected bool
	}{
		{"legacy scalar", "connect:\n  allow-cidr: 192.168.1.0/24\n", CIDRs{"192.168.1.0/24"}, false},
		{"YAML list", "connect:\n  allow-cidr:\n    - 192.168.1.0/24\n    - 172.18.0.0/16\n", CIDRs{"192.168.1.0/24", "172.18.0.0/16"}, false},
		{"empty list", "connect:\n  allow-cidr: []\n", CIDRs{}, false},
		{"invalid item", "connect:\n  allow-cidr:\n    - 192.168.1.0/24\n    - 123\n", nil, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(tc.raw), 0o600); err != nil {
				t.Fatal(err)
			}
			got, err := Load(dir)
			if tc.errorExpected {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got.Connect.AllowCIDRs, tc.want) {
				t.Fatalf("got %v want %v", got.Connect.AllowCIDRs, tc.want)
			}
		})
	}
}

func TestLoadNonExistent(t *testing.T) {
	c, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if c.Connect.SocksPort != 10808 {
		t.Errorf("default socks-port = %d", c.Connect.SocksPort)
	}
}

func TestLoadCorrupted(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("{bad"), 0o600)

	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected error for corrupted yaml")
	}
}
