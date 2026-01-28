#!/usr/bin/env python3
"""
Seed Priming Words to Qdrant
Carrega pares prime→target para a EVA entender associações semânticas
USA OLLAMA LOCAL (GRÁTIS)
"""

import json
import time
import requests
from qdrant_client import QdrantClient
from qdrant_client.http import models

# ============ CONFIGURAÇÃO ============
QDRANT_HOST = "104.248.219.200"
QDRANT_PORT = 6333
COLLECTION_NAME = "context_priming"

# Ollama LOCAL (GRÁTIS)
OLLAMA_URL = "http://localhost:11434"
EMBEDDING_MODEL = "nomic-embed-text"
VECTOR_SIZE = 768

# Arquivo de priming
PRIMING_FILE = r"D:\dev\priming_front_end_v6\src\assets\json\word.json"

# Batching
BATCH_SIZE = 50  # Menor para evitar timeout
MAX_RETRIES = 3

# ======================================


def generate_embedding(text: str) -> list[float]:
    """Gera embedding usando Ollama LOCAL (grátis)"""
    try:
        response = requests.post(
            f"{OLLAMA_URL}/api/embeddings",
            json={
                "model": EMBEDDING_MODEL,
                "prompt": text
            },
            timeout=30
        )
        if response.status_code == 200:
            return response.json()["embedding"]
        return None
    except Exception as e:
        print(f"   ERRO embedding: {e}")
        return None


def upsert_with_retry(client, collection_name, points, max_retries=3):
    """Upsert com retry em caso de falha"""
    for attempt in range(max_retries):
        try:
            client.upsert(collection_name=collection_name, points=points)
            return True
        except Exception as e:
            print(f"      Tentativa {attempt + 1}/{max_retries} falhou: {e}")
            if attempt < max_retries - 1:
                time.sleep(2)  # Espera antes de retry
    return False


def main():
    print("=" * 70)
    print("PRIMING SEMANTIC LOADER - EVA Memory System")
    print("Usando OLLAMA LOCAL (grátis)")
    print("=" * 70)
    print()

    # 1. Testar Ollama
    print("[1/5] Testando Ollama...")
    test = generate_embedding("teste")
    if test is None:
        print("      ERRO: Ollama não está rodando!")
        print("      Execute: ollama serve")
        return
    print(f"      OK! ({len(test)} dimensões)")
    print()

    # 2. Conectar ao Qdrant (com compatibilidade desabilitada)
    print("[2/5] Conectando ao Qdrant...")
    print(f"      Host: {QDRANT_HOST}:{QDRANT_PORT}")
    client = QdrantClient(
        host=QDRANT_HOST,
        port=QDRANT_PORT,
        timeout=60,
        prefer_grpc=False
    )
    print("      Conectado!")
    print()

    # 3. Verificar se collection existe e quantos pontos tem
    print(f"[3/5] Verificando collection '{COLLECTION_NAME}'...")

    try:
        info = client.get_collection(collection_name=COLLECTION_NAME)
        existing_points = info.points_count
        print(f"      Collection existe com {existing_points} pontos")

        if existing_points > 0:
            print(f"      Continuando de onde parou...")
            start_from = existing_points
        else:
            start_from = 0
    except:
        print("      Collection não existe, criando...")
        client.create_collection(
            collection_name=COLLECTION_NAME,
            vectors_config=models.VectorParams(
                size=VECTOR_SIZE,
                distance=models.Distance.COSINE
            ),
            on_disk_payload=True
        )
        print("      Collection criada!")
        start_from = 0
    print()

    # 4. Carregar arquivo
    print(f"[4/5] Carregando priming...")
    with open(PRIMING_FILE, 'r', encoding='utf-8') as f:
        priming_pairs = json.load(f)

    total = len(priming_pairs)
    print(f"      {total} pares no arquivo")
    print(f"      Começando do índice {start_from}")
    print()

    # 5. Processar
    print("[5/5] Gerando embeddings e inserindo...")
    print()

    inserted = 0
    errors = 0
    points = []

    for idx in range(start_from, total):
        pair = priming_pairs[idx]

        # Texto para embedding: prime + target
        text = f"{pair['prime']} {pair['target']}"

        # Gerar embedding
        embedding = generate_embedding(text)

        if embedding is None:
            errors += 1
            continue

        # Criar ponto
        points.append(models.PointStruct(
            id=idx + 1,
            vector=embedding,
            payload={
                "prime": pair["prime"],
                "target": pair["target"],
                "pair": f"{pair['prime']} → {pair['target']}",
                "type": "semantic_priming"
            }
        ))

        inserted += 1

        # Progresso
        if (idx + 1) % 100 == 0:
            print(f"   {idx + 1}/{total} processados...")

        # Upload em batches
        if len(points) >= BATCH_SIZE:
            if upsert_with_retry(client, COLLECTION_NAME, points):
                points = []
                time.sleep(0.3)  # Pequena pausa entre batches
            else:
                print(f"   ERRO: Falha após {MAX_RETRIES} tentativas no índice {idx}")
                print(f"   Execute novamente para continuar de onde parou.")
                return

    # Upload resto
    if points:
        upsert_with_retry(client, COLLECTION_NAME, points)

    print()
    print("=" * 70)
    print("RESUMO")
    print("=" * 70)
    print()
    print(f"Total de pares: {total}")
    print(f"Inseridos nesta execução: {inserted}")
    print(f"Erros: {errors}")
    print()

    # Verificar
    collection_info = client.get_collection(collection_name=COLLECTION_NAME)
    print(f"Total no Qdrant: {collection_info.points_count}")
    print()

    if collection_info.points_count >= total:
        # Teste de busca
        print("Teste: buscando 'doctor hospital'...")
        test_vector = generate_embedding("doctor hospital")
        if test_vector:
            results = client.search(
                collection_name=COLLECTION_NAME,
                query_vector=test_vector,
                limit=5
            )
            print("Resultados:")
            for r in results:
                print(f"   {r.payload['pair']} (score: {r.score:.3f})")

        print()
        print("=" * 70)
        print("PRONTO! EVA agora entende associações semânticas.")
        print("=" * 70)
    else:
        print("Execute novamente para continuar...")


if __name__ == "__main__":
    main()
