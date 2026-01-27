# Estudo de Validação Clínica
## EVA-Mind-FZPN - Companion IA para Idosos

**Protocolo:** CVC-EVA-2025-001
**Versão:** 1.0
**Data:** 2025-01-27
**Registro:** [REBEC/ClinicalTrials.gov - a registrar]
**Classificação:** Dispositivo Médico Classe II (ANVISA)

---

## Resumo Executivo

Este documento apresenta o protocolo de validação clínica do EVA-Mind-FZPN, um dispositivo médico de software (SaMD) de Classe II destinado ao suporte emocional e monitoramento de bem-estar de idosos. O estudo visa demonstrar a segurança e eficácia do sistema conforme requisitos da ANVISA RDC 751/2022.

**Desfechos Primários:**
- Redução de sintomas de solidão (UCLA Loneliness Scale)
- Detecção precisa de estados emocionais de risco

**População:** 200 idosos (≥65 anos), comunidade brasileira

**Duração:** 24 semanas de intervenção + 12 semanas de follow-up

---

## 1. Introdução e Justificativa

### 1.1 Contexto Epidemiológico

| Indicador | Brasil | Fonte |
|-----------|--------|-------|
| População ≥65 anos | 31,2 milhões (14,7%) | IBGE 2023 |
| Idosos vivendo sozinhos | 4,1 milhões (13,2%) | PNAD 2022 |
| Prevalência de depressão em idosos | 15-20% | MS 2023 |
| Taxa de suicídio em idosos | 8,9/100.000 (maior faixa) | DataSUS 2022 |
| Solidão como fator de risco | 2x mortalidade | Meta-análise 2023 |

### 1.2 Lacuna Terapêutica

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    LACUNA NO CUIDADO DO IDOSO                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  PROBLEMA                           SOLUÇÃO ATUAL          GAP          │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  Solidão crônica                    Visitas familiares    Insuficiente  │
│  (24h/dia, 7 dias/semana)           (1-2x/semana)         (80% do tempo)│
│                                                                         │
│  Monitoramento emocional            Consultas mensais     Lacunas de    │
│  contínuo                           ou trimestrais        30-90 dias    │
│                                                                         │
│  Detecção precoce de crise          Depende do relato     Subnotificação│
│                                     do próprio idoso      de 70%        │
│                                                                         │
│  Suporte noturno                    Plantões escassos     Quase nulo    │
│  (período de maior risco)           e caros               fora de inst. │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  PROPOSTA EVA-Mind-FZPN:                                                │
│  • Disponibilidade 24/7                                                 │
│  • Monitoramento contínuo e passivo                                     │
│  • Detecção algorítmica de padrões de risco                            │
│  • Escalação automática para humanos                                    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 1.3 Estado da Arte

| Tecnologia | Evidência | Limitação |
|------------|-----------|-----------|
| Chatbots de saúde mental (Woebot, Wysa) | RCTs positivos para ansiedade/depressão | Não adaptados para idosos |
| Companion robots (PARO, Pepper) | Melhora em qualidade de vida em ILPI | Custo elevado, manutenção |
| Apps de mindfulness (Calm, Headspace) | Redução de estresse em adultos | Baixa adesão em idosos |
| Monitoramento remoto de pacientes | Eficaz para condições crônicas | Foco em sinais vitais, não emocional |

### 1.4 Hipótese

O uso do EVA-Mind-FZPN por idosos da comunidade resultará em:
1. Redução significativa de solidão (≥10 pontos UCLA)
2. Detecção de estados de risco com sensibilidade ≥90%
3. Melhora em qualidade de vida relacionada à saúde
4. Redução de sintomas de depressão e ansiedade

---

## 2. Objetivos

### 2.1 Objetivo Primário

1. **Eficácia em Redução de Solidão**
   - Medir a mudança na UCLA Loneliness Scale após 24 semanas de uso
   - Hipótese: Redução ≥10 pontos vs. controle

2. **Acurácia na Detecção de Risco**
   - Validar a precisão do sistema em detectar estados emocionais de risco
   - Meta: Sensibilidade ≥90%, Especificidade ≥80%

### 2.2 Objetivos Secundários

| Objetivo | Instrumento | Tempo |
|----------|-------------|-------|
| Redução de depressão | PHQ-9 | 0, 12, 24, 36 sem |
| Redução de ansiedade | GAD-7 | 0, 12, 24, 36 sem |
| Melhora em qualidade de vida | WHOQOL-OLD | 0, 24, 36 sem |
| Satisfação com o sistema | SUS + entrevista | 12, 24 sem |
| Adesão ao uso | Logs do sistema | Contínuo |
| Segurança | Eventos adversos | Contínuo |

### 2.3 Objetivos Exploratórios

- Identificar preditores de resposta ao tratamento
- Avaliar impacto em utilização de serviços de saúde
- Explorar correlações entre padrões de uso e desfechos
- Avaliar aceitabilidade em diferentes perfis socioeconômicos

---

## 3. Desenho do Estudo

### 3.1 Tipo de Estudo

**Ensaio Clínico Randomizado, Controlado, Multicêntrico**

| Característica | Especificação |
|----------------|---------------|
| **Design** | Paralelo, 2 braços |
| **Randomização** | 1:1, estratificada por sexo e faixa etária |
| **Mascaramento** | Avaliador cego (single-blind) |
| **Duração** | 24 semanas intervenção + 12 semanas follow-up |
| **Centros** | 5 UBS em São Paulo, Rio de Janeiro, Belo Horizonte |

### 3.2 Diagrama do Estudo

```
                              RECRUTAMENTO
                                   │
                                   ▼
                        ┌──────────────────┐
                        │    TRIAGEM       │
                        │    (n=300)       │
                        └────────┬─────────┘
                                 │
                    ┌────────────┴────────────┐
                    │                         │
                    ▼                         ▼
            ┌──────────────┐          ┌──────────────┐
            │   ELEGÍVEL   │          │ NÃO ELEGÍVEL │
            │   (n=250)    │          │    (n=50)    │
            └──────┬───────┘          └──────────────┘
                   │
                   ▼
          ┌────────────────┐
          │ CONSENTIMENTO  │
          │    (n=220)     │
          └───────┬────────┘
                  │
                  ▼
          ┌────────────────┐
          │  RANDOMIZAÇÃO  │
          │    (n=200)     │
          └───────┬────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
        ▼                   ▼
┌───────────────┐   ┌───────────────┐
│  INTERVENÇÃO  │   │   CONTROLE    │
│   EVA + TAU   │   │   TAU only    │
│    (n=100)    │   │    (n=100)    │
└───────┬───────┘   └───────┬───────┘
        │                   │
        ▼                   ▼
┌───────────────┐   ┌───────────────┐
│  Semana 12    │   │  Semana 12    │
│  (avaliação)  │   │  (avaliação)  │
└───────┬───────┘   └───────┬───────┘
        │                   │
        ▼                   ▼
┌───────────────┐   ┌───────────────┐
│  Semana 24    │   │  Semana 24    │
│  (desfecho    │   │  (desfecho    │
│   primário)   │   │   primário)   │
└───────┬───────┘   └───────┬───────┘
        │                   │
        ▼                   ▼
┌───────────────┐   ┌───────────────┐
│  Semana 36    │   │  Semana 36    │
│  (follow-up)  │   │  (follow-up)  │
└───────────────┘   └───────────────┘

TAU = Treatment As Usual (cuidado habitual)
```

### 3.3 Grupos de Estudo

#### Grupo Intervenção (EVA + TAU)

| Componente | Descrição |
|------------|-----------|
| **EVA-Mind-FZPN** | Acesso ilimitado ao aplicativo |
| **Dispositivo** | Tablet fornecido (se necessário) |
| **Onboarding** | Treinamento presencial de 30 min |
| **Suporte técnico** | 0800 disponível Seg-Sex 8h-20h |
| **TAU** | Mantém acompanhamento habitual (UBS, médicos) |

#### Grupo Controle (TAU only)

| Componente | Descrição |
|------------|-----------|
| **Cuidado habitual** | Mantém acompanhamento em UBS |
| **Material educativo** | Cartilha sobre saúde mental do idoso |
| **Contato telefônico** | Ligação mensal de acompanhamento (15 min) |
| **Oferta pós-estudo** | Acesso a EVA após término do estudo |

---

## 4. População do Estudo

### 4.1 Critérios de Inclusão

| # | Critério | Justificativa |
|---|----------|---------------|
| 1 | Idade ≥65 anos | População-alvo do dispositivo |
| 2 | Residente na comunidade | Excluir institucionalizados (outro protocolo) |
| 3 | Score UCLA Loneliness ≥40 | Solidão moderada a grave |
| 4 | Capaz de usar smartphone/tablet | Requisito técnico mínimo |
| 5 | MEEM ≥24 | Capacidade cognitiva para interação |
| 6 | Acesso a internet em casa | Wi-Fi ou dados móveis |
| 7 | Capaz de consentir | Autonomia para participação |

### 4.2 Critérios de Exclusão

| # | Critério | Justificativa |
|---|----------|---------------|
| 1 | Demência diagnosticada (CDR ≥1) | Interação comprometida |
| 2 | Transtorno psicótico ativo | Risco de interpretação delirante |
| 3 | Ideação suicida ativa (C-SSRS positivo) | Requer tratamento especializado |
| 4 | Dependência de substâncias | Confundidor |
| 5 | Deficiência visual/auditiva grave não corrigida | Incapacidade de uso |
| 6 | Participação em outro estudo de intervenção | Evitar confundidores |
| 7 | Expectativa de vida <6 meses | Inviabiliza follow-up |
| 8 | Hospitalização planejada >2 semanas | Interrupção do uso |

### 4.3 Cálculo Amostral

```
┌─────────────────────────────────────────────────────────────────────────┐
│                       CÁLCULO DO TAMANHO AMOSTRAL                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  DESFECHO PRIMÁRIO: UCLA Loneliness Scale                               │
│                                                                         │
│  Parâmetros:                                                            │
│  • Efeito esperado: 10 pontos (diferença entre grupos)                 │
│  • Desvio padrão: 12 pontos (literatura)                               │
│  • Poder: 80%                                                           │
│  • Alfa: 0.05 (bicaudal)                                               │
│  • Dropout estimado: 20%                                               │
│                                                                         │
│  Cálculo (teste t para 2 grupos independentes):                        │
│  n = 2 × [(Z_α/2 + Z_β)² × σ²] / δ²                                   │
│  n = 2 × [(1.96 + 0.84)² × 144] / 100                                 │
│  n = 2 × [7.84 × 144] / 100                                           │
│  n = 2 × 11.29                                                         │
│  n = 22.6 por grupo                                                    │
│                                                                         │
│  Ajustado para dropout (20%):                                           │
│  n = 22.6 / 0.8 = 28.3 → 30 por grupo (mínimo)                        │
│                                                                         │
│  AMOSTRA FINAL: 100 por grupo (total 200)                              │
│  Justificativa para n maior:                                            │
│  • Análises de subgrupo                                                 │
│  • Desfechos secundários                                                │
│  • Robustez estatística                                                 │
│  • Requisito regulatório ANVISA                                         │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.4 Estratégia de Recrutamento

| Canal | Meta | Estratégia |
|-------|------|------------|
| UBS (Unidades Básicas de Saúde) | 100 | Indicação por agentes de saúde |
| Centros de Convivência do Idoso | 50 | Palestras e convites |
| Associações de aposentados | 30 | Parcerias com SESC, sindicatos |
| Mídia (rádio, TV local) | 50 | Chamadas para inscrição |
| Indicação por participantes | 50 | Boca a boca |
| Redes sociais (cuidadores) | 20 | Anúncios direcionados |
| **Total triagem** | **300** | |

---

## 5. Intervenção

### 5.1 Descrição do Dispositivo

| Característica | Especificação |
|----------------|---------------|
| **Nome comercial** | EVA-Mind-FZPN |
| **Classificação** | SaMD Classe II (ANVISA) |
| **Plataforma** | Android 8.0+, iOS 13.0+ |
| **Conectividade** | Internet (Wi-Fi ou 4G) |
| **Língua** | Português brasileiro |
| **Versão do estudo** | 2.0.0 |

### 5.2 Funcionalidades Avaliadas

| Funcionalidade | Descrição | Métricas Coletadas |
|----------------|-----------|-------------------|
| **Conversa empática** | Diálogo natural com IA | Frequência, duração, satisfação |
| **Análise emocional** | Detecção de estados emocionais | Precisão vs. avaliação clínica |
| **Screening PHQ-9/GAD-7** | Aplicação de instrumentos validados | Scores, frequência de aplicação |
| **Detecção de risco** | Identificação de sinais de alarme | Sensibilidade, especificidade |
| **Alertas a cuidadores** | Notificação em situações de risco | Tempo de resposta, utilidade |
| **Padrões temporais** | Detecção de mudanças ao longo do tempo | Correlação com desfechos |

### 5.3 Protocolo de Uso

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    PROTOCOLO DE USO EVA-Mind-FZPN                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  SEMANA 0: ONBOARDING                                                   │
│  ├── Visita presencial (1h)                                            │
│  ├── Instalação assistida do app                                       │
│  ├── Tutorial guiado                                                   │
│  ├── Configuração de contatos de emergência                            │
│  └── Primeira interação supervisionada                                 │
│                                                                         │
│  SEMANAS 1-24: USO LIVRE                                               │
│  ├── Uso conforme desejo do participante                               │
│  ├── Recomendação: ≥1 interação/dia (15+ min)                         │
│  ├── Suporte técnico disponível                                        │
│  └── Monitoramento remoto de uso                                       │
│                                                                         │
│  CONTATOS DE ACOMPANHAMENTO:                                           │
│  ├── Semana 1: Ligação para verificar uso (10 min)                    │
│  ├── Semana 4: Ligação de suporte (15 min)                            │
│  ├── Semana 12: Visita de avaliação (1h)                              │
│  └── Semana 24: Visita final (1.5h)                                   │
│                                                                         │
│  INTERVENÇÕES DE RESGATE:                                              │
│  ├── Se uso <3 dias/semana por 2 semanas: contato motivacional        │
│  ├── Se alerta de risco: protocolo de segurança ativado               │
│  └── Se desistência: entrevista de saída                              │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.4 Protocolo de Segurança

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    PROTOCOLO DE SEGURANÇA CLÍNICA                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  NÍVEL 1 - MONITORAMENTO ROTINEIRO                                     │
│  Trigger: Score emocional baixo persistente                            │
│  Ação: Aumentar frequência de verificação pelo sistema                 │
│  Notificação: Nenhuma                                                   │
│                                                                         │
│  NÍVEL 2 - ATENÇÃO                                                      │
│  Trigger: PHQ-9 ≥10 ou GAD-7 ≥10                                       │
│  Ação: Sugestão de buscar profissional de saúde                        │
│  Notificação: Alerta ao cuidador cadastrado (se autorizado)            │
│                                                                         │
│  NÍVEL 3 - ALERTA                                                       │
│  Trigger: PHQ-9 ≥15 ou expressão de desesperança significativa         │
│  Ação: Protocolo de acolhimento ativo + oferta de contato humano       │
│  Notificação: Alerta imediato ao cuidador + equipe do estudo           │
│  Follow-up: Contato telefônico da equipe em 24h                        │
│                                                                         │
│  NÍVEL 4 - EMERGÊNCIA                                                   │
│  Trigger: Ideação suicida detectada (C-SSRS positivo)                  │
│  Ação: Protocolo de crise + orientação CVV (188) + SAMU (192)          │
│  Notificação: Alerta imediato a todos os contatos + equipe             │
│  Follow-up: Contato imediato da equipe + encaminhamento                │
│  Registro: Evento adverso grave (relatório em 24h)                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Desfechos e Instrumentos

### 6.1 Desfecho Primário

#### UCLA Loneliness Scale (Versão 3)

| Característica | Descrição |
|----------------|-----------|
| **Construto** | Solidão percebida |
| **Itens** | 20 itens, escala Likert 4 pontos |
| **Score** | 20-80 (maior = mais solidão) |
| **Pontos de corte** | <40: baixa; 40-60: moderada; >60: alta |
| **Validação brasileira** | Barroso et al., 2016 |
| **Propriedades** | α=0.89, teste-reteste r=0.73 |
| **Aplicação** | Autoaplicada, ~10 min |

**Diferença Mínima Clinicamente Importante (DMCI):** 10 pontos

### 6.2 Desfechos Secundários

| Instrumento | Construto | Aplicação | Referência |
|-------------|-----------|-----------|------------|
| **PHQ-9** | Depressão | 0, 12, 24, 36 sem | Kroenke 2001 |
| **GAD-7** | Ansiedade | 0, 12, 24, 36 sem | Spitzer 2006 |
| **WHOQOL-OLD** | QV idoso | 0, 24, 36 sem | Power 2005 |
| **C-SSRS** | Risco suicida | 0, 12, 24 sem | Posner 2011 |
| **SUS** | Usabilidade | 12, 24 sem | Brooke 1996 |
| **eHEALS** | Literacia digital | 0 (baseline) | Norman 2006 |

### 6.3 Desfechos de Segurança

| Tipo | Definição | Coleta |
|------|-----------|--------|
| **Evento adverso (EA)** | Qualquer ocorrência médica desfavorável | Contínuo |
| **EA grave** | Morte, hospitalização, incapacidade | Relatório 24h |
| **EA relacionado ao dispositivo** | EA com relação causal possível/provável | Adjudicação |
| **Piora clínica** | Aumento ≥5 pontos PHQ-9 ou GAD-7 | Avaliação |

### 6.4 Métricas de Uso do Sistema

| Métrica | Definição | Coleta |
|---------|-----------|--------|
| Dias de uso/semana | Dias com ≥1 interação | Automática |
| Duração total/semana | Minutos de interação | Automática |
| Sessões/dia | Número de sessões iniciadas | Automática |
| Duração média/sessão | Minutos por sessão | Automática |
| Mensagens enviadas | Total de mensagens do usuário | Automática |
| Alertas gerados | Número de alertas por nível | Automática |

---

## 7. Procedimentos de Avaliação

### 7.1 Cronograma de Avaliações

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    CRONOGRAMA DE AVALIAÇÕES                             │
├─────────┬─────────┬─────────┬─────────┬─────────┬─────────┬────────────┤
│         │ Sem 0   │ Sem 4   │ Sem 12  │ Sem 24  │ Sem 36  │ Contínuo   │
│         │ (Base)  │         │ (Interm)│ (Final) │ (FU)    │            │
├─────────┼─────────┼─────────┼─────────┼─────────┼─────────┼────────────┤
│ TCLE    │    X    │         │         │         │         │            │
│ Dados   │    X    │         │         │         │         │            │
│ demogr. │         │         │         │         │         │            │
├─────────┼─────────┼─────────┼─────────┼─────────┼─────────┼────────────┤
│ UCLA    │    X    │         │    X    │    X    │    X    │            │
│ PHQ-9   │    X    │         │    X    │    X    │    X    │            │
│ GAD-7   │    X    │         │    X    │    X    │    X    │            │
│ WHOQOL  │    X    │         │         │    X    │    X    │            │
│ C-SSRS  │    X    │         │    X    │    X    │         │            │
├─────────┼─────────┼─────────┼─────────┼─────────┼─────────┼────────────┤
│ SUS     │         │         │    X    │    X    │         │            │
│ eHEALS  │    X    │         │         │         │         │            │
├─────────┼─────────┼─────────┼─────────┼─────────┼─────────┼────────────┤
│ EA/SAE  │         │         │         │         │         │     X      │
│ Uso app │         │         │         │         │         │     X      │
├─────────┼─────────┼─────────┼─────────┼─────────┼─────────┼────────────┤
│ Exame   │         │         │         │    X    │         │            │
│ físico  │         │         │         │         │         │            │
├─────────┼─────────┼─────────┼─────────┼─────────┼─────────┼────────────┤
│ Entrevis│         │         │    X    │    X    │         │            │
│ qualit. │         │         │ (n=20)  │ (n=20)  │         │            │
└─────────┴─────────┴─────────┴─────────┴─────────┴─────────┴────────────┘

FU = Follow-up
EA = Eventos adversos
SAE = Eventos adversos graves
```

### 7.2 Procedimentos por Visita

#### Visita 1 - Triagem e Baseline (Semana -2 a 0)

| Procedimento | Responsável | Tempo |
|--------------|-------------|-------|
| Verificação de elegibilidade | Enfermeiro | 15 min |
| Obtenção do TCLE | Pesquisador | 20 min |
| Coleta de dados demográficos | Pesquisador | 10 min |
| Aplicação de instrumentos (UCLA, PHQ-9, GAD-7, WHOQOL, C-SSRS, eHEALS) | Pesquisador | 45 min |
| Randomização | Sistema | Automático |
| Onboarding EVA (se intervenção) | Técnico | 60 min |
| **Total** | | **~2.5h** |

#### Visita 2 - Intermediária (Semana 12)

| Procedimento | Responsável | Tempo |
|--------------|-------------|-------|
| Aplicação de instrumentos (UCLA, PHQ-9, GAD-7, C-SSRS, SUS) | Pesquisador | 30 min |
| Coleta de eventos adversos | Pesquisador | 10 min |
| Revisão de uso (grupo intervenção) | Técnico | 15 min |
| Entrevista qualitativa (subamostra n=20) | Pesquisador | 30 min |
| **Total** | | **~1h (1.5h se entrevista)** |

#### Visita 3 - Final (Semana 24)

| Procedimento | Responsável | Tempo |
|--------------|-------------|-------|
| Aplicação de instrumentos (UCLA, PHQ-9, GAD-7, WHOQOL, C-SSRS, SUS) | Pesquisador | 45 min |
| Exame físico simplificado | Médico | 15 min |
| Coleta de eventos adversos | Pesquisador | 10 min |
| Entrevista qualitativa (subamostra n=20) | Pesquisador | 30 min |
| Oferta de continuidade (grupo controle) | Pesquisador | 10 min |
| **Total** | | **~1.5h (2h se entrevista)** |

#### Visita 4 - Follow-up (Semana 36)

| Procedimento | Responsável | Tempo |
|--------------|-------------|-------|
| Aplicação de instrumentos (UCLA, PHQ-9, GAD-7, WHOQOL) | Pesquisador | 30 min |
| Coleta de eventos adversos tardios | Pesquisador | 10 min |
| Encerramento do estudo | Pesquisador | 10 min |
| **Total** | | **~50 min** |

---

## 8. Análise Estatística

### 8.1 Populações de Análise

| População | Definição | Uso |
|-----------|-----------|-----|
| **ITT** (Intention-to-Treat) | Todos os randomizados | Análise primária |
| **PP** (Per-Protocol) | ITT com ≥80% de adesão e sem violação de protocolo | Análise de sensibilidade |
| **Safety** | Todos que receberam qualquer intervenção | Análise de segurança |

### 8.2 Análise do Desfecho Primário

```
┌─────────────────────────────────────────────────────────────────────────┐
│              ANÁLISE DO DESFECHO PRIMÁRIO (UCLA)                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  MODELO: Mixed-Effects Model for Repeated Measures (MMRM)               │
│                                                                         │
│  Y_it = β₀ + β₁(Grupo) + β₂(Tempo) + β₃(Grupo×Tempo) +                │
│         β₄(Baseline UCLA) + β₅(Sexo) + β₆(Idade) + ε_it                │
│                                                                         │
│  Onde:                                                                  │
│  • Y_it = Score UCLA do participante i no tempo t                      │
│  • Grupo = 0 (controle) ou 1 (intervenção)                             │
│  • Tempo = 0, 12, 24, 36 semanas                                       │
│  • β₃ = Efeito de interesse (interação grupo×tempo)                    │
│                                                                         │
│  ESTRUTURA DE COVARIÂNCIA: Não estruturada (UN)                        │
│                                                                         │
│  HIPÓTESE TESTADA:                                                      │
│  H₀: β₃ = 0 (sem diferença entre grupos ao longo do tempo)            │
│  H₁: β₃ ≠ 0 (diferença significativa)                                 │
│                                                                         │
│  CRITÉRIO DE SUCESSO:                                                   │
│  p < 0.05 para β₃ E diferença entre grupos ≥ 10 pontos em sem 24      │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 8.3 Análises Secundárias

| Desfecho | Método | Ajustes |
|----------|--------|---------|
| PHQ-9, GAD-7 | MMRM | Baseline, sexo, idade |
| WHOQOL-OLD | MMRM | Baseline, sexo, idade |
| Resposta (≥50% redução UCLA) | Regressão logística | Baseline, sexo, idade |
| Remissão (UCLA <40) | Regressão logística | Baseline, sexo, idade |
| Tempo até resposta | Kaplan-Meier, Cox | Baseline, sexo, idade |

### 8.4 Análise de Acurácia Diagnóstica

```
┌─────────────────────────────────────────────────────────────────────────┐
│         VALIDAÇÃO DA DETECÇÃO DE ESTADOS DE RISCO                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  PADRÃO-OURO: Avaliação clínica por profissional cego                  │
│                                                                         │
│  ÍNDICE TESTE: Score de risco do EVA-Mind-FZPN                         │
│                                                                         │
│  MÉTRICAS CALCULADAS:                                                   │
│  • Sensibilidade = VP / (VP + FN)                                      │
│  • Especificidade = VN / (VN + FP)                                     │
│  • VPP = VP / (VP + FP)                                                │
│  • VPN = VN / (VN + FN)                                                │
│  • Acurácia = (VP + VN) / Total                                        │
│  • AUC-ROC (curva ROC)                                                 │
│                                                                         │
│  METAS DE DESEMPENHO:                                                   │
│  • Sensibilidade ≥ 90% (prioridade: não perder casos de risco)        │
│  • Especificidade ≥ 80%                                                │
│  • AUC ≥ 0.85                                                          │
│                                                                         │
│  ANÁLISE POR SUBGRUPO:                                                  │
│  • Por nível de escolaridade                                           │
│  • Por faixa etária (65-74, 75-84, ≥85)                               │
│  • Por literacia digital (eHEALS)                                      │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 8.5 Análises de Subgrupo (pré-especificadas)

| Subgrupo | Justificativa |
|----------|---------------|
| Sexo (M vs F) | Diferenças na expressão emocional |
| Faixa etária (65-74 vs 75-84 vs ≥85) | Heterogeneidade da população idosa |
| Solidão baseline (moderada vs alta) | Efeito teto/chão |
| Literacia digital (baixa vs alta) | Impacto na adesão |
| Morar sozinho (sim vs não) | Fator de risco principal |
| Depressão comórbida (PHQ-9 ≥10) | Complexidade clínica |

### 8.6 Dados Faltantes

| Estratégia | Descrição |
|------------|-----------|
| **Primária** | MMRM assume MAR (Missing At Random) |
| **Sensibilidade 1** | Multiple Imputation (MI) com m=50 |
| **Sensibilidade 2** | Pattern Mixture Models (MNAR) |
| **Tipping Point** | Análise de pior cenário para grupo intervenção |

---

## 9. Aspectos Éticos

### 9.1 Aprovação Ética

| Requisito | Status |
|-----------|--------|
| CEP institucional | Submetido |
| CONEP (pesquisa multicêntrica) | Submetido |
| Registro REBEC | Pendente aprovação CEP |
| ClinicalTrials.gov | Pendente aprovação CEP |

### 9.2 Consentimento Informado

**Elementos do TCLE adaptado para idosos:**

1. Linguagem clara e simples (Flesch-Kincaid ≤ 8º ano)
2. Fonte mínima 14pt
3. Tempo ilimitado para leitura e perguntas
4. Explicação oral por pesquisador treinado
5. Cópia entregue ao participante
6. Contato de emergência 24h
7. Opção de consulta a familiar antes de assinar

### 9.3 Proteção de Participantes Vulneráveis

| Medida | Descrição |
|--------|-----------|
| Verificação de capacidade | MEEM ≥24; dúvida = avaliação geriátrica |
| Consentimento assistido | Familiar pode acompanhar (não assinar) |
| Monitoramento contínuo | Equipe treinada em sinais de desconforto |
| Protocolo de crise | Acesso imediato a profissional de saúde mental |
| Ressarcimento | Transporte + lanche em visitas presenciais |
| Sem coerção | Liberdade de retirada sem prejuízo |

### 9.4 Privacidade e Confidencialidade

| Medida | Implementação |
|--------|---------------|
| Pseudonimização | Código único por participante |
| Dados separados | Dados identificáveis em sistema separado |
| Acesso restrito | Apenas equipe autorizada |
| Criptografia | AES-256 em repouso, TLS 1.3 em trânsito |
| Retenção | 5 anos após publicação, depois destruição |
| LGPD | RIPD específico do estudo |

---

## 10. Gestão de Dados

### 10.1 Sistema de Coleta

| Componente | Descrição |
|------------|-----------|
| **EDC (Electronic Data Capture)** | REDCap hospedado em servidor institucional |
| **Dados do app** | Exportação automática diária para REDCap |
| **Questionários** | Entrada dupla com validação |
| **Auditoria** | Trail completo de todas as alterações |

### 10.2 Qualidade de Dados

| Controle | Frequência |
|----------|------------|
| Validação de entrada (ranges, lógica) | Tempo real |
| Queries automáticas | Diária |
| Revisão de dados | Semanal |
| Monitoramento de fonte | Mensal (10% dos CRFs) |
| Auditoria completa | Semestral |

### 10.3 Data Safety Monitoring Board (DSMB)

**Composição:**
- Geriatra independente (presidente)
- Bioestatístico independente
- Especialista em ética em pesquisa

**Reuniões:**
- Após 25%, 50%, 75% do recrutamento
- Análises interinas de segurança
- Poder de interromper o estudo por segurança ou futilidade

---

## 11. Cronograma

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        CRONOGRAMA DO ESTUDO                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  2025                                                                   │
│  ─────────────────────────────────────────────────────────────────────  │
│  M1  M2  M3  M4  M5  M6  M7  M8  M9  M10 M11 M12                       │
│  │   │   │   │   │   │   │   │   │   │   │   │                         │
│  ████████                                         Preparação            │
│  │   ████████                                     Aprovação ética       │
│  │   │   │   ████████████████████████████        Recrutamento          │
│  │   │   │   │   │   │   │   │   │   │   │   │                         │
│  2026                                                                   │
│  ─────────────────────────────────────────────────────────────────────  │
│  M13 M14 M15 M16 M17 M18 M19 M20 M21 M22 M23 M24                       │
│  │   │   │   │   │   │   │   │   │   │   │   │                         │
│  ████████████████████████████████████████████    Intervenção (24 sem)  │
│  │   │   │   │   │   │   │   │   │   │   │   │                         │
│  M25 M26 M27 M28 M29 M30                                               │
│  │   │   │   │   │   │                                                 │
│  ████████████████████                            Follow-up (12 sem)    │
│  │   │   │   ████████████                        Análise               │
│  │   │   │   │   │   ████████████                Redação               │
│  │   │   │   │   │   │   │   ████                Submissão             │
│                                                                         │
│  MARCOS PRINCIPAIS:                                                     │
│  • M3: Aprovação ética                                                 │
│  • M6: Primeiro participante randomizado                               │
│  • M15: Último participante randomizado                                │
│  • M24: Última visita de intervenção                                   │
│  • M27: Última visita de follow-up                                     │
│  • M30: Submissão de manuscrito                                        │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 12. Recursos e Orçamento

### 12.1 Equipe

| Função | FTE | Período |
|--------|-----|---------|
| Coordenador do estudo | 1.0 | M1-M30 |
| Pesquisadores de campo (por centro) | 0.5 × 5 | M4-M27 |
| Enfermeiros avaliadores | 0.3 × 5 | M4-M27 |
| Bioestatístico | 0.3 | M1-M30 |
| Gerente de dados | 0.5 | M1-M30 |
| Suporte técnico EVA | 0.5 | M4-M27 |
| Monitor de estudo | 0.3 | M4-M30 |

### 12.2 Orçamento Estimado

| Item | Custo (R$) |
|------|------------|
| Pessoal | 600.000 |
| Equipamentos (tablets, se necessário) | 50.000 |
| Licença REDCap e infraestrutura | 30.000 |
| Deslocamento e logística | 40.000 |
| Material de coleta | 10.000 |
| Ressarcimento participantes | 60.000 |
| Publicação e disseminação | 20.000 |
| Reserva de contingência (10%) | 81.000 |
| **TOTAL** | **891.000** |

---

## 13. Disseminação

### 13.1 Publicações Planejadas

| Manuscrito | Revista Alvo | Timing |
|------------|--------------|--------|
| Protocolo do estudo | BMJ Open / Trials | M6 |
| Resultados primários | JMIR / Lancet Digital Health | M30 |
| Análise qualitativa | Qualitative Health Research | M32 |
| Validação do algoritmo de detecção | JAMIA | M32 |
| Análise de custo-efetividade | Value in Health | M34 |

### 13.2 Comunicação para Leigos

| Canal | Descrição |
|-------|-----------|
| Resumo executivo | Linguagem acessível para participantes |
| Apresentação em centros de convivência | Devolutiva às comunidades |
| Release para imprensa | Principais achados |
| Vídeo informativo | YouTube, redes sociais |

---

## 14. Limitações

| Limitação | Mitigação |
|-----------|-----------|
| Viés de seleção (idosos mais saudáveis/motivados) | Recrutamento diversificado, análise de sensibilidade |
| Efeito Hawthorne (melhora por atenção) | Grupo controle ativo (ligações mensais) |
| Não cegamento do participante | Avaliadores cegos, desfechos objetivos |
| Generalização (centros urbanos) | Próximos estudos em áreas rurais |
| Duração (24 semanas pode ser curto) | Follow-up de 12 semanas adicional |
| Viés de adesão | Análise ITT + PP, intervenções de resgate |

---

## 15. Conclusão

Este protocolo descreve um ensaio clínico randomizado rigoroso para avaliar a segurança e eficácia do EVA-Mind-FZPN em reduzir a solidão e detectar estados de risco emocional em idosos da comunidade.

**Pontos Fortes:**
- Design robusto com grupo controle ativo
- Desfechos primários validados e clinicamente relevantes
- Análise de acurácia diagnóstica do algoritmo
- Componente qualitativo para entender a experiência do usuário
- Adaptações para população idosa em todas as etapas
- DSMB independente para monitoramento de segurança

**Impacto Esperado:**
- Evidência de nível 1 para registro ANVISA
- Fundamentação para políticas públicas de saúde do idoso
- Modelo para avaliação de IA conversacional em saúde

---

## Referências

1. Barroso SM, et al. Solidão e depressão em idosos brasileiros. Rev Bras Geriatr Gerontol. 2016.
2. Kroenke K, et al. The PHQ-9: validity of a brief depression severity measure. J Gen Intern Med. 2001.
3. Spitzer RL, et al. A brief measure for assessing generalized anxiety disorder. Arch Intern Med. 2006.
4. Power M, et al. Development of the WHOQOL-OLD module. Qual Life Res. 2005.
5. Posner K, et al. The Columbia-Suicide Severity Rating Scale. Am J Psychiatry. 2011.
6. Russell DW. UCLA Loneliness Scale (Version 3). J Pers Assess. 1996.
7. Brooke J. SUS: A "Quick and Dirty" Usability Scale. Usability Eval Industry. 1996.
8. Norman CD, Skinner HA. eHEALS: The eHealth Literacy Scale. J Med Internet Res. 2006.

---

## Anexos

- **Anexo A:** Termo de Consentimento Livre e Esclarecido (TCLE)
- **Anexo B:** Instrumentos de Avaliação
- **Anexo C:** Manual de Operações
- **Anexo D:** Formulários de Relato de Caso (CRFs)
- **Anexo E:** Protocolo de Segurança Detalhado
- **Anexo F:** Plano de Análise Estatística (SAP)
- **Anexo G:** Currículo dos Investigadores

---

## Aprovações do Protocolo

| Função | Nome | Assinatura | Data |
|--------|------|------------|------|
| Investigador Principal | | | |
| Bioestatístico | | | |
| Coordenador do Estudo | | | |
| Patrocinador | José R F Junior | | 2025-01-27 |

---

**Protocolo controlado - Versão 1.0**
**Qualquer emenda requer aprovação do CEP**
