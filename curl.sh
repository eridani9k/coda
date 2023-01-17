#!/bin/bash

printf "POST with valid JSON:\n"
curl -i -X POST 127.0.0.1:10000/echo \
	-H 'Content-Type: application/json' \
	-d '{"game":"Mobile Legends", "gamerID":"GYUTDTE", "points":20}'

printf "POST with invalid JSON:\n"
curl -i -X POST 127.0.0.1:10000/echo \
	-H 'Content-Type: application/json' \
	-d '{{"game":"Mobile Legends", "gamerID":"GYUTDTE", "points":20}'

printf "Unsupported GET:\n"
curl -i 127.0.0.1:10000/echo
