# Network Proxy (Linux, SOCKS5/HTTP)

LAN sharing is opt-in. Normal `happ connect 1` continues to listen on
`127.0.0.1:10808` (SOCKS5) and `127.0.0.1:10809` (HTTP).

## Enable on Debian

First find the real LAN IPv4 address and subnet with `ip -4 addr` and
`ip -4 route`; do not assume the sample addresses below apply to your host.

```sh
happ connect 1 --mode proxy \
  --allow-lan \
  --listen 192.168.1.210 \
  --allow-cidr 192.168.1.0/24
```

Clients on the allowed subnet can use:

- SOCKS5: `192.168.1.210:10808`
- HTTP proxy: `192.168.1.210:10809`

Local proxy listeners stay active on `127.0.0.1` with the same ports.
The feature does not change Linux routing, start a TUN device, or require root.

## YAML configuration

With an existing `config.yaml` in the happ-cli config directory:

```yaml
connect:
  mode: proxy
  socks-port: 10808
  http-port: 10809
  listen: 192.168.1.210
  allow-lan: true
  allow-cidr: 192.168.1.0/24
  dns-proxy: true
```

## Verify

```sh
ss -lntp | grep -E '10808|10809'
# From a separate client on the permitted LAN:
curl --socks5-hostname 192.168.1.210:10808 https://api.ipify.org
```

## Security

This proxy has **no authentication**. A CIDR-based Xray routing rule blocks
forwarding from addresses outside the specified range, but does not prevent
those clients from opening TCP connections. Use the Debian firewall as an
additional boundary. Never configure router port forwarding to these ports,
and do not expose them on a public interface.

The CLI rejects `0.0.0.0`, public IP addresses, invalid/noncanonical CIDRs
and LAN configuration in TUN modes. Only private IPv4 addresses are supported.
No change is made to the external subscription fetching behavior or telemetry
of the upstream program.
