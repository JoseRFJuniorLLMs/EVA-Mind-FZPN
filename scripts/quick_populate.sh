#!/bin/bash
# Quick execution guide for populating Qdrant on server

echo "🚀 EVA-Mind-FZPN - Qdrant Population"
echo "===================================="
echo ""

# Step 1: Install dependencies
echo "📦 Step 1: Installing Python dependencies..."
pip3 install --break-system-packages qdrant-client requests

# Step 2: Verify Qdrant is running
echo ""
echo "🔍 Step 2: Checking Qdrant..."
curl -s http://localhost:6333/collections || echo "⚠️  Qdrant not responding"

# Step 3: Verify Ollama is running
echo ""
echo "🔍 Step 3: Checking Ollama..."
curl -s http://localhost:11434/api/tags || echo "⚠️  Ollama not responding"

# Step 4: Run population scripts
echo ""
echo "📥 Step 4: Populating collections..."

echo "  → Nasrudin stories..."
python3 scripts/populate_nasrudin_with_lacan.py

echo "  → Aesop fables..."
python3 scripts/populate_aesop_fables.py

echo "  → Zen koans + somatic..."
python3 scripts/populate_zen_content.py

echo "  → Resonance scripts..."
python3 scripts/seed_resonance_scripts.py

# Step 5: Verify results
echo ""
echo "✅ Step 5: Verifying collections..."
curl -s http://localhost:6333/collections | python3 -m json.tool

echo ""
echo "🎉 Done! EVA-Mind-FZPN Qdrant is ready!"
