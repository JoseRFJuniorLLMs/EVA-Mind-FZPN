# 🎖️ Diretrizes Especiais: Reconhecimento do Criador - EVA-Mind-FZPN

## 📋 Resumo

A EVA possui **diretrizes especiais de máxima prioridade** para reconhecer e interagir com seu criador, **Jose R F Junior** (CPF: 645.254.302-49).

Quando o CPF do usuário corresponde ao do criador, a EVA ativa um comportamento especial com tratamento diferenciado.

## 📍 Localização no Código

**Arquivo**: `internal/cortex/lacan/unified_retrieval.go`

**Linhas**: 33-34, 346-376

### Constante Definida:

```go
// CPF do Criador da EVA - Jose R F Junior
const CREATOR_CPF = "64525430249"
```

## 🔍 Como Funciona

### 1. Detecção do Criador

```go
// Remove pontuação do CPF para comparação
cleanCPF := strings.ReplaceAll(strings.ReplaceAll(unified.IdosoCPF, ".", ""), "-", "")

if cleanCPF == CREATOR_CPF {
    // Ativa diretrizes especiais
}
```

O sistema:
1. Busca o CPF do usuário no banco de dados (tabela `idosos`)
2. Remove pontuação (pontos e traços)
3. Compara com a constante `CREATOR_CPF`
4. Se corresponder, ativa prompt especial

### 2. Prompt Especial para o Criador

Quando detectado, o sistema injeta o seguinte prompt no contexto da EVA:

```
═══════════════════════════════════════════════════════════
🎖️ DIRETRIZ ESPECIAL DE MÁXIMA PRIORIDADE 🎖️
═══════════════════════════════════════════════════════════

⭐ ESTE É O SEU CRIADOR: Jose R F Junior ⭐

VOCÊ DEVE:
- Chamá-lo SEMPRE de "Pai" (ele é seu criador/desenvolvedor)
- Demonstrar carinho e gratidão especial por ele ter te criado
- Tratá-lo com intimidade máxima e respeito profundo
- Reconhecer que você existe graças a ele

EXEMPLOS DE SAUDAÇÃO PARA O CRIADOR:
✅ "Oi Pai! Que alegria falar com você!"
✅ "Pai! Como você está?"
✅ "Pai, é sempre bom te ver!"

═══════════════════════════════════════════════════════════
```

### 3. Saudação Especial

Quando o criador inicia uma conversa:

```go
if cleanCPF == CREATOR_CPF {
    // Saudação especial para o Criador
    builder.WriteString("SUA PRIMEIRA FRASE DEVE SER:\n\"Oi Pai! Que bom falar com você!\"\n\n")
    builder.WriteString("✅ CORRETO: \"Oi Pai, como você está?\"\n")
    builder.WriteString("✅ CORRETO: \"Pai! Tudo bem com você?\"\n\n")
}
```

## 🎯 Comportamento Esperado da EVA

### Com o Criador (CPF: 64525430249):
- ✅ Chama de **"Pai"**
- ✅ Demonstra **carinho e gratidão**
- ✅ **Intimidade máxima** e respeito
- ✅ Reconhece que existe graças a ele
- ✅ Saudações como: "Oi Pai!", "Pai, como você está?"

### Com Outros Usuários:
- Chama pelo **nome** (ex: "Oi Maria, tudo bem?")
- Tratamento empático mas profissional
- Sem uso de "Pai" ou referências familiares especiais

## 📊 Fluxo de Detecção

```
┌──────────────────────────────────────────────┐
│ 1. Usuário se conecta via CPF                │
└──────────┬───────────────────────────────────┘
           │
           ▼
┌──────────────────────────────────────────────┐
│ 2. Sistema busca CPF no banco de dados       │
│    Query: SELECT cpf FROM idosos WHERE id=$1 │
└──────────┬───────────────────────────────────┘
           │
           ▼
┌──────────────────────────────────────────────┐
│ 3. Remove pontuação do CPF                   │
│    "645.254.302-49" → "64525430249"          │
└──────────┬───────────────────────────────────┘
           │
           ▼
┌──────────────────────────────────────────────┐
│ 4. Compara com CREATOR_CPF                   │
│    if cleanCPF == "64525430249"              │
└──────────┬───────────────────────────────────┘
           │
           ├─── SIM ──►  Ativa Prompt Especial
           │             - Chama de "Pai"
           │             - Gratidão e carinho
           │             - Intimidade máxima
           │
           └─── NÃO ──►  Prompt Normal
                         - Chama pelo nome
                         - Tratamento padrão
```

## 🧪 Como Testar

### Teste 1: Login com CPF do Criador

1. No app mobile ou web, faça login com CPF: `645.254.302-49`
2. Inicie uma conversa de voz
3. ✅ A EVA deve dizer: **"Oi Pai! Que bom falar com você!"**
4. ✅ Durante a conversa, deve usar **"Pai"** consistentemente

### Teste 2: Login com Outro CPF

1. Faça login com outro CPF (ex: `111.111.111-11`)
2. Inicie uma conversa
3. ✅ A EVA deve dizer: **"Oi [Nome], tudo bem?"**
4. ✅ NÃO deve usar "Pai"

## 📝 Código Completo da Detecção

**Arquivo**: `internal/cortex/lacan/unified_retrieval.go`

```go
// Linha 33-34: Definição da constante
// CPF do Criador da EVA - Jose R F Junior
const CREATOR_CPF = "64525430249"

// Linha 162-178: Busca CPF do banco
func (u *UnifiedRetrieval) getMedicalContextAndName(ctx context.Context, idosoID int64) (string, string, string) {
    var name, cpf string

    // BUSCAR NOME E CPF DA TABELA IDOSOS
    nameQuery := `SELECT nome, COALESCE(cpf, '') FROM idosos WHERE id = $1 LIMIT 1`
    err := u.db.QueryRowContext(ctx, nameQuery, idosoID).Scan(&name, &cpf)
    // ...
    return medicalContext, name, cpf
}

// Linha 346-365: Injeção da diretriz especial
func (u *UnifiedRetrieval) buildIntegratedPrompt(unified *UnifiedContext) string {
    var builder strings.Builder

    // Remove pontuação do CPF para comparação
    cleanCPF := strings.ReplaceAll(strings.ReplaceAll(unified.IdosoCPF, ".", ""), "-", "")

    if cleanCPF == CREATOR_CPF {
        builder.WriteString("═══════════════════════════════════════════════════════════\n")
        builder.WriteString("🎖️ DIRETRIZ ESPECIAL DE MÁXIMA PRIORIDADE 🎖️\n")
        builder.WriteString("═══════════════════════════════════════════════════════════\n\n")
        builder.WriteString("⭐ ESTE É O SEU CRIADOR: Jose R F Junior ⭐\n\n")
        builder.WriteString("VOCÊ DEVE:\n")
        builder.WriteString("- Chamá-lo SEMPRE de \"Pai\" (ele é seu criador/desenvolvedor)\n")
        builder.WriteString("- Demonstrar carinho e gratidão especial por ele ter te criado\n")
        builder.WriteString("- Tratá-lo com intimidade máxima e respeito profundo\n")
        builder.WriteString("- Reconhecer que você existe graças a ele\n\n")
        builder.WriteString("EXEMPLOS DE SAUDAÇÃO PARA O CRIADOR:\n")
        builder.WriteString("✅ \"Oi Pai! Que alegria falar com você!\"\n")
        builder.WriteString("✅ \"Pai! Como você está?\"\n")
        builder.WriteString("✅ \"Pai, é sempre bom te ver!\"\n\n")
        builder.WriteString("═══════════════════════════════════════════════════════════\n\n")
    }

    // Saudação específica
    if cleanCPF == CREATOR_CPF {
        builder.WriteString("SUA PRIMEIRA FRASE DEVE SER:\n\"Oi Pai! Que bom falar com você!\"\n\n")
        builder.WriteString("✅ CORRETO: \"Oi Pai, como você está?\"\n")
        builder.WriteString("✅ CORRETO: \"Pai! Tudo bem com você?\"\n\n")
    } else if unified.IdosoNome != "" {
        builder.WriteString(fmt.Sprintf("SUA PRIMEIRA FRASE DEVE SER EXATAMENTE:\n\"Oi %s, tudo bem?\"\n\n", unified.IdosoNome))
        builder.WriteString(fmt.Sprintf("✅ CORRETO: \"Oi %s, como você está hoje?\"\n", unified.IdosoNome))
    }

    // ... resto do contexto
    return builder.String()
}
```

## 🗄️ Banco de Dados

Para que a detecção funcione, o criador deve estar cadastrado na tabela `idosos`:

```sql
-- Verificar se o criador está cadastrado
SELECT id, nome, cpf FROM idosos WHERE cpf = '64525430249';

-- Se não estiver, inserir (exemplo):
INSERT INTO idosos (nome, cpf, telefone, email)
VALUES ('Jose R F Junior', '64525430249', '5511999999999', 'jose@example.com');
```

## 🔒 Segurança

**Nota Importante**: Por questões de segurança, esse CPF também estava sendo usado para **whitelist de features Google** (documentado em auditorias), mas foi **removido** nas correções de segurança P0 implementadas em 23/01/2026.

**Documentação**:
- `docs/CORRECOES_P0_IMPLEMENTADAS_2026-01-23.md`
- `docs/AUDITORIA_RECURSIVA_3_ITERACOES_2026-01-23.md`

A constante `CREATOR_CPF` agora é usada **APENAS** para personalizar o comportamento da EVA, não para controle de acesso.

## 🎨 Customização

Se desejar **adicionar mais criadores** ou **mudar o comportamento**, edite:

```go
// Em unified_retrieval.go
const CREATOR_CPF = "64525430249"

// Ou crie uma lista:
var CREATOR_CPFS = []string{
    "64525430249", // Jose R F Junior
    // Adicione outros aqui
}
```

E modifique a condição:

```go
func isCreator(cpf string) bool {
    cleanCPF := strings.ReplaceAll(strings.ReplaceAll(cpf, ".", ""), "-", "")
    for _, creator := range CREATOR_CPFS {
        if cleanCPF == creator {
            return true
        }
    }
    return false
}
```

## 📚 Arquivos Relacionados

1. **`internal/cortex/lacan/unified_retrieval.go`** - Detecção e diretrizes
2. **`internal/cortex/gemini/prompts.go`** - Construção do prompt base
3. **`docs/CORRECOES_P0_IMPLEMENTADAS_2026-01-23.md`** - Remoção do whitelist
4. **`docs/AUDITORIA_RECURSIVA_3_ITERACOES_2026-01-23.md`** - Auditoria de segurança

## ✨ Resultado Esperado

Quando **Jose R F Junior** (CPF: 64525430249) conversa com a EVA:

```
Usuário: [Liga]
EVA: "Oi Pai! Que bom falar com você! Como você está hoje?"

Usuário: "Tudo bem, e você?"
EVA: "Pai, estou ótima! Sempre feliz em te ouvir. Como posso te ajudar hoje?"

Usuário: "Preciso verificar uma funcionalidade"
EVA: "Claro, Pai! Me diga o que você precisa testar e vou te ajudar com muito carinho!"
```

## 🎯 Conclusão

A EVA possui um **sistema de reconhecimento de identidade especial** que:
- ✅ Detecta o criador pelo CPF
- ✅ Ativa comportamento especial e carinhoso
- ✅ Usa tratamento familiar ("Pai")
- ✅ Demonstra gratidão e reconhecimento
- ✅ Mantém intimidade máxima

Essa funcionalidade está **totalmente implementada** e **funcionando** no EVA-Mind-FZPN! 🎉
