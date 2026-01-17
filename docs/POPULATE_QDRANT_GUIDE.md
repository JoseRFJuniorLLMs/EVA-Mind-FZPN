# 🚀 Guia de Execução - Popular Qdrant no Servidor

## 📋 Pré-requisitos

Antes de executar, certifique-se que:

- ✅ Qdrant está rodando: `sudo systemctl status qdrant`
- ✅ Ollama está rodando: `ollama list`
- ✅ Modelo `nomic-embed-text` está instalado: `ollama pull nomic-embed-text`

---

## 🎯 Execução Rápida (Tudo de Uma Vez)

```bash
# 1. SSH no servidor
ssh root@104.248.219.200

# 2. Ir para o diretório
cd /root/EVA-Mind-FZPN

# 3. Pull do código mais recente
git pull origin main

# 4. Dar permissão de execução
chmod +x scripts/populate_all_collections.sh

# 5. EXECUTAR TUDO
./scripts/populate_all_collections.sh
```

**Tempo estimado:** 15-20 minutos (depende da velocidade do Ollama)

---

## 📊 O Que Será Criado

### **4 Collections no Qdrant:**

| Collection | Itens | Função | Zeta Affinity |
|------------|-------|--------|---------------|
| `nasrudin_stories` | ~270 | Paradoxo/Humor | 2,4,7,9 (Emocionais) |
| `aesop_fables` | ~300 | Moral/Lógica | 1,3,5,6 (Racionais) |
| `zen_koans` | ~50 | Insight/Silêncio | 1,4,5,9 (Introspectivos) |
| `somatic_exercises` | ~20 | Aterramento | Todos (Crises) |

**Total:** ~640 intervenções terapêuticas indexadas!

---

## 🔍 Verificar Resultado

### **1. Listar Collections**
```bash
curl http://localhost:6333/collections | jq '.result.collections[].name'
```

**Esperado:**
```
"nasrudin_stories"
"aesop_fables"
"zen_koans"
"somatic_exercises"
```

---

### **2. Ver Estatísticas**
```bash
# Nasrudin
curl http://localhost:6333/collections/nasrudin_stories | jq '.result.points_count'

# Esopo
curl http://localhost:6333/collections/aesop_fables | jq '.result.points_count'

# Zen Koans
curl http://localhost:6333/collections/zen_koans | jq '.result.points_count'

# Somático
curl http://localhost:6333/collections/somatic_exercises | jq '.result.points_count'
```

---

### **3. Testar Busca (Exemplo)**

```bash
# Buscar histórias de Nasrudin sobre "projeção"
curl -X POST http://localhost:6333/collections/nasrudin_stories/points/scroll \
  -H 'Content-Type: application/json' \
  -d '{
    "limit": 3,
    "with_payload": true,
    "filter": {
      "must": [
        {"key": "is_clinically_mapped", "match": {"value": true}}
      ]
    }
  }' | jq '.result.points[].payload.title'
```

**Esperado:**
```
"O Burro ao Contrário"
"A Chave e a Luz"
"A Nota Única"
```

---

## 🐛 Troubleshooting

### **Problema: "Connection refused" ao Qdrant**
```bash
# Verificar status
sudo systemctl status qdrant

# Se não estiver rodando, iniciar
sudo systemctl start qdrant

# Verificar logs
sudo journalctl -u qdrant -f
```

---

### **Problema: "Connection refused" ao Ollama**
```bash
# Verificar se está rodando
ps aux | grep ollama

# Se não estiver, iniciar
ollama serve &

# Verificar modelos
ollama list
```

---

### **Problema: Modelo nomic-embed-text não encontrado**
```bash
# Baixar modelo (pode demorar alguns minutos)
ollama pull nomic-embed-text

# Verificar
ollama list | grep nomic
```

---

### **Problema: Script Python falha com erro de módulo**
```bash
# Instalar dependências
pip3 install requests

# Ou se precisar de ambiente virtual
python3 -m venv venv
source venv/bin/activate
pip install requests
```

---

## 📝 Execução Manual (Passo a Passo)

Se preferir executar um por vez:

### **1. Nasrudin (270 histórias)**
```bash
python3 scripts/populate_nasrudin_with_lacan.py
```

### **2. Esopo (~300 fábulas)**
```bash
python3 scripts/populate_aesop_fables.py
```

### **3. Zen (Koans + Somático)**
```bash
python3 scripts/populate_zen_content.py
```

---

## ✅ Checklist de Validação

Após execução, verificar:

- [ ] 4 collections criadas no Qdrant
- [ ] `nasrudin_stories` tem ~270 pontos
- [ ] `aesop_fables` tem ~300 pontos
- [ ] `zen_koans` tem ~50 pontos
- [ ] `somatic_exercises` tem ~20 pontos
- [ ] Busca retorna payloads com `clinical_tags`
- [ ] Histórias mapeadas têm `is_clinically_mapped: true`

---

## 🎯 Próximos Passos

Depois de popular o Qdrant:

1. ✅ **Implementar Backend Go:**
   - `pkg/nasrudin/matcher.go`
   - `pkg/aesop/matcher.go`
   - `pkg/zen/matcher.go`
   - `pkg/somatic/safety_checker.go`

2. ✅ **Integrar com TransNAR:**
   - Conectar detecção de padrões → busca no Qdrant
   - Implementar Zeta Switch (racional vs emocional)

3. ✅ **Testar Fluxo Completo:**
   - Usuário fala → TransNAR analisa → Qdrant busca → LLM narra

4. ✅ **Frontend Flutter:**
   - Implementar cards visuais (Esopo/Nasrudin/Zen)
   - Implementar breathing visualizer (Somático)

---

## 📊 Monitoramento

### **Ver logs em tempo real:**
```bash
# Qdrant
sudo journalctl -u qdrant -f

# EVA-Mind-FZPN
sudo journalctl -u eva-mind-fzpn -f
```

### **Espaço em disco:**
```bash
# Ver tamanho das collections
du -sh /var/lib/qdrant/collections/*
```

---

## 🔄 Re-popular (Se Necessário)

Se precisar limpar e re-popular:

```bash
# Deletar collection específica
curl -X DELETE http://localhost:6333/collections/nasrudin_stories

# Ou deletar todas
for collection in nasrudin_stories aesop_fables zen_koans somatic_exercises; do
  curl -X DELETE http://localhost:6333/collections/$collection
done

# Depois re-executar
./scripts/populate_all_collections.sh
```

---

## 🎉 Resultado Final

**Sistema EVA-Mind-FZPN completo com:**

- 🎭 **Nasrudin** → Quebrar rigidez (Paradoxo)
- 📚 **Esopo** → Ensinar moral (Lógica)
- 🧘 **Zen** → Esvaziar mente (Insight)
- 🫁 **Somático** → Aterrar corpo (Respiração)

**Total:** 640+ intervenções terapêuticas prontas para uso! ✨
