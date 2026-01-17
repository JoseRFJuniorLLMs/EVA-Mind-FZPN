# Instruções para Executar no Servidor

## 1. Fazer Upload do Script

```bash
# No seu computador local (PowerShell)
scp d:\dev\EVA\EVA-Mind-FZPN\scripts\populate_nasrudin_with_lacan.py root@104.248.219.200:/root/EVA-Mind-FZPN/scripts/
```

## 2. No Servidor, Executar

```bash
# SSH no servidor
ssh root@104.248.219.200

# Ir para o diretório
cd /root/EVA-Mind-FZPN

# Executar o script
python3 scripts/populate_nasrudin_with_lacan.py
```

## 3. O Que o Script Faz

✅ Lê as 270 histórias de `docs/book1.txt`
✅ Gera embeddings usando **Ollama** (nomic-embed-text)
✅ Aplica **Schema Lacaniano** nas histórias-chave:
   - História 215: A Chave e a Luz → negation_as_desire
   - História 250: A Nota Única → compulsive_repetition
   - História 208: O Burro ao Contrário → projection
   - História 206: O Gato e a Carne → internal_contradiction
   - História 233: A Lua no Poço → reactive_formation

✅ Insere no Qdrant com payload completo:
   - `transnar_rule`: Qual regra TransNAR ativa
   - `trigger_condition`: Quando usar a história
   - `eva_followup`: Frase pós-história
   - `clinical_tags`: Conceitos lacanianos

## 4. Resultado Esperado

```
======================================================================
🧠 PONTE LACAN-NASRUDIN → QDRANT
======================================================================

📖 Lendo histórias de Nasrudin...
✅ Encontradas 270 histórias

📊 Total: 270 histórias (5 com mapeamento Lacaniano)

🔧 Configurando Qdrant...
✅ Collection 'nasrudin_stories' criada

📥 Inserindo no Qdrant...

Progresso (270/270): |████████████████████████████████████████| 100% ✅ 270 | ❌ 0

======================================================================

📊 RESULTADO:
   ✅ Inseridas: 270
   ❌ Falhas: 0
   🧠 Com Schema Lacaniano: 5
   📦 Points no Qdrant: 270

✨ Ponte Lacan-Nasrudin estabelecida!
======================================================================
```

## 5. Testar a Busca

Depois de popular, teste se a busca semântica funciona:

```bash
# Buscar história para "culpar os outros"
curl -X POST http://localhost:6333/collections/nasrudin_stories/points/search \
  -H 'Content-Type: application/json' \
  -d '{
    "vector": [0.1, 0.2, ...],  # Embedding de "User blames others"
    "limit": 3,
    "with_payload": true,
    "filter": {
      "must": [
        {"key": "is_clinically_mapped", "match": {"value": true}}
      ]
    }
  }'
```

## 6. Próximos Passos

Depois de popular o Qdrant:

1. ✅ Implementar `pkg/nasrudin/matcher.go` (busca no Qdrant)
2. ✅ Integrar com TransNAR (detector.go)
3. ✅ Criar narrator.go (LLM conta a história)
4. ✅ Testar fluxo completo

---

**Nota:** O script usa apenas 5 histórias mapeadas manualmente como prova de conceito. 
Depois podemos expandir o mapeamento Lacaniano para mais histórias conforme necessário.
