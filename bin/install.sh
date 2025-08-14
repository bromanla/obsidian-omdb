#!/bin/bash
# Install systemd service

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
USER_NAME="$USER"

source "$(dirname "$0")/common.sh"

# Build binary
bash "$SCRIPT_DIR/build.sh"

# Create unit file
cat <<EOF | sudo tee "$UNIT_FILE" > /dev/null
[Unit]
Description=OMDB Telegram Bot
After=network.target

[Service]
Type=simple
ExecStart=${SCRIPT_DIR}/${SERVICE_NAME}
WorkingDirectory=${SCRIPT_DIR}
Restart=always
User=${USER_NAME}

[Install]
WantedBy=multi-user.target
EOF

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable "$SERVICE_NAME"
sudo systemctl restart "$SERVICE_NAME"

success "Service '${SERVICE_NAME}' installed and started as ${USER_NAME}."