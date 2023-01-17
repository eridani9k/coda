#!/bin/bash

curl -X POST 127.0.0.1:10000/echo \
	-H 'Content-Type: application/json' \
	-d '{"game":"Mobile Legends", "gamerID":"GYUTDTE", "points":20}'
