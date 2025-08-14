#!/bin/bash
# Remove systemd service

source "$(dirname "$0")/common.sh"

# Stop and disable the service
sudo systemctl stop ${SERVICE_NAME}
sudo systemctl disable ${SERVICE_NAME}

# Remove the unit file
sudo rm -f ${UNIT_FILE}

# Reload systemd configuration
sudo systemctl daemon-reload

# Print success message in green
success "Service '${SERVICE_NAME}' removed."
