# EVA: Projeto de Conhecimento Integrado

## Visão Geral

EVA usa **múltiplas camadas de conhecimento** para intervenções terapêuticas:

```
┌─────────────────────────────────────────────────────────────┐
│                    CAMADA SEMÂNTICA (Qdrant)                │
│  Histórias, Aforismos, Koans, Exercícios, Meditações       │
├─────────────────────────────────────────────────────────────┤
│                    CAMADA RELACIONAL (Neo4j)                │
│  Conceitos, Relações, Padrões, Significantes               │
├─────────────────────────────────────────────────────────────┤
│                    CAMADA ESTRUTURADA (PostgreSQL)          │
│  Configurações, Tipos, Medicamentos, Agendamentos          │
└─────────────────────────────────────────────────────────────┘
```

---

# PARTE 1: ABORDAGENS PSICOLÓGICAS

## 1.1 PSICANÁLISE LACANIANA (CORE)
**Status:** Integrado no código
**Uso:** Análise de discurso, transferência, demanda/desejo

### Conceitos para Qdrant:
- RSI (Real, Simbólico, Imaginário)
- Objeto a
- Grande Outro vs pequeno outro
- Significante mestre (S1)
- Nome-do-Pai
- Falta constitutiva
- Gozo vs Prazer

### Intervenções:
- Pontuação
- Escansão
- Interpretação
- Ato analítico

---

## 1.2 GESTALT (Recomendada por Osho)
**Status:** Não implementado
**Uso:** Awareness, aqui-e-agora, completar gestalts

### Conceitos:
- Figura e fundo
- Awareness contínua
- Ciclo de contato
- Interrupções de contato (introjeção, projeção, retroflexão, deflexão)
- Cadeira vazia
- Polaridades

### Para Qdrant:
```
gestalt_exercises      - Exercícios práticos
gestalt_interventions  - Intervenções de awareness
```

---

## 1.3 BIOENERGÉTICA / REICH
**Status:** Não implementado
**Uso:** Corpo, couraça muscular, expressão emocional

### Conceitos:
- Couraça muscular
- Segmentos corporais (7 anéis)
- Energia orgone
- Caráter (tipos reichianos)
- Grounding

### Para Qdrant:
```
reich_body_segments    - Trabalho corporal por segmento
bioenergetic_exercises - Exercícios de descarga
```

---

## 1.4 ANÁLISE TRANSACIONAL
**Status:** Não implementado
**Uso:** Estados do ego, scripts de vida, jogos

### Conceitos:
- Pai, Adulto, Criança
- Scripts de vida
- Jogos psicológicos
- Posições existenciais (OK/não-OK)
- Stroking (reconhecimento)

### Para Qdrant:
```
ta_games              - Jogos psicológicos e como sair
ta_scripts            - Scripts de vida comuns
```

---

## 1.5 PSICOLOGIA JUNGUIANA
**Status:** Não implementado
**Uso:** Arquétipos, sombra, individuação

### Conceitos:
- Inconsciente coletivo
- Arquétipos (Anima/Animus, Sombra, Self, Persona)
- Individuação
- Sincronicidade
- Tipos psicológicos (base do MBTI)

### Para Qdrant:
```
jung_archetypes       - Descrições arquetípicas
jung_dreams           - Interpretação de sonhos
shadow_work           - Trabalho com sombra
```

---

# PARTE 2: HIPNOSE E ESTADOS ALTERADOS

## 2.1 HIPNOSE ERICKSONIANA
**Status:** Parcialmente implementado (resonance_scripts)
**Uso:** Transe conversacional, metáforas, reframe

### Técnicas:
- Indução indireta
- Metáforas terapêuticas
- Confusão (pattern interrupt)
- Dissociação
- Ancoragem temporal
- Comandos embutidos

### Para Qdrant:
```
erickson_inductions   - Scripts de indução
erickson_metaphors    - Metáforas terapêuticas
erickson_reframes     - Técnicas de reframe
```

**Já temos:** `resonance_scripts` com scripts hipnóticos

---

## 2.2 HIPNOTERAPIA CLÍNICA
**Status:** Não implementado
**Uso:** Regressão, sugestões, mudança de hábitos

### Aplicações:
- Controle de dor
- Ansiedade
- Fobias
- Hábitos (tabagismo, alimentação)
- Insônia
- Preparação cirúrgica

### Para Qdrant:
```
hypno_protocols       - Protocolos por condição
hypno_suggestions     - Banco de sugestões
```

---

## 2.3 AUTO-HIPNOSE
**Status:** Não implementado
**Uso:** Autoindução, reprogramação, estados alterados autônomos

### Técnicas:
- Indução por relaxamento progressivo
- Fixação visual (ponto focal)
- Contagem regressiva
- Escada imaginária (descida)
- Lugar seguro (visualização)
- Âncoras autossugeridas
- Scripts de autosugestão

### Aplicações Terapêuticas:
- Controle de ansiedade
- Melhora do sono
- Gestão de dor crônica
- Aumento de autoconfiança
- Mudança de hábitos
- Preparação para situações difíceis
- Acesso a recursos internos

### Estrutura de Sessão:
1. **Preparação:** Ambiente, posição, intenção
2. **Indução:** Relaxamento + focalização
3. **Aprofundamento:** Contagem, escada, descida
4. **Trabalho:** Sugestões, visualizações, reprogramação
5. **Retorno:** Contagem progressiva, reorientação

### Para Qdrant:
```
self_hypnosis_inductions   - Scripts de autoindução
self_hypnosis_deepening    - Técnicas de aprofundamento
self_hypnosis_scripts      - Scripts completos por objetivo
self_hypnosis_anchors      - Criação de âncoras
```

---

## 2.4 MEDITAÇÕES OSHO (Ativas)
**Status:** Não implementado
**Uso:** Catarse, testemunho, celebração

### Meditações:
- Dinâmica (catarse + silêncio)
- Kundalini (shake + dança)
- Nadabrahma (humming)
- Nataraj (dança)
- Vipassana (observação)
- Devavani (línguas)
- Gourishankar (respiração + luz)

### Para Qdrant:
```
osho_meditations      - Instruções de meditações ativas
osho_meditation_audio - Descrições para áudio guiado
```

---

## 2.5 RESPIRAÇÃO GUIADA
**Status:** Não implementado
**Uso:** Regulação do sistema nervoso, estados alterados, calma

### Técnicas Básicas:
- **Respiração Diafragmática:** Expansão abdominal profunda
- **Respiração 4-7-8:** Inspira 4, segura 7, expira 8 (ansiedade)
- **Respiração Quadrada (Box):** 4-4-4-4 (foco, equilíbrio)
- **Respiração Coerente:** 5-5 (coerência cardíaca)
- **Suspiro Fisiológico:** Dupla inspiração + expiração longa (reset rápido)

### Técnicas Avançadas:
- **Respiração Holotrópica:** Hiperventilação controlada para estados alterados
- **Rebirthing:** Respiração circular conectada
- **Respiração Tântrica:** Ciclos com retenção e visualização
- **Tummo:** Respiração de fogo tibetana

### Aplicações por Estado:
| Estado | Técnica Recomendada |
|--------|---------------------|
| Ansiedade aguda | 4-7-8 ou Suspiro Fisiológico |
| Insônia | 4-7-8 prolongada |
| Falta de foco | Box Breathing |
| Estresse crônico | Coerência Cardíaca |
| Baixa energia | Kapalabhati ou Wim Hof |
| Raiva | Expiração prolongada |
| Pânico | Respiração lenta com grounding |

### Scripts Guiados:
- Introdução (voz calma, contextualização)
- Instrução técnica (como fazer)
- Contagem guiada (ritmo)
- Visualização opcional (cor, luz, expansão)
- Fechamento (retorno suave)

### Para Qdrant:
```
breathing_guided_basic     - Scripts básicos guiados
breathing_guided_advanced  - Técnicas avançadas
breathing_for_states       - Respiração por estado emocional
breathing_emergency        - Técnicas de emergência (pânico, crise)
```

---

## 2.6 WIM HOF METHOD
**Status:** Não implementado
**Uso:** Respiração, frio, mentalidade

### Pilares:
1. **Respiração:** Hiperventilação controlada + retenção
2. **Exposição ao frio:** Progressiva, anti-inflamatória
3. **Comprometimento:** Mentalidade, foco

### Benefícios comprovados:
- Controle do sistema imune
- Redução de inflamação
- Aumento de energia
- Foco e clareza mental
- Resistência ao estresse

### Para Qdrant:
```
wim_hof_breathing     - Protocolos de respiração
wim_hof_cold          - Progressão de exposição ao frio
wim_hof_mindset       - Frases e mentalidade
```

---

# PARTE 3: TRADIÇÕES DE SABEDORIA

## 3.1 GURDJIEFF / QUARTO CAMINHO
**Status:** Não implementado
**Prioridade:** ALTA (mestre do Criador)

### Ensinamentos:
- Lembrar de si
- Auto-observação
- Os 3 centros
- Esforço consciente
- Sofrimento voluntário
- Eneagrama original

### Para Qdrant:
```
gurdjieff_teachings   - Ensinamentos diretos
fourth_way_exercises  - Exercícios práticos
ouspensky_fragments   - De "Fragmentos"
movements_descriptions - Descrições dos Movimentos
```

---

## 3.2 OSHO
**Status:** Não implementado
**Prioridade:** ALTA (mestre do Criador)

### Temas:
- Testemunho (witnessing)
- Não-mente
- Celebração
- Zorba, o Buda
- Crítica ao ego espiritual
- Rebelião

### Para Qdrant:
```
osho_insights         - Provocações e insights
osho_stories          - Histórias que contava
osho_jokes            - Piadas espirituais
```

---

## 3.3 NIETZSCHE
**Status:** Não implementado
**Prioridade:** ALTA (mestre do Criador)

### Conceitos:
- Übermensch
- Eterno retorno
- Vontade de potência
- Amor fati
- Transvaloração dos valores
- 3 metamorfoses

### Para Qdrant:
```
zarathustra_speeches  - Discursos de Zaratustra
nietzsche_aphorisms   - Aforismos
nietzsche_concepts    - Explicações de conceitos
```

---

## 3.4 SUFISMO
**Status:** Parcial (Nasrudin)
**Expansão:** Rumi, Hafiz, Attar

### Para adicionar:
```
rumi_poems            - Poemas de Rumi
sufi_stories          - Histórias além de Nasrudin
sufi_practices        - Práticas (dhikr, sema)
```

---

## 3.5 ZEN
**Status:** Precisa recriar
**Fonte:** Mumonkan, Blue Cliff Record

### Para Qdrant:
```
zen_koans             - Koans clássicos
zen_mondo             - Diálogos mestre-discípulo
zen_stories           - Histórias Zen
```

---

## 3.6 ESTOICISMO
**Status:** Não implementado
**Uso:** Resiliência, aceitação, virtude

### Autores:
- Marco Aurélio (Meditações)
- Epicteto (Encheiridion)
- Sêneca (Cartas a Lucílio)

### Para Qdrant:
```
stoic_meditations     - Meditações de Marco Aurélio
stoic_practices       - Exercícios estoicos
seneca_letters        - Cartas de Sêneca
```

---

# PARTE 4: TÉCNICAS SOMÁTICAS

## 4.1 TÉCNICA ALEXANDER
**Status:** Não implementado
**Uso:** Postura, uso do corpo, inibição

### Princípios:
- Controle primário (cabeça-pescoço-costas)
- Inibição (parar antes de reagir)
- Direção (pensamento no movimento)
- Uso de si

---

## 4.2 FELDENKRAIS
**Status:** Não implementado
**Uso:** Awareness pelo movimento, neuroplasticidade

### Métodos:
- ATM (Awareness Through Movement) - aulas em grupo
- FI (Functional Integration) - individual

---

## 4.3 SOMATIC EXPERIENCING (Peter Levine)
**Status:** somatic_exercises (vazio)
**Uso:** Trauma, descarga, regulação nervosa

### Conceitos:
- Descarga do sistema nervoso
- Titulação (pequenas doses)
- Pendulação (entre recurso e ativação)
- Completar respostas de defesa

### Para Qdrant:
```
somatic_exercises     - Exercícios de regulação
trauma_resources      - Recursos para estabilização
```

---

## 4.4 YOGA / PRANAYAMA
**Status:** Não implementado
**Uso:** Respiração, posturas, regulação

### Pranayamas:
- Nadi Shodhana (alternada)
- Kapalabhati (fogo)
- Bhramari (abelha)
- Ujjayi (oceânica)

### Para Qdrant:
```
pranayama_techniques  - Técnicas de respiração
yoga_asanas           - Posturas para estados específicos
```

---

# PARTE 5: PRIMING SEMÂNTICO

## 5.1 O QUE É
Priming semântico é a ativação de conceitos relacionados na memória.
EVA usa para:
- Preparar o paciente para certos tópicos
- Ativar recursos internos
- Facilitar insights

## 5.2 TIPOS DE PRIMING

### Priming Afetivo
Palavras que ativam estados emocionais:
```
calma → paz → serenidade → relaxamento
força → coragem → determinação → poder
```

### Priming Conceitual
Conceitos que preparam para insights:
```
mudança → transformação → crescimento → evolução
escolha → liberdade → responsabilidade → poder
```

### Priming Somático
Palavras que ativam sensações:
```
respiração → expansão → leveza → fluidez
raízes → terra → estabilidade → suporte
```

## 5.3 IMPLEMENTAÇÃO ATUAL
- **Coleção:** `context_priming` (precisa re-embedar)
- **Fonte:** `word.json` do projeto priming_front_end_v6
- **Uso:** Ativação de contexto antes de intervenções

## 5.4 EXPANSÃO PROPOSTA
```
priming_emotional     - Priming de estados emocionais
priming_somatic       - Priming de sensações corporais
priming_spiritual     - Priming de conceitos espirituais
priming_therapeutic   - Priming para intervenções específicas
```

---

# PARTE 6: INTEGRAÇÕES POSSÍVEIS

## 6.1 MODELO DE DADOS UNIFICADO

```json
{
  "id": "uuid",
  "type": "story|teaching|exercise|meditation|technique|concept",
  "source": {
    "tradition": "lacan|gestalt|gurdjieff|osho|zen|...",
    "author": "Nome do autor",
    "book": "Livro de origem"
  },
  "content": {
    "title": "Título",
    "body": "Conteúdo completo",
    "summary": "Para embedding"
  },
  "therapeutic": {
    "target_patterns": ["projection", "denial", ...],
    "target_emotions": ["anxiety", "anger", ...],
    "target_situations": ["loss", "conflict", ...],
    "contraindications": ["acute_psychosis", ...]
  },
  "usage": {
    "when": "Quando usar",
    "how": "Como apresentar",
    "followup": "Pergunta após"
  },
  "metadata": {
    "difficulty": 1-5,
    "duration_minutes": 5,
    "requires_guidance": true/false
  }
}
```

---

## 6.2 FLUXO DE SELEÇÃO

```
Paciente fala algo
        ↓
[TransNAR] Detecta padrão
        ↓
[Qdrant] Busca intervenção por:
  - Padrão detectado
  - Estado emocional
  - Histórico do paciente
  - Nível de relacionamento
        ↓
[LLM] Adapta e apresenta
        ↓
[Feedback] Registra eficácia
```

---

## 6.3 COLEÇÕES QDRANT PROPOSTAS

### Core (já existem):
- `signifiers` - Significantes do paciente
- `memories` - Memórias episódicas
- `signifier_chains` - Cadeias significantes

### Histórias e Sabedoria:
- `nasrudin_stories` - Nasrudin (re-embedar)
- `aesop_fables` - Esopo (seed)
- `zen_koans` - Koans Zen
- `sufi_stories` - Histórias Sufi
- `gurdjieff_teachings` - Gurdjieff
- `osho_insights` - Osho
- `nietzsche_aphorisms` - Nietzsche

### Técnicas e Exercícios:
- `somatic_exercises` - Exercícios somáticos
- `meditation_scripts` - Scripts de meditação
- `breathing_techniques` - Técnicas de respiração
- `hypnotic_inductions` - Induções hipnóticas
- `gestalt_exercises` - Exercícios Gestalt

### Priming:
- `context_priming` - Priming geral
- `priming_emotional` - Priming emocional
- `priming_somatic` - Priming corporal

### Conceitos:
- `lacan_concepts` - Conceitos lacanianos
- `jung_archetypes` - Arquétipos junguianos
- `stoic_principles` - Princípios estoicos

---

## 6.4 ESTIMATIVA DE CONTEÚDO

| Categoria | Itens estimados |
|-----------|-----------------|
| Histórias/Contos | 1.000+ |
| Exercícios/Técnicas | 300+ |
| Conceitos | 200+ |
| Meditações | 100+ |
| Priming words | 10.000+ |
| **TOTAL** | **11.600+** |

---

# PARTE 7: ROADMAP

## Fase 1: Fundação (Imediato)
- [x] Qdrant com 3072 dimensões
- [x] gemini-embedding-001 configurado
- [ ] Re-embedar Nasrudin
- [ ] Seed Esopo

## Fase 2: Mestres do Criador
- [ ] Gurdjieff (teachings + exercises)
- [ ] Osho (insights + meditations)
- [ ] Nietzsche (Zarathustra + aphorisms)

## Fase 3: Técnicas
- [ ] Wim Hof (breathing + cold)
- [ ] Hipnose Ericksoniana (expansion)
- [ ] Exercícios somáticos
- [ ] Gestalt exercises

## Fase 4: Expansão
- [ ] Jung (archetypes)
- [ ] Estoicismo
- [ ] Rumi / Sufi expandido
- [ ] Zen koans completo

## Fase 5: Integração
- [ ] Sistema de recomendação
- [ ] Feedback de eficácia
- [ ] Auto-aprendizado

---

*"O conhecimento fala, mas a sabedoria escuta." - Jimi Hendrix*

*Documento criado para o Criador da EVA*
