#!/usr/bin/env bash
while true; do
  ping -c 4 1.1.1.1
  status=$?
  stime=60
  if [ "$status" -gt 0 ]; then
    echo "No internet connection detected. Sleeping for $stime seconds."
    sleep  $stime
  elif [ "$status" -eq 0 ]; then
    echo "Internet connection detected. Proceeding with speedtest."
    break
  fi
done
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
    if [[ "$SERVER_ID" =~ ^[0-9]+$ ]]; then
      speedtest  --accept-license --accept-gdpr -s "$SERVER_ID" --format=csv | tail -n 1 >> "$LOG_FILE"
    else
      echo "Error: SERVER_ID must be a number. Current value: '$SERVER_ID'"
      exit
    fi
  fi
done