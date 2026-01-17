#!/usr/bin/env python3
"""
Seed Resonance Scripts to Qdrant
Loads hypnotic scripts with SSML and voice settings
"""

import json
import sys
from qdrant_client import QdrantClient
from qdrant_client.http import models

# Configuration
QDRANT_URL = "localhost"
QDRANT_PORT = 6333
OLLAMA_URL = "http://localhost:11434"
COLLECTION_NAME = "resonance_scripts"
SCRIPTS_FILE = "/root/EVA-Mind-FZPN/data/resonance_scripts_core.json"

def generate_embedding_ollama(text: str):
    """Generate embedding using Ollama"""
    import requests
    
    try:
        response = requests.post(
            f"{OLLAMA_URL}/api/embeddings",
            json={
                "model": "nomic-embed-text",
                "prompt": text[:2000]
            },
            timeout=30
        )
        if response.status_code == 200:
            return response.json()["embedding"]
        return None
    except Exception as e:
        print(f"❌ Erro embedding: {e}")
        return None

def main():
    print("=" * 70)
    print("🌊 DEEP RESONANCE ENGINE - SEED SCRIPTS")
    print("=" * 70)
    print()
    
    # 1. Connect to Qdrant
    print("🔗 Conectando ao Qdrant...")
    client = QdrantClient(host=QDRANT_URL, port=QDRANT_PORT)
    
    # 2. Create collection
    print(f"🔧 Criando collection '{COLLECTION_NAME}'...")
    
    try:
        client.delete_collection(collection_name=COLLECTION_NAME)
        print("   ⚠️  Collection existente deletada")
    except:
        pass
    
    client.create_collection(
        collection_name=COLLECTION_NAME,
        vectors_config=models.VectorParams(
            size=768,
            distance=models.Distance.COSINE
        ),
        on_disk_payload=True
    )
    print("   ✅ Collection criada")
    print()
    
    # 3. Load scripts
    print(f"📖 Carregando scripts de {SCRIPTS_FILE}...")
    with open(SCRIPTS_FILE, 'r', encoding='utf-8') as f:
        scripts = json.load(f)
    
    print(f"   ✅ {len(scripts)} scripts carregados")
    print()
    
    # 4. Process and insert
    print("📥 Gerando embeddings e inserindo...")
    print()
    
    points = []
    for idx, script in enumerate(scripts, 1):
        payload = script['payload']
        
        # Generate text for embedding
        embed_text = (
            f"{payload['title']} "
            f"{payload['category']} "
            f"{' '.join(payload['target_symptom'])}"
        )
        
        print(f"   [{idx}/{len(scripts)}] {payload['title']}")
        print(f"       Categoria: {payload['category']}")
        print(f"       Sintomas: {', '.join(payload['target_symptom'][:2])}...")
        
        # Generate embedding
        vector = generate_embedding_ollama(embed_text)
        
        if vector is None:
            print(f"       ❌ Falha no embedding")
            continue
        
        # Create point
        points.append(models.PointStruct(
            id=idx,
            vector=vector,
            payload=payload
        ))
        
        print(f"       ✅ Embedding gerado ({len(vector)} dims)")
        print()
    
    # 5. Upload to Qdrant
    print("📤 Fazendo upload para Qdrant...")
    client.upsert(
        collection_name=COLLECTION_NAME,
        points=points
    )
    print("   ✅ Upload concluído")
    print()
    
    # 6. Verify
    print("🔍 Verificando...")
    collection_info = client.get_collection(collection_name=COLLECTION_NAME)
    print(f"   📊 Points no Qdrant: {collection_info.points_count}")
    print()
    
    print("=" * 70)
    print("✨ Deep Resonance Engine - Scripts Carregados!")
    print()
    print("Categorias disponíveis:")
    categories = set(s['payload']['category'] for s in scripts)
    for cat in sorted(categories):
        count = sum(1 for s in scripts if s['payload']['category'] == cat)
        print(f"   • {cat}: {count} scripts")
    print("=" * 70)

if __name__ == "__main__":
    main()
