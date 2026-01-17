#!/usr/bin/env python3
"""
Zen Stories → Qdrant com Schema Lacaniano
Duas Collections:
1. zen_koans: Histórias para esvaziar a mente (overthinking)
2. somatic_exercises: Exercícios de centramento (pânico/ansiedade)

Uso no servidor:
    python3 populate_zen_content.py
"""

import requests
import json
import re
import time
from typing import List, Dict, Optional

# ============================================================================
# CONFIGURAÇÕES
# ============================================================================
QDRANT_URL = "http://localhost:6333"
OLLAMA_URL = "http://localhost:11434"
COLLECTION_KOANS = "zen_koans"
COLLECTION_SOMATIC = "somatic_exercises"
BOOK_PATH = "/root/EVA-Mind-FZPN/docs/zen.txt"

# ============================================================================
# MAPEAMENTO LACANIANO MANUAL (Koans-Chave)
# ============================================================================
KOAN_MAPPING = {
    "xicara_cha": {
        "title": "Uma Xícara de Chá",
        "transnar_rule": "intellectualization",
        "target_state": "mental_saturation",
        "intervention_type": "shock_insight",
        "zeta_affinity": [1, 5, 6],  # Tipos analíticos
        "trigger_condition": "User overthinks, analyzes excessively, cannot stop mental chatter",
        "eva_followup": "Sua mente parece essa xícara de chá. Estamos tentando colocar mais chá numa xícara que já transbordou. Que tal pararmos de pensar por um minuto?"
    },
    "diamante": {
        "title": "Encontrando um Diamante",
        "transnar_rule": "attachment",
        "target_state": "material_obsession",
        "intervention_type": "detachment",
        "zeta_affinity": [3, 8],
        "trigger_condition": "User obsesses over material loss or gain",
        "eva_followup": "Como o homem do diamante, às vezes nos apegamos tanto ao que temos que esquecemos de viver."
    }
}

# ============================================================================
# MAPEAMENTO SOMATIC (Exercícios de Centramento)
# ============================================================================
SOMATIC_MAPPING = {
    "equilibrio": {
        "title": "Centralização pelo Equilíbrio",
        "instruction": "Tente permanecer igualmente em ambos os pés; então, imagine que você está alternando levemente o seu equilíbrio de um pé para o outro.",
        "symptoms": ["panic_attack", "dizziness", "disassociation"],
        "action": "grounding",
        "duration_seconds": 60,
        "eva_voice_command": "Pare tudo. Fique em pé. Sinta os dois pés no chão. Agora, imagine o peso mudando levemente de um pé para o outro."
    },
    "respiracao": {
        "title": "Atenção na Respiração",
        "instruction": "Quando sua respiração entrar, sinta que você está entrando. Quando sair, sinta que você está saindo.",
        "symptoms": ["anxiety", "hyperventilation", "stress"],
        "action": "breath_awareness",
        "duration_seconds": 120,
        "eva_voice_command": "Feche os olhos. Sinta o ar entrando. Você está entrando com ele. Sinta o ar saindo. Você está saindo com ele."
    },
    "corpo_presente": {
        "title": "Sentir o Corpo Inteiro",
        "instruction": "Sinta seu corpo inteiro como uma única presença, sem divisões.",
        "symptoms": ["fragmentation", "disassociation", "numbness"],
        "action": "body_integration",
        "duration_seconds": 90,
        "eva_voice_command": "Sinta seu corpo todo de uma vez. Não pense nas partes. Sinta a presença inteira, como uma unidade."
    }
}

# ============================================================================
# FUNÇÕES AUXILIARES
# ============================================================================

def print_progress(current, total, prefix='', suffix=''):
    """Barra de progresso visual"""
    percent = int(100 * current / total)
    filled = int(40 * current / total)
    bar = '█' * filled + '-' * (40 - filled)
    print(f'\r{prefix} |{bar}| {percent}% {suffix}', end='\r')
    if current == total:
        print()

def generate_embedding_ollama(text: str) -> Optional[List[float]]:
    """Gera embedding usando Ollama"""
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
        print(f"\n❌ Erro embedding: {e}")
        return None

def parse_zen_content(file_path: str) -> tuple[List[Dict], List[Dict]]:
    """Parse do arquivo zen.txt - separa koans e exercícios"""
    print("\n📖 Lendo conteúdo Zen...")
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    koans = []
    exercises = []
    
    # Detectar seção de Centering (Vigyan Bhairav Tantra)
    centering_start = content.find("Ó Shiva")
    
    if centering_start > 0:
        koan_section = content[:centering_start]
        centering_section = content[centering_start:]
    else:
        koan_section = content
        centering_section = ""
    
    # Parse Koans (histórias curtas)
    # Pattern: Título em maiúsculas ou número seguido de texto
    koan_pattern = r'([A-ZÁÀÂÃÉÈÊÍÏÓÔÕÖÚÇÑ\s]+)\n\n(.*?)(?=\n\n[A-ZÁÀÂÃÉÈÊÍÏÓÔÕÖÚÇÑ\s]+\n\n|\Z)'
    koan_matches = re.findall(koan_pattern, koan_section, re.DOTALL)
    
    for idx, (title, text) in enumerate(koan_matches[:50], 1):  # Limitar a 50 koans
        title = title.strip()
        text = text.strip()
        
        if len(text) < 100 or len(text) > 2000:
            continue
        
        koan_id = f"zen_koan_{idx:03d}"
        
        payload = {
            "koan_id": koan_id,
            "title": title,
            "text": text,
            "source": "A Carne e os Ossos do Zen - Paul Reps",
            "language": "pt-BR",
            "type": "narrative_koan"
        }
        
        # Adicionar mapeamento se existir
        key = title.lower().replace(" ", "_")[:20]
        if key in KOAN_MAPPING:
            mapping = KOAN_MAPPING[key]
            payload["clinical_tags"] = {
                "transnar_rule": mapping["transnar_rule"],
                "target_state": mapping["target_state"],
                "intervention_type": mapping["intervention_type"],
                "zeta_affinity": mapping["zeta_affinity"]
            }
            payload["trigger_condition"] = mapping["trigger_condition"]
            payload["eva_followup"] = mapping["eva_followup"]
            payload["is_clinically_mapped"] = True
        else:
            payload["clinical_tags"] = {
                "zeta_affinity": [1, 4, 5, 9]  # Default: Tipos introspectivos
            }
            payload["is_clinically_mapped"] = False
        
        koans.append(payload)
    
    # Parse Exercícios Somáticos (Centering)
    if centering_section:
        # Extrair instruções (frases que começam com verbos imperativos)
        exercise_pattern = r'([A-Z][^.!?]*(?:sinta|imagine|permaneça|tente|observe)[^.!?]*\.)'
        exercise_matches = re.findall(exercise_pattern, centering_section, re.IGNORECASE)
        
        for idx, instruction in enumerate(exercise_matches[:20], 1):  # Limitar a 20
            instruction = instruction.strip()
            
            if len(instruction) < 30:
                continue
            
            exercise_id = f"somatic_{idx:03d}"
            
            payload = {
                "exercise_id": exercise_id,
                "title": f"Exercício de Centramento {idx}",
                "instruction": instruction,
                "source": "Vigyan Bhairav Tantra (via Paul Reps)",
                "type": "somatic_exercise",
                "duration_seconds": 60,
                "is_clinically_mapped": False
            }
            
            # Adicionar mapeamento manual se existir
            for key, mapping in SOMATIC_MAPPING.items():
                if key in instruction.lower():
                    payload.update(mapping)
                    payload["is_clinically_mapped"] = True
                    break
            
            if not payload["is_clinically_mapped"]:
                payload["symptoms"] = ["anxiety", "stress"]
                payload["action"] = "mindfulness"
            
            exercises.append(payload)
    
    print(f"✅ Encontrados {len(koans)} koans e {len(exercises)} exercícios\n")
    return koans, exercises

def create_collection(name: str):
    """Cria collection no Qdrant"""
    print(f"\n🔧 Criando collection '{name}'...")
    
    # Deletar se existir
    response = requests.get(f"{QDRANT_URL}/collections/{name}")
    if response.status_code == 200:
        requests.delete(f"{QDRANT_URL}/collections/{name}")
        time.sleep(1)
    
    # Criar nova
    payload = {
        "vectors": {
            "size": 768,
            "distance": "Cosine"
        },
        "on_disk_payload": True
    }
    
    response = requests.put(
        f"{QDRANT_URL}/collections/{name}",
        json=payload
    )
    
    if response.status_code == 200:
        print(f"✅ Collection '{name}' criada\n")
    else:
        print(f"❌ Erro: {response.text}")
        exit(1)

def insert_item(collection: str, item: Dict, point_id: int) -> bool:
    """Insere item no Qdrant"""
    
    # Gerar texto para embedding
    if item.get("type") == "narrative_koan":
        embed_text = f"{item.get('trigger_condition', '')}. {item['title']}. {item['text'][:500]}"
    else:  # somatic_exercise
        embed_text = f"{item.get('symptoms', [])}. {item['instruction']}"
    
    vector = generate_embedding_ollama(embed_text)
    if vector is None:
        return False
    
    point = {
        "id": point_id,
        "vector": vector,
        "payload": item
    }
    
    response = requests.put(
        f"{QDRANT_URL}/collections/{collection}/points",
        json={"points": [point]}
    )
    
    return response.status_code == 200

# ============================================================================
# MAIN
# ============================================================================

def main():
    print("=" * 70)
    print("🧘 ZEN (PRESENÇA/CENTRAMENTO) → QDRANT")
    print("   Collection 1: zen_koans (Esvaziar a Mente)")
    print("   Collection 2: somatic_exercises (Aterramento)")
    print("=" * 70)
    
    # 1. Parse
    koans, exercises = parse_zen_content(BOOK_PATH)
    
    # 2. Criar collections
    create_collection(COLLECTION_KOANS)
    create_collection(COLLECTION_SOMATIC)
    
    # 3. Inserir Koans
    print(f"📥 Inserindo {len(koans)} koans...\n")
    success_koans = 0
    for idx, koan in enumerate(koans, 1):
        print_progress(idx, len(koans), prefix='Koans:', suffix=f'✅ {success_koans}')
        if insert_item(COLLECTION_KOANS, koan, idx):
            success_koans += 1
        time.sleep(0.3)
    
    # 4. Inserir Exercícios
    print(f"\n📥 Inserindo {len(exercises)} exercícios...\n")
    success_exercises = 0
    for idx, exercise in enumerate(exercises, 1):
        print_progress(idx, len(exercises), prefix='Exercícios:', suffix=f'✅ {success_exercises}')
        if insert_item(COLLECTION_SOMATIC, exercise, idx):
            success_exercises += 1
        time.sleep(0.3)
    
    # 5. Verificar
    print("\n" + "=" * 70)
    print(f"\n📊 RESULTADO:")
    print(f"   ✅ Koans: {success_koans}/{len(koans)}")
    print(f"   ✅ Exercícios: {success_exercises}/{len(exercises)}")
    print("\n✨ Zen (Presença/Centramento) estabelecido!")
    print("   → Para Overthinking, Ansiedade, Pânico")
    print("=" * 70)

if __name__ == "__main__":
    main()
