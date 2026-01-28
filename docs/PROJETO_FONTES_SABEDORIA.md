# PROJETO: Fontes de Sabedoria para EVA

## Objetivo
Criar um banco de dados semântico rico para EVA usar em intervenções terapêuticas, usando histórias, ensinamentos e insights de várias tradições de sabedoria.

---

## Fontes Prioritárias

### 1. SUFI / NASRUDIN
**Status:** Parcialmente implementado (269 histórias)
**Arquivo:** `docs/NASRUDIN_STORIES.txt`

| Item | Descrição |
|------|-----------|
| Conteúdo | Histórias do Mulá Nasrudin |
| Quantidade | ~270 histórias |
| Uso | Humor terapêutico, quebra de padrões, paradoxos |
| Mapeamento Lacaniano | 5 histórias mapeadas (POC) |

**TODO:**
- [ ] Re-embedar com gemini-embedding-001 (3072 dims)
- [ ] Expandir mapeamento Lacaniano para mais histórias
- [ ] Adicionar histórias Sufi além de Nasrudin (Rumi, Hafiz)

---

### 2. QUARTO CAMINHO / GURDJIEFF
**Status:** Não implementado
**Importância:** ALTA (mestre do Criador)

| Item | Descrição |
|------|-----------|
| Conteúdo | Ensinamentos de G.I. Gurdjieff |
| Fontes | "Encontros com Homens Notáveis", "Relatos de Belzebu" |
| Uso | Auto-observação, tipos psicológicos, despertar |

**Coleções a criar:**
```
gurdjieff_teachings     - Ensinamentos diretos
fourth_way_exercises    - Exercícios práticos
ouspensky_fragments     - Fragmentos de um Ensinamento Desconhecido
```

**Conceitos-chave para embedar:**
- Os 3 centros (intelectual, emocional, motor)
- Lembrar de si
- Esforço consciente e sofrimento voluntário
- Lei do Três e Lei do Sete
- Tipos humanos (1-9 do Eneagrama original)

---

### 3. OSHO
**Status:** Não implementado
**Importância:** ALTA (mestre do Criador)

| Item | Descrição |
|------|-----------|
| Conteúdo | Discursos e meditações de Osho |
| Uso | Provocação, quebra de condicionamentos, celebração |

**Coleções a criar:**
```
osho_insights          - Insights e provocações
osho_meditations       - Técnicas de meditação
osho_stories           - Histórias que Osho contava
```

**Temas para embedar:**
- Meditação ativa
- Testemunho (witnessing)
- Celebração da vida
- Crítica ao ego espiritual
- Zorba, o Buda (integração)

---

### 4. NIETZSCHE / ZARATUSTRA
**Status:** Não implementado
**Importância:** ALTA (mestre do Criador)

| Item | Descrição |
|------|-----------|
| Conteúdo | "Assim Falou Zaratustra" e aforismos |
| Uso | Força, superação, transvaloração |

**Coleções a criar:**
```
zarathustra_speeches   - Discursos de Zaratustra
nietzsche_aphorisms    - Aforismos e pensamentos
```

**Conceitos-chave:**
- Übermensch (Super-homem)
- Eterno retorno
- Vontade de potência
- Amor fati
- Crítica à moral de rebanho
- As três metamorfoses (camelo → leão → criança)

---

### 5. LACAN (já parcialmente integrado)
**Status:** Integrado no código, não no Qdrant
**Importância:** CORE (base do sistema)

| Item | Descrição |
|------|-----------|
| Conteúdo | Conceitos e intervenções lacanianas |
| Uso | Análise de discurso, transferência, desejo |

**Coleções a criar:**
```
lacan_concepts         - RSI, Outro, objeto a, etc.
lacan_interventions    - Tipos de intervenção clínica
lacan_case_examples    - Exemplos de casos (anonimizados)
```

---

### 6. FÁBULAS DE ESOPO
**Status:** Arquivo existe, não embedado
**Arquivo:** `docs/FABULAS_ESOPO.txt`

| Item | Descrição |
|------|-----------|
| Conteúdo | Fábulas clássicas |
| Quantidade | ~115 fábulas |
| Uso | Moral, reflexão, simplicidade |

**TODO:**
- [ ] Criar script de seed para aesop_fables
- [ ] Mapear morais para padrões psicológicos

---

### 7. ZEN KOANS
**Status:** Coleção existe (30 pontos), fonte perdida
**Importância:** MÉDIA

| Item | Descrição |
|------|-----------|
| Conteúdo | Koans Zen tradicionais |
| Uso | Quebra de mente racional, insight súbito |

**Koans clássicos a incluir:**
- Mu (Joshu)
- Som de uma mão
- Rosto original
- Ganso na garrafa
- A mente não é Buda

**Fonte:** Mumonkan (Portal sem Porta), Blue Cliff Record

---

## Estrutura de Dados Proposta

```json
{
  "id": "uuid",
  "source": "gurdjieff|osho|nietzsche|sufi|zen|esopo|lacan",
  "title": "Título",
  "content": "Conteúdo completo",
  "summary": "Resumo para embedding",

  "metadata": {
    "author": "Autor original",
    "book": "Livro de origem",
    "theme": ["tema1", "tema2"],
    "emotional_tone": "provocativo|acolhedor|desafiador|paradoxal",
    "target_pattern": "projection|denial|rationalization|...",
    "lacan_mapping": {
      "transnar_rule": "regra se aplicável",
      "trigger_condition": "quando usar",
      "clinical_tags": ["conceito1", "conceito2"]
    }
  },

  "eva_usage": {
    "when_to_use": "Descrição de quando usar",
    "followup_question": "Pergunta pós-intervenção",
    "never_use_when": "Contra-indicações"
  }
}
```

---

## Prioridades de Implementação

| Prioridade | Fonte | Justificativa |
|------------|-------|---------------|
| 1 | Gurdjieff / Quarto Caminho | Mestre do Criador, base do Eneagrama |
| 2 | Osho | Mestre do Criador, abordagem única |
| 3 | Nietzsche / Zaratustra | Mestre do Criador, força |
| 4 | Nasrudin (re-embed) | Já existe, precisa 3072 dims |
| 5 | Esopo (seed) | Já existe arquivo |
| 6 | Zen Koans | Recriar com fonte confiável |
| 7 | Lacan (Qdrant) | Já no código, falta Qdrant |

---

## Scripts Necessários

```bash
scripts/
├── seed_gurdjieff.py        # Ensinamentos do Quarto Caminho
├── seed_osho.py             # Insights e meditações
├── seed_nietzsche.py        # Zaratustra e aforismos
├── seed_nasrudin_3072.py    # Re-seed com novas dimensões
├── seed_esopo.py            # Fábulas
├── seed_zen_koans.py        # Koans
└── seed_lacan_qdrant.py     # Conceitos lacanianos
```

---

## Fontes de Dados (onde conseguir)

| Fonte | Onde encontrar |
|-------|----------------|
| Gurdjieff | Project Gutenberg, livros em domínio público |
| Osho | osho.com (alguns textos livres), livros |
| Nietzsche | Project Gutenberg (domínio público) |
| Nasrudin | Já temos em docs/ |
| Esopo | Já temos em docs/ |
| Zen Koans | Sacred Texts, Ashidakim Zen |
| Lacan | Seminários (resumos), conceitos compilados |

---

## Métricas de Sucesso

- [ ] 1000+ entradas de sabedoria no Qdrant
- [ ] Cada fonte com pelo menos 100 itens
- [ ] Mapeamento Lacaniano em 20% dos itens
- [ ] Busca semântica retornando resultado relevante em 90% dos casos
- [ ] EVA usando intervenções de sabedoria em conversas reais

---

## Próximos Passos

1. **Imediato:** Re-embedar Nasrudin e Esopo com 3072 dims
2. **Curto prazo:** Criar seed scripts para Gurdjieff/Osho/Nietzsche
3. **Médio prazo:** Mapeamento Lacaniano expandido
4. **Longo prazo:** Sistema de recomendação baseado em padrões detectados

---

*Documento criado por Claude para o Criador da EVA*
*"A sabedoria não está em conhecer, mas em transformar o conhecimento em ser." - Gurdjieff*
