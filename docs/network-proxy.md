# Network Proxy (Linux, SOCKS5/HTTP)

LAN sharing is **disabled by default**. Normal `happ connect 1` listens on
`127.0.0.1:10808` (SOCKS5) and `127.0.0.1:10809` (HTTP).

## Enable LAN access on Debian

Find the real server IP with `ip -4 addr`.

```sh
happ connect 1 --mode proxy --allow-lan --listen 192.168.1.210
```

This intentionally exposes unauthenticated SOCKS5/HTTP proxy ports on the
specified private IP, including to Docker containers that can reach that IP.
**Protect ports 10808 and 10809 with host firewall rules**, and never forward
them from your router.

The local `127.0.0.1` listeners remain available.

## Optional allowlist

For access limited to home LAN and Docker clients:

```sh
happ connect 1 --mode proxy \
  --allow-lan --listen 192.168.1.210 \
  --allow-cidr 192.168.1.0/24 \
  --allow-cidr 172.18.0.0/16
```

`--allow-cidr` is optional and repeatable. If it is absent, source filtering
is disabled; the host firewall and interface reachability are the only access
barriers. If it is present, Xray routes allowed client addresses to the proxy
and blocks others on the LAN listener. This filtering happens **after** the TCP
connection, not at the firewall boundary.

Allowed networks need not contain the server IP. Only canonical private
IPv4 CIDRs with prefix length /16 or longer are accepted.

## YAML

Legacy single-string syntax is supported:

```yaml
connect:
  mode: proxy
  listen: 192.168.1.210
  allow-lan: true
  allow-cidr: 192.168.1.0/24
```

For multiple networks use a YAML list:

```yaml
connect:
  mode: proxy
  listen: 192.168.1.210
  allow-lan: true
  allow-cidr:
    - 192.168.1.0/24
    - 172.18.0.0/16
```

For unrestricted LAN omit `allow-cidr` entirely.

## Verification

```sh
ss -lntp | grep -E '10808|10809'
curl --socks5-hostname 127.0.0.1:10808 https://api.ipify.org
```

From Docker or a permitted LAN client:

```sh
curl --socks5-hostname 192.168.1.210:10808 https://api.ipify.org
```

The command does not alter Linux routes or start TUN and does not require root.
It rejects `0.0.0.0`, public IP listeners, IPv6 listeners, and LAN flags in
TUN modes.
