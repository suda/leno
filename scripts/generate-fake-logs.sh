#!/usr/bin/env bash

levels=("debug" "info" "warn" "silly" "trace")
messages=(
    "loading redis lua script"
    "Initializing stream processor to receive feature flag updates"
    "Enabling Sentry"
    "Gracefully shutting down in response to signal"
    "All active requests complete"
    "Completed 200 OK in 4ms"
    "Checking Gmail for CSVs"
    "Timed out"
    "No matching hosted zones"
    "Scoped order and limit are ignored, it's forced to be batch order and batch size"
    "passing message to callback"
    "Housekeeping found ${RANDOM} registered hooks"
    "Device event emitted to recipient"
    "User function errored"
)

while true
do
    level=${levels[$(( RANDOM % 5 ))]}
    message=${messages[$(( RANDOM % ${#messages[@]} ))]}
    echo "{\"level\": \"${level}\",\"time\":$(date +%s%N | cut -b1-13),\"msg\":\"${message}\",\"pid\":${$},\"hostname\":\"$(hostname)\",\"v\":1}"
    sleep 1s
done