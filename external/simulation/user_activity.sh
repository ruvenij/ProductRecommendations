#!/bin/bash

echo -e "\n\nFetching all products..."
curl http://localhost:8080/api/recommendations?user_id=U1004
echo

echo "Creating sample product..."
curl -X POST http://localhost:8080/api/activity \
  -H "Content-Type: application/json" \
  -d '{"user_id": "U1004", "type": "view", "product_id": "P1009"}'

curl -X POST http://localhost:8080/api/activity \
  -H "Content-Type: application/json" \
  -d '{"user_id": "U1004", "type": "wishlist", "product_id": "P1018"}'

curl -X POST http://localhost:8080/api/activity \
  -H "Content-Type: application/json" \
  -d '{"user_id": "U1001", "type": "wishlist", "product_id": "P1020"}'

curl -X POST http://localhost:8080/api/activity \
  -H "Content-Type: application/json" \
  -d '{"user_id": "U1007", "type": "wishlist", "product_id": "P1022"}'

curl -X POST http://localhost:8080/api/activity \
  -H "Content-Type: application/json" \
  -d '{"user_id": "U1004", "type": "purchase", "product_id": "P1009", "quantity": 1, "price": "999.99"}'

curl -X POST http://localhost:8080/api/activity \
  -H "Content-Type: application/json" \
  -d '{"user_id": "U1010", "type": "wishlist", "product_id": "P1029"}'

echo -e "\n\nFetching all products..."
curl http://localhost:8080/api/recommendations?user_id=U1004
echo

sleep 2

echo "Creating sample product..."
curl -X POST http://localhost:8080/api/activity \
  -H "Content-Type: application/json" \
  -d '{"user_id": "U1004", "type": "view", "product_id": "P1027"}'

curl -X POST http://localhost:8080/api/activity \
  -H "Content-Type: application/json" \
  -d '{"user_id": "U1010", "type": "wishlist", "product_id": "P1056"}'

curl -X POST http://localhost:8080/api/activity \
  -H "Content-Type: application/json" \
  -d '{"user_id": "U1009", "type": "wishlist", "product_id": "P1052"}'

curl -X POST http://localhost:8080/api/activity \
  -H "Content-Type: application/json" \
  -d '{"user_id": "U1009", "type": "wishlist", "product_id": "P1053"}'

curl -X POST http://localhost:8080/api/activity \
  -H "Content-Type: application/json" \
  -d '{"user_id": "U1004", "type": "purchase", "product_id": "P1031", "quantity": 1, "price": "999.99"}'

curl -X POST http://localhost:8080/api/activity \
  -H "Content-Type: application/json" \
  -d '{"user_id": "U1005", "type": "wishlist", "product_id": "P1090"}'

echo -e "\n\nFetching all products..."
curl http://localhost:8080/api/recommendations?user_id=U1004
echo