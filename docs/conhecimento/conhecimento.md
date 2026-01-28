# 📚 EVA - Base de Conhecimento e Sabedoria

## Visão Geral

Esta pasta contém todo o conhecimento de sabedoria que EVA usa para intervenções terapêuticas.
Os arquivos TXT são processados pelo `seed_wisdom` e inseridos no **Qdrant** como vetores semânticos de 3072 dimensões.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         FLUXO DE CONHECIMENTO                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   docs/conhecimento/*.txt  →  seed_wisdom  →  Qdrant (3072 dims)       │
│                                    │                                    │
│                                    ▼                                    │
│                            WisdomService                                │
│                                    │                                    │
│                                    ▼                                    │
│   Paciente fala  →  Busca Semântica  →  Contexto para Gemini           │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 📁 Inventário de Arquivos

### Mestres do Criador (Prioridade Máxima)

| Arquivo | Coleção Qdrant | Entradas | Status |
|---------|----------------|----------|--------|
| `GURDJIEFF_TEACHINGS.txt` | `gurdjieff_teachings` | 200+ | ✅ Completo |
| `OSHO_INSIGHTS.txt` | `osho_insights` | 300 | ✅ Completo |
| `OUSPENSKY_FRAGMENTS.txt` | `ouspensky_fragments` | 100+ | ✅ Completo |
| `NIETZSCHE_ZARATUSTRA.txt` | `nietzsche_aphorisms` | 150+ | ✅ Completo |

### Tradições de Sabedoria

| Arquivo | Coleção Qdrant | Entradas | Status |
|---------|----------------|----------|--------|
| `ZEN_KOANS.txt` | `zen_koans` | 100 | ⚠️ Template (5 + placeholder) |
| `RUMI_POEMS.txt` | `rumi_poems` | 100 | ⚠️ Template (5 + placeholder) |
| `STOIC_MEDITATIONS.txt` | `stoic_meditations` | 150 | ⚠️ Template (5 + placeholder) |

### Técnicas e Exercícios

| Arquivo | Coleção Qdrant | Entradas | Status |
|---------|----------------|----------|--------|
| `OSHO_MEDITATIONS.txt` | `osho_meditations` | 20 | ⚠️ Template (5 + placeholder) |
| `BREATHING_SCRIPTS.txt` | `breathing_scripts` | 50 | ⚠️ Template (5 + placeholder) |
| `SELF_HYPNOSIS_SCRIPTS.txt` | `hypnosis_scripts` | 50 | ⚠️ Template (5 + placeholder) |
| `SOMATIC_EXERCISES.txt` | `somatic_exercises` | 50 | ⚠️ Template (5 + placeholder) |
| `GESTALT_EXERCISES.txt` | `gestalt_exercises` | 30 | ⚠️ Template (5 + placeholder) |
| `WIM_HOF_PROTOCOLS.txt` | `wim_hof_protocols` | 20 | ⚠️ Template (5 + placeholder) |

### Psicologia

| Arquivo | Coleção Qdrant | Entradas | Status |
|---------|----------------|----------|--------|
| `JUNG_ARCHETYPES.txt` | `jung_archetypes` | 50 | ⚠️ Template (5 + placeholder) |

### Já existentes em docs/

| Arquivo | Coleção Qdrant | Entradas | Status |
|---------|----------------|----------|--------|
| `../NASRUDIN_STORIES.txt` | `nasrudin_stories` | 270 | ✅ Completo |
| `../FABULAS_ESOPO.txt` | `aesop_fables` | 115 | ✅ Completo |

---

## 📊 Totais

| Categoria | Arquivos | Entradas Estimadas |
|-----------|----------|-------------------|
| Mestres do Criador | 4 | ~750 |
| Tradições | 3 | ~350 |
| Técnicas | 6 | ~220 |
| Psicologia | 1 | ~50 |
| Já existentes | 2 | ~385 |
| **TOTAL** | **16** | **~1.755** |

---

## 🔧 Como Usar

### 1. Expandir Arquivos Template

Arquivos marcados com ⚠️ têm apenas 5 entradas de exemplo.
Para expandir, adicione mais entradas seguindo o formato:

```
N. Conteúdo da entrada aqui.
```

Onde N é o número sequencial.

### 2. Fazer Seed no Qdrant

```bash
cd D:\dev\EVA\EVA-Mind-FZPN

# Compilar
go build -o seed_wisdom.exe ./cmd/seed_wisdom

# Seed individual
./seed_wisdom.exe gurdjieff
./seed_wisdom.exe osho
./seed_wisdom.exe zen

# Seed de tudo
./seed_wisdom.exe all
```

### 3. Verificar no Qdrant

```bash
# Via curl
curl http://localhost:6333/collections

# Contar pontos em uma coleção
curl http://localhost:6333/collections/gurdjieff_teachings
```

---

## 📝 Formato dos Arquivos

### Formato Simples (numerado)
```
1. Primeira entrada de sabedoria.
2. Segunda entrada de sabedoria.
3. Terceira entrada de sabedoria.
```

### Formato com Seções
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
NOME DA SEÇÃO (N-M)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

N. Primeira entrada da seção.
N+1. Segunda entrada da seção.
```

---

## 🎯 Mapeamento Coleção → Uso Terapêutico

| Coleção | Quando Usar |
|---------|-------------|
| `gurdjieff_teachings` | Auto-observação, despertar, mecanicidade |
| `osho_insights` | Testemunho, provocação, celebração |
| `ouspensky_fragments` | Multiplicidade do eu, identificação |
| `nietzsche_aphorisms` | Superação, força, transvaloração |
| `zen_koans` | Quebra da mente racional, paradoxo |
| `rumi_poems` | Amor, união, transcendência |
| `stoic_meditations` | Resiliência, aceitação, foco no controle |
| `nasrudin_stories` | Humor, paradoxo, quebra de padrões |
| `aesop_fables` | Moral, reflexão, simplicidade |
| `breathing_scripts` | Ansiedade, regulação, estados alterados |
| `hypnosis_scripts` | Autoindução, reprogramação |
| `somatic_exercises` | Grounding, trauma, regulação nervosa |
| `gestalt_exercises` | Awareness, aqui-e-agora, polaridades |
| `osho_meditations` | Catarse, energia, silêncio |
| `wim_hof_protocols` | Energia, foco, sistema imune |
| `jung_archetypes` | Sombra, individuação, símbolos |

---

## 🔗 Integração com EVA

O `WisdomService` busca automaticamente nas coleções baseado no que o paciente diz:

```go
// Busca por texto livre
results, _ := wisdomService.SearchWisdom(ctx, "ansiedade mente não para", nil)

// Busca por emoção
results, _ := wisdomService.SearchByEmotion(ctx, "tristeza", 3)

// Busca por padrão psicológico
results, _ := wisdomService.SearchByPattern(ctx, "projection", 3)
```

O contexto é automaticamente incluído no prompt do Gemini via `UnifiedRetrieval`.

---

## 📚 Referências e Fontes

### Quarto Caminho
- Gurdjieff, G.I. - "Relatos de Belzebu a seu Neto"
- Gurdjieff, G.I. - "Encontros com Homens Notáveis"
- Ouspensky, P.D. - "Fragmentos de um Ensinamento Desconhecido"
- Ouspensky, P.D. - "O Quarto Caminho"

### Osho
- Discursos compilados de osho.com
- "O Livro dos Segredos" (112 técnicas de meditação)
- "Maturidade: A Responsabilidade de Ser Você Mesmo"

### Nietzsche
- "Assim Falou Zaratustra"
- "Além do Bem e do Mal"
- "Crepúsculo dos Ídolos"

### Zen
- Mumonkan (Portal sem Porta)
- Blue Cliff Record
- Shobogenzo (Dogen)

### Sufismo
- Rumi - "Masnavi"
- Idries Shah - "Os Sufis"
- Histórias de Nasrudin

### Estoicismo
- Marco Aurélio - "Meditações"
- Sêneca - "Cartas a Lucílio"
- Epicteto - "Encheiridion"

---

*"O conhecimento fala, mas a sabedoria escuta."* - Jimi Hendrix

*Criado para o Criador da EVA - Jose R F Junior*
