#!/bin/bash

# add performance testing

mintimeout=200

load_test(){
    for i in {1..10}
    do
        timeout=`expr $mintimeout + $(($RANDOM%1000))`
        curl -X GET --header "Accept: */*" "http://localhost:4001/v1/api/smart?timeout=${timeout}"
        echo ""
    done
}

for i in {1..100}
do
    load_test &
done

wait

echo ""
echo "==========================================================="
echo "=========== finished 1.000 concurrent requests ============"
