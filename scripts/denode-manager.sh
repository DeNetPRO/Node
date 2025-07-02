#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

DEFAULT_INSTALL_DIR="$HOME/.denode-manager"
INSTALL_DIR="${1:-$DEFAULT_INSTALL_DIR}"
SERVER_START="$INSTALL_DIR/server"
SERVER_PROCESS="server"
PID_FILE="$INSTALL_DIR/server.pid"
LOGS_DIR="$INSTALL_DIR/logs"
FRONTEND_URL="http://localhost:1111"

if [ ! -d "$INSTALL_DIR" ]; then
  echo -e "${RED}ERROR: Installation directory ${BOLD}$INSTALL_DIR${NC} does not exist${NC}"
  exit 1
fi

if [ ! -f "$SERVER_START" ]; then
  echo -e "${RED}ERROR: Server binary ${BOLD}$SERVER_START${NC} not found${NC}"
  exit 1
fi

is_server_running() {
  if [ -f "$PID_FILE" ]; then
    pid=$(cat "$PID_FILE")
    if ps -p "$pid" > /dev/null 2>&1; then
      return 0
    fi
  fi
  ps aux | grep -v grep | grep -E "[[:space:]]$SERVER_PROCESS[[:space:]].*$SERVER_START" > /dev/null
  return $?
}

get_server_pid() {
  if [ -f "$PID_FILE" ] && ps -p "$(cat "$PID_FILE")" > /dev/null 2>&1; then
    cat "$PID_FILE"
  else
    ps aux | grep -v grep | grep -E "[[:space:]]$SERVER_PROCESS[[:space:]].*$SERVER_START" | awk '{print $2}' | head -n 1
  fi
}

echo -e "${BLUE}Denode Manager Control${NC}"
echo "1. Start application"
echo "2. Stop application"
echo "3. Restart application"
echo "4. Check status"
echo "5. Open web interface"
echo "6. View logs"
echo "7. Exit"
read -p "Select action (1-7): " choice

case $choice in
  1)
    if is_server_running; then
      echo -e "${RED}Application already running! (PID: $(get_server_pid))${NC}"
    else
      cd "$INSTALL_DIR" || {
        echo -e "${RED}ERROR: Failed to change to ${BOLD}$INSTALL_DIR${NC}"
        exit 1
      }
      ./server > "$LOGS_DIR/app.log" 2>&1 &
      echo $! > "$PID_FILE"
      sleep 2
      if is_server_running; then
        echo -e "${GREEN}Application started successfully (PID: $(cat "$PID_FILE")).${NC}"
      else
        echo -e "${RED}ERROR: Failed to start application. Check ${BOLD}$LOGS_DIR/app.log${NC}"
        cat "$LOGS_DIR/app.log"
      fi
    fi
    ;;
  2)
    if is_server_running; then
      pid=$(get_server_pid)
      kill -TERM $pid
      rm -f "$PID_FILE"
      echo -e "${GREEN}Application stopped!${NC}"
    else
      echo -e "${RED}Application not running!${NC}"
    fi
    ;;
  3)
    if is_server_running; then
      pid=$(get_server_pid)
      kill -TERM $pid
      rm -f "$PID_FILE"
      sleep 1
    fi
    cd "$INSTALL_DIR" || {
      echo -e "${RED}ERROR: Failed to change to ${BOLD}$INSTALL_DIR${NC}"
      exit 1
    }
    ./server > "$LOGS_DIR/app.log" 2>&1 &
    echo $! > "$PID_FILE"
    sleep 2
    if is_server_running; then
      echo -e "${GREEN}Application restarted (PID: $(cat "$PID_FILE")).${NC}"
    else
      echo -e "${RED}ERROR: Failed to restart application. Check ${BOLD}$LOGS_DIR/app.log${NC}"
      cat "$LOGS_DIR/app.log"
    fi
    ;;
  4)
    if is_server_running; then
      echo -e "${GREEN}Application running (PID: $(get_server_pid)).${NC}"
    else
      echo -e "${RED}Application not running.${NC}"
    fi
    ;;
  5)
    if is_server_running; then
      case "$(uname -s | tr '[:upper:]' '[:lower:]')" in
        darwin)
          open "$FRONTEND_URL"
          ;;
        linux)
          xdg-open "$FRONTEND_URL" >/dev/null 2>&1
          ;;
        *)
          echo -e "${RED}ERROR: Unsupported OS for opening browser${NC}"
          ;;
      esac
    else
      echo -e "${RED}Application not running. Start application before opening web interface.${NC}"
    fi
    ;;
  6)
    if [ -f "$LOGS_DIR/app.log" ]; then
      echo -e "${YELLOW}Recent log entries (${BOLD}$LOGS_DIR/app.log${NC}):"
      tail -n 20 "$LOGS_DIR/app.log"
    else
      echo -e "${RED}Log file not found in ${BOLD}$LOGS_DIR/app.log${NC}"
    fi
    ;;
  7)
    echo -e "${BLUE}Exit.${NC}"
    exit 0
    ;;
  *)
    echo -e "${RED}Invalid choice!${NC}"
    ;;
esac
