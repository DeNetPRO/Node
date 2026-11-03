# Access DeNode Manager Web remotely

By default, DeNode Manager is available only on the computer where it runs. To open it from another device, use an SSH tunnel or a private VPN.

**Keep Manager's TCP port 1111 private.** Manager uses HTTP, so opening this port to the internet can expose the interface and its traffic. SSH or a VPN provides a private connection.

The Datakeeper node uses a separate port for peer communication. Do not close that port or place it behind the VPN when configuring remote Manager access.

Choose a method:

- **SSH tunnel:** use it if you have SSH access to the computer running Manager.
- **Private VPN:** use it to connect from a phone or another device. If you need a VPN, see the [VPN setup guide](./private-vpn-setup.md).

## Option 1: SSH tunnel

On the client computer, run:

```sh
ssh -N -o ExitOnForwardFailure=yes -L 127.0.0.1:1111:127.0.0.1:1111 SSH_USER@MANAGER_SSH_HOST
```

Replace `SSH_USER` and `MANAGER_SSH_HOST` with the credentials and address you normally use to SSH into the **Manager computer**. For a custom SSH port add `-p SSH_PORT`. Verify the host-key fingerprint through a trusted channel on first connection.

Keep the terminal open and open [http://localhost:1111](http://localhost:1111) on the same client computer. Press `Ctrl+C` to disconnect.

## Option 2: VPN

### 1. Connect the devices and choose the Manager address

Connect your phone or laptop to the VPN. The computer running Manager must also be reachable through that VPN. Choose the address based on where the VPN connection ends:

| VPN setup | Address to use for Manager | Example |
| --- | --- | --- |
| VPN server runs on the Manager computer | That computer's VPN address | `10.90.90.1` |
| VPN server is a different computer; Manager computer connects as a VPN client | The VPN address assigned to the Manager computer | `10.90.90.2` |
| VPN gateway routes to the Manager computer's local network | The Manager computer's local network address | `192.168.1.20` |

The second row means the **computer running Manager joins the VPN as a client** and receives its own VPN address. The VPN server's public address is only used to connect to the VPN; it is not the Manager address. In the gateway setup, the VPN must route traffic to the Manager computer and allow its replies back to the client.

Use the chosen Manager address in both places below: the `--add-listen` command on the Manager computer, and the URL opened on your phone. The phone's VPN address is different; use it only when writing a firewall rule that allows trusted VPN clients.

On the Manager computer, find its addresses with `ip -br address` (Linux), `ifconfig` (macOS), or `ipconfig` (Windows). `--add-listen` must use an address assigned to an active network interface on this computer. For example, if the Manager computer's VPN address is `10.90.90.2`, use `10.90.90.2` below. Do not enter the phone's address or the VPN server's public address.

When the VPN connection ends at a separate gateway, traffic from that gateway to the Manager computer is outside the VPN encryption.

### 2. Allow access privately

Allow TCP **1111** from trusted VPN clients to the Manager computer. Remove public access to this port in the computer firewall, cloud firewall and router, including IPv6 where applicable. Keep your administrative connection open while changing firewall rules.

### 3. Start Manager on that address

On the Manager computer, stop Manager using your normal procedure. Make sure the interface with the selected address is active. Then run, replacing `10.90.90.2` with the address you chose for the Manager computer:

```sh
cd ~/.denode-manager
./server --add-listen 10.90.90.2:1111
```

Alternatively, start or restart it using the management script:

```sh
DENODE_MANAGER_ADD_LISTEN=10.90.90.2:1111 bash denode-manager.sh
```

Use your usual installation directory and account. If Manager runs as a system service, add the flag or environment variable to the service configuration so it applies on restart. The installed version must support `--add-listen`. By default, Manager also remains available locally at `127.0.0.1:1111`.

### 4. Open Manager

On the VPN-connected client, open `http://10.90.90.2:1111`, using the same Manager address you passed to `--add-listen`. On the Manager computer itself, use `http://localhost:1111`.

## Installing as a PWA

A VPN does not turn HTTP into HTTPS. The remote HTTP address can open in a browser, but an installable PWA generally needs HTTPS with a certificate trusted by the client and must meet the app's installability requirements. A home-screen shortcut may not be an installed PWA. This guide does not configure HTTPS. See [MDN's PWA installation requirements](https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps/Guides/Making_PWAs_installable).

An SSH tunnel opened on the client provides a `localhost` URL, which browsers treat as secure for local development.

## Troubleshooting

- **No VPN connection:** check the VPN endpoint, UDP port, keys and firewalls. Generate traffic and check for a recent WireGuard handshake.
- **VPN connects but Manager does not open:** check the route to the Manager address, the Manager listener and the TCP 1111 firewall rule. With a separate gateway, check the return route too.
- **Manager fails to start:** confirm the selected address belongs to an active interface on the Manager computer and your version supports `--add-listen`.
- **SSH fails:** check SSH access and forwarding. Do not expose port 1111 to fix it.

To check for accidental public access, test each public IPv4 and IPv6 address from outside the VPN and local network, for example `nc -vz -w 3 MANAGER_PUBLIC_IP 1111`. Replace the placeholder with the public IP of the Manager computer or its router. A successful connection means the port is exposed; a failed test alone does not prove all paths are closed.

For the separate tunnel used by the Datakeeper **node**, see the [SSH Tunnel Guide](./ssh-tunnel.md).
