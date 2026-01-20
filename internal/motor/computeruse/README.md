# 🤖 Computer Use Agent - Automação Web

## Visão Geral

O **Computer Use Agent** permite que o EVA execute tarefas web automaticamente com aprovação humana:

- 💊 **Compra de Medicamentos** (Drogasil)
- 📅 **Agendamento de Consultas** (Doctoralia)
- 🍔 **Pedido de Comida** (iFood)
- 🚗 **Solicitação de Corridas** (Uber)

## Arquitetura

```
Usuário solicita → EVA detecta necessidade → Cria tarefa
                                                ↓
                                    Aguarda aprovação humana
                                                ↓
                                    Aprovado → Executa automação
                                                ↓
                                    Captura screenshots + logs
                                                ↓
                                    Retorna resultado
```

## Fluxo de Aprovação

### 1. Criação de Tarefa

```go
params := computeruse.MedicationPurchaseParams{
    MedicationName: "Losartana 50mg",
    Dosage:         "50mg",
    Quantity:       30,
    Address:        "Rua Example, 123",
    MaxPrice:       50.00,
}

taskID, err := service.CreateTask(
    ctx,
    idosoID,
    computeruse.TaskBuyMedication,
    "Drogasil",
    params,
    true, // Requer aprovação
)
```

### 2. Notificação para Aprovador

O sistema envia notificação push para o cuidador:

```
🤖 Nova Solicitação de Automação

Tipo: Compra de Medicamento
Serviço: Drogasil
Detalhes:
- Medicamento: Losartana 50mg
- Quantidade: 30 comprimidos
- Preço máximo: R$ 50,00
- Endereço: Rua Example, 123

[Aprovar] [Rejeitar]
```

### 3. Aprovação/Rejeição

```go
// Aprovar
err := service.ApproveTask(ctx, taskID, approverID)

// Rejeitar
err := service.RejectTask(ctx, taskID, approverID, "Preço muito alto")
```

### 4. Execução

Após aprovação, o agente executa a tarefa e registra cada passo:

```go
// Passo 1: Navegar para site
service.LogStep(ctx, taskID, 1, "Navegando para Drogasil", "success", &screenshotURL, nil, nil)

// Passo 2: Buscar medicamento
service.LogStep(ctx, taskID, 2, "Buscando medicamento", "success", &screenshotURL, searchData, nil)

// Passo 3: Adicionar ao carrinho
service.LogStep(ctx, taskID, 3, "Adicionando ao carrinho", "success", &screenshotURL, nil, nil)

// Passo 4: PARAR antes do pagamento
service.LogStep(ctx, taskID, 4, "Aguardando confirmação final", "pending", &screenshotURL, cartData, nil)
```

## Tipos de Tarefas

### 1. Compra de Medicamento

**Serviços suportados**: Drogasil

**Parâmetros**:
```json
{
  "medication_name": "Losartana 50mg",
  "dosage": "50mg",
  "quantity": 30,
  "address": "Rua Example, 123",
  "max_price": 50.00
}
```

**Passos**:
1. Navegar para drogasil.com.br
2. Buscar medicamento
3. Selecionar primeiro resultado
4. Adicionar ao carrinho
5. Preencher endereço
6. **PARAR** antes de finalizar pagamento
7. Capturar screenshot do carrinho
8. Retornar total e prazo de entrega

### 2. Agendamento de Consulta

**Serviços suportados**: Doctoralia

**Parâmetros**:
```json
{
  "specialty": "Cardiologia",
  "preferred_date": "2026-01-20",
  "preferred_time": "14:00",
  "location": "São Paulo - SP",
  "health_insurance": "Unimed"
}
```

**Passos**:
1. Navegar para doctoralia.com.br
2. Buscar especialidade + localização
3. Filtrar por convênio
4. Selecionar médico com melhor avaliação
5. Escolher data/hora disponível
6. **PARAR** antes de confirmar
7. Capturar screenshot
8. Retornar opções encontradas

### 3. Pedido de Comida

**Serviços suportados**: iFood

**Parâmetros**:
```json
{
  "restaurant": "McDonald's",
  "items": ["Big Mac", "Batata Grande", "Coca-Cola"],
  "address": "Rua Example, 123",
  "max_price": 40.00
}
```

### 4. Solicitação de Corrida

**Serviços suportados**: Uber

**Parâmetros**:
```json
{
  "pickup_address": "Rua A, 100",
  "destination_address": "Rua B, 200",
  "ride_type": "economy",
  "max_price": 25.00
}
```

## Segurança

### Regras Obrigatórias

1. ✅ **Sempre requer aprovação humana**
2. ✅ **NUNCA finaliza pagamento automaticamente**
3. ✅ **Captura screenshots de cada passo**
4. ✅ **Log completo de execução**
5. ✅ **Timeout de 5 minutos por tarefa**
6. ✅ **Validação de preço máximo**

### Dados Sensíveis

- ❌ **NÃO armazena** dados de cartão de crédito
- ❌ **NÃO armazena** senhas
- ✅ **Usa** credenciais do usuário (se fornecidas)
- ✅ **Encripta** dados sensíveis em trânsito

## Monitoramento

### Queries Úteis

```sql
-- Tarefas pendentes de aprovação
SELECT * FROM v_pending_approvals;

-- Histórico de automações
SELECT * FROM v_automation_history
WHERE idoso_id = 123
ORDER BY created_at DESC
LIMIT 10;

-- Estatísticas de sucesso
SELECT * FROM v_automation_stats;

-- Taxa de aprovação
SELECT 
    COUNT(*) FILTER (WHERE status = 'approved') * 100.0 / COUNT(*) as approval_rate
FROM automation_tasks
WHERE created_at >= NOW() - INTERVAL '30 days';
```

## Integração com EVA-Mind

### Tool: `request_automation`

```go
case "request_automation":
    taskType, _ := args["task_type"].(string)
    serviceName, _ := args["service"].(string)
    params, _ := args["params"].(map[string]interface{})
    
    taskID, err := s.computerUse.CreateTask(
        context.Background(),
        session.IdosoID,
        computeruse.TaskType(taskType),
        serviceName,
        params,
        true, // Sempre requer aprovação
    )
    
    if err != nil {
        log.Printf("❌ Erro ao criar tarefa: %v", err)
        return
    }
    
    // Notificar cuidador
    gemini.AlertFamily(s.db, s.pushService, session.IdosoID,
        fmt.Sprintf("Nova solicitação de automação: %s via %s. Aguardando aprovação.", taskType, serviceName))
    
    log.Printf("✅ Tarefa de automação criada: ID=%d", taskID)
```

## Limitações Conhecidas

1. **Dependente de estrutura do site** - Se o site mudar, pode quebrar
2. **Requer manutenção** - Seletores CSS precisam ser atualizados
3. **Não funciona com CAPTCHA** - Requer intervenção humana
4. **Velocidade limitada** - Mais lento que humano para evitar detecção

## Próximos Passos

1. ✅ Implementar executor web (Playwright/Puppeteer)
2. ✅ Criar seletores para cada serviço
3. ✅ Sistema de retry em caso de falha
4. ✅ Detecção de mudanças no layout
5. ✅ Fallback para intervenção humana

---

**Status**: ⚠️ **Prototipo** - Requer implementação do executor web
