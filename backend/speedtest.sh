#!/usr/bin/env bash
LOG_FILE="/app/log/ookla_speedtest_log.csv"
echo "Writing to log file at: $LOG_FILE"
while true; do
  current_minute=$(date +"%M")
  current_second=$(date +"%S")
  
  total_seconds_in_hour=$(( (10#$current_minute * 60) + 10#$current_second ))
  
  sleep_seconds=$(( 3600 - total_seconds_in_hour ))
  
  echo "Current time: $(date). Sleeping for $sleep_seconds seconds to sync with top of hour."
  sleep "$sleep_seconds"
  if [ -z "$SERVER_ID" ]; then
    speedtest  --accept-license --accept-gdpr --format=csv | tail -n 1 >> "$LOG_FILE"
  else
    speedtest  --accept-license --accept-gdpr -s "$SERVER_ID" --format=csv | tail -n 1 >> "$LOG_FILE"
  fi
done