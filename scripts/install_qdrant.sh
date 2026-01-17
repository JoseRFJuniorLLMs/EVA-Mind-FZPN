#!/bin/bash

# ============================================
# Qdrant Installation Script
# EVA-Mind FZPN Production Server
# ============================================

set -e  # Exit on error

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 Installing Qdrant Vector Database"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 1. Create data directory
echo "📁 Creating data directory..."
sudo mkdir -p /var/lib/qdrant
sudo chown -R $USER:$USER /var/lib/qdrant
echo "✅ Data directory created: /var/lib/qdrant"
echo ""

# 2. Pull Qdrant Docker image
echo "📦 Pulling Qdrant Docker image..."
docker pull qdrant/qdrant:latest
echo "✅ Image pulled"
echo ""

# 3. Stop existing Qdrant container (if any)
echo "🛑 Stopping existing Qdrant container..."
docker stop qdrant 2>/dev/null || true
docker rm qdrant 2>/dev/null || true
echo "✅ Cleanup complete"
echo ""

# 4. Start Qdrant container
echo "🚀 Starting Qdrant container..."
docker run -d \
  --name qdrant \
  --restart always \
  -p 6333:6333 \
  -p 6334:6334 \
  -v /var/lib/qdrant:/qdrant/storage \
  qdrant/qdrant:latest

echo "✅ Qdrant container started"
echo ""

# 5. Wait for Qdrant to be ready
echo "⏳ Waiting for Qdrant to be ready..."
sleep 5

# 6. Health check
echo "🏥 Checking Qdrant health..."
HEALTH=$(curl -s http://localhost:6333/health)
echo "Response: $HEALTH"
echo ""

# 7. Create collections
echo "📚 Creating collections..."

# Collection: memories
echo "  Creating 'memories' collection..."
curl -X PUT http://localhost:6333/collections/memories \
  -H 'Content-Type: application/json' \
  -d '{
    "vectors": {
      "size": 768,
      "distance": "Cosine"
    }
  }'
echo ""

# Collection: signifiers
echo "  Creating 'signifiers' collection..."
curl -X PUT http://localhost:6333/collections/signifiers \
  -H 'Content-Type: application/json' \
  -d '{
    "vectors": {
      "size": 768,
      "distance": "Cosine"
    }
  }'
echo ""

# Collection: context_priming
echo "  Creating 'context_priming' collection..."
curl -X PUT http://localhost:6333/collections/context_priming \
  -H 'Content-Type: application/json' \
  -d '{
    "vectors": {
      "size": 768,
      "distance": "Cosine"
    }
  }'
echo ""

echo "✅ Collections created"
echo ""

# 8. Verify collections
echo "📋 Verifying collections..."
curl -s http://localhost:6333/collections | jq .
echo ""

# 9. Display info
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Qdrant Installation Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📊 Service Information:"
echo "  HTTP API:  http://localhost:6333"
echo "  gRPC API:  localhost:6334"
echo "  Dashboard: http://localhost:6333/dashboard"
echo "  Data Dir:  /var/lib/qdrant"
echo ""
echo "🔧 Useful Commands:"
echo "  docker logs qdrant -f     # View logs"
echo "  docker restart qdrant     # Restart service"
echo "  docker stop qdrant        # Stop service"
echo ""
echo "📝 Next Steps:"
echo "  1. Add to .env: QDRANT_HOST=localhost"
echo "  2. Add to .env: QDRANT_PORT=6334"
echo "  3. Update EVA-Mind code to use Qdrant"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
