# AGENTS.md - DeNet Node Documentation Repository

## Repository Overview
Documentation repository for DeNet Node - a decentralized data storage protocol. No compilation, linting, or testing required. Changes deploy directly from main branch.

## Critical Information
- **Current Version:** v4.0.1-rc13 (node binaries), v1.0.7 (Desktop Manager)
- **Platforms:** Windows, macOS, Linux (AMD64, ARM64, ARMv6)
- **Binary Location:** Downloads from GitHub releases
- **Default Install Directory:** `~/.denode-manager` (Manager), `~/denet` (CLI)
- **Default Port:** 55050
- **Web Interface:** http://localhost:1111 (Desktop Manager)

## File Structure
- `/readme.md` - Main README with installation instructions and download links
- `/guides/` - Platform-specific installation and configuration guides
- `/assets/` & `/assets/webp/` - Images (prefer WebP format)
- `/scripts/` - `install.sh` (manager installer), `denode-manager.sh` (control script)

## Key Commands & Scripts

### install.sh
Downloads and installs DeNode Manager to `~/.denode-manager`. Detects OS/architecture automatically. Requires the zip archive in current directory.

### denode-manager.sh
Interactive menu script to start/stop/restart the server, view logs, and open web interface. Requires installation directory as argument (default: `~/.denode-manager`).

### CLI Node Commands
```
# Interactive mode (prompts for config)
./denode

# Non-interactive mode (all flags required)
./denode --address <wallet> --license <id> --storage <path> --share <GiB> --rpc <url> --ip <addr> --port <num>
```

### Config Management
```
# Generate new config (interactive or non-interactive)
./denode config generate --address 0x... --license 12345 --storage /path --share 200 --rpc https://rpc-peaq.peaq.network --ip 0.0.0.0 --port 55050

# View config
./denode config get --address 0x... [--license 12345]

# Modify config (interactive only)
./denode config set --address 0x...

# Delete all config files
./denode config clear
```

### Running Multiple Nodes
Each node requires separate directory, unique license, and different port:
```bash
mkdir -p ~/denet-node1 ~/denet-node2
cp denode-binary ~/denet-node1/denode
cp denode-binary ~/denet-node2/denode
# Configure each separately with different license/port
```

## Documentation Guidelines
- **Naming:** lowercase with hyphens (`install-denode-linux.md`)
- **Links:** relative paths (`./guides/monitoring.md`), kebab-case anchors
- **Images:** store in `/assets/`, create WebP versions in `/assets/webp/`
- **Tables:** use for version matrices, platform support, download links
- **Code blocks:** always specify language (```bash, ```powershell)

## Git Workflow
```bash
git checkout -b docs/update-guide-name
# Make changes
git add .
git commit -m "docs: update installation guide"
git push origin docs/update-guide-name
```
Commit prefix: `docs:` for documentation changes.

## Common Tasks
- **Add new guide:** Create in `/guides/`, add to readme.md TOC, include screenshots
- **Update version:** Update all download links across files (check readme.md and all install guides)
- **Verify links:** All internal links use relative paths, external links include `https://`

## Troubleshooting
- **Broken links:** Check anchor casing (GitHub uses lowercase kebab-case)
- **Image paths:** Use relative paths from guide location (`../assets/image.png`)
- **Multiple licenses:** Launch one at a time to avoid transaction issues
- **Manager reinstallation:** Ensure `denode` process is fully stopped before reinstalling

## External Resources
- Official site: https://denet.pro
- Discord: https://discord.gg/cPz9m4cSWv
- GitHub releases: https://github.com/DeNetPRO/Node/releases
