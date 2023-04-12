#!/bin/bash

for i in {1..15}; do
  ffmpeg -i noise_raw.mkv \
  -c:v libx264 -vf scale=1280:720 \
  -profile:v high -preset medium -r 24 \
  -force_key_frames "expr:gte(t,n_forced*2)" \
  -b:v 2000k \
  -minrate 2000k -maxrate 2000k \
  -bufsize 16500k \
  -t 120 -f mpegts noise_720p_2k_$i.ts
done
