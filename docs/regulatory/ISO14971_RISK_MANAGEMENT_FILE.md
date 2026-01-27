# ARQUIVO DE GESTÃO DE RISCOS
## EVA-Mind-FZPN - Conforme ISO 14971:2019

**Documento:** RMF-001
**Versão:** 0.1 (DRAFT)
**Data:** 2026-01-27
**Status:** Em elaboração

---

## 1. ESCOPO E CONTEXTO

### 1.1 Identificação do Produto

| Campo | Valor |
|-------|-------|
| Nome do Produto | EVA-Mind |
| Versão | 1.0 |
| Classificação | SaMD Classe II (RDC 751/2022) |
| Fabricante | [Nome da empresa] |
| Responsável Técnico | José R F Junior |

### 1.2 Uso Pretendido

**Descrição:** Sistema de inteligência artificial para acompanhamento e suporte emocional de idosos, com capacidade de:
- Realizar avaliações clínicas padronizadas (PHQ-9, GAD-7, C-SSRS)
- Detectar sinais de risco psicológico e suicida
- Alertar cuidadores e profissionais de saúde em situações de emergência
- Fornecer suporte emocional através de conversas

**População alvo:**
- Idosos (65+ anos) em acompanhamento domiciliar
- Cuidadores familiares ou profissionais
- Profissionais de saúde mental

**Ambiente de uso:**
- Residencial (domicílio do idoso)
- Dispositivos: smartphones, tablets, computadores
- Conectividade: Internet (Wi-Fi, 4G/5G)

### 1.3 Indicações de Uso

O EVA-Mind é indicado para:
1. Triagem de sintomas de depressão (PHQ-9)
2. Triagem de sintomas de ansiedade (GAD-7)
3. Avaliação de risco suicida (C-SSRS)
4. Suporte emocional complementar ao tratamento profissional
5. Monitoramento de bem-estar entre consultas

### 1.4 Contraindicações

O EVA-Mind NÃO deve ser usado como:
1. Substituto de atendimento profissional de saúde mental
2. Ferramenta de diagnóstico definitivo
3. Única fonte de avaliação em situações de crise
4. Sistema para pacientes com demência severa

### 1.5 Limitações Conhecidas

1. Requer conectividade com internet
2. Dependente da capacidade cognitiva do usuário para interação
3. Não substitui avaliação clínica presencial
4. Escalas validadas mas não diagnósticas

---

## 2. PROCESSO DE GESTÃO DE RISCO

### 2.1 Responsabilidades

| Função | Responsabilidade |
|--------|------------------|
| Alta Direção | Aprovar política de gestão de riscos |
| Engenheiro de Qualidade | Conduzir análises de risco |
| Equipe de Desenvolvimento | Implementar controles de risco |
| Validação | Verificar eficácia dos controles |

### 2.2 Critérios de Aceitabilidade de Risco

#### Matriz de Probabilidade

| Nível | Descrição | Critério |
|-------|-----------|----------|
| 1 | Muito improvável | <0.01% (1 em 10.000) |
| 2 | Improvável | 0.01-0.1% (1 em 1.000-10.000) |
| 3 | Possível | 0.1-1% (1 em 100-1.000) |
| 4 | Provável | 1-10% (1 em 10-100) |
| 5 | Frequente | >10% (>1 em 10) |

#### Matriz de Severidade

| Nível | Descrição | Impacto |
|-------|-----------|---------|
| 1 | Insignificante | Inconveniência temporária |
| 2 | Menor | Lesão menor sem tratamento |
| 3 | Sério | Lesão requerendo tratamento |
| 4 | Crítico | Lesão permanente |
| 5 | Catastrófico | Morte ou risco de vida |

#### Matriz de Aceitabilidade

|          | Sev 1 | Sev 2 | Sev 3 | Sev 4 | Sev 5 |
|----------|-------|-------|-------|-------|-------|
| Prob 5   | 🟡    | 🟡    | 🔴    | 🔴    | 🔴    |
| Prob 4   | 🟢    | 🟡    | 🟡    | 🔴    | 🔴    |
| Prob 3   | 🟢    | 🟢    | 🟡    | 🟡    | 🔴    |
| Prob 2   | 🟢    | 🟢    | 🟢    | 🟡    | 🟡    |
| Prob 1   | 🟢    | 🟢    | 🟢    | 🟢    | 🟡    |

- 🟢 **Aceitável** - Risco aceitável sem ação
- 🟡 **ALARP** - Tão baixo quanto razoavelmente praticável
- 🔴 **Inaceitável** - Requer redução de risco

---

## 3. IDENTIFICAÇÃO DE PERIGOS

### 3.1 Fontes de Perigos por Categoria

#### 3.1.1 Perigos Relacionados à Energia (N/A)
*Software não possui componentes de energia direta*

#### 3.1.2 Perigos Biológicos (N/A)
*Software não possui componentes biológicos*

#### 3.1.3 Perigos Relacionados ao Uso

| ID | Perigo | Descrição |
|----|--------|-----------|
| H-USE-001 | Erro de interpretação do usuário | Usuário interpreta resultado como diagnóstico definitivo |
| H-USE-002 | Falha de comunicação | Alerta não entregue ao cuidador |
| H-USE-003 | Dependência excessiva | Usuário deixa de buscar ajuda profissional |
| H-USE-004 | Uso por população não indicada | Uso por crianças ou pessoas com demência severa |
| H-USE-005 | Interface não compreendida | Idoso não consegue usar o sistema |

#### 3.1.4 Perigos Relacionados à Informação

| ID | Perigo | Descrição |
|----|--------|-----------|
| H-INFO-001 | Score incorreto | Algoritmo calcula score errado |
| H-INFO-002 | Classificação incorreta | Risco subestimado ou superestimado |
| H-INFO-003 | Dados corrompidos | Perda de histórico clínico |
| H-INFO-004 | Vazamento de dados | Exposição de dados sensíveis |
| H-INFO-005 | Atraso na informação | Alerta de emergência atrasado |

#### 3.1.5 Perigos Funcionais

| ID | Perigo | Descrição |
|----|--------|-----------|
| H-FUN-001 | Sistema indisponível | Paciente não consegue acessar em crise |
| H-FUN-002 | Falha no processamento | Conversa não processada corretamente |
| H-FUN-003 | Integração falha | Canal de alerta indisponível |
| H-FUN-004 | Modelo de IA incorreto | Resposta inadequada do LLM |
| H-FUN-005 | Escalação falha | Todos os canais de alerta falham |

---

## 4. ANÁLISE DE RISCO

### 4.1 Análise Preliminar de Perigos (PHA)

| ID | Perigo | Causa | Sequência | Situação Perigosa | Dano | P | S | Risco |
|----|--------|-------|-----------|-------------------|------|---|---|-------|
| R-001 | H-INFO-001 | Bug no código | Score calculado incorretamente | Risco suicida não detectado | Suicídio | 1 | 5 | 🟡 |
| R-002 | H-INFO-002 | Limiar incorreto | Paciente classificado como baixo risco | Falta de intervenção | Dano psicológico | 2 | 4 | 🟡 |
| R-003 | H-USE-002 | Falha de rede | Alerta não entregue | Cuidador não informado | Atraso no socorro | 2 | 5 | 🟡 |
| R-004 | H-FUN-005 | Todas APIs falharam | Nenhum canal funciona | Isolamento em crise | Suicídio | 1 | 5 | 🟡 |
| R-005 | H-INFO-004 | Vulnerabilidade | Dados expostos | Violação de privacidade | Dano moral/legal | 2 | 3 | 🟢 |
| R-006 | H-USE-001 | Falta de orientação | Usuário acredita em diagnóstico | Tratamento inadequado | Piora do quadro | 3 | 3 | 🟡 |
| R-007 | H-FUN-001 | Servidor down | Sistema indisponível | Paciente sem suporte | Ansiedade aumentada | 3 | 2 | 🟢 |

---

## 5. AVALIAÇÃO E CONTROLE DE RISCO

### 5.1 Controles de Risco Implementados

#### R-001: Score calculado incorretamente

| Aspecto | Detalhe |
|---------|---------|
| **Risco Original** | P:1 × S:5 = 🟡 ALARP |
| **Controles Implementados** | |
| 1. Design inerentemente seguro | Algoritmos baseados em escalas validadas (PHQ-9, GAD-7, C-SSRS) |
| 2. Medidas de proteção | 25 testes unitários para escalas clínicas |
| 3. Informação de segurança | Disclaimer que não substitui avaliação profissional |
| **Verificação** | `internal/cortex/scales/clinical_scales_test.go` - 25 testes passando |
| **Risco Residual** | P:1 × S:5 = 🟡 ALARP (aceitável) |

```go
// Controle: Testes unitários cobrindo todos os cenários
func TestCSSRSRiskLevels(t *testing.T) {
    // Suicidal behavior (Q6) = SEMPRE crítico
    result := CalculateCSSRS([]int{0, 0, 0, 0, 0, 1}) // Apenas Q6
    assert.Equal(t, "critical", result.RiskLevel)
    assert.True(t, result.HasSuicidalBehavior)
}
```

---

#### R-002: Risco subestimado

| Aspecto | Detalhe |
|---------|---------|
| **Risco Original** | P:2 × S:4 = 🟡 ALARP |
| **Controles Implementados** | |
| 1. Design inerentemente seguro | Comportamento suicida (Q6 C-SSRS) = CRÍTICO automático |
| 2. Medidas de proteção | Limiar conservador (qualquer ideação ≥ alerta) |
| 3. Informação de segurança | Recursos de emergência sempre fornecidos (CVV 188, SAMU 192) |
| **Verificação** | Teste automatizado + código auditável |
| **Risco Residual** | P:1 × S:4 = 🟢 Aceitável |

```go
// Controle: Qualquer comportamento suicida = CRÍTICO
if result.HasSuicidalBehavior {
    result.RiskLevel = "critical"
    result.RequiresIntervention = true
    result.InterventionPlan = &InterventionPlan{
        Priority: "immediate",
        Actions: []string{"SAMU 192", "CVV 188", "Contato de emergência"},
    }
}
```

---

#### R-003: Alerta não entregue

| Aspecto | Detalhe |
|---------|---------|
| **Risco Original** | P:2 × S:5 = 🟡 ALARP |
| **Controles Implementados** | |
| 1. Design inerentemente seguro | Múltiplos canais redundantes (Push→WhatsApp→SMS→Email→Ligação) |
| 2. Medidas de proteção | Escalação automática com timeout por prioridade |
| 3. Informação de segurança | Log de todas as tentativas para auditoria |
| **Verificação** | `internal/cortex/alert/escalation_test.go` - 17 testes passando |
| **Risco Residual** | P:1 × S:5 = 🟡 ALARP (aceitável) |

```go
// Controle: Escalação multi-canal
type EscalationConfig struct {
    Channels []AlertChannel // Push, WhatsApp, SMS, Email, Call
    Timeouts map[AlertPriority]time.Duration{
        PriorityCritical: 30 * time.Second,  // Escalação rápida
        PriorityHigh:     2 * time.Minute,
        PriorityMedium:   5 * time.Minute,
        PriorityLow:      15 * time.Minute,
    }
}
```

---

#### R-004: Todos os canais falham

| Aspecto | Detalhe |
|---------|---------|
| **Risco Original** | P:1 × S:5 = 🟡 ALARP |
| **Controles Implementados** | |
| 1. Design inerentemente seguro | 5 canais independentes com providers diferentes |
| 2. Medidas de proteção | Recursos locais exibidos (CVV 188, SAMU 192) mesmo sem conexão |
| 3. Informação de segurança | Orientação para buscar ajuda presencial em caso de falha |
| **Verificação** | Teste de falha total implementado |
| **Risco Residual** | P:1 × S:5 = 🟡 ALARP (aceitável - probabilidade muito baixa) |

---

#### R-005: Vazamento de dados

| Aspecto | Detalhe |
|---------|---------|
| **Risco Original** | P:2 × S:3 = 🟢 Aceitável |
| **Controles Implementados** | |
| 1. Design inerentemente seguro | Autenticação JWT, senhas bcrypt (cost 14) |
| 2. Medidas de proteção | HTTPS, trilha de auditoria LGPD |
| 3. Informação de segurança | Política de privacidade, consentimento explícito |
| **Verificação** | Testes de autenticação, auditoria LGPD implementada |
| **Risco Residual** | P:1 × S:3 = 🟢 Aceitável |

---

### 5.2 Resumo de Riscos Residuais

| ID | Risco | Residual | Status |
|----|-------|----------|--------|
| R-001 | Score incorreto | 🟡 ALARP | ✅ Aceitável |
| R-002 | Risco subestimado | 🟢 Aceitável | ✅ OK |
| R-003 | Alerta não entregue | 🟡 ALARP | ✅ Aceitável |
| R-004 | Todos canais falham | 🟡 ALARP | ✅ Aceitável |
| R-005 | Vazamento de dados | 🟢 Aceitável | ✅ OK |
| R-006 | Diagnóstico incorreto | 🟡 ALARP | ✅ Aceitável |
| R-007 | Sistema indisponível | 🟢 Aceitável | ✅ OK |

**Conclusão:** Todos os riscos identificados estão em níveis aceitáveis ou ALARP após a implementação dos controles.

---

## 6. AVALIAÇÃO DE RISCO RESIDUAL GERAL

### 6.1 Risco-Benefício

**Benefícios esperados:**
1. Detecção precoce de risco suicida
2. Suporte emocional contínuo entre consultas
3. Redução de isolamento social
4. Alertas rápidos para cuidadores
5. Triagem inicial para otimizar recursos de saúde mental

**Riscos residuais:**
1. Possibilidade remota de score incorreto (mitigado por testes)
2. Possibilidade remota de falha total de alertas (mitigado por redundância)

**Avaliação:** Os benefícios superam significativamente os riscos residuais. O produto é considerado seguro para o uso pretendido.

---

## 7. INFORMAÇÕES PARA PRODUÇÃO E PÓS-PRODUÇÃO

### 7.1 Monitoramento Pós-Mercado

| Métrica | Monitoramento | Ação se Threshold |
|---------|---------------|-------------------|
| Taxa de falha de alertas | Prometheus `alerts_failed_total` | Investigar se >1% |
| Scores fora do padrão | Auditoria de `clinical_assessments` | Investigar anomalias |
| Reclamações de usuários | Sistema de tickets | Análise de causa raiz |
| Eventos adversos | Relatório obrigatório ANVISA | Notificação em 72h |

### 7.2 Critérios para Revisão

A gestão de riscos deve ser revisada quando:
1. Mudança significativa no software (nova funcionalidade clínica)
2. Evento adverso reportado
3. Feedback negativo recorrente
4. Mudança regulatória aplicável
5. Anualmente (revisão programada)

---

## 8. HISTÓRICO DE REVISÕES

| Versão | Data | Autor | Descrição |
|--------|------|-------|-----------|
| 0.1 | 2026-01-27 | Auto-gerado | Versão inicial (draft) |

---

## 9. APROVAÇÕES

| Função | Nome | Assinatura | Data |
|--------|------|------------|------|
| Elaborado por | | | |
| Revisado por | | | |
| Aprovado por | | | |

---

## ANEXOS

### A. Referências
- ISO 14971:2019 - Medical devices - Application of risk management to medical devices
- RDC 751/2022 - ANVISA
- IEC 62304:2006/Amd1:2015 - Medical device software — Software life cycle processes

### B. Documentos Relacionados
- Arquivo de Usabilidade IEC 62366-1 (a elaborar)
- Plano de Validação de Software (a elaborar)
- Dossiê Técnico ANVISA (a elaborar)
