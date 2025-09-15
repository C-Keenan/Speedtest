#!/usr/bin/env bash
LOG_FILE="/app/log/ookla_speedtest_log.csv"
HEADERS="server name,server id,idle latency,idle jitter,packet loss,download,upload,download bytes,upload bytes,share url,download server count,download latency,download latency jitter,download latency low,download latency high,upload latency,upload latency jitter,upload latency low,upload latency high,idle latency low,idle latency high,timestamp"
SPDTSTLOC="/usr/local/bin/speedtest"
if [ ! -x "$SPDTSTLOC" ]; then
  echo "Error: Speedtest binary not found or not executable at $SPDTSTLOC. Please ensure the binary is present and has execute permissions. If you're running this as a docker container, ensure the volume is mounted correctly."
  exit 0
fi
if [ ! -f "$LOG_FILE" ]; then
  echo "$HEADERS" > "$LOG_FILE"
  chmod 664 "$LOG_FILE"
  echo "Log file created with header."
fi
echo "Writing to log file at: $LOG_FILE"
if [ "$ACCEPT_EULA" == "true" ] && [ "$ACCEPT_PRIVACY" == "true" ] && [ "$ACCEPT_TERMS" == "true" ]; then
  echo "All required agreements accepted. Starting speedtest script."
  while true; do
    ping -c 4 1.1.1.1 > /dev/null 2>&1
    status=$?
    stime=60
    if [ "$status" -gt 0 ]; then
      echo "No internet connection detected. Sleeping for $stime seconds."
      sleep "$stime"
    elif [ "$status" -eq 0 ]; then
      echo "Internet connection detected. Proceeding with speedtest."
      break
    fi
  done
  if [ -n "$SERVER_ID" ] && ! [[ "$SERVER_ID" =~ ^[0-9]+$ ]]; then
    echo "Error: SERVER_ID must be a number if set. Current value: '$SERVER_ID'"
    exit 0
  fi
  ./speedtest.sh
else
  echo "You must accept the EULA, Privacy, and Terms to proceed. The links to those respective documents are below."
  echo "EULA: https://www.speedtest.net/about/eula"
  echo "Privacy: https://www.speedtest.net/about/privacy"
  echo "Terms: https://www.speedtest.net/about/terms"
  echo "Set the following environment variables to true to accept: ACCEPT_EULA, ACCEPT_PRIVACY, ACCEPT_TERMS"
  echo "Stopping container and exiting now."
  exit 0
fi
