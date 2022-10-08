#!/bin/bash

# add performance testing

mintimeout=200

for i in {0..100}
do
    timeout=`expr $mintimeout + $(($RANDOM%500))`
    echo "\n-- TIMEOUT: $timeout -- ITERATION: $i"
    curl -X GET --header "Accept: */*" "http://localhost:4001/v1/api/smart?timeout=${timeout}"
done
echo "\n-- TESTING DONE --"
