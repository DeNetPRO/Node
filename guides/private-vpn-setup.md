# Set up a private VPN for DeNode Manager

This guide sets up WireGuard with the VPN server **on the same computer as Manager**. A phone or laptop connects to it as a client. If you already have a VPN, skip to [remote Manager access](./remote-manager-access.md#option-2-vpn).

If your VPN server is on a separate computer, connect the Manager computer and your client to that VPN. The VPN must route traffic between them, including the return path; a WireGuard server forwarding between peers needs IP forwarding and firewall rules. If a gateway routes to Manager's LAN, configure a return route or appropriate NAT. Then follow [remote Manager access](./remote-manager-access.md#option-2-vpn) using the Manager computer's own reachable address. The steps below do not configure a separate gateway.

## Before you start

- You need administrator access to the Manager computer and access to the phone or laptop.
- The Manager computer or its router needs a reachable public IP or DNS name. Allow UDP **51820** to the Manager computer in its firewall and cloud firewall; if behind a router, forward UDP 51820 to it. Carrier-grade NAT or blocked inbound traffic can prevent this setup.
- Keep TCP **1111** closed to the public internet. Keep your administrative connection open while changing firewall rules.
- Choose a VPN subnet that does not overlap with existing networks. This example uses:

| What | Example | Purpose |
| --- | --- | --- |
| Manager VPN address | `10.90.90.1` | Address used to open Manager |
| First client VPN address | `10.90.90.2` | Address assigned to the phone or laptop |
| VPN endpoint | `VPN_SERVER_ENDPOINT_HOST:51820` | Public IP or DNS name and UDP port used to connect to WireGuard |

Replace `VPN_SERVER_ENDPOINT_HOST` with the public IP or DNS name of the VPN server (or its router), without the port. It is different from the private VPN address. Keep private keys and configurations containing them secret; exchange only public keys.

## 1. Install WireGuard

Install WireGuard on the Manager computer and client. See the [official installation page](https://www.wireguard.com/install/). On Ubuntu or Debian, install the server tools with:

```sh
sudo apt update
sudo apt install wireguard
```

The Linux server steps below use `wg-quick` and `systemd`. On Windows or macOS, use the WireGuard app instead.

## 2. Configure the VPN server on the Manager computer

### Linux

Create the server key pair:

```sh
sudo mkdir -p /etc/wireguard
sudo sh -c 'umask 077; wg genkey > /etc/wireguard/manager-server.key; wg pubkey < /etc/wireguard/manager-server.key > /etc/wireguard/manager-server.pub'
sudo cat /etc/wireguard/manager-server.pub
```

Copy the displayed **public key** for the client. Create `/etc/wireguard/wg0.conf` using an administrator-owned editor:

```ini
[Interface]
Address = 10.90.90.1/24
ListenPort = 51820
PrivateKey = VPN_SERVER_PRIVATE_KEY
```

Replace `VPN_SERVER_PRIVATE_KEY` locally with the contents of `/etc/wireguard/manager-server.key`. Do not share it. Protect the configuration:

```sh
sudo chmod 600 /etc/wireguard/wg0.conf
```

### Windows or macOS

In WireGuard, create an empty tunnel. The app generates a private key and displays the corresponding public key. Keep the generated private key, set `Address = 10.90.90.1/24` and `ListenPort = 51820`, and copy the public key for the client. Configure your firewall to allow UDP 51820. The computer and tunnel must stay on for remote access.

## 3. Configure a client

Install the WireGuard app on your phone or computer and create a new tunnel with its own key pair. Keep its private key on the client; copy its **public key** for step 4. Set these fields:

```ini
[Interface]
PrivateKey = CLIENT_PRIVATE_KEY
Address = 10.90.90.2/32

[Peer]
PublicKey = VPN_SERVER_PUBLIC_KEY
Endpoint = VPN_SERVER_ENDPOINT_HOST:51820
AllowedIPs = 10.90.90.1/32
PersistentKeepalive = 25
```

Use the client's generated private key, the server public key from step 2, and your VPN endpoint. Some apps display these as separate fields. Leave DNS empty. `AllowedIPs` sends only traffic to the Manager VPN address through this tunnel.

On Linux, generate the client key pair on the client using `umask 077; wg genkey > client-private.key` and `wg pubkey < client-private.key > client-public.key`. Insert the private key locally into `/etc/wireguard/manager.conf`, protect the file with `sudo chmod 600 /etc/wireguard/manager.conf`, and use `sudo wg-quick up manager` to connect after step 4.

For each additional device, create a separate key pair and use a different address, such as `10.90.90.3/32`. Do not share client configurations or QR codes: they contain private keys.

## 4. Add the client to the VPN server

Add this peer to the server tunnel configuration, replacing the placeholder with the client public key:

```ini
[Peer]
PublicKey = CLIENT_PUBLIC_KEY
AllowedIPs = 10.90.90.2/32
```

Add a separate peer with its own public key and address for each additional client.

On Linux, add the peer to `/etc/wireguard/wg0.conf` and start the tunnel:

```sh
sudo systemctl enable --now wg-quick@wg0
sudo wg show wg0
```

If it was already running, apply the edit with `sudo systemctl restart wg-quick@wg0`; existing VPN clients will briefly disconnect. On Windows or macOS, save the tunnel and activate it; if already active, deactivate and reactivate it.

## 5. Connect and open Manager

Activate the client tunnel. Generate traffic to `10.90.90.1` and check for a recent handshake and transferred data in WireGuard, or run `sudo wg show wg0` on the Linux server. If there is no handshake, check the endpoint, UDP port, router and firewalls, public keys, and whether both tunnels are active.

Then follow [remote Manager access](./remote-manager-access.md#option-2-vpn) to allow trusted VPN clients to reach TCP 1111 and configure Manager to listen on `10.90.90.1:1111`. Open `http://10.90.90.1:1111` on the connected client.

This direct setup does not require IP forwarding: Manager and the VPN server are on the same computer. A VPN does not itself provide HTTPS or guarantee PWA installation; see the [PWA note](./remote-manager-access.md#installing-as-a-pwa).
