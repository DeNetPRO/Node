# Running the Node as a systemd Service

> **Advanced setup.** This guide is for the standalone CLI node on Linux. Distributions and file layouts differ, so treat the unit files here as a reference and adjust the paths to your machine. We can't provide individual support for setup-specific differences.

A systemd service starts the node at boot and brings it back if it exits unexpectedly, so nobody has to keep a terminal open. If you use DeNode Manager instead of the CLI binary, skip this: the Manager runs its own server process, see [DeNode Manager Web for Linux](./install-denode-manager-web-linux.md).

**Before you start:** the node must already be installed, configured, and confirmed to run when you start it by hand ([Linux install guide](./install-denode-linux.md), [CLI configuration](./configuring-cli.md)). A service has no keyboard, so it cannot answer configuration or password prompts — anything left unanswered makes the node exit immediately.

## Table of Contents
- [Step 1: Keep the Wallet Password Out of the Unit File](#step-1-keep-the-wallet-password-out-of-the-unit-file)
- [Step 2: Write the Unit File](#step-2-write-the-unit-file)
- [Step 3: Enable and Verify](#step-3-enable-and-verify)
- [Day-to-Day Commands](#day-to-day-commands)
- [Where the Logs Are](#where-the-logs-are)
- [Multiple Nodes on One Machine](#multiple-nodes-on-one-machine)
- [Running Under a Dedicated Account (Optional)](#running-under-a-dedicated-account-optional)
- [Replacing the Binary After a Release](#replacing-the-binary-after-a-release)
- [Removing the Service](#removing-the-service)
- [If the Node Does Not Stay Running](#if-the-node-does-not-stay-running)

## Step 1: Keep the Wallet Password Out of the Unit File

Files in `/etc/systemd/system` are world-readable, so every account on the machine could read the wallet password if it were written there. Put it in a separate file that only root can read; systemd reads it as root before it drops privileges to the service user.

```bash
sudo mkdir -p /etc/denode
sudo chmod 700 /etc/denode
sudo tee /etc/denode/denode.env > /dev/null <<'EOF'
DENODE_PASSWORD=your_password
EOF
sudo chmod 600 /etc/denode/denode.env
```

If the password contains spaces, put the whole value in double quotes (`DENODE_PASSWORD="my pass word"`) — systemd strips the quotes and keeps the rest as one value.

## Step 2: Write the Unit File

```bash
sudo nano /etc/systemd/system/denode.service
```

```ini
[Unit]
Description=DeNet Datakeeper Node
Documentation=https://github.com/DeNetPRO/Node/blob/Dev/guides/systemd-service.md
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=YOUR_USER
WorkingDirectory=/home/YOUR_USER/denet
EnvironmentFile=/etc/denode/denode.env
ExecStart=/home/YOUR_USER/denet/denode --address 0xYOUR_ADDRESS --license YOUR_LICENSE
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

| What to replace | Value |
|-----------------|-------|
| `YOUR_USER` | the Linux account you configured the node under — `whoami` in the terminal where you ran `./denode` |
| `0xYOUR_ADDRESS` | your Datakeeper wallet address, the one used in [configuration](./configuring-cli.md) |
| `YOUR_LICENSE` | the license ID this node runs |
| paths | absolute paths to your node directory. No `~`: systemd does not expand it |

Why the settings are the way they are:

| Setting | Reason |
|---------|--------|
| `After=`/`Wants=network-online.target` | plain `network.target` is reached before the machine has real connectivity, so the node would start against a down interface. `network-online` waits for the connection, and `Wants=` is what pulls that target in |
| `User=` | the config and encrypted key live in `~/.denode/<address>/`, so the service must run as the account where you configured the node. systemd sets `HOME` from this account, and that is how the node finds its config |
| `WorkingDirectory=` | the directory the node starts in; relative paths in your config resolve from here |
| `EnvironmentFile=` | passes `DENODE_PASSWORD` without storing it in a world-readable file |
| `Restart=always`, `RestartSec=5` | the node is restarted after any exit, clean or crashed, with 5 seconds between attempts |

## Step 3: Enable and Verify

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now denode.service
sudo systemctl status denode.service     # want: active (running)
journalctl -u denode.service -f          # live output of the node
```

`enable --now` starts the service now and makes it start on every boot. If the status is `failed` a second after the first start, the node most likely tried to ask a question — see [troubleshooting](#if-the-node-does-not-stay-running).

## Day-to-Day Commands

| Task | Command |
|------|---------|
| Stop, start, restart | `sudo systemctl stop denode` / `start` / `restart` |
| Is it running, is it enabled for boot | `sudo systemctl status denode` / `systemctl is-enabled denode` |
| Follow the log | `journalctl -u denode -f` |
| Log since this boot / last 100 lines | `journalctl -u denode -b` / `journalctl -u denode -n 100` |
| Re-read an edited unit file | `sudo systemctl daemon-reload && sudo systemctl restart denode` |
| Check the unit file for typos | `systemd-analyze verify /etc/systemd/system/denode.service` |

## Where the Logs Are

A CLI node prints to the terminal, and a service has no terminal, so journald captures that output: `journalctl -u denode` is the log. The `denode.log` file from the [nohup instructions](./install-denode-linux.md#run-in-background) stays empty for a service, and `~/.denode/logs` is written by DeNode Manager, not by the CLI node. When support asks for logs:

```bash
journalctl -u denode --since "-2 hours" --no-pager > denode.log
```

## Multiple Nodes on One Machine

One template unit covers all instances instead of a copy per node, so a change is made once. The instance name is the license ID.

```bash
sudo nano /etc/systemd/system/denode@.service
```

```ini
[Unit]
Description=DeNet Datakeeper Node (license %i)
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=YOUR_USER
WorkingDirectory=/home/YOUR_USER/denet
EnvironmentFile=/etc/denode/%i.env
ExecStart=/home/YOUR_USER/denet/denode --address ${DENODE_ADDRESS} --license ${DENODE_LICENSE}
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Each instance gets its own env file, named after its license. `${DENODE_ADDRESS}` and `${DENODE_LICENSE}` are substituted by systemd from that file — no shell is involved:

```bash
sudo tee /etc/denode/12345.env > /dev/null <<'EOF'
DENODE_PASSWORD=your_password
DENODE_ADDRESS=0xYOUR_ADDRESS
DENODE_LICENSE=12345
EOF
sudo chmod 600 /etc/denode/12345.env

sudo systemctl daemon-reload
sudo systemctl enable --now denode@12345.service
```

Give each instance its own license ID and its own port — two nodes on port `55050` cannot both listen. Set the port with [`config set` or `config generate --port`](./denode-command.md#config-management-commands), and forward each port separately ([Port Forwarding & DDNS](./port-forwarding-ddns.md)).

## Running Under a Dedicated Account (Optional)

A separate system account keeps the node away from your desktop session, but the config and key live in the home of whoever configured the node, so the setup has to be redone under that account. Only switch if you are ready to configure the node again.

```bash
sudo adduser --system --group --home /var/lib/denode --shell /usr/sbin/nologin denet
```

On RHEL/Fedora the same account is created with `sudo useradd -r -g denet -d /var/lib/denode -s /sbin/nologin`, and the `nologin` path there is `/sbin/nologin`.

```bash
sudo mkdir -p /var/lib/denode/denet && sudo chown -R denet:denet /var/lib/denode
sudo -u denet -H bash                      # from here on you are denet, HOME=/var/lib/denode
curl -LO https://github.com/DeNetPRO/Node/releases/download/v4.1.3/denode-linux-amd64
mv denode-linux-amd64 /var/lib/denode/denet/denode && chmod +x /var/lib/denode/denet/denode
/var/lib/denode/denet/denode               # configure interactively, as denet
```

Then use `User=denet` and those paths in the unit file. If the shared storage directory is on a mounted disk, the `denet` account must be able to write there (`sudo chown` it to `denet`, or mount it with the right uid), otherwise the node starts fine and fails only when it first accepts files.

## Replacing the Binary After a Release

The config stays in `~/.denode`, so an upgrade is a file swap with no reconfiguration. Stop the service first: `cp` into a running binary fails with "text file busy", and renaming the file under a live process leaves that process running the old version.

```bash
sudo systemctl stop denode
curl -LO https://github.com/DeNetPRO/Node/releases/download/vX.Y.Z/denode-linux-amd64
chmod +x denode-linux-amd64
mv denode-linux-amd64 /home/YOUR_USER/denet/denode    # as the node's own user; sudo -u denet if it runs under the dedicated account
sudo systemctl start denode
journalctl -u denode -n 20 --no-pager                 # confirm it came up
```

## Removing the Service

```bash
sudo systemctl disable --now denode.service
sudo rm /etc/systemd/system/denode.service
sudo systemctl daemon-reload
```

The node's wallet and configuration are untouched by this — they are in `~/.denode`. To wipe those as well, see [Reinstalling Your Node](./install-denode-linux.md#reinstalling-your-node).

## If the Node Does Not Stay Running

| What the log shows | What it means | What to do |
|--------------------|---------------|------------|
| Exits within a second, mentioning a password or a prompt | `DENODE_PASSWORD` is missing or wrong, or the address/license don't match a saved config | `sudo cat /etc/denode/denode.env`, then run `./denode config get --address 0x...` as `YOUR_USER` to see what is actually configured |
| Start error about a missing file or directory | a wrong path in `ExecStart`, or `User=` is not the account that owns `~/.denode` | Compare with `ls -la` on the binary and with `ls ~/.denode` under that account |
| Fails at load, "unit file has syntax error" | a typo in the unit file | `systemd-analyze verify /etc/systemd/system/denode.service` |
| Restarts and reports the address already in use | an old `nohup` process of the same node is still running | `ps aux \| grep denode`, then `kill <PID>` of the leftover process |
| Runs with no errors, but no transactions appear | the service is fine; peers can't reach the node | Check [Public IP](./public-ip.md) and [port forwarding](./port-forwarding-ddns.md) |

One caution if you add hardening options such as `ProtectSystem=` or `ProtectHome=`: the node writes to its storage directory and to `~/.denode`. Either allow those with `ReadWritePaths=` or leave the hardening options out, otherwise the node runs until the first file arrives and fails then.
