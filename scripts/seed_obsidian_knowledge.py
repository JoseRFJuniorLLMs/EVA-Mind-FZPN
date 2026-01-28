#!/usr/bin/env python3
"""
Seed Obsidian Knowledge to Qdrant
Loads markdown files from Obsidian vault as user knowledge base
"""

import os
import re
import hashlib
from pathlib import Path
from qdrant_client import QdrantClient
from qdrant_client.http import models
import requests

# Configuration
QDRANT_URL = "localhost"
QDRANT_PORT = 6333
OLLAMA_URL = "http://localhost:11434"
COLLECTION_NAME = "user_knowledge"

# User data
USER_CPF = "64525430249"
USER_NAME = "Jose R F Junior"
USER_ID = 1  # ID do usuário no sistema (ajustar conforme necessário)

# Obsidian vault path
OBSIDIAN_VAULT = r"D:\dev\AI\AI-900_Obsidian"


def generate_embedding_ollama(text: str):
    """Generate embedding using Ollama"""
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
        print(f"   Erro embedding: {e}")
        return None


def clean_markdown(content: str) -> str:
    """Remove markdown formatting for better embedding"""
    # Remove code blocks
    content = re.sub(r'```[\s\S]*?```', '', content)
    # Remove inline code
    content = re.sub(r'`[^`]+`', '', content)
    # Remove links but keep text
    content = re.sub(r'\[([^\]]+)\]\([^\)]+\)', r'\1', content)
    # Remove images
    content = re.sub(r'!\[([^\]]*)\]\([^\)]+\)', '', content)
    # Remove headers markers
    content = re.sub(r'^#+\s*', '', content, flags=re.MULTILINE)
    # Remove bold/italic
    content = re.sub(r'\*+([^*]+)\*+', r'\1', content)
    content = re.sub(r'_+([^_]+)_+', r'\1', content)
    # Clean up whitespace
    content = re.sub(r'\n{3,}', '\n\n', content)
    return content.strip()


def extract_metadata(filepath: Path, content: str) -> dict:
    """Extract metadata from file path and content"""
    # Get category from folder
    parts = filepath.relative_to(OBSIDIAN_VAULT).parts
    category = parts[0] if len(parts) > 1 else "General"

    # Get title from filename
    title = filepath.stem.replace('_', ' ')

    # Extract tags from content (Obsidian format: #tag)
    tags = re.findall(r'#(\w+)', content)

    # Extract first heading as summary
    first_heading = re.search(r'^#\s*(.+)$', content, re.MULTILINE)
    summary = first_heading.group(1) if first_heading else title

    return {
        "title": title,
        "category": category,
        "tags": list(set(tags))[:10],
        "summary": summary,
        "source": "obsidian",
        "vault": "AI-900_Obsidian"
    }


def get_all_markdown_files(vault_path: str) -> list:
    """Get all markdown files from vault"""
    vault = Path(vault_path)
    files = []

    for md_file in vault.rglob("*.md"):
        # Skip hidden folders and files
        if any(part.startswith('.') for part in md_file.parts):
            continue
        files.append(md_file)

    return sorted(files)


def generate_point_id(filepath: str) -> int:
    """Generate unique ID from filepath"""
    hash_obj = hashlib.md5(filepath.encode())
    return int(hash_obj.hexdigest()[:15], 16)


def main():
    print("=" * 70)
    print("OBSIDIAN KNOWLEDGE LOADER - EVA Memory System")
    print("=" * 70)
    print()
    print(f"Usuario: {USER_NAME}")
    print(f"CPF: {USER_CPF}")
    print(f"Vault: {OBSIDIAN_VAULT}")
    print()

    # 1. Connect to Qdrant
    print("[1/5] Conectando ao Qdrant...")
    client = QdrantClient(host=QDRANT_URL, port=QDRANT_PORT)
    print("      Conectado!")
    print()

    # 2. Create or recreate collection
    print(f"[2/5] Preparando collection '{COLLECTION_NAME}'...")

    try:
        # Check if collection exists
        collections = client.get_collections()
        collection_names = [c.name for c in collections.collections]

        if COLLECTION_NAME in collection_names:
            print(f"      Collection '{COLLECTION_NAME}' ja existe, adicionando dados...")
        else:
            client.create_collection(
                collection_name=COLLECTION_NAME,
                vectors_config=models.VectorParams(
                    size=768,
                    distance=models.Distance.COSINE
                ),
                on_disk_payload=True
            )
            print("      Collection criada!")
    except Exception as e:
        print(f"      Erro: {e}")
        return
    print()

    # 3. Get all markdown files
    print(f"[3/5] Buscando arquivos markdown...")
    md_files = get_all_markdown_files(OBSIDIAN_VAULT)
    print(f"      Encontrados: {len(md_files)} arquivos")
    print()

    # 4. Process and insert
    print("[4/5] Processando arquivos...")
    print()

    points = []
    categories = {}

    for idx, filepath in enumerate(md_files, 1):
        try:
            # Read file
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()

            if len(content.strip()) < 50:
                print(f"   [{idx}/{len(md_files)}] {filepath.stem} - SKIP (muito curto)")
                continue

            # Extract metadata
            metadata = extract_metadata(filepath, content)

            # Clean content for embedding
            clean_content = clean_markdown(content)

            # Text for embedding (title + content)
            embed_text = f"{metadata['title']}\n\n{clean_content}"

            print(f"   [{idx}/{len(md_files)}] {metadata['title']}")
            print(f"       Categoria: {metadata['category']}")

            # Generate embedding
            vector = generate_embedding_ollama(embed_text)

            if vector is None:
                print(f"       FALHA no embedding")
                continue

            # Track categories
            categories[metadata['category']] = categories.get(metadata['category'], 0) + 1

            # Create point with user association
            point_id = generate_point_id(str(filepath))

            points.append(models.PointStruct(
                id=point_id,
                vector=vector,
                payload={
                    "user_id": USER_ID,
                    "user_cpf": USER_CPF,
                    "user_name": USER_NAME,
                    "title": metadata['title'],
                    "category": metadata['category'],
                    "tags": metadata['tags'],
                    "summary": metadata['summary'],
                    "content": clean_content[:5000],  # Limit content size
                    "source": metadata['source'],
                    "vault": metadata['vault'],
                    "file_path": str(filepath.relative_to(OBSIDIAN_VAULT)),
                    "content_length": len(clean_content)
                }
            ))

            print(f"       OK ({len(vector)} dims)")

        except Exception as e:
            print(f"   [{idx}/{len(md_files)}] ERRO: {e}")
            continue

    print()

    # 5. Upload to Qdrant
    print(f"[5/5] Fazendo upload para Qdrant ({len(points)} documentos)...")

    if points:
        # Upload in batches
        batch_size = 50
        for i in range(0, len(points), batch_size):
            batch = points[i:i + batch_size]
            client.upsert(
                collection_name=COLLECTION_NAME,
                points=batch
            )
            print(f"      Batch {i // batch_size + 1}: {len(batch)} pontos")

        print("      Upload concluido!")
    else:
        print("      Nenhum ponto para upload!")

    print()

    # Summary
    print("=" * 70)
    print("RESUMO")
    print("=" * 70)
    print()
    print(f"Usuario: {USER_NAME} (CPF: {USER_CPF})")
    print(f"Documentos carregados: {len(points)}")
    print()
    print("Categorias:")
    for cat, count in sorted(categories.items(), key=lambda x: -x[1]):
        print(f"   {cat}: {count} docs")
    print()

    # Verify
    collection_info = client.get_collection(collection_name=COLLECTION_NAME)
    print(f"Total de pontos na collection: {collection_info.points_count}")
    print()
    print("=" * 70)
    print("Conhecimento carregado com sucesso!")
    print()
    print("Para buscar, use filtro:")
    print(f'   filter: {{"user_cpf": "{USER_CPF}"}}')
    print("=" * 70)


if __name__ == "__main__":
    main()
