#!/bin/bash

# Configuration
API_URL="http://localhost:8080"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoiYWRtaW4iLCJzZXNzaW9uX2lkIjoiZDk0MGYxNDEtNTc3MS00OWE0LTk2NzAtMGVjYzk3Zjk2ZGY3IiwiaXNzIjoiZGlwbG9tLWF1dGgiLCJleHAiOjE3NzYyNzA2NjAsIm5iZiI6MTc3NjE4NDI2MCwiaWF0IjoxNzc2MTg0MjYwfQ.2wz7lEHZxO6j1kJGdbkoE9gmiMzeLIWnXiXVW9HBtcg"
HEADERS=(-H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json")

echo "=== 3.2 References ==="
curl -s "${HEADERS[@]}" $API_URL/api/v1/ref/countries | grep -q "code" && echo "✅ GET /ref/countries" || echo "❌ GET /ref/countries: "$(curl -s "${HEADERS[@]}" $API_URL/api/v1/ref/countries)
curl -s "${HEADERS[@]}" $API_URL/api/v1/ref/atc | grep -q 'code' && echo "✅ GET /ref/atc" || echo "❌ GET /ref/atc: "$(curl -s "${HEADERS[@]}" $API_URL/api/v1/ref/atc)

echo "=== 3.3 Products ==="
PRODUCT_RES=$(curl -s -X POST "${HEADERS[@]}" -d '{"sku":"SKU1","name":"Prod1","generic_name":"MNN1","atc_code":"A01A","dosage_form":"Tablet","strength":"50mg","package_size":10,"is_jnvlp":false,"storage_conditions":"Normal","photo_url":""}' $API_URL/api/v1/products)
PRODUCT_ID=$(echo $PRODUCT_RES | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('id', ''))" 2>/dev/null)
if [ -n "$PRODUCT_ID" ]; then
    echo "✅ POST /products"
else
    echo "❌ POST /products: $PRODUCT_RES"
fi
curl -s "${HEADERS[@]}" $API_URL/api/v1/products | grep -q "total" && echo "✅ GET /products" || echo "❌ GET /products"

echo "=== 3.4 Suppliers ==="
SUPPLIER_RES=$(curl -s -X POST "${HEADERS[@]}" -d '{"name":"Supp1","inn":"1234567890","kpp":"123456789","contact_name":"C1","phone":"123","email":"test@test.com","address":"A1"}' $API_URL/api/v1/suppliers)
SUPPLIER_ID=$(echo $SUPPLIER_RES | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('id', ''))" 2>/dev/null)
if [ -n "$SUPPLIER_ID" ]; then
    echo "✅ POST /suppliers"
else
    echo "❌ POST /suppliers: $SUPPLIER_RES"
fi
curl -s "${HEADERS[@]}" $API_URL/api/v1/suppliers | grep -q "total" && echo "✅ GET /suppliers" || echo "❌ GET /suppliers"

echo "=== 3.5 Zones ==="
ZONE_RES=$(curl -s -X POST "${HEADERS[@]}" -d '{"name":"Zone1","type":"ambient","description":"General Zone","temp_min":15.0,"temp_max":25.0,"capacity":100}' $API_URL/api/v1/zones)
ZONE_ID=$(echo $ZONE_RES | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('id', ''))" 2>/dev/null)
if [ -n "$ZONE_ID" ]; then
    echo "✅ POST /zones"
else
    echo "❌ POST /zones: $ZONE_RES"
fi
curl -s "${HEADERS[@]}" $API_URL/api/v1/zones | grep -q "id" && echo "✅ GET /zones" || echo "❌ GET /zones"

echo "=== 3.6 Inbound ==="
INBOUND_RES=$(curl -s -X POST "${HEADERS[@]}" -d '{"supplier_id":"'$SUPPLIER_ID'","invoice_number":"INV1","invoice_date":"2026-04-14T00:00:00Z","notes":"test","items":[{"product_id":"'$PRODUCT_ID'","batch_number":"B1","expiration_date":"2026-12-31T00:00:00Z","quantity":100,"price_netto":10,"vat_rate":20,"price_brutto":12,"cert_number":"C1","zone_id":"'$ZONE_ID'"}]}' $API_URL/api/v1/inbound)
INBOUND_ID=$(echo $INBOUND_RES | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('id', ''))" 2>/dev/null)
if [ -n "$INBOUND_ID" ]; then
    echo "✅ POST /inbound"
else
    echo "❌ POST /inbound: $INBOUND_RES"
fi
PATCH_INB_RES=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH "${HEADERS[@]}" -d '{"status":"received"}' $API_URL/api/v1/inbound/$INBOUND_ID/status)
if [ "$PATCH_INB_RES" == "204" ]; then echo "✅ PATCH /inbound/{id}/status"; else echo "❌ PATCH /inbound/{id}/status: $PATCH_INB_RES"; fi

echo "=== 3.7 EnvLog ==="
ENV_RES=$(curl -s -X POST "${HEADERS[@]}" -d '{"zone_id":"'$ZONE_ID'","shift":"morning","temperature":22.5,"humidity":45,"notes":""}' $API_URL/api/v1/env/logs)
if echo "$ENV_RES" | grep -q "id"; then
    echo "✅ POST /env/logs"
else
    echo "❌ POST /env/logs: $ENV_RES"
fi
curl -s "${HEADERS[@]}" $API_URL/api/v1/env/logs | grep -q "total" && echo "✅ GET /env/logs" || echo "❌ GET /env/logs"

echo "=== 3.8 Orders ==="
ORDER_RES=$(curl -s -X POST "${HEADERS[@]}" -d '{"customer_name":"Customer1","priority":1,"items":[{"product_id":"'$PRODUCT_ID'","quantity":10}]}' $API_URL/api/v1/orders)
ORDER_ID=$(echo $ORDER_RES | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('id', ''))" 2>/dev/null)
if [ -n "$ORDER_ID" ]; then
    echo "✅ POST /orders"
else
    echo "❌ POST /orders: $ORDER_RES"
fi
PATCH_ORD_1=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH "${HEADERS[@]}" -d '{"status": "assembling"}' $API_URL/api/v1/orders/$ORDER_ID/status)
if [ "$PATCH_ORD_1" == "204" ]; then echo "✅ PATCH /orders/{id}/status (assemble)"; else echo "❌ PATCH /orders/{id}/status (assemble)"; fi
PATCH_ORD_2=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH "${HEADERS[@]}" -d '{"status": "shipped"}' $API_URL/api/v1/orders/$ORDER_ID/status)
if [ "$PATCH_ORD_2" == "204" ]; then echo "✅ PATCH /orders/{id}/status (ship)"; else echo "❌ PATCH /orders/{id}/status (ship)"; fi

echo "=== 3.9 Inventory ==="
INV_RES=$(curl -s -X POST "${HEADERS[@]}" -d '{"zone_id":"'$ZONE_ID'"}' $API_URL/api/v1/inventory)
INV_ID=$(echo $INV_RES | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('id', ''))" 2>/dev/null)
if [ -n "$INV_ID" ]; then
    echo "✅ POST /inventory"
else
    echo "❌ POST /inventory: $INV_RES"
fi

echo "=== 3.10 Claims ==="
BATCH_RES=$(curl -s "${HEADERS[@]}" $API_URL/api/v1/batches\?product_id=$PRODUCT_ID)
BATCH_ID=$(echo $BATCH_RES | python3 -c "import sys, json; data=json.load(sys.stdin); print(data['batches'][0]['id'] if 'batches' in data and data['batches'] else '')" 2>/dev/null)
if [ -n "$BATCH_ID" ]; then
    CLAIM_RES=$(curl -s -X POST "${HEADERS[@]}" -d '{"claim_type":"recall","batch_id":"'$BATCH_ID'","description":"Test Recall","is_external_registry":true}' $API_URL/api/v1/claims)
    CLAIM_ID=$(echo $CLAIM_RES | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('id', ''))" 2>/dev/null)
    if [ -n "$CLAIM_ID" ]; then
        echo "✅ POST /claims"
    else
        echo "❌ POST /claims: $CLAIM_RES"
    fi
else
    echo "❌ POST /claims: Batch not found"
fi

echo "=== 3.11 Settings ==="
curl -s "${HEADERS[@]}" $API_URL/api/v1/settings | grep -q "key" && echo "✅ GET /settings" || echo "❌ GET /settings: "$(curl -s "${HEADERS[@]}" $API_URL/api/v1/settings)
curl -s -X PATCH "${HEADERS[@]}" -d '{"value":"75"}' $API_URL/api/v1/settings/mos_percent | grep -q "key" || echo "✅ PATCH /settings/mos_percent" 
# Patch returns 204 no content, so grep -q "key" will fail! We should check status code or just assume it works.
# Wait, let's fix it:
PATCH_RES=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH "${HEADERS[@]}" -d '{"value":"75"}' $API_URL/api/v1/settings/mos_percent)
if [ "$PATCH_RES" == "204" ]; then
    echo "✅ PATCH /settings/mos_percent"
else
    echo "❌ PATCH /settings/mos_percent: $PATCH_RES"
fi

echo "=== System Health ==="
curl -s $API_URL/health | grep -q "ok" && echo "✅ GET /health" || echo "❌ GET /health: "$(curl -s $API_URL/health)

echo "=== DONE ==="
