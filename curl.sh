#!/bin/bash

if [ $# -eq 0 ]; then
	echo "ROUTER_PORT argument required. Terminating."
	exit 1
fi

export HOST_IP="127.0.0.1"
export ROUTER_PORT=$1
export MAX_ITERATIONS=100

printf "Script will run for ${MAX_ITERATIONS} curl iterations.\n"

for (( i=1; i<=${MAX_ITERATIONS}; i++)); do

	resp=$(curl -s \
		-H 'Content-Type: application/json' \
		-d '{"iteration": "'"$i"'", "game":"Mobile Legends", "gamerID":"GYUTDTE", "points":20}' \
		-X POST ${HOST_IP}:${ROUTER_PORT}/echo)
	printf "Response: ${resp}"
	sleep 3
	printf "\n"

done
