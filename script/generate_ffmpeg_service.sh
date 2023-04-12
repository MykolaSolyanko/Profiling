#!/bin/bash

base_content='[Unit]
Description=RTSP streaming service

[Service]
ExecStart=/usr/bin/ffmpeg -re -stream_loop -1 -i /home/ubuntu/noise_video/noise_1080p_4k_%%NUMBER%%.ts -c copy  -f rtsp -rtsp_transport tcp  rtsp://localhost:%%PORT%%/axis-media/media.amp

[Install]
WantedBy=multi-user.target'

start_port=30004
end_port=30018
start_number=1

for ((port=start_port, number=start_number; port<=end_port; port++, number++))
do
    file_content="${base_content//%%PORT%%/$port}"
    file_content="${file_content//%%NUMBER%%/$number}"
    file_name="/lib/systemd/system/ffmpeg_1080_${port}.service"
    echo "$file_content" > $file_name
done
