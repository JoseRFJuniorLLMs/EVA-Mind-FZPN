# 🏥 MedGemma - Análise de Imagens Médicas

## Visão Geral

O módulo **MedGemma** adiciona capacidades de análise de imagens médicas ao EVA-Mind, permitindo:
- 📋 Análise de receitas médicas com extração automática de medicamentos
- 🩹 Análise de feridas e lesões com avaliação de gravidade
- 🔬 Análise de resultados de exames (futuro)

## Componentes

### 1. **Service** (`service.go`)
Cliente principal do MedGemma para análise de imagens.

**Principais Funções**:
- `NewMedGemmaService(apiKey)`: Cria cliente configurado
- `AnalyzePrescription(ctx, imageData, mimeType)`: Analisa receita médica
- `AnalyzeWound(ctx, imageData, mimeType)`: Analisa ferida/lesão

### 2. **Audit Logger** (`audit.go`)
Gerencia auditoria e persistência no banco de dados.

**Funções**:
- `LogPrescriptionAnalysis()`: Salva análise de receita
- `LogWoundAnalysis()`: Salva análise de ferida
- `SaveMedicationsFromPrescription()`: Extrai e salva medicamentos
- `MarkNotified()`: Marca notificação enviada
- `GetPendingAlerts()`: Busca alertas pendentes

## Uso

### Análise de Receita Médica

```go
// Criar serviço
medgemma, err := medgemma.NewMedGemmaService(apiKey)
if err != nil {
    log.Fatal(err)
}

// Analisar imagem
analysis, err := medgemma.AnalyzePrescription(ctx, imageBytes, "image/jpeg")
if err != nil {
    log.Fatal(err)
}

// Resultado
fmt.Printf("Médico: %s (CRM: %s)\n", analysis.DoctorName, analysis.DoctorCRM)
fmt.Printf("Medicamentos encontrados: %d\n", len(analysis.Medications))

for _, med := range analysis.Medications {
    fmt.Printf("- %s %s - %s\n", med.Name, med.Dosage, med.Frequency)
}
```

### Análise de Ferida

```go
// Analisar ferida
analysis, err := medgemma.AnalyzeWound(ctx, imageBytes, "image/jpeg")
if err != nil {
    log.Fatal(err)
}

// Resultado
fmt.Printf("Tipo: %s\n", analysis.Type)
fmt.Printf("Gravidade: %s\n", analysis.Severity)
fmt.Printf("Requer atendimento: %v\n", analysis.SeekMedicalCare)

if analysis.SeekMedicalCare {
    fmt.Printf("Urgência: %s\n", analysis.Urgency)
    fmt.Printf("Recomendações: %v\n", analysis.Recommendations)
}
```

## Integração via WebSocket

### Tool: `analyze_medical_image`

**Parâmetros**:
- `image` (string): Imagem em base64
- `type` (string): Tipo de análise (`prescription`, `wound`, `lab_result`)

**Exemplo de Chamada**:
```json
{
  "type": "tool_call",
  "tool": "analyze_medical_image",
  "args": {
    "image": "data:image/jpeg;base64,/9j/4AAQSkZJRg...",
    "type": "prescription"
  }
}
```

## Banco de Dados

### Tabela: `medical_image_analysis`

```sql
CREATE TABLE medical_image_analysis (
    id BIGSERIAL PRIMARY KEY,
    idoso_id BIGINT NOT NULL,
    image_type VARCHAR(50),
    analysis_result JSONB,
    severity VARCHAR(20),
    requires_medical_attention BOOLEAN,
    caregiver_notified BOOLEAN,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Views Úteis

- `v_analyzed_prescriptions`: Receitas analisadas
- `v_analyzed_wounds`: Feridas analisadas
- `v_medical_image_alerts`: Alertas pendentes

## Fluxo de Análise

### Receita Médica

```
1. Mobile envia imagem via WebSocket
2. EVA-Mind decodifica base64
3. MedGemma analisa imagem
4. Extrai medicamentos, médico, data
5. Salva no banco de dados
6. Atualiza tabela de medicamentos
7. Notifica cuidador se controlados
8. Retorna resultado para mobile
```

### Ferida/Lesão

```
1. Mobile envia foto da ferida
2. EVA-Mind decodifica base64
3. MedGemma analisa imagem
4. Avalia tipo, tamanho, gravidade
5. Detecta sinais de infecção
6. Salva análise no banco
7. Se grave → Notifica cuidador
8. Retorna recomendações
```

## Níveis de Gravidade (Feridas)

| Nível | Critérios | Ação |
|-------|-----------|------|
| **CRÍTICO** | Sangramento intenso, queimadura 3º grau, infecção severa | Notificação imediata + pronto-socorro |
| **ALTO** | Ferida profunda, sinais moderados de infecção | Notificação + consulta urgente |
| **MÉDIO** | Ferida superficial com sinais leves de infecção | Orientação + monitoramento |
| **BAIXO** | Ferida superficial limpa | Orientação de cuidados |

## Segurança e Compliance

### Disclaimers Obrigatórios

✅ Todas as análises incluem:
- "Esta é uma análise automatizada e não substitui avaliação médica profissional"
- Recomendação de consulta médica para casos graves
- Não fornece diagnósticos definitivos

### Auditoria Completa

- ✅ Todas as análises são registradas
- ✅ Imagens podem ser armazenadas (opcional)
- ✅ Notificações são rastreadas
- ✅ Logs sanitizados (sem PII desnecessário)

## Monitoramento

### Queries Úteis

```sql
-- Análises das últimas 24h
SELECT 
    image_type,
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE requires_medical_attention) as alertas
FROM medical_image_analysis
WHERE created_at >= NOW() - INTERVAL '24 hours'
GROUP BY image_type;

-- Alertas pendentes
SELECT * FROM v_medical_image_alerts;

-- Receitas analisadas hoje
SELECT * FROM v_analyzed_prescriptions
WHERE data_analise::date = CURRENT_DATE;

-- Feridas graves não notificadas
SELECT * FROM v_analyzed_wounds
WHERE gravidade IN ('ALTO', 'CRÍTICO')
  AND cuidador_notificado = false;
```

## Limitações Conhecidas

1. **Qualidade da Imagem**: Requer foto nítida e bem iluminada
2. **Caligrafia**: Receitas manuscritas podem ter baixa precisão
3. **Idioma**: Otimizado para português brasileiro
4. **Tipos de Lesão**: Melhor performance em feridas superficiais
5. **Não é Diagnóstico**: Sempre recomenda consulta médica

## Próximos Passos

1. ✅ Adicionar suporte para resultados de exames laboratoriais
2. ✅ Implementar OCR especializado para receitas manuscritas
3. ✅ Integração com banco de dados de medicamentos (ANVISA)
4. ✅ Análise de bulas e embalagens de medicamentos
5. ✅ Detecção de interações medicamentosas

## Troubleshooting

### Erro: "JSON não encontrado na resposta"
- Gemini pode retornar formato diferente
- Verificar logs para ver resposta completa
- Ajustar prompts se necessário

### Erro: "Erro ao decodificar imagem"
- Verificar se imagem está em base64 válido
- Confirmar formato (JPEG, PNG)
- Verificar tamanho da imagem (<10MB)

### Medicamentos não extraídos
- Verificar qualidade da foto
- Receita pode estar manuscrita (baixa precisão)
- Tentar com foto mais nítida

---

**Criado em**: 15 de janeiro de 2026  
**Versão**: 1.0
