#!/bin/bash
cd 1-catalog-service
go run main.go & P1=$!
cd ../2-offer-service
go run main.go & P2=$!
cd ../3-order-service
go run main.go & P3=$!
wait $P1 $P2 $P3