#!/bin/bash

for (( i=0; i<100; i++)); do

	printf "POST with valid JSON:\n"
	curl -i -X POST 127.0.0.1:8080/echo \
		-H 'Content-Type: application/json' \
		-d '{"game":"Mobile Legends", "gamerID":"GYUTDTE", "points":20}'
	sleep 2

	printf "POST with invalid JSON:\n"
	curl -i -X POST 127.0.0.1:8080/echo \
		-H 'Content-Type: application/json' \
		-d '{{"game":"Mobile Legends", "gamerID":"GYUTDTE", "points":20}'
	sleep 2

	printf "Unsupported GET:\n"
	curl -i 127.0.0.1:8080/echo
	sleep 2

	printf "Iteration complete...\n"
done
