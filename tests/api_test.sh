#!/bin/bash

echo -n "Password hash call: "
RESULT=$(curl -s --data "password=angryMonkey" http://localhost:8080/hash)

re='^[0-9]+$'
if ! [[ $RESULT =~ $re ]] ; then
   echo "Error"; 
else
    echo "Passed"; 
fi

x=5
while [ $x -gt 0 ]
    do
    sleep 1s
    echo -n "."
    x=$(( $x - 1 ))
done

echo -e 

echo -n "Password get call: "
GETRESULT=$(curl -s http://localhost:8080/hash/$RESULT)

if [[ -n "$GETRESULT" ]]; then
   echo "Passed"; 
else
    echo "Error"; 
fi

echo "API stats call: "
STATS=$(curl -s http://localhost:8080/stats)


echo $STATS | python -m json.tool
