#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" = "darwin" ]; then
  if sysctl -n hw.optional.arm64 2>/dev/null | grep -q 1; then
    ARCH="arm64"
  fi
fi

case "$OS" in
  linux)
    case "$ARCH" in
      x86_64) ARCHIVE="denode-manager-linux-amd64" ;;
      armv6*) ARCHIVE="denode-manager-linux-armv6" ;;
      aarch64) ARCHIVE="denode-manager-linux-arm64" ;;
      *) echo -e "${RED}ERROR: Unsupported architecture: ${BOLD}$ARCH${NC}"; exit 1 ;;
    esac
    ;;
  darwin)
    case "$ARCH" in
      x86_64) ARCHIVE="denode-manager-darwin-amd64" ;;
      arm64) ARCHIVE="denode-manager-darwin-arm64" ;;
      *) echo -e "${RED}ERROR: Unsupported architecture: ${BOLD}$ARCH${NC}"; exit 1 ;;
    esac
    ;;
  *) echo -e "${RED}ERROR: Unsupported OS: ${BOLD}$OS${NC}"; exit 1 ;;
esac

DEFAULT_INSTALL_DIR="$HOME/.denode-manager"
INSTALL_DIR="${1:-$DEFAULT_INSTALL_DIR}"
STATIC_DIR="$INSTALL_DIR/static"
LOGS_DIR="$INSTALL_DIR/logs"

echo -e "\n${BLUE}=== Starting Denode Manager Installation ===${NC}\n"
echo -e "Detected OS: ${BOLD}$OS${NC}, Architecture: ${BOLD}$ARCH${NC}"
echo -e "Selected archive: ${BOLD}$ARCHIVE.zip${NC}"
echo -e "Installation directory: ${BOLD}$INSTALL_DIR${NC}\n"

if [ ! -f "$ARCHIVE.zip" ]; then
  echo -e "${RED}ERROR: Archive ${BOLD}$ARCHIVE.zip${NC} not found in current directory"
  exit 1
fi

TEMP_DIR=$(mktemp -d)
echo -e "Extracting zip archive to ${BOLD}$TEMP_DIR${NC}..."
unzip -q "$ARCHIVE.zip" -d "$TEMP_DIR"
if [ $? -ne 0 ]; then
  echo -e "${RED}ERROR: Failed to extract zip archive${NC}"
  rm -rf "$TEMP_DIR"
  exit 1
fi

# Check for a .tar.gz file in the extracted contents
TAR_FILE=$(find "$TEMP_DIR" -name "*.tar.gz" -type f)
if [ -n "$TAR_FILE" ]; then
  echo -e "Found tar.gz archive: ${BOLD}$TAR_FILE${NC}"
  echo -e "Extracting tar.gz to ${BOLD}$TEMP_DIR${NC}..."
  tar -xzf "$TAR_FILE" -C "$TEMP_DIR"
  if [ $? -ne 0 ]; then
    echo -e "${RED}ERROR: Failed to extract tar.gz archive${NC}"
    rm -rf "$TEMP_DIR"
    exit 1
  fi
  # Remove the tar.gz file after extraction to clean up
  rm -f "$TAR_FILE"
  echo -e "${GREEN}Removed intermediate tar.gz file${NC}"
else
  echo -e "No tar.gz file found, assuming direct extraction from zip${NC}"
fi

# Define paths to expected files in the extracted archive
NODE_BIN="$TEMP_DIR/denode"
SERVER_BIN="$TEMP_DIR/server"
STATIC_SRC="$TEMP_DIR/static"

if [ ! -f "$NODE_BIN" ]; then
  echo -e "${RED}ERROR: Node binary not found in archive${NC}"
  rm -rf "$TEMP_DIR"
  exit 1
fi

if [ ! -f "$SERVER_BIN" ]; then
  echo -e "${RED}ERROR: Server binary not found in archive${NC}"
  rm -rf "$TEMP_DIR"
  exit 1
fi

if [ ! -d "$STATIC_SRC" ]; then
  echo -e "${RED}ERROR: Static directory not found in archive${NC}"
  rm -rf "$TEMP_DIR"
  exit 1
fi

mkdir -p "$INSTALL_DIR" "$LOGS_DIR"
echo -e "${GREEN}Created directories: ${BOLD}$INSTALL_DIR${NC}, ${BOLD}$LOGS_DIR${NC}"

echo -e "Installing ${BOLD}denode${NC} to ${BOLD}$INSTALL_DIR${NC}..."
cp "$NODE_BIN" "$INSTALL_DIR/denode"
chmod +x "$INSTALL_DIR/denode"
chown -R "$(whoami)" "$INSTALL_DIR/denode"
if [ "$OS" = "darwin" ]; then
  xattr -d com.apple.quarantine "$INSTALL_DIR/denode" 2>/dev/null
fi

echo -e "Installing server to ${BOLD}$INSTALL_DIR${NC}..."
cp "$SERVER_BIN" "$INSTALL_DIR/server"
chmod +x "$INSTALL_DIR/server"
chown -R "$(whoami)" "$INSTALL_DIR/server"
if [ "$OS" = "darwin" ]; then
  xattr -d com.apple.quarantine "$INSTALL_DIR/server" 2>/dev/null
fi

echo -e "Installing static files to ${BOLD}$STATIC_DIR${NC}..."
rm -rf "$STATIC_DIR"
cp -r "$STATIC_SRC" "$STATIC_DIR"

echo -e "Checking denode installation...${NC}"
NODE_VERSION=$("$INSTALL_DIR/denode" --version 2>/dev/null)
if [ $? -eq 0 ]; then
  echo -e "${GREEN}Denode installed successfully: ${BOLD}$NODE_VERSION${NC}"
else
  echo -e "${RED}ERROR: Denode installation failed${NC}"
  rm -rf "$TEMP_DIR"
  exit 1
fi

rm -rf "$TEMP_DIR"
echo -e "\n${GREEN}=== Installation Complete! ===${NC}"
echo -e "${BLUE}Use denode-manager.sh to start and manage the server.${NC}"
printf "\n${RED}WARNING: Manager Web now listens on localhost only by default.\n"
printf "It cannot be opened from another device unless you configure remote access.\n"
printf "Do not expose TCP port 1111 to the internet. Guide: https://github.com/DeNetPRO/Node/blob/Dev/guides/remote-manager-access.md${NC}\n"
echo -e "${BLUE}Installation directory: ${BOLD}$INSTALL_DIR${NC}"
