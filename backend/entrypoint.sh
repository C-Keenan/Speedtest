#!/usr/bin/env bash
LOG_FILE="/app/log/ookla_speedtest_log.csv"
HEADERS="server name,server id,idle latency,idle jitter,packet loss,download,upload,download bytes,upload bytes,share url,download server count,download latency,download latency jitter,download latency low,download latency high,upload latency,upload latency jitter,upload latency low,upload latency high,idle latency low,idle latency high"
if [ ! -f "$LOG_FILE" ]; then
  echo "$HEADERS" > "$LOG_FILE"
  chmod 664 "$LOG_FILE"
  echo "Log file created with header."
fi
if [ "$ACCEPT_EULA" == "true" ] && [ "$ACCEPT_GDPR" == "true" ] && [ "$ACCEPT_TERMS" == "true" ]; then
  echo "All required agreements accepted. Starting speedtest script."
  ./speedtest.sh
else
  echo "You must accept the EULA, GDPR, and Terms to proceed. The links to those respective documents are below."
  echo "EULA: https://www.speedtest.net/about/eula"
  echo "GDPR: https://www.speedtest.net/about/privacy"
  echo "Terms: https://www.speedtest.net/about/terms"
  echo "Set the following environment variables to true to accept: ACCEPT_EULA, ACCEPT_GDPR, ACCEPT_TERMS"
  echo "Stopping container and exiting now."
  exit 0
fi
