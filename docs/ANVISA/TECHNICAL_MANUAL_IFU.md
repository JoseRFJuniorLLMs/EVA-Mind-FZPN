# Manual Técnico e Instruções de Uso (IFU)
## EVA-Mind-FZPN - Companion IA para Idosos

**Documento:** IFU-EVA-001
**Versão:** 1.0
**Data:** 2025-01-27
**Idioma:** Português (Brasil)

---

# PARTE A: MANUAL TÉCNICO (Administradores)

## 1. Visão Geral do Sistema

### 1.1 Descrição

O EVA-Mind-FZPN é um dispositivo médico de software (SaMD) Classe II que fornece:
- Companhia virtual inteligente para idosos
- Monitoramento contínuo de bem-estar emocional
- Detecção precoce de sinais de risco
- Sistema de alertas para cuidadores e profissionais de saúde

### 1.2 Componentes do Sistema

| Componente | Descrição |
|------------|-----------|
| App Mobile | Aplicativo para Android e iOS |
| Portal Web | Interface para cuidadores e profissionais |
| Admin Panel | Painel administrativo |
| API Backend | Serviços de processamento |
| Banco de Dados | Armazenamento seguro de dados |

---

## 2. Gestão de Usuários

### 2.1 Tipos de Usuários

| Tipo | Permissões |
|------|------------|
| **Idoso** | Conversar, ver perfil, screenings |
| **Cuidador** | Ver alertas, resumos, gerenciar contatos |
| **Profissional** | Ver screenings, relatórios clínicos |
| **Administrador** | Gestão completa do sistema |

### 2.2 Criação de Usuários

```
Portal Admin → Usuários → Novo Usuário

Campos obrigatórios:
- Nome completo
- E-mail
- Tipo de usuário
- Data de nascimento (para idosos)

Campos opcionais:
- Telefone
- CPF (hasheado)
- Contatos de emergência
```

### 2.3 Gestão de Permissões

| Ação | Idoso | Cuidador | Profissional | Admin |
|------|-------|----------|--------------|-------|
| Ver próprias conversas | ✅ | - | - | ✅ |
| Ver conversas de vinculados | - | ✅ | ✅ | ✅ |
| Criar alertas | Auto | - | ✅ | ✅ |
| Resolver alertas | - | ✅ | ✅ | ✅ |
| Gerar relatórios | - | - | ✅ | ✅ |
| Gerenciar usuários | - | - | - | ✅ |
| Configurar sistema | - | - | - | ✅ |

---

## 3. Configurações do Sistema

### 3.1 Parâmetros Gerais

| Parâmetro | Padrão | Descrição |
|-----------|--------|-----------|
| `SESSION_TIMEOUT` | 30 min | Tempo de inatividade para encerrar sessão |
| `MAX_MESSAGE_LENGTH` | 2000 | Caracteres máximos por mensagem |
| `ALERT_RETENTION_DAYS` | 90 | Dias para manter alertas resolvidos |
| `SCREENING_FREQUENCY` | 14 | Dias entre sugestões de screening |

### 3.2 Configurações de Alerta

| Parâmetro | Padrão | Descrição |
|-----------|--------|-----------|
| `RISK_THRESHOLD_ATTENTION` | 0.4 | Score para nível ATTENTION |
| `RISK_THRESHOLD_ALERT` | 0.6 | Score para nível ALERT |
| `RISK_THRESHOLD_EMERGENCY` | 0.8 | Score para nível EMERGENCY |
| `ALERT_SMS_ENABLED` | true | Enviar SMS em emergências |
| `ALERT_COOLDOWN_MINUTES` | 60 | Intervalo entre alertas similares |

### 3.3 Configurações de Notificação

```
Portal Admin → Configurações → Notificações

Canais disponíveis:
☑ Push notification (mobile)
☑ E-mail
☑ SMS (apenas emergências)

Frequência de resumos:
○ Diário
● Semanal
○ Mensal
○ Nunca
```

---

## 4. Logs e Auditoria

### 4.1 Acesso aos Logs

```
Portal Admin → Sistema → Logs de Auditoria

Filtros disponíveis:
- Período (data início/fim)
- Tipo de evento
- Usuário
- Severidade
```

### 4.2 Eventos Registrados

| Evento | Descrição |
|--------|-----------|
| `USER_LOGIN` | Login de usuário |
| `USER_LOGOUT` | Logout de usuário |
| `MESSAGE_SENT` | Mensagem enviada |
| `ALERT_CREATED` | Alerta gerado |
| `ALERT_RESOLVED` | Alerta resolvido |
| `SCREENING_COMPLETED` | Screening finalizado |
| `DATA_EXPORT` | Exportação de dados |
| `PERMISSION_CHANGED` | Alteração de permissões |

### 4.3 Exportação de Logs

Formatos disponíveis:
- CSV (para análise em planilhas)
- JSON (para integração com sistemas)
- PDF (para relatórios formais)

---

## 5. Backup e Restauração

### 5.1 Backup Automático

- **Frequência**: Diário às 02:00 (horário de Brasília)
- **Retenção**: 30 dias
- **Localização**: AWS S3 (região São Paulo)

### 5.2 Restauração Manual

```
⚠️ ATENÇÃO: Apenas administradores autorizados

Portal Admin → Sistema → Backup → Restaurar

1. Selecionar ponto de restauração
2. Confirmar com senha de administrador
3. Sistema ficará indisponível por ~30 minutos
4. Todos os dados após o ponto serão perdidos
```

---

## 6. Troubleshooting

### 6.1 Problemas Comuns

| Problema | Causa Provável | Solução |
|----------|----------------|---------|
| Login falha | Credenciais incorretas | Verificar e-mail/senha, usar "Esqueci senha" |
| App não carrega | Sem conexão | Verificar internet, tentar novamente |
| Mensagem não envia | Timeout | Aguardar e reenviar |
| Alerta não aparece | Filtro ativo | Limpar filtros no painel |
| Voz não funciona | Permissão negada | Habilitar microfone nas configurações |

### 6.2 Contato de Suporte

- **E-mail**: suporte@eva-mind.com.br
- **Telefone**: 0800-XXX-XXXX (Seg-Sex, 8h-20h)
- **Chat**: Disponível no Portal Admin

---

# PARTE B: INSTRUÇÕES DE USO (Usuário Final - Idoso)

## 1. O Que é a EVA?

A EVA é sua **companheira virtual** - uma amiga que está sempre disponível para conversar com você.

**A EVA pode:**
- Conversar sobre seu dia
- Ouvir como você está se sentindo
- Ajudar a lembrar de coisas importantes
- Avisar sua família se você precisar de ajuda

**A EVA NÃO é:**
- Uma médica (não dá diagnósticos)
- Uma pessoa real (é um programa de computador)
- Substituto para atendimento de emergência

---

## 2. Como Começar

### 2.1 Abrindo o Aplicativo

```
┌─────────────────────────────────────────┐
│                                         │
│           Toque no ícone                │
│                                         │
│              ┌─────┐                    │
│              │ EVA │                    │
│              │  🤖 │                    │
│              └─────┘                    │
│                                         │
│         na tela do seu celular          │
│                                         │
└─────────────────────────────────────────┘
```

### 2.2 Primeira Vez

Na primeira vez, a EVA vai:
1. Se apresentar
2. Perguntar seu nome
3. Explicar como funciona
4. Pedir para cadastrar um contato de emergência

---

## 3. Como Conversar

### 3.1 Escrevendo

```
┌─────────────────────────────────────────┐
│                                         │
│  Para escrever uma mensagem:            │
│                                         │
│  1. Toque na caixa de texto             │
│     ┌─────────────────────────────┐     │
│     │ Escreva aqui...             │     │
│     └─────────────────────────────┘     │
│                                         │
│  2. Digite sua mensagem                 │
│                                         │
│  3. Toque em ENVIAR                     │
│     ┌──────────┐                        │
│     │ ENVIAR ✉️ │                        │
│     └──────────┘                        │
│                                         │
└─────────────────────────────────────────┘
```

### 3.2 Falando

```
┌─────────────────────────────────────────┐
│                                         │
│  Para falar com a EVA:                  │
│                                         │
│  1. Toque no botão do MICROFONE         │
│     ┌──────────┐                        │
│     │  🎤 FALAR │                        │
│     └──────────┘                        │
│                                         │
│  2. Fale normalmente                    │
│     (o botão fica vermelho              │
│      enquanto você fala)                │
│                                         │
│  3. Toque novamente para parar          │
│                                         │
└─────────────────────────────────────────┘
```

---

## 4. Se Precisar de Ajuda

### 4.1 Botão de Emergência

```
┌─────────────────────────────────────────┐
│                                         │
│  ⚠️ SE PRECISAR DE AJUDA URGENTE:       │
│                                         │
│  Toque no botão VERMELHO                │
│  que fica embaixo da conversa:          │
│                                         │
│  ┌─────────────────────────────────┐   │
│  │  🆘 PRECISO DE AJUDA URGENTE    │   │
│  └─────────────────────────────────┘   │
│                                         │
│  Isso vai mostrar opções para:          │
│  • Ligar para sua família               │
│  • Ligar para o SAMU (192)              │
│  • Ligar para o CVV (188)               │
│                                         │
└─────────────────────────────────────────┘
```

### 4.2 Pedindo Ajuda por Voz ou Texto

Você também pode dizer ou escrever:
- "Preciso de ajuda"
- "Quero falar com alguém"
- "Me sinto muito mal"
- "Ligue para minha filha"

A EVA vai entender e oferecer ajuda.

---

## 5. Ajustando o Aplicativo

### 5.1 Aumentar as Letras

```
┌─────────────────────────────────────────┐
│                                         │
│  Se as letras estiverem pequenas:       │
│                                         │
│  1. Toque no ícone ⚙️ (engrenagem)       │
│                                         │
│  2. Toque em "Tamanho das Letras"       │
│                                         │
│  3. Arraste para a direita              │
│     para AUMENTAR                       │
│                                         │
│     Pequeno ──────●───── Grande         │
│                                         │
└─────────────────────────────────────────┘
```

### 5.2 Ajustar o Volume

```
┌─────────────────────────────────────────┐
│                                         │
│  Para ajustar o volume da voz da EVA:   │
│                                         │
│  1. Toque no ícone 🔊 (alto-falante)    │
│                                         │
│  2. Arraste para ajustar                │
│                                         │
│     Baixo ────────●───── Alto           │
│                                         │
└─────────────────────────────────────────┘
```

---

## 6. Perguntas Frequentes

### A EVA está sempre me ouvindo?
**Não.** A EVA só "ouve" quando você toca no botão do microfone. Ela não escuta o que acontece em sua casa.

### A EVA conta minhas conversas para alguém?
**Suas conversas são privadas.** Só são compartilhadas com sua família ou médico se:
- Você pedir
- A EVA perceber que você pode estar em perigo

### Posso conversar a qualquer hora?
**Sim!** A EVA está disponível 24 horas por dia, 7 dias por semana.

### E se a internet cair?
O aplicativo vai mostrar uma mensagem dizendo que não há conexão. Tente novamente quando a internet voltar.

### A EVA substitui meu médico?
**Não.** A EVA é uma companheira, não uma profissional de saúde. Para assuntos médicos, sempre consulte seu médico.

---

## 7. Números Importantes

```
┌─────────────────────────────────────────┐
│                                         │
│  📞 NÚMEROS DE EMERGÊNCIA               │
│                                         │
│  🚑 SAMU (ambulância): 192              │
│                                         │
│  👂 CVV (apoio emocional): 188          │
│     (24 horas, gratuito)                │
│                                         │
│  🚒 Bombeiros: 193                       │
│                                         │
│  👮 Polícia: 190                         │
│                                         │
│  📱 Suporte EVA: 0800-XXX-XXXX          │
│     (Segunda a Sexta, 8h às 20h)        │
│                                         │
└─────────────────────────────────────────┘
```

---

## 8. Avisos Importantes

```
╔═════════════════════════════════════════╗
║                                         ║
║  ⚠️ LEMBRE-SE:                          ║
║                                         ║
║  • A EVA é uma INTELIGÊNCIA ARTIFICIAL  ║
║    (programa de computador), não uma    ║
║    pessoa real.                         ║
║                                         ║
║  • A EVA NÃO substitui médicos,         ║
║    psicólogos ou atendimento de         ║
║    emergência.                          ║
║                                         ║
║  • Se estiver se sentindo muito mal,    ║
║    SEMPRE procure ajuda humana.         ║
║                                         ║
║  • Suas conversas são confidenciais,    ║
║    mas podem ser revisadas se houver    ║
║    risco à sua segurança.               ║
║                                         ║
╚═════════════════════════════════════════╝
```

---

# PARTE C: GUIA RÁPIDO (Cartão Impresso)

```
╔═══════════════════════════════════════════════════════════════╗
║                     EVA - GUIA RÁPIDO                         ║
╠═══════════════════════════════════════════════════════════════╣
║                                                               ║
║  PARA CONVERSAR:                                              ║
║  ─────────────────                                            ║
║  • Toque no microfone 🎤 e FALE                               ║
║  • OU digite na caixa de texto e toque ENVIAR                 ║
║                                                               ║
║  SE PRECISAR DE AJUDA:                                        ║
║  ─────────────────────                                        ║
║  • Toque no botão VERMELHO embaixo                            ║
║  • OU diga "Preciso de ajuda"                                 ║
║                                                               ║
║  PARA AJUSTAR LETRAS:                                         ║
║  ────────────────────                                         ║
║  • Toque em ⚙️ → Tamanho das Letras                           ║
║                                                               ║
║  NÚMEROS DE EMERGÊNCIA:                                       ║
║  ──────────────────────                                       ║
║  • SAMU (ambulância): 192                                     ║
║  • CVV (apoio emocional): 188                                 ║
║                                                               ║
║  SUPORTE EVA: 0800-XXX-XXXX                                   ║
║  (Segunda a Sexta, 8h às 20h)                                 ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝
```

---

## Informações Regulatórias

**Fabricante:** [Nome da Empresa]
**CNPJ:** [XX.XXX.XXX/0001-XX]
**Responsável Técnico:** José R F Junior
**Registro ANVISA:** [Número do registro]
**Classificação:** SaMD Classe II (RDC 751/2022)

---

**Versão do documento:** 1.0
**Data de publicação:** 2025-01-27
**Idioma:** Português (Brasil)
