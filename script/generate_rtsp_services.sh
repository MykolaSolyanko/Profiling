#!/bin/bash

base_content='[Unit]
Description=RTSP streaming service

[Service]
ExecStart=/home/ubuntu/mediamtx /home/ubuntu/mediamtx.yml
Environment=RTSP_RTSPADDRESS=":%%PORT%%"

[Install]
WantedBy=multi-user.target'

start_port=30004
end_port=30018

for ((port=start_port; port<=end_port; port++))
do
    file_content="${base_content//%%PORT%%/$port}"
    file_name="/lib/systemd/system/rtsp_1080_${port}.service"
    echo "$file_content" > $file_name
done

