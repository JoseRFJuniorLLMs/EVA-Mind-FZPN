#!/bin/bash

# ============================================
# Simple Qdrant Migration via HTTP API
# ============================================

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔄 PostgreSQL → Qdrant Migration (HTTP API)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check Qdrant
if ! curl -s http://localhost:6333/health > /dev/null; then
    echo "❌ Qdrant not running!"
    exit 1
fi
echo "✅ Qdrant is running"

# Example: Insert test point
echo "📦 Inserting test data..."

curl -X PUT http://localhost:6333/collections/memories/points \
  -H 'Content-Type: application/json' \
  -d '{
    "points": [
      {
        "id": 1,
        "vector": [0.1, 0.2, 0.3],
        "payload": {
          "user_id": 123,
          "content": "Test memory",
          "timestamp": "2026-01-17T00:00:00Z"
        }
      }
    ]
  }'

echo ""
echo "✅ Test data inserted"
echo ""

# Verify
echo "🔍 Verifying..."
curl -s http://localhost:6333/collections/memories | jq .

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Migration test complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📝 Note: For full migration, you need to:"
echo "  1. Extract embeddings from PostgreSQL"
echo "  2. Convert to JSON format"
echo "  3. POST to Qdrant HTTP API"
echo ""
echo "📖 API Docs: https://qdrant.tech/documentation/concepts/points/"
echo ""
