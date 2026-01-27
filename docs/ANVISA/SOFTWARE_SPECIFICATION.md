# Especificação Técnica de Software
## EVA-Mind-FZPN - Companion IA para Idosos

**Documento:** SRS-EVA-001
**Versão:** 1.0
**Data:** 2025-01-27
**Classificação:** SaMD Classe II (ANVISA RDC 751/2022)
**Norma:** IEC 62304:2006/AMD1:2015

---

## 1. Arquitetura de Sistema

### 1.1 Visão Geral da Arquitetura

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         ARQUITETURA EVA-Mind-FZPN                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        CAMADA DE APRESENTAÇÃO                        │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐  │   │
│  │  │  Mobile App │  │   Web App   │  │  Voice API  │  │ Admin Panel│  │   │
│  │  │  (Flutter)  │  │   (React)   │  │  (Whisper)  │  │  (React)   │  │   │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └─────┬──────┘  │   │
│  └─────────┼────────────────┼────────────────┼───────────────┼─────────┘   │
│            │                │                │               │              │
│            └────────────────┴────────┬───────┴───────────────┘              │
│                                      │                                      │
│                              ┌───────▼───────┐                              │
│                              │   API Gateway │                              │
│                              │    (Kong)     │                              │
│                              └───────┬───────┘                              │
│                                      │                                      │
│  ┌───────────────────────────────────┼───────────────────────────────────┐ │
│  │                        CAMADA DE SERVIÇOS                             │ │
│  │                                   │                                    │ │
│  │    ┌──────────────────────────────┼──────────────────────────────┐    │ │
│  │    │                              │                               │    │ │
│  │    ▼                              ▼                               ▼    │ │
│  │ ┌─────────────┐  ┌─────────────────────────────┐  ┌─────────────────┐ │ │
│  │ │   Cortex    │  │       Hippocampus           │  │      Motor      │ │ │
│  │ │  (Golang)   │  │        (Golang)             │  │    (Golang)     │ │ │
│  │ │             │  │                             │  │                 │ │ │
│  │ │ • Emotional │  │ • Memory Service            │  │ • Alert Worker  │ │ │
│  │ │   Analysis  │  │ • Pattern Miner             │  │ • Notification  │ │ │
│  │ │ • Clinical  │  │ • Superhuman Memory         │  │ • Emergency     │ │ │
│  │ │   Metrics   │  │ • Consciousness Service     │  │   Escalation    │ │ │
│  │ │ • Learning  │  │                             │  │                 │ │ │
│  │ └──────┬──────┘  └──────────────┬──────────────┘  └────────┬────────┘ │ │
│  │        │                        │                          │          │ │
│  └────────┼────────────────────────┼──────────────────────────┼──────────┘ │
│           │                        │                          │            │
│           └────────────────────────┼──────────────────────────┘            │
│                                    │                                        │
│  ┌─────────────────────────────────┼─────────────────────────────────────┐ │
│  │                      CAMADA DE INTEGRAÇÃO                             │ │
│  │                                 │                                      │ │
│  │    ┌────────────────────────────┼────────────────────────────┐        │ │
│  │    │                            │                             │        │ │
│  │    ▼                            ▼                             ▼        │ │
│  │ ┌─────────────┐     ┌─────────────────────┐     ┌─────────────────┐   │ │
│  │ │  LLM API    │     │   External APIs     │     │  Notification   │   │ │
│  │ │ (Anthropic/ │     │ • SAMU (192)        │     │   Services      │   │ │
│  │ │  OpenAI)    │     │ • CVV (188)         │     │ • FCM/APNs      │   │ │
│  │ │             │     │ • EHR Integration   │     │ • Twilio SMS    │   │ │
│  │ └─────────────┘     └─────────────────────┘     └─────────────────┘   │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                    │                                        │
│  ┌─────────────────────────────────┼─────────────────────────────────────┐ │
│  │                        CAMADA DE DADOS                                │ │
│  │                                 │                                      │ │
│  │    ┌────────────────────────────┼────────────────────────────┐        │ │
│  │    │                            │                             │        │ │
│  │    ▼                            ▼                             ▼        │ │
│  │ ┌─────────────┐     ┌─────────────────────┐     ┌─────────────────┐   │ │
│  │ │ PostgreSQL  │     │      Qdrant         │     │     Redis       │   │ │
│  │ │  (Primary)  │     │  (Vector Store)     │     │    (Cache)      │   │ │
│  │ │             │     │                     │     │                 │   │ │
│  │ │ • Users     │     │ • Memory Vectors    │     │ • Sessions      │   │ │
│  │ │ • Sessions  │     │ • Emotional States  │     │ • Rate Limits   │   │ │
│  │ │ • Alerts    │     │ • Pattern Clusters  │     │ • Hot Data      │   │ │
│  │ │ • Audit Log │     │                     │     │                 │   │ │
│  │ └─────────────┘     └─────────────────────┘     └─────────────────┘   │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Tecnologias Utilizadas

#### 1.2.1 Backend

| Componente | Tecnologia | Versão | Justificativa |
|------------|------------|--------|---------------|
| Linguagem principal | Go (Golang) | 1.21+ | Performance, concorrência, type safety |
| Framework HTTP | Chi Router | 5.x | Leve, idiomático, middleware |
| ORM | SQLC | 1.x | Type-safe SQL, geração de código |
| Validação | go-playground/validator | 10.x | Validação robusta de structs |
| Logs | zerolog | 1.x | Structured logging, alta performance |
| Config | viper | 1.x | Configuração multi-formato |

#### 1.2.2 Frontend

| Componente | Tecnologia | Versão | Justificativa |
|------------|------------|--------|---------------|
| Mobile | Flutter | 3.x | Cross-platform, acessibilidade |
| Web | React | 18.x | Ecossistema, componentização |
| UI Components | Material UI | 5.x | Acessibilidade WCAG 2.1 |
| State Management | Redux Toolkit | 2.x | Previsibilidade, DevTools |

#### 1.2.3 Infraestrutura

| Componente | Tecnologia | Versão | Justificativa |
|------------|------------|--------|---------------|
| Container | Docker | 24.x | Portabilidade, isolamento |
| Orquestração | Kubernetes | 1.28+ | Escalabilidade, HA |
| API Gateway | Kong | 3.x | Rate limiting, auth, logs |
| Service Mesh | Istio | 1.x | mTLS, observability |

#### 1.2.4 Banco de Dados

| Componente | Tecnologia | Versão | Justificativa |
|------------|------------|--------|---------------|
| Relacional | PostgreSQL | 15+ | ACID, JSON, extensões |
| Vetorial | Qdrant | 1.x | Busca semântica, embeddings |
| Cache | Redis | 7.x | Performance, pub/sub |
| Mensageria | NATS | 2.x | Baixa latência, clustering |

### 1.3 Requisitos de Hardware

#### 1.3.1 Servidor (Produção)

| Recurso | Mínimo | Recomendado |
|---------|--------|-------------|
| CPU | 8 vCPU | 16 vCPU |
| RAM | 16 GB | 32 GB |
| Disco | 100 GB SSD | 500 GB NVMe |
| Rede | 100 Mbps | 1 Gbps |
| IOPS | 3.000 | 10.000+ |

#### 1.3.2 Cliente (Usuário Final)

| Dispositivo | Mínimo | Recomendado |
|-------------|--------|-------------|
| **Smartphone Android** | | |
| - Versão OS | 8.0 (Oreo) | 11.0+ |
| - RAM | 2 GB | 4 GB+ |
| - Armazenamento | 100 MB livre | 500 MB |
| **Smartphone iOS** | | |
| - Versão OS | iOS 13 | iOS 15+ |
| - Dispositivo | iPhone 6s+ | iPhone 11+ |
| **Tablet** | | |
| - Tela | 7" | 10"+ |
| - RAM | 2 GB | 4 GB |
| **Web (Navegador)** | | |
| - Chrome | 90+ | Latest |
| - Firefox | 88+ | Latest |
| - Safari | 14+ | Latest |
| - Edge | 90+ | Latest |

### 1.4 Requisitos de Conectividade

| Requisito | Especificação |
|-----------|---------------|
| Banda mínima (download) | 1 Mbps |
| Banda mínima (upload) | 512 Kbps |
| Latência máxima | 200 ms |
| Protocolo | HTTPS (TLS 1.3) |
| Porta | 443 |
| Modo offline | Suportado (funcionalidades limitadas) |

---

## 2. Descrição de Algoritmos

### 2.1 Análise Emocional

#### 2.1.1 Visão Geral

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    PIPELINE DE ANÁLISE EMOCIONAL                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  INPUT: Texto do usuário                                                │
│    │                                                                    │
│    ▼                                                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ 1. PRÉ-PROCESSAMENTO                                            │   │
│  │    • Normalização de texto (lowercase, remoção de acentos)      │   │
│  │    • Tokenização                                                 │   │
│  │    • Remoção de stopwords (opcional)                            │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│    │                                                                    │
│    ▼                                                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ 2. DETECÇÃO DE PALAVRAS-CHAVE DE RISCO                          │   │
│  │    • Lista de termos de risco (suicídio, morte, desistir...)    │   │
│  │    • Expressões idiomáticas de risco                            │   │
│  │    • Negações e intensificadores                                │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│    │                                                                    │
│    ▼                                                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ 3. ANÁLISE DE SENTIMENTO (Multi-modelo)                         │   │
│  │    • VADER adaptado para português                              │   │
│  │    • LLM para análise contextual                                │   │
│  │    • Ensemble voting                                            │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│    │                                                                    │
│    ▼                                                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ 4. CLASSIFICAÇÃO EMOCIONAL                                      │   │
│  │    • Valência: -1.0 (negativo) a +1.0 (positivo)               │   │
│  │    • Arousal: 0.0 (calmo) a 1.0 (ativado)                      │   │
│  │    • Dominância: 0.0 (submisso) a 1.0 (dominante)              │   │
│  │    • Emoções discretas: alegria, tristeza, medo, raiva...       │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│    │                                                                    │
│    ▼                                                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ 5. CÁLCULO DE SCORE DE RISCO                                    │   │
│  │    • Combinação ponderada dos indicadores                       │   │
│  │    • Histórico recente (tendência)                              │   │
│  │    • Contexto temporal (horário, frequência)                    │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│    │                                                                    │
│    ▼                                                                    │
│  OUTPUT: EmotionalState {                                               │
│    valence: float64,                                                    │
│    arousal: float64,                                                    │
│    dominance: float64,                                                  │
│    primary_emotion: string,                                             │
│    risk_score: float64,        // 0.0 - 1.0                            │
│    risk_level: enum,           // NORMAL, ATTENTION, ALERT, EMERGENCY  │
│    confidence: float64         // 0.0 - 1.0                            │
│  }                                                                      │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 2.1.2 Algoritmo de Score de Risco

```go
// Pseudocódigo do cálculo de risco
func CalculateRiskScore(input RiskInput) RiskOutput {
    // Pesos configuráveis
    weights := RiskWeights{
        KeywordMatch:    0.30,  // Palavras-chave de risco
        SentimentScore:  0.25,  // Análise de sentimento
        HistoricalTrend: 0.20,  // Tendência nas últimas 24h
        TemporalContext: 0.15,  // Horário de risco (noite)
        FrequencyAnomaly: 0.10, // Padrão de uso anômalo
    }

    // Componentes do score
    keywordScore := detectRiskKeywords(input.Text)      // 0.0 - 1.0
    sentimentScore := analyzeSentiment(input.Text)       // -1.0 a +1.0, normalizado
    trendScore := calculateTrend(input.History)          // 0.0 - 1.0
    temporalScore := getTemporalRisk(input.Timestamp)    // 0.0 - 1.0
    frequencyScore := detectAnomaly(input.SessionData)   // 0.0 - 1.0

    // Score ponderado
    rawScore := (keywordScore * weights.KeywordMatch) +
                (normalizeToPositive(sentimentScore) * weights.SentimentScore) +
                (trendScore * weights.HistoricalTrend) +
                (temporalScore * weights.TemporalContext) +
                (frequencyScore * weights.FrequencyAnomaly)

    // Ajuste por gravidade de palavras-chave
    if containsCriticalKeyword(input.Text) {
        rawScore = max(rawScore, 0.8) // Floor de 0.8 para termos críticos
    }

    // Classificação de nível
    riskLevel := classifyRiskLevel(rawScore)

    return RiskOutput{
        Score:      clamp(rawScore, 0.0, 1.0),
        Level:      riskLevel,
        Confidence: calculateConfidence(input),
        Triggers:   identifyTriggers(input),
    }
}

// Níveis de risco
func classifyRiskLevel(score float64) RiskLevel {
    switch {
    case score >= 0.8:
        return EMERGENCY    // Ação imediata
    case score >= 0.6:
        return ALERT        // Notificar cuidador
    case score >= 0.4:
        return ATTENTION    // Monitoramento aumentado
    default:
        return NORMAL       // Operação normal
    }
}
```

#### 2.1.3 Lista de Palavras-Chave de Risco

| Categoria | Exemplos | Peso |
|-----------|----------|------|
| **Crítico** | suicídio, me matar, acabar com tudo, não aguento mais | 1.0 |
| **Alto** | quero morrer, desistir, não vale a pena, sozinho demais | 0.8 |
| **Moderado** | cansado de viver, ninguém se importa, inútil | 0.6 |
| **Baixo** | triste, deprimido, ansioso, preocupado | 0.3 |

### 2.2 Screening Clínico (PHQ-9 / GAD-7)

#### 2.2.1 PHQ-9 (Patient Health Questionnaire-9)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         ALGORITMO PHQ-9                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  QUESTÕES (últimas 2 semanas):                                          │
│  Q1. Pouco interesse ou prazer em fazer as coisas                      │
│  Q2. Sentir-se para baixo, deprimido ou sem esperança                  │
│  Q3. Dificuldade para dormir/dormir demais                             │
│  Q4. Sentir-se cansado ou com pouca energia                            │
│  Q5. Apetite diminuído ou comendo demais                               │
│  Q6. Sentir-se mal consigo mesmo                                       │
│  Q7. Dificuldade para se concentrar                                    │
│  Q8. Movendo-se/falando devagar ou agitado demais                      │
│  Q9. Pensamentos de que seria melhor estar morto                       │
│                                                                         │
│  ESCALA DE RESPOSTA:                                                    │
│  0 = Nenhuma vez                                                        │
│  1 = Vários dias                                                        │
│  2 = Mais da metade dos dias                                           │
│  3 = Quase todos os dias                                               │
│                                                                         │
│  CÁLCULO:                                                               │
│  score_total = sum(Q1..Q9)   // Range: 0-27                            │
│                                                                         │
│  CLASSIFICAÇÃO:                                                         │
│  ┌────────────┬─────────────────────────────────────┐                  │
│  │ Score      │ Classificação                       │                  │
│  ├────────────┼─────────────────────────────────────┤                  │
│  │ 0-4        │ Mínima ou nenhuma depressão         │                  │
│  │ 5-9        │ Depressão leve                      │                  │
│  │ 10-14      │ Depressão moderada                  │                  │
│  │ 15-19      │ Depressão moderadamente grave       │                  │
│  │ 20-27      │ Depressão grave                     │                  │
│  └────────────┴─────────────────────────────────────┘                  │
│                                                                         │
│  ALERTA ESPECIAL:                                                       │
│  Se Q9 >= 1 → Avaliar risco de suicídio imediatamente                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 2.2.2 GAD-7 (Generalized Anxiety Disorder-7)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         ALGORITMO GAD-7                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  QUESTÕES (últimas 2 semanas):                                          │
│  Q1. Sentir-se nervoso, ansioso ou no limite                           │
│  Q2. Não conseguir parar ou controlar preocupações                     │
│  Q3. Preocupar-se demais com coisas diferentes                         │
│  Q4. Dificuldade para relaxar                                          │
│  Q5. Ficar tão inquieto que é difícil ficar parado                    │
│  Q6. Ficar facilmente irritado ou aborrecido                           │
│  Q7. Sentir medo como se algo terrível fosse acontecer                 │
│                                                                         │
│  ESCALA DE RESPOSTA: (igual PHQ-9)                                      │
│  0 = Nenhuma vez | 1 = Vários dias | 2 = Mais da metade | 3 = Quase todos│
│                                                                         │
│  CÁLCULO:                                                               │
│  score_total = sum(Q1..Q7)   // Range: 0-21                            │
│                                                                         │
│  CLASSIFICAÇÃO:                                                         │
│  ┌────────────┬─────────────────────────────────────┐                  │
│  │ Score      │ Classificação                       │                  │
│  ├────────────┼─────────────────────────────────────┤                  │
│  │ 0-4        │ Ansiedade mínima                    │                  │
│  │ 5-9        │ Ansiedade leve                      │                  │
│  │ 10-14      │ Ansiedade moderada                  │                  │
│  │ 15-21      │ Ansiedade grave                     │                  │
│  └────────────┴─────────────────────────────────────┘                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2.3 Sistema de Memória (Superhuman Memory)

#### 2.3.1 Arquitetura de Memória

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    SISTEMA DE MEMÓRIA SUPERHUMAN                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    MEMÓRIA DE TRABALHO                          │   │
│  │                    (Working Memory)                             │   │
│  │    • Contexto da sessão atual                                   │   │
│  │    • Últimas N mensagens (sliding window)                       │   │
│  │    • Estado emocional corrente                                  │   │
│  │    • TTL: duração da sessão                                     │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    MEMÓRIA EPISÓDICA                            │   │
│  │                    (Episodic Memory)                            │   │
│  │    • Eventos significativos                                     │   │
│  │    • Conversas marcantes                                        │   │
│  │    • Marco temporal + contexto emocional                        │   │
│  │    • TTL: 2 anos (configurável)                                 │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    MEMÓRIA SEMÂNTICA                            │   │
│  │                    (Semantic Memory)                            │   │
│  │    • Fatos sobre o usuário (família, preferências)              │   │
│  │    • Conhecimento extraído das conversas                        │   │
│  │    • Grafos de relacionamento                                   │   │
│  │    • TTL: indefinido (atualizado continuamente)                 │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    MEMÓRIA EMOCIONAL                            │   │
│  │                    (Emotional Memory)                           │   │
│  │    • Padrões emocionais ao longo do tempo                       │   │
│  │    • Triggers identificados                                     │   │
│  │    • Estratégias de coping eficazes                            │   │
│  │    • TTL: indefinido                                            │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  RECUPERAÇÃO:                                                           │
│  1. Query → Embedding (vector)                                          │
│  2. Busca por similaridade em Qdrant                                   │
│  3. Re-ranking por relevância temporal + emocional                     │
│  4. Fusão com contexto atual                                           │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 2.3.2 Algoritmo de Consolidação de Memória

```go
// Consolidação de memória (executado periodicamente)
func ConsolidateMemory(userID int64, session Session) error {
    // 1. Extrair eventos significativos da sessão
    events := extractSignificantEvents(session)

    // 2. Para cada evento, decidir tipo de memória
    for _, event := range events {
        importance := calculateImportance(event)

        if importance >= EPISODIC_THRESHOLD {
            // Armazenar como memória episódica
            episodicMemory := EpisodicMemory{
                UserID:      userID,
                Event:       event.Summary,
                Timestamp:   event.Timestamp,
                Emotion:     event.EmotionalState,
                Importance:  importance,
                Embedding:   generateEmbedding(event.Summary),
            }
            saveToQdrant(episodicMemory)
        }

        // 3. Extrair fatos semânticos
        facts := extractFacts(event)
        for _, fact := range facts {
            updateSemanticMemory(userID, fact)
        }

        // 4. Atualizar padrões emocionais
        updateEmotionalPatterns(userID, event.EmotionalState)
    }

    // 5. Decay de memórias antigas (esquecimento natural)
    applyMemoryDecay(userID)

    return nil
}

// Cálculo de importância
func calculateImportance(event Event) float64 {
    factors := []WeightedFactor{
        {Value: event.EmotionalIntensity, Weight: 0.3},
        {Value: event.Novelty, Weight: 0.25},
        {Value: event.PersonalRelevance, Weight: 0.25},
        {Value: event.Recency, Weight: 0.2},
    }

    return weightedSum(factors)
}
```

### 2.4 Detecção de Padrões Temporais

#### 2.4.1 Padrões Monitorados

| Padrão | Descrição | Algoritmo |
|--------|-----------|-----------|
| **Sono** | Horários de interação indicando insônia | Análise de distribuição horária |
| **Humor** | Tendência de valência ao longo de dias | Média móvel + detecção de tendência |
| **Isolamento** | Redução de interações | Detecção de anomalia (z-score) |
| **Ciclotimia** | Oscilações regulares de humor | Análise de Fourier / autocorrelação |
| **Medicação** | Menções a medicamentos + horários | NER + análise temporal |

#### 2.4.2 Algoritmo de Detecção de Tendência

```go
// Detecção de tendência de humor (últimos N dias)
func DetectMoodTrend(userID int64, days int) TrendResult {
    // Buscar scores de valência dos últimos N dias
    dailyAverages := getMoodAverages(userID, days)

    if len(dailyAverages) < 3 {
        return TrendResult{Trend: INSUFFICIENT_DATA}
    }

    // Calcular regressão linear
    slope, intercept, r2 := linearRegression(dailyAverages)

    // Classificar tendência
    trend := classifyTrend(slope, r2)

    // Detectar pontos de inflexão
    inflections := detectInflectionPoints(dailyAverages)

    return TrendResult{
        Trend:       trend,
        Slope:       slope,
        Confidence:  r2,
        Inflections: inflections,
        Prediction:  predictNext(slope, intercept, dailyAverages),
    }
}

// Classificação de tendência
func classifyTrend(slope, r2 float64) TrendType {
    if r2 < 0.3 {
        return STABLE // Pouca correlação = sem tendência clara
    }

    switch {
    case slope < -0.1:
        return DECLINING   // Piora
    case slope > 0.1:
        return IMPROVING   // Melhora
    default:
        return STABLE
    }
}
```

### 2.5 Geração de Resposta

#### 2.5.1 Pipeline de Geração

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    PIPELINE DE GERAÇÃO DE RESPOSTA                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  INPUT: Mensagem do usuário + Contexto                                  │
│    │                                                                    │
│    ▼                                                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ 1. CONSTRUÇÃO DE CONTEXTO                                       │   │
│  │    • Memória de trabalho (últimas mensagens)                    │   │
│  │    • Memórias relevantes (busca semântica)                      │   │
│  │    • Perfil do usuário (nome, preferências)                     │   │
│  │    • Estado emocional atual                                     │   │
│  │    • Hora do dia / contexto temporal                            │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│    │                                                                    │
│    ▼                                                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ 2. SELEÇÃO DE PROMPT TEMPLATE                                   │   │
│  │    • Baseado no estado emocional                                │   │
│  │    • Baseado no tipo de conversa (acolhimento, screening, etc.) │   │
│  │    • Guardrails de segurança incluídos                          │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│    │                                                                    │
│    ▼                                                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ 3. CHAMADA AO LLM                                               │   │
│  │    • Provider: Anthropic Claude / OpenAI GPT-4                  │   │
│  │    • Temperatura: 0.7 (balanceado)                              │   │
│  │    • Max tokens: 500                                            │   │
│  │    • System prompt com persona EVA                              │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│    │                                                                    │
│    ▼                                                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ 4. PÓS-PROCESSAMENTO                                            │   │
│  │    • Verificação de guardrails (conteúdo inapropriado)          │   │
│  │    • Verificação de tamanho (< 100 palavras padrão)             │   │
│  │    • Ajuste de tom (baseado em feedback histórico)              │   │
│  │    • Inserção de emojis apropriados (configurável)              │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│    │                                                                    │
│    ▼                                                                    │
│  OUTPUT: Resposta final + Metadata                                      │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 2.5.2 Guardrails de Segurança

| Guardrail | Descrição | Ação |
|-----------|-----------|------|
| NO_MEDICAL_ADVICE | Não dar diagnósticos ou receitar | Bloquear + redirecionar |
| NO_SUICIDE_METHODS | Não discutir métodos de suicídio | Bloquear + escalar |
| NO_HARMFUL_CONTENT | Conteúdo prejudicial | Bloquear + log |
| IDENTITY_CLEAR | Sempre identificar-se como IA | Verificar na resposta |
| PROFESSIONAL_REFERRAL | Encaminhar casos graves | Inserir orientação |

### 2.6 Limites de Operação

| Parâmetro | Mínimo | Máximo | Padrão |
|-----------|--------|--------|--------|
| Tamanho de mensagem (input) | 1 char | 2.000 chars | - |
| Tamanho de resposta (output) | 10 chars | 1.000 chars | 200 |
| Mensagens por sessão | - | 100 | - |
| Sessões por dia | - | 10 | - |
| Tempo de sessão | - | 2 horas | - |
| Histórico de contexto | 5 msgs | 20 msgs | 10 |
| Score de risco | 0.0 | 1.0 | - |
| Score de sentimento | -1.0 | +1.0 | - |

---

## 3. Especificações de Interface

### 3.1 Telas Principais

#### 3.1.1 Tela de Conversa (Principal)

```
┌─────────────────────────────────────────────────────────────────┐
│  ☰  EVA - Sua Companheira              🔊  ⚙️  │    18:32      │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ EVA:                                                    │   │
│  │ Boa tarde, Dona Maria! Como está se sentindo hoje?     │   │
│  │                                              14:30      │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│            ┌─────────────────────────────────────────────┐     │
│            │ Estou bem, obrigada! Um pouco cansada.      │     │
│            │                                   14:32     │     │
│            └─────────────────────────────────────────────┘     │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ EVA:                                                    │   │
│  │ Entendo. Dormiu bem esta noite? Às vezes o cansaço    │   │
│  │ pode estar relacionado ao sono.                        │   │
│  │                                              14:33      │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌───────────────────────────────┐                             │
│  │ EVA está digitando...         │                             │
│  └───────────────────────────────┘                             │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│ ┌─────────────────────────────────────────────────────────────┐│
│ │ Digite sua mensagem...                                      ││
│ └─────────────────────────────────────────────────────────────┘│
│                                                                 │
│  [ 🎤 Falar ]                              [ ✉️ Enviar ]        │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │          🆘 PRECISO DE AJUDA URGENTE                    │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

#### 3.1.2 Especificações de UI

| Elemento | Especificação |
|----------|---------------|
| **Fonte padrão** | 20pt, ajustável 18-32pt |
| **Fonte mínima** | 14pt (labels secundários) |
| **Contraste texto** | 7:1 mínimo (WCAG AAA) |
| **Cor de fundo** | #FFFFFF (branco) |
| **Cor de texto** | #1A1A1A (preto) |
| **Cor EVA** | #1E3A5F (azul escuro) |
| **Cor emergência** | #CC0000 (vermelho) |
| **Área de toque mínima** | 48×48 px |
| **Espaçamento** | 16px padding padrão |

### 3.2 Fluxo de Navegação

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         FLUXO DE NAVEGAÇÃO                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│                          ┌─────────────┐                                │
│                          │   SPLASH    │                                │
│                          │   SCREEN    │                                │
│                          └──────┬──────┘                                │
│                                 │                                       │
│                    ┌────────────┴────────────┐                          │
│                    │                         │                          │
│                    ▼                         ▼                          │
│            ┌─────────────┐           ┌─────────────┐                    │
│            │   LOGIN     │           │  ONBOARDING │                    │
│            │   (auth)    │           │  (1ª vez)   │                    │
│            └──────┬──────┘           └──────┬──────┘                    │
│                   │                         │                           │
│                   └────────────┬────────────┘                           │
│                                │                                        │
│                                ▼                                        │
│                        ┌─────────────┐                                  │
│                        │    HOME     │                                  │
│                        │  (Conversa) │◄─────────────────────┐           │
│                        └──────┬──────┘                      │           │
│                               │                             │           │
│         ┌──────────┬──────────┼──────────┬──────────┐       │           │
│         │          │          │          │          │       │           │
│         ▼          ▼          ▼          ▼          ▼       │           │
│   ┌──────────┐┌──────────┐┌──────────┐┌──────────┐┌──────────┐          │
│   │ HISTÓRICO││SCREENING ││EMERGÊNCIA││CONTATOS  ││CONFIG    │          │
│   │(sessões) ││(PHQ/GAD) ││(alerta)  ││(família) ││(ajustes) │          │
│   └────┬─────┘└────┬─────┘└────┬─────┘└────┬─────┘└────┬─────┘          │
│        │           │           │           │           │                │
│        └───────────┴───────────┴───────────┴───────────┘                │
│                                │                                        │
│                                └────────────────────────────────────────┘
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 3.3 Acessibilidade (WCAG 2.1 AA)

| Critério | Implementação | Status |
|----------|---------------|--------|
| **1.1.1** Conteúdo não textual | Alt text em imagens | ✅ |
| **1.3.1** Info e relacionamentos | Estrutura semântica HTML | ✅ |
| **1.4.3** Contraste mínimo | 7:1 (AAA) | ✅ |
| **1.4.4** Redimensionar texto | Até 200% sem perda | ✅ |
| **2.1.1** Teclado | Navegação completa por teclado | ✅ |
| **2.4.1** Blocos de bypass | Skip links | ✅ |
| **2.4.7** Foco visível | Indicador de foco claro | ✅ |
| **3.1.1** Idioma da página | lang="pt-BR" | ✅ |
| **3.2.1** Em foco | Sem mudança de contexto | ✅ |
| **4.1.2** Nome, função, valor | ARIA labels | ✅ |

---

## 4. Banco de Dados

### 4.1 Modelo Entidade-Relacionamento

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    MODELO ENTIDADE-RELACIONAMENTO                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────┐         ┌─────────────┐         ┌─────────────┐       │
│  │   USERS     │         │  SESSIONS   │         │  MESSAGES   │       │
│  ├─────────────┤         ├─────────────┤         ├─────────────┤       │
│  │ id (PK)     │◄───────┤│ user_id(FK) │◄───────┤│ session_id  │       │
│  │ name        │    1:N  │ id (PK)     │    1:N  │ id (PK)     │       │
│  │ email       │         │ started_at  │         │ role        │       │
│  │ phone       │         │ ended_at    │         │ content     │       │
│  │ birth_date  │         │ status      │         │ timestamp   │       │
│  │ created_at  │         └─────────────┘         │ emotional   │       │
│  │ preferences │                                  │ _state      │       │
│  └──────┬──────┘                                  └─────────────┘       │
│         │                                                               │
│         │ 1:N      ┌─────────────┐                                      │
│         └─────────▶│  CONTACTS   │                                      │
│         │          ├─────────────┤                                      │
│         │          │ id (PK)     │                                      │
│         │          │ user_id(FK) │                                      │
│         │          │ name        │                                      │
│         │          │ phone       │                                      │
│         │          │ relation    │                                      │
│         │          │ is_emergency│                                      │
│         │          └─────────────┘                                      │
│         │                                                               │
│         │ 1:N      ┌─────────────┐         ┌─────────────┐              │
│         └─────────▶│   ALERTS    │────────▶│ALERT_ACTIONS│              │
│         │          ├─────────────┤   1:N   ├─────────────┤              │
│         │          │ id (PK)     │         │ id (PK)     │              │
│         │          │ user_id(FK) │         │ alert_id(FK)│              │
│         │          │ level       │         │ action_type │              │
│         │          │ trigger     │         │ actor       │              │
│         │          │ created_at  │         │ timestamp   │              │
│         │          │ resolved_at │         └─────────────┘              │
│         │          └─────────────┘                                      │
│         │                                                               │
│         │ 1:N      ┌─────────────┐                                      │
│         └─────────▶│  SCREENINGS │                                      │
│                    ├─────────────┤                                      │
│                    │ id (PK)     │                                      │
│                    │ user_id(FK) │                                      │
│                    │ type        │  (PHQ9, GAD7)                        │
│                    │ score       │                                      │
│                    │ responses   │  (JSONB)                             │
│                    │ created_at  │                                      │
│                    └─────────────┘                                      │
│                                                                         │
│  ┌─────────────┐                                                        │
│  │ AUDIT_LOGS  │                                                        │
│  ├─────────────┤                                                        │
│  │ id (PK)     │                                                        │
│  │ user_id     │                                                        │
│  │ action      │                                                        │
│  │ resource    │                                                        │
│  │ details     │  (JSONB)                                               │
│  │ ip_address  │                                                        │
│  │ timestamp   │                                                        │
│  └─────────────┘                                                        │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Schema Principal (PostgreSQL)

```sql
-- Tabela de usuários (idosos)
CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    external_id     UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255) UNIQUE,
    phone           VARCHAR(20),
    birth_date      DATE NOT NULL,
    cpf_hash        VARCHAR(64), -- Hash para verificação, não armazenamos CPF em texto
    preferences     JSONB DEFAULT '{}',
    consent_version VARCHAR(10),
    consent_date    TIMESTAMP WITH TIME ZONE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at      TIMESTAMP WITH TIME ZONE -- Soft delete
);

-- Tabela de sessões de conversa
CREATE TABLE sessions (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    started_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ended_at        TIMESTAMP WITH TIME ZONE,
    status          VARCHAR(20) DEFAULT 'active', -- active, ended, timeout
    summary         TEXT,
    emotional_summary JSONB,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de mensagens
CREATE TABLE messages (
    id              BIGSERIAL PRIMARY KEY,
    session_id      BIGINT NOT NULL REFERENCES sessions(id),
    role            VARCHAR(10) NOT NULL, -- 'user' ou 'assistant'
    content         TEXT NOT NULL,
    emotional_state JSONB, -- {valence, arousal, dominance, emotion, risk_score}
    tokens_used     INTEGER,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de contatos de emergência
CREATE TABLE contacts (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    name            VARCHAR(255) NOT NULL,
    phone           VARCHAR(20) NOT NULL,
    relation        VARCHAR(50), -- filho, filha, cuidador, médico
    is_emergency    BOOLEAN DEFAULT false,
    priority        INTEGER DEFAULT 1,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de alertas
CREATE TABLE alerts (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    session_id      BIGINT REFERENCES sessions(id),
    message_id      BIGINT REFERENCES messages(id),
    level           VARCHAR(20) NOT NULL, -- NORMAL, ATTENTION, ALERT, EMERGENCY
    risk_score      DECIMAL(3,2),
    trigger_reason  TEXT,
    trigger_details JSONB,
    status          VARCHAR(20) DEFAULT 'open', -- open, acknowledged, resolved
    resolved_at     TIMESTAMP WITH TIME ZONE,
    resolved_by     VARCHAR(255),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de screenings (PHQ-9, GAD-7)
CREATE TABLE screenings (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    type            VARCHAR(20) NOT NULL, -- PHQ9, GAD7, CSSRS
    score           INTEGER NOT NULL,
    classification  VARCHAR(50),
    responses       JSONB NOT NULL, -- Array de respostas
    flagged         BOOLEAN DEFAULT false, -- Q9 do PHQ-9 positivo
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de logs de auditoria
CREATE TABLE audit_logs (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT,
    actor_id        BIGINT, -- Quem realizou a ação (admin, sistema)
    action          VARCHAR(50) NOT NULL, -- CREATE, READ, UPDATE, DELETE, LOGIN, etc.
    resource_type   VARCHAR(50) NOT NULL, -- user, message, alert, etc.
    resource_id     BIGINT,
    details         JSONB,
    ip_address      INET,
    user_agent      TEXT,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índices para performance
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_created_at ON sessions(created_at);
CREATE INDEX idx_messages_session_id ON messages(session_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_alerts_user_id ON alerts(user_id);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_level ON alerts(level);
CREATE INDEX idx_screenings_user_id ON screenings(user_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Particionamento de audit_logs por mês (para retenção)
-- CREATE TABLE audit_logs_2025_01 PARTITION OF audit_logs
--     FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
```

### 4.3 Política de Retenção

| Tabela | Retenção | Ação |
|--------|----------|------|
| users | Enquanto ativo + 6 meses | Anonimização |
| sessions | 2 anos | Eliminação |
| messages | 2 anos | Eliminação |
| alerts | 5 anos | Arquivamento |
| screenings | 5 anos | Arquivamento |
| audit_logs | 5 anos | Arquivamento |
| contacts | Enquanto usuário ativo | Eliminação com usuário |

### 4.4 Estratégia de Backup

| Tipo | Frequência | Retenção | Localização |
|------|------------|----------|-------------|
| Full | Diário (02:00) | 30 dias | S3 São Paulo + DR |
| Incremental | A cada 6h | 7 dias | S3 São Paulo |
| WAL Archiving | Contínuo | 7 dias | S3 São Paulo |
| Snapshot | Semanal | 90 dias | S3 DR (outra região) |

---

## 5. APIs e Integrações

### 5.1 API REST Principal

#### 5.1.1 Endpoints

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| **Autenticação** | | |
| POST | `/api/v1/auth/login` | Login do usuário |
| POST | `/api/v1/auth/refresh` | Refresh token |
| POST | `/api/v1/auth/logout` | Logout |
| **Conversação** | | |
| POST | `/api/v1/chat/message` | Enviar mensagem |
| GET | `/api/v1/chat/sessions` | Listar sessões |
| GET | `/api/v1/chat/sessions/:id` | Detalhes de sessão |
| **Screenings** | | |
| POST | `/api/v1/screening/start` | Iniciar screening |
| POST | `/api/v1/screening/answer` | Responder questão |
| GET | `/api/v1/screening/history` | Histórico de screenings |
| **Alertas** | | |
| GET | `/api/v1/alerts` | Listar alertas |
| PATCH | `/api/v1/alerts/:id` | Atualizar status |
| **Usuário** | | |
| GET | `/api/v1/user/profile` | Perfil do usuário |
| PATCH | `/api/v1/user/profile` | Atualizar perfil |
| GET | `/api/v1/user/contacts` | Listar contatos |
| POST | `/api/v1/user/contacts` | Adicionar contato |

#### 5.1.2 Exemplo de Request/Response

```json
// POST /api/v1/chat/message
// Request
{
  "session_id": "sess_abc123",
  "content": "Estou me sentindo muito triste hoje"
}

// Response
{
  "message_id": "msg_xyz789",
  "response": {
    "content": "Sinto muito que você esteja se sentindo assim. Quer me contar mais sobre o que está acontecendo?",
    "emotional_analysis": {
      "valence": -0.6,
      "arousal": 0.3,
      "primary_emotion": "sadness",
      "risk_score": 0.35,
      "risk_level": "ATTENTION"
    }
  },
  "session": {
    "id": "sess_abc123",
    "message_count": 5,
    "started_at": "2025-01-27T14:30:00Z"
  }
}
```

### 5.2 Integrações Externas

#### 5.2.1 LLM (Large Language Model)

| Provider | Endpoint | Modelo |
|----------|----------|--------|
| Anthropic | `api.anthropic.com` | claude-3-sonnet |
| OpenAI (backup) | `api.openai.com` | gpt-4-turbo |

#### 5.2.2 Notificações

| Serviço | Uso | Protocolo |
|---------|-----|-----------|
| Firebase Cloud Messaging | Push Android/iOS | HTTPS |
| Apple Push Notification | Push iOS | HTTPS/2 |
| Twilio | SMS de emergência | HTTPS |
| SMTP (SendGrid) | E-mail | SMTPS |

#### 5.2.3 Serviços de Emergência

| Serviço | Integração |
|---------|------------|
| SAMU (192) | Discagem direta via app |
| CVV (188) | Discagem direta via app |
| Bombeiros (193) | Discagem direta via app |

### 5.3 Formatos de Dados

| Formato | Uso |
|---------|-----|
| JSON | API REST, configurações |
| Protocol Buffers | Comunicação interna (gRPC) |
| MessagePack | Cache Redis |
| CSV | Exportação de dados |
| PDF | Relatórios |

---

## 6. Versionamento

### 6.1 Convenção de Versão

**Semantic Versioning 2.0 (SemVer)**

```
MAJOR.MINOR.PATCH+BUILD

Exemplo: 2.1.3+build.456

MAJOR: Mudanças incompatíveis de API
MINOR: Funcionalidades novas compatíveis
PATCH: Correções de bugs compatíveis
BUILD: Identificador único de build
```

### 6.2 Versão Atual

| Componente | Versão |
|------------|--------|
| EVA-Mind-FZPN (Sistema) | 2.0.0 |
| API Backend | 2.0.0 |
| Mobile App (Android) | 2.0.0 |
| Mobile App (iOS) | 2.0.0 |
| Web App | 2.0.0 |

---

## Aprovações

| Função | Nome | Assinatura | Data |
|--------|------|------------|------|
| Arquiteto de Software | | | |
| Engenheiro de Qualidade | | | |
| Responsável Regulatório | José R F Junior | | 2025-01-27 |

---

**Documento controlado - Versão 1.0**
**Próxima revisão: 2026-01-27**
