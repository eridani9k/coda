#!/bin/bash

for (( i=0; i<999; i++)); do

	printf "POST with valid JSON:\n"
	curl -i -X POST 127.0.0.1:8080/echo \
		-H 'Content-Type: application/json' \
		-d '{"messageOrder": 1, "game":"Mobile Legends", "gamerID":"GYUTDTE", "points":20}'
	sleep 3

	printf "POST with valid JSON:\n"
	curl -i -X POST 127.0.0.1:8080/echo \
		-H 'Content-Type: application/json' \
		-d '{"messageOrder": 2, "game":"Mobile Legends", "gamerID":"GYUTDTE", "points":20}'
	sleep 3

	printf "POST with valid JSON:\n"
	curl -i -X POST 127.0.0.1:8080/echo \
		-H 'Content-Type: application/json' \
		-d '{"messageOrder": 3, "game":"Mobile Legends", "gamerID":"GYUTDTE", "points":20}'
	sleep 3

	printf "Iteration complete...\n"
done
