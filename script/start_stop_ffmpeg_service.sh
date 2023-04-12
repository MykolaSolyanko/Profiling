#!/bin/bash

# Check if the count argument is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <count>"
  exit 1
fi

# Set the starting port and read the count from the command-line argument
start_port=30004
command=$1
count=$2

# Calculate the ending port based on the count
end_port=$((start_port + count - 1))

# Iterate through the ports and run the systemctl cat command for each service file
for ((port=start_port; port<=end_port; port++))
do
    service_name="ffmpeg_1080_${port}.service"
    systemctl $command  $service_name
done
