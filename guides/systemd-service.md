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
- [Running Under a Dedicated Account](#running-under-a-dedicated-account)
- [Replacing the Binary After a Release](#replacing-the-binary-after-a-release)
- [Removing the Service](#removing-the-service)
- [If the Node Does Not Stay Running](#if-the-node-does-not-stay-running)
- [Sandboxing Options](#sandboxing-options)

## Step 1: Keep the Wallet Password Out of the Unit File

Files in `/etc/systemd/system` are world-readable, so every account on the machine could read the wallet password if it were written there. Put it in a file that only root can read; systemd reads it as root before it drops privileges to the service user.

```bash
sudo mkdir -p /etc/denode && sudo chmod 700 /etc/denode
sudo install -o root -g root -m 600 /dev/null /etc/denode/denode.env
sudo nano /etc/denode/denode.env
```

The file holds one line:

```
DENODE_PASSWORD=your_password
```

Edit the file rather than writing the password inside a command: the shell saves everything you type into its history.

systemd substitutes nothing inside these files, so a `$` in the password is taken literally. Backslashes are different: unquoted, they act as escape characters, so if the password contains one, wrap the whole value in single quotes (`DENODE_PASSWORD='my pass word'`) — a single-quoted value is read verbatim. A password containing a single quote can't be expressed either way, so keep the wallet password to letters, digits and `- _ .`.

What this protects against and what it does not: other unprivileged accounts on the machine can no longer read the password, but root can, and so can anything running as the `User=` account — the node must hold the password to unlock the key, so a program running as that user reads it from the process environment. That is the main reason to give the node its own account, see [Running Under a Dedicated Account](#running-under-a-dedicated-account).

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
StartLimitIntervalSec=300
StartLimitBurst=10

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
| `StartLimitIntervalSec`/`StartLimitBurst` | 10 start attempts in 5 minutes — restarts included — put the unit into `failed` instead of retrying forever and filling the journal, which is what a wrong password would do. After fixing the cause: `sudo systemctl reset-failed denode` |

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
umask 077
journalctl -u denode --since "-2 hours" --no-pager > denode-journal.log
```

The journal contains your wallet address, license IDs and the paths the node serves, so look through the file before sending it, and send it privately rather than pasting it into a public channel.

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
StartLimitIntervalSec=300
StartLimitBurst=10

[Service]
Type=simple
User=YOUR_USER
WorkingDirectory=/home/YOUR_USER/denet
EnvironmentFile=/etc/denode/denode.env
EnvironmentFile=/etc/denode/%i.env
ExecStart=/home/YOUR_USER/denet/denode --address ${DENODE_ADDRESS} --license ${DENODE_LICENSE}
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Here `%i` is the instance name — the license ID. The password stays in the shared `denode.env` from Step 1, so it exists in one place on disk; the per-license file holds only the address and license:

```bash
sudo tee /etc/denode/12345.env > /dev/null <<'EOF'
DENODE_ADDRESS=0xYOUR_ADDRESS
DENODE_LICENSE=12345
EOF
sudo chmod 600 /etc/denode/12345.env

sudo systemctl daemon-reload
sudo systemctl enable --now denode@12345.service
```

`${DENODE_ADDRESS}` and `${DENODE_LICENSE}` are substituted by systemd from those files — no shell is involved, which is what makes this safe. Two things follow from that: use `%i`, never `%I`, because `%I` hands back the instance name with special characters unescaped, and never wrap `ExecStart` in `/bin/sh -c` — inside a shell the same values get split and substituted.

Give each instance its own license ID and its own port — two nodes on port `55050` cannot both listen. Set the port with [`config set` or `config generate --port`](./denode-command.md#config-management-commands), and forward each port separately ([Port Forwarding & DDNS](./port-forwarding-ddns.md)).

## Running Under a Dedicated Account

Two reasons to run the node as its own system account instead of your own: someone who gets into the node finds nothing outside `/var/lib/denode`, and someone who gets into your desktop session does not get the wallet password, which any program running as your account could otherwise read from the node's environment (Step 1). The cost is that the config and the key live in the home of whoever configured the node, so choosing this path means configuring the node again under the new account.

```bash
sudo adduser --system --group --home /var/lib/denode --shell /usr/sbin/nologin denet
```

On RHEL/Fedora the same account is created with `sudo useradd -r -g denet -d /var/lib/denode -s /sbin/nologin`, and `nologin` there lives in `/sbin`.

```bash
sudo mkdir -p /var/lib/denode/denet && sudo chown -R denet:denet /var/lib/denode
sudo -u denet -H bash                      # from here on you are denet, HOME=/var/lib/denode
curl -LO https://github.com/DeNetPRO/Node/releases/download/v4.1.3/denode-linux-amd64
mv denode-linux-amd64 /var/lib/denode/denet/denode && chmod +x /var/lib/denode/denet/denode
/var/lib/denode/denet/denode               # configure interactively, as denet
```

Download the binary only from the official release URL over HTTPS, never from a link posted in chat and never into `/tmp` or another world-writable directory.

Then use `User=denet` and those paths in the unit file. The shared storage directory must be writable by `denet`, and how to achieve that depends on the filesystem: `chown` it before mounting any data on it — `chown -R` on an already mounted data disk rewrites ownership of everything stored there and can break other programs using that disk — or mount with `uid=denet,gid=denet`, which is the only option for vfat, exfat and CIFS where `chown` does nothing, or grant access without touching owners: `sudo setfacl -R -m u:denet:rwx /mnt/data`. Without that access the node starts fine and fails only when it first accepts files.

## Replacing the Binary After a Release

The config stays in `~/.denode`, so an upgrade is a file swap with no reconfiguration. Stop the service first: `cp` into a running binary fails with "text file busy", and renaming the file under a live process leaves that process running the old version. Download from the official release URL under `github.com/DeNetPRO/Node/releases` only, as in the previous section.

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
| Exits within a second, mentioning a password or a prompt | `DENODE_PASSWORD` is missing or wrong, or the address/license don't match a saved config | Check the variable is there without printing it — `sudo grep -c '^DENODE_PASSWORD=' /etc/denode/denode.env` should print `1` — then run `./denode config get --address 0x...` as `YOUR_USER` to see what is configured. Never paste or screenshot the contents of `denode.env` into a ticket |
| Start error about a missing file or directory | a wrong path in `ExecStart`, or `User=` is not the account that owns `~/.denode` | Compare with `ls -la` on the binary and with `ls ~/.denode` under that account |
| Fails at load, "unit file has syntax error" | a typo in the unit file | `systemd-analyze verify /etc/systemd/system/denode.service` |
| Restarts and reports the address already in use | an old `nohup` process of the same node is still running | `ps aux \| grep denode`, then `kill <PID>` of the leftover process |
| Runs with no errors, but no transactions appear | the service is fine; peers can't reach the node | Check [Public IP](./public-ip.md) and [port forwarding](./port-forwarding-ddns.md) |

## Sandboxing Options

`NoNewPrivileges=yes`, `PrivateTmp=yes` and `ProtectSystem=full` don't restrict anything this node writes — they never cover home directories or mounted data disks — so they are free to add. `ProtectSystem=strict` makes the whole filesystem read-only, and then every writable path has to be allow-listed or the node dies on the first file it accepts: `ReadWritePaths=/var/lib/denode /mnt/data1`. A disk mounted read-only stays read-only, because `ReadWritePaths=` doesn't change the mount itself. `ProtectHome=yes` breaks this node outright: it hides `~/.denode` completely, and `ReadWritePaths=` cannot bring it back. Use `ProtectHome=read-only` together with `ReadWritePaths=` pointing at your node directory, or run the node under the dedicated account above, whose files sit outside `/home` and stay unaffected.
