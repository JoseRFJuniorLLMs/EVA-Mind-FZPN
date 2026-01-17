#!/bin/bash

# ============================================
# Qdrant Migration Script
# Migrate embeddings from PostgreSQL to Qdrant
# ============================================

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔄 Starting Qdrant Migration"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if Qdrant is running
echo "🔍 Checking Qdrant status..."
if ! curl -s http://localhost:6333/health > /dev/null; then
    echo "❌ Qdrant is not running!"
    echo "   Start it with: sudo systemctl start qdrant"
    exit 1
fi
echo "✅ Qdrant is running"
echo ""

# Check if PostgreSQL is accessible
echo "🔍 Checking PostgreSQL status..."
if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
    echo "❌ PostgreSQL is not accessible!"
    exit 1
fi
echo "✅ PostgreSQL is accessible"
echo ""

# Verify collections exist
echo "🔍 Verifying Qdrant collections..."
COLLECTIONS=$(curl -s http://localhost:6333/collections | jq -r '.result.collections[].name')

if ! echo "$COLLECTIONS" | grep -q "memories"; then
    echo "⚠️ Collection 'memories' not found. Creating..."
    curl -X PUT http://localhost:6333/collections/memories \
      -H 'Content-Type: application/json' \
      -d '{"vectors": {"size": 768, "distance": "Cosine"}}'
    echo ""
fi

if ! echo "$COLLECTIONS" | grep -q "signifiers"; then
    echo "⚠️ Collection 'signifiers' not found. Creating..."
    curl -X PUT http://localhost:6333/collections/signifiers \
      -H 'Content-Type: application/json' \
      -d '{"vectors": {"size": 768, "distance": "Cosine"}}'
    echo ""
fi

echo "✅ Collections verified"
echo ""

# Build migration tool
echo "🔨 Building migration tool..."
cd "$(dirname "$0")/.."
go build -o migrate_to_qdrant cmd/migrate_to_qdrant.go
echo "✅ Build complete"
echo ""

# Run migration
echo "🚀 Running migration..."
./migrate_to_qdrant

# Cleanup
rm -f migrate_to_qdrant

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Migration complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📊 Verify migration:"
echo "  curl http://localhost:6333/collections/memories"
echo "  curl http://localhost:6333/collections/signifiers"
echo ""
echo "🎨 View in dashboard:"
echo "  http://104.248.219.200:8888"
echo ""
