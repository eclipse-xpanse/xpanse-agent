#!/bin/bash

#
# SPDX-License-Identifier: Apache-2.0
#  SPDX-FileCopyrightText: Huawei Inc.
#

serviceId=""
resourceName=""
pollingInterval=30
xpanseApiEndpoint=""
agentVersion=""

# Function to display usage
usage() {
    echo "Usage: $0 --serviceId <serviceId> --resourceName <resourceName> --pollingInterval <pollingInterval> --xpanseApiEndpoint <xpanseApiEndpoint> --agentVersion <agentVersion>"
    exit 1
}

create_agent_user() {
  NEW_USER="xpanse-agent"
  if id "$NEW_USER" &>/dev/null; then
      echo "User $NEW_USER already exists. Skipping user creation."
  else
    useradd -m $NEW_USER
    if [[ $? -eq 0 ]]; then
        echo "user $NEW_USER created."
    else
      echo "user $NEW_USER creation failed."
      exit 1
    fi
  fi
}

add_user_to_sudo_file() {
  # Add the new user to the sudoers file
  echo "$NEW_USER ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
  if [[ $? -eq 0 ]]; then
      echo "user $NEW_USER added to sudoers"
  else
      echo "user $NEW_USER could not be added to sudoers file"
      exit 1
  fi
}

create_agent_working_directory() {
  # Create application directory
  mkdir -p /home/$NEW_USER/xpanse-agent
  if [[ $? -eq 0 ]]; then
      echo "working directory for agent created successfully"
  else
      echo "working directory creation failed"
      exit 1
  fi
}

download_and_extract_agent_application() {
  # Download your application
  wget -O /home/$NEW_USER/xpanse-agent/xpanse-agent_"${agentVersion}"_linux_amd64.tar.gz https://github.com/eclipse-xpanse/xpanse-agent/releases/download/v"${agentVersion}"/xpanse-agent_"${agentVersion}"_linux_amd64.tar.gz
  if [[ $? -eq 0 ]]; then
    tar xvf /home/$NEW_USER/xpanse-agent/xpanse-agent_"${agentVersion}"_linux_amd64.tar.gz -C /home/$NEW_USER/xpanse-agent
    if [[ $? -eq 0 ]]; then
          echo "agent with version $agentVersion downloaded and extracted"
      else
          echo "failed to extract agent from downloaded file"
          exit 1
      fi
  else
    echo "failed to download agent from github"
          exit 1
  fi
}

create_xpanse_agent_service() {
cat <<EOL > /etc/systemd/system/xpanse-agent.service
[Unit]
Description=xpanse-agent service
After=network.target

[Service]
Type=simple
User=$NEW_USER
WorkingDirectory=/home/$NEW_USER/xpanse-agent
ExecStart=/home/$NEW_USER/xpanse-agent/xpanse-agent start --serviceId $serviceId	--resourceName $resourceName --xpanseApiEndpoint $xpanseApiEndpoint --jobFrequency $pollingInterval
Restart=always
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOL

if [[ $? -eq 0 ]]; then
      echo "xpanse agent service configuration file created"
  else
      echo "failed to create xpanse agent service configuration file"
      exit 1
  fi
}

enable_and_start_agent_service() {
  systemctl daemon-reload
  if [[ $? -eq 0 ]]; then
      systemctl enable xpanse-agent.service
  fi
  if [[ $? -eq 0 ]]; then
      systemctl start xpanse-agent.service
  fi
  if [[ $? -eq 0 ]]; then
      echo "xpanse agent service enabled and started successfully"
    else
      echo "failed to enable and start the xpanse agent service"
      exit 1
    fi
}


# Parse named arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --serviceId)
            serviceId="$2"
            shift 2
            ;;
        --resourceName)
            resourceName="$2"
            shift 2
            ;;
        --pollingInterval)
            pollingInterval="$2"
            shift 2
            ;;
        --xpanseApiEndpoint)
            xpanseApiEndpoint="$2"
            shift 2
            ;;
        --agentVersion)
            agentVersion="$2"
            shift 2
            ;;
        *)
            usage
            ;;
    esac
done

# Check if required arguments are provided
if [[ -z "$serviceId" || -z "$resourceName"  || -z "$pollingInterval" || -z "$xpanseApiEndpoint" || -z "$agentVersion" ]]; then
    usage
fi

create_agent_user
add_user_to_sudo_file
create_agent_working_directory
download_and_extract_agent_application
create_xpanse_agent_service
enable_and_start_agent_service

echo "service created and started successfully"