#!/bin/bash

if [ $# -eq 0 ]; then
	echo "ROUTER_PORT argument required. Terminating."
	exit 1
fi

export ROUTER_PORT=$1

for (( i=0; i<20; i++)); do

	printf "POST with valid JSON:\n"
	curl -i -X POST 127.0.0.1:${ROUTER_PORT}/echo \
		-H 'Content-Type: application/json' \
		-d '{"messageOrder": 1, "game":"Mobile Legends", "gamerID":"GYUTDTE", "points":20}'
	sleep 3

	printf "POST with valid JSON:\n"
	curl -i -X POST 127.0.0.1:${ROUTER_PORT}/echo \
		-H 'Content-Type: application/json' \
		-d '{"messageOrder": 2, "game":"Mobile Legends", "gamerID":"GYUTDTE", "points":20}'
	sleep 3

	printf "POST with valid JSON:\n"
	curl -i -X POST 127.0.0.1:${ROUTER_PORT}/echo \
		-H 'Content-Type: application/json' \
		-d '{"messageOrder": 3, "game":"Mobile Legends", "gamerID":"GYUTDTE", "points":20}'
	sleep 3

	printf "Iteration complete...\n"
done
