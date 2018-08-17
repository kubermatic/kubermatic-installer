#!/usr/bin/env bash
set -e

sonobuoy run

# sonobuoy states that tests can run up to 60 min. They lied. Mine took _WAY_ longer...
TIMEOUT=14400 # 4 hours
INTERVAL=45
ELAPSED=0

while true
do
	if [ "$ELAPSED" -gt "$TIMEOUT" ]; then
		echo "sonobuoy timeoutted."
		exit 2
	fi

	STATUS=$(sonobuoy status 2>&1 || true)

	if [[ $STATUS = *"Sonobuoy is still running"* ]]; then
		echo "waiting for sonobuoy to finish since $((ELAPSED/60)) minutes"
	elif [[ $STATUS = *'error attempting to run sonobuoy: pod has status "Pending"'*  ]]; then
		echo "waiting for sonobuoy pod to start running..."
	elif [[ $STATUS = *"Sonobuoy has failed"* ]]; then
		echo "sonobuoy failed."
		exit 1
	elif [[ $STATUS = *"Sonobuoy has completed"* ]]; then
		echo "sonobuoy completed."
		exit 0
	else
		echo "couldnt parse sonobuoys result: $STATUS"
	fi

	sleep $INTERVAL
	ELAPSED=$((ELAPSED+INTERVAL))
done