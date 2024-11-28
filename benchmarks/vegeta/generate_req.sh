#!/bin/bash
echo "" > req.txt
while read -r uid; do
    echo "GET http://localhost:8080/api/v1/order/$uid" >> req.txt
done < ../uuids.txt
