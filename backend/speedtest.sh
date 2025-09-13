#!/usr/bin/env bash
LOG_FILE="/app/log/ookla_speedtest_log.csv"

while true; do
  current_minute=$(date +"%M")
  current_second=$(date +"%S")

  total_seconds_in_hour=$(( (10#$current_minute * 60) + 10#$current_second ))
  
  sleep_seconds=$(( 3600 - total_seconds_in_hour ))
  
  echo "Current time: $(date). Sleeping for $sleep_seconds seconds to sync with top of hour."
  sleep "$sleep_seconds"

  SPEEDTEST_OUTPUT=""
  if [ -z "$SERVER_ID" ]; then
    SPEEDTEST_OUTPUT=$(speedtest --accept-license --accept-gdpr --format=csv)
  else
    SPEEDTEST_OUTPUT=$(speedtest --accept-license --accept-gdpr -s "$SERVER_ID" --format=csv)
  fi

  STATUS=$?
  CURRENT_TIME=$(date -u +"%Y-%m-%dT%H:%M:%S.%NZ")

  if [ $STATUS -eq 0 ] && [ -n "$SPEEDTEST_OUTPUT" ]; then
    TEMP_OUTPUT=$(echo "$SPEEDTEST_OUTPUT" | tail -n 1)
    echo "$TEMP_OUTPUT,$CURRENT_TIME" >> "$LOG_FILE"
  else
    ZERO_VALUES="Failed Server,,0,0,0,0,0,0,0,,0,0,0,0,0,0,0,0,0,0,0"
    echo "$ZERO_VALUES,$CURRENT_TIME" >> "$LOG_FILE"
    echo "Speedtest failed. Logged zero values."
  fi
done
