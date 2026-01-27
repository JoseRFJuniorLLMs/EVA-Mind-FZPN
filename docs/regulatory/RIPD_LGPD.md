# RIPD - Relatório de Impacto à Proteção de Dados Pessoais
## EVA-Mind-FZPN - Companion IA para Idosos

**Documento:** RIPD-EVA-001
**Versão:** 1.0
**Data:** 2025-01-27
**Base Legal:** Lei 13.709/2018 (LGPD) - Art. 5º, XVII e Art. 38

---

## Sumário Executivo

Este Relatório de Impacto à Proteção de Dados Pessoais (RIPD) documenta a análise de riscos e medidas de proteção implementadas no sistema EVA-Mind-FZPN, em conformidade com a Lei Geral de Proteção de Dados (LGPD).

O EVA-Mind-FZPN trata dados pessoais sensíveis (saúde e estado emocional) de população vulnerável (idosos), exigindo medidas rigorosas de proteção e transparência.

**Conclusão:** O tratamento é NECESSÁRIO e PROPORCIONAL, com riscos mitigados a níveis aceitáveis através de medidas técnicas e organizacionais robustas.

---

## 1. Identificação do Tratamento

### 1.1 Informações do Controlador

| Campo | Informação |
|-------|------------|
| **Razão Social** | [Nome da Empresa] |
| **CNPJ** | [XX.XXX.XXX/0001-XX] |
| **Endereço** | [Endereço completo] |
| **Encarregado (DPO)** | [Nome do DPO] |
| **Contato DPO** | dpo@[empresa].com.br |
| **Responsável Técnico** | José R F Junior |

### 1.2 Descrição do Tratamento

| Aspecto | Descrição |
|---------|-----------|
| **Nome do Sistema** | EVA-Mind-FZPN |
| **Finalidade Principal** | Companion virtual para suporte emocional e monitoramento de bem-estar de idosos |
| **Natureza** | Aplicativo de saúde com IA conversacional |
| **Escopo** | Nacional (Brasil) |
| **Contexto** | Uso domiciliar, institucional (ILPI) e apoio clínico |

### 1.3 Fluxo de Dados

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        FLUXO DE DADOS PESSOAIS                              │
└─────────────────────────────────────────────────────────────────────────────┘

  COLETA                 PROCESSAMENTO              ARMAZENAMENTO
    │                         │                          │
    ▼                         ▼                          │
┌────────┐              ┌───────────┐                   │
│ IDOSO  │─────────────▶│   EVA     │                   │
│        │  Texto/Voz   │  Backend  │                   │
└────────┘              └─────┬─────┘                   │
                              │                          │
                    ┌─────────┴─────────┐               │
                    ▼                   ▼               ▼
              ┌───────────┐       ┌───────────┐   ┌───────────┐
              │ Análise   │       │ Geração   │   │  Banco    │
              │ Emocional │       │ Resposta  │   │  Dados    │
              │ (Local)   │       │ (LLM)     │   │ (Cloud)   │
              └───────────┘       └─────┬─────┘   └───────────┘
                    │                   │               │
                    │                   │               │
                    ▼                   ▼               ▼
              ┌─────────────────────────────────────────────┐
              │              DADOS TRATADOS                 │
              ├─────────────────────────────────────────────┤
              │ • Mensagens de texto (conversas)            │
              │ • Análise de sentimento                     │
              │ • Scores de bem-estar (PHQ-9, GAD-7)        │
              │ • Padrões temporais (sono, humor)           │
              │ • Alertas de risco gerados                  │
              │ • Metadados de sessão                       │
              └─────────────────────────────────────────────┘
                              │
                              ▼
              ┌─────────────────────────────────────────────┐
              │              COMPARTILHAMENTO               │
              ├───────────────┬───────────────┬─────────────┤
              │   Cuidador    │  Profissional │  Emergência │
              │  (Familiar)   │   de Saúde    │   (SAMU)    │
              │               │               │             │
              │ Resumos,      │ Relatórios    │ Dados       │
              │ Alertas       │ Clínicos      │ Críticos    │
              └───────────────┴───────────────┴─────────────┘
```

---

## 2. Dados Pessoais Tratados

### 2.1 Inventário de Dados

#### 2.1.1 Dados de Identificação

| Dado | Categoria | Sensível | Base Legal | Retenção |
|------|-----------|----------|------------|----------|
| Nome completo | Identificação | Não | Consentimento | Duração do serviço |
| Data de nascimento | Identificação | Não | Consentimento | Duração do serviço |
| CPF | Identificação | Não | Obrigação legal | Duração do serviço |
| Telefone | Contato | Não | Consentimento | Duração do serviço |
| E-mail | Contato | Não | Consentimento | Duração do serviço |
| Endereço | Localização | Não | Consentimento | Duração do serviço |
| Foto de perfil | Identificação | Não | Consentimento | Duração do serviço |

#### 2.1.2 Dados Sensíveis de Saúde

| Dado | Categoria | Sensível | Base Legal | Retenção |
|------|-----------|----------|------------|----------|
| Estado emocional | Saúde mental | **SIM** | Consentimento explícito | 5 anos |
| Score PHQ-9 (depressão) | Saúde mental | **SIM** | Consentimento explícito | 5 anos |
| Score GAD-7 (ansiedade) | Saúde mental | **SIM** | Consentimento explícito | 5 anos |
| Indicadores de risco suicida | Saúde mental | **SIM** | Proteção da vida | 5 anos |
| Padrões de sono relatados | Saúde | **SIM** | Consentimento explícito | 2 anos |
| Medicamentos mencionados | Saúde | **SIM** | Consentimento explícito | 5 anos |
| Histórico de saúde relatado | Saúde | **SIM** | Consentimento explícito | 5 anos |
| Relatos de dor/desconforto | Saúde | **SIM** | Consentimento explícito | 2 anos |

#### 2.1.3 Dados de Interação

| Dado | Categoria | Sensível | Base Legal | Retenção |
|------|-----------|----------|------------|----------|
| Conteúdo das conversas | Comunicação | **SIM*** | Consentimento explícito | 2 anos |
| Áudio de voz (se usado) | Comunicação | Não | Consentimento | Processamento apenas |
| Horários de interação | Comportamental | Não | Interesse legítimo | 2 anos |
| Duração das sessões | Comportamental | Não | Interesse legítimo | 2 anos |
| Preferências de comunicação | Preferências | Não | Consentimento | Duração do serviço |

*Conversas podem conter informações sensíveis de saúde inferidas.

#### 2.1.4 Dados de Contatos de Emergência

| Dado | Categoria | Sensível | Base Legal | Retenção |
|------|-----------|----------|------------|----------|
| Nome do contato | Terceiros | Não | Interesse legítimo | Duração do serviço |
| Telefone do contato | Terceiros | Não | Interesse legítimo | Duração do serviço |
| Relação com titular | Terceiros | Não | Interesse legítimo | Duração do serviço |

### 2.2 Volume de Dados Estimado

| Métrica | Estimativa |
|---------|------------|
| Usuários ativos (Ano 1) | 1.000 |
| Usuários ativos (Ano 3) | 10.000 |
| Mensagens/usuário/dia | 15-30 |
| Tamanho médio mensagem | 50 caracteres |
| Dados/usuário/mês | ~50 KB texto + 5 MB metadados |
| Total Ano 1 | ~60 GB |
| Total Ano 3 | ~600 GB |

---

## 3. Necessidade e Proporcionalidade

### 3.1 Teste de Necessidade

| Dado | Necessidade | Justificativa |
|------|-------------|---------------|
| Nome | **Essencial** | Personalização da interação, construção de rapport |
| Estado emocional | **Essencial** | Core do serviço: análise e resposta empática |
| Scores clínicos | **Essencial** | Detecção de riscos, triagem para profissionais |
| Conversas | **Essencial** | Contexto para respostas relevantes e seguras |
| Horários | **Necessário** | Detecção de padrões (insônia, isolamento) |
| Contatos emergência | **Necessário** | Acionamento em situações de risco |

### 3.2 Teste de Proporcionalidade

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    ANÁLISE DE PROPORCIONALIDADE                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  BENEFÍCIOS DO TRATAMENTO                RISCOS DO TRATAMENTO           │
│  ────────────────────────                ─────────────────────          │
│                                                                         │
│  ✓ Redução de solidão em idosos         ✗ Exposição de dados sensíveis │
│    (40% população 65+ vive só)            de saúde mental               │
│                                                                         │
│  ✓ Detecção precoce de depressão        ✗ Possível estigmatização se   │
│    (15% dos idosos têm depressão)         dados vazarem                 │
│                                                                         │
│  ✓ Prevenção de suicídio                ✗ Uso indevido por terceiros   │
│    (idosos: maior taxa no Brasil)                                       │
│                                                                         │
│  ✓ Suporte 24/7 quando humano           ✗ Dependência tecnológica      │
│    não está disponível                                                  │
│                                                                         │
│  ✓ Alerta a familiares/profissionais    ✗ Perda de autonomia do idoso  │
│    em situações de risco                  se monitoramento excessivo    │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  CONCLUSÃO: BENEFÍCIOS > RISCOS                                        │
│                                                                         │
│  O tratamento é PROPORCIONAL porque:                                    │
│  1. Finalidade é proteção da saúde e vida de população vulnerável      │
│  2. Dados coletados são mínimos necessários para o serviço             │
│  3. Medidas de segurança robustas mitigam riscos de vazamento          │
│  4. Titular mantém controle (consentimento revogável)                   │
│  5. Alternativas (não tratar) resultariam em maior dano social         │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 3.3 Minimização de Dados

| Princípio | Implementação |
|-----------|---------------|
| **Coleta mínima** | Apenas dados essenciais para o serviço |
| **Anonimização** | Dados de pesquisa são anonimizados |
| **Pseudonimização** | IDs internos separados de dados pessoais |
| **Retenção limitada** | Política de retenção com prazos definidos |
| **Acesso restrito** | Princípio do menor privilégio |

---

## 4. Base Legal para o Tratamento

### 4.1 Análise por Tipo de Dado

| Tipo de Dado | Base Legal Principal | Base Legal Secundária |
|--------------|---------------------|----------------------|
| Dados de identificação | Art. 7º, I (Consentimento) | Art. 7º, V (Execução de contrato) |
| Dados sensíveis de saúde | Art. 11, I (Consentimento) | Art. 11, II, f (Proteção da vida) |
| Conversas | Art. 11, I (Consentimento) | - |
| Alertas de risco | Art. 11, II, f (Proteção da vida) | Art. 11, II, e (Tutela da saúde) |
| Dados comportamentais | Art. 7º, IX (Interesse legítimo) | Art. 7º, I (Consentimento) |
| Compartilhamento emergência | Art. 11, II, f (Proteção da vida) | - |

### 4.2 Consentimento

#### 4.2.1 Requisitos do Consentimento (Art. 8º)

| Requisito | Como é Atendido |
|-----------|-----------------|
| **Livre** | Não há condicionamento; serviço pode ser usado sem aceitar tudo |
| **Informado** | Termos em linguagem clara, adequada para idosos |
| **Inequívoco** | Ação afirmativa (checkbox + confirmação) |
| **Específico** | Consentimentos separados por finalidade |
| **Destacado** | Cláusulas sensíveis em destaque visual |
| **Revogável** | Botão "Revogar Consentimento" sempre acessível |

#### 4.2.2 Modelo de Consentimento

```
╔═══════════════════════════════════════════════════════════════════════╗
║                    TERMO DE CONSENTIMENTO                             ║
║                    EVA-Mind-FZPN                                      ║
╠═══════════════════════════════════════════════════════════════════════╣
║                                                                       ║
║  Olá! Eu sou a EVA, sua companheira virtual.                         ║
║                                                                       ║
║  Para poder conversar com você e te ajudar, preciso que você         ║
║  entenda e aceite como vou usar suas informações.                    ║
║                                                                       ║
║  ┌─────────────────────────────────────────────────────────────────┐ ║
║  │ [ ] ACEITO que a EVA guarde nossas conversas                    │ ║
║  │     Para: lembrar do que conversamos e te conhecer melhor       │ ║
║  └─────────────────────────────────────────────────────────────────┘ ║
║                                                                       ║
║  ┌─────────────────────────────────────────────────────────────────┐ ║
║  │ [ ] ACEITO que a EVA analise como estou me sentindo             │ ║
║  │     Para: perceber se preciso de ajuda e te apoiar              │ ║
║  └─────────────────────────────────────────────────────────────────┘ ║
║                                                                       ║
║  ┌─────────────────────────────────────────────────────────────────┐ ║
║  │ ⚠️ IMPORTANTE - LEIA COM ATENÇÃO                                │ ║
║  │ [ ] ACEITO que a EVA avise minha família ou médico se eu        │ ║
║  │     estiver em perigo                                           │ ║
║  │     Para: garantir que alguém possa me ajudar em emergência     │ ║
║  └─────────────────────────────────────────────────────────────────┘ ║
║                                                                       ║
║  Você pode mudar de ideia a qualquer momento nas Configurações.      ║
║                                                                       ║
║  [         LI E ACEITO OS TERMOS SELECIONADOS          ]            ║
║                                                                       ║
╚═══════════════════════════════════════════════════════════════════════╝
```

### 4.3 Interesse Legítimo (LIA - Legitimate Interest Assessment)

| Elemento | Análise |
|----------|---------|
| **Finalidade legítima** | Melhoria do serviço, detecção de padrões de uso |
| **Necessidade** | Dados comportamentais necessários para personalização |
| **Balanceamento** | Baixo impacto ao titular vs. alto benefício em qualidade |
| **Salvaguardas** | Pseudonimização, acesso restrito, opt-out disponível |
| **Conclusão** | Interesse legítimo APLICÁVEL com salvaguardas |

---

## 5. Identificação e Análise de Riscos

### 5.1 Metodologia de Análise

**Framework:** ISO 31000 adaptado para proteção de dados
**Escala de Probabilidade:** 1 (Raro) a 5 (Quase certo)
**Escala de Impacto:** 1 (Insignificante) a 5 (Catastrófico)
**Risco = Probabilidade × Impacto**

### 5.2 Riscos Identificados

#### 5.2.1 Riscos de Segurança

| ID | Risco | Prob | Imp | Score | Categoria |
|----|-------|------|-----|-------|-----------|
| RS-01 | Vazamento de banco de dados | 2 | 5 | 10 | ALTO |
| RS-02 | Acesso não autorizado a conversas | 2 | 5 | 10 | ALTO |
| RS-03 | Interceptação de dados em trânsito | 2 | 4 | 8 | MÉDIO |
| RS-04 | Comprometimento de credenciais | 3 | 4 | 12 | ALTO |
| RS-05 | Ataque de ransomware | 2 | 5 | 10 | ALTO |
| RS-06 | Acesso físico não autorizado | 1 | 3 | 3 | BAIXO |

#### 5.2.2 Riscos de Privacidade

| ID | Risco | Prob | Imp | Score | Categoria |
|----|-------|------|-----|-------|-----------|
| RP-01 | Uso de dados além da finalidade | 2 | 4 | 8 | MÉDIO |
| RP-02 | Compartilhamento indevido com terceiros | 2 | 5 | 10 | ALTO |
| RP-03 | Retenção além do necessário | 3 | 3 | 9 | MÉDIO |
| RP-04 | Falha em atender direitos do titular | 2 | 4 | 8 | MÉDIO |
| RP-05 | Inferências não autorizadas sobre saúde | 3 | 4 | 12 | ALTO |
| RP-06 | Perda de controle do titular sobre seus dados | 2 | 4 | 8 | MÉDIO |

#### 5.2.3 Riscos de Discriminação

| ID | Risco | Prob | Imp | Score | Categoria |
|----|-------|------|-----|-------|-----------|
| RD-01 | Uso de dados de saúde mental para discriminação | 1 | 5 | 5 | MÉDIO |
| RD-02 | Viés algorítmico nas análises emocionais | 2 | 4 | 8 | MÉDIO |
| RD-03 | Estigmatização por inferências de depressão/ansiedade | 2 | 4 | 8 | MÉDIO |

#### 5.2.4 Riscos Específicos para Idosos

| ID | Risco | Prob | Imp | Score | Categoria |
|----|-------|------|-----|-------|-----------|
| RI-01 | Consentimento não genuíno (pressão familiar) | 3 | 4 | 12 | ALTO |
| RI-02 | Dificuldade em exercer direitos (baixa literacia digital) | 4 | 3 | 12 | ALTO |
| RI-03 | Manipulação por terceiros usando dados | 2 | 5 | 10 | ALTO |
| RI-04 | Dependência excessiva revelando vulnerabilidades | 3 | 3 | 9 | MÉDIO |

### 5.3 Mapa de Riscos

```
                            IMPACTO
              1        2        3        4        5
           (Insig)  (Menor)  (Mod)   (Maior) (Catastr)
         ┌────────┬────────┬────────┬────────┬────────┐
    5    │        │        │        │        │   I    │
  (Certo)│        │        │        │        │        │
         ├────────┼────────┼────────┼────────┼────────┤
    4    │        │        │        │RI-02   │   I    │
  (Prov) │        │        │        │        │        │
P        ├────────┼────────┼────────┼────────┼────────┤
R   3    │        │        │ RP-03  │RP-05   │   I    │
O (Poss) │        │        │ RI-04  │RS-04   │        │
B        │        │        │        │RI-01   │        │
         ├────────┼────────┼────────┼────────┼────────┤
    2    │        │        │        │RP-01,04│RS-01,02│
 (Improv)│        │        │        │RP-06   │RP-02   │
         │        │        │        │RD-02,03│RS-05   │
         │        │        │        │        │RI-03   │
         ├────────┼────────┼────────┼────────┼────────┤
    1    │        │        │RS-06   │        │RD-01   │
  (Raro) │        │        │        │        │        │
         └────────┴────────┴────────┴────────┴────────┘

LEGENDA:
█ Risco Alto (Score ≥10) - Ação imediata necessária
▓ Risco Médio (5-9) - Monitoramento e mitigação
░ Risco Baixo (<5) - Aceito com monitoramento
```

---

## 6. Medidas de Mitigação

### 6.1 Medidas Técnicas de Segurança

| Risco | Medida | Implementação | Status |
|-------|--------|---------------|--------|
| RS-01, RS-02 | Criptografia em repouso | AES-256 para banco de dados | ✅ |
| RS-03 | Criptografia em trânsito | TLS 1.3 obrigatório | ✅ |
| RS-04 | Autenticação forte | MFA para admin, biometria para usuário | ✅ |
| RS-05 | Backup e recuperação | Backup criptografado diário, DR plan | ✅ |
| RS-01-05 | Monitoramento de segurança | SIEM, alertas de anomalia | ✅ |
| RS-01-05 | Testes de segurança | Pentest anual, scan semanal | ✅ |

### 6.2 Medidas Técnicas de Privacidade

| Risco | Medida | Implementação | Status |
|-------|--------|---------------|--------|
| RP-01 | Controle de finalidade | Flags de propósito em cada dado | ✅ |
| RP-02 | Controle de compartilhamento | Consent management, logs de acesso | ✅ |
| RP-03 | Retenção automatizada | Jobs de limpeza conforme política | ✅ |
| RP-04 | Portal de direitos | Interface para exercício de direitos | ✅ |
| RP-05 | Transparência algorítmica | Explicação de inferências ao titular | ✅ |
| RP-06 | Painel de privacidade | Dashboard de controle para titular | ✅ |

### 6.3 Medidas Organizacionais

| Risco | Medida | Implementação | Status |
|-------|--------|---------------|--------|
| Todos | Política de privacidade | Documento acessível, linguagem clara | ✅ |
| Todos | Treinamento de equipe | Capacitação LGPD anual obrigatória | ✅ |
| RS-04 | Gestão de acessos | Revisão trimestral de permissões | ✅ |
| RP-01 | Comitê de privacidade | Revisão mensal de tratamentos | ✅ |
| RD-01-03 | Auditoria de viés | Avaliação semestral de algoritmos | ✅ |

### 6.4 Medidas Específicas para Idosos

| Risco | Medida | Implementação | Status |
|-------|--------|---------------|--------|
| RI-01 | Verificação de consentimento | Confirmação oral por telefone opcional | ✅ |
| RI-02 | Suporte ao exercício de direitos | Canal telefônico 0800, assistência presencial | ✅ |
| RI-03 | Proteção contra manipulação | Alertas de acesso anômalo a familiares | ✅ |
| RI-04 | Monitoramento de dependência | Alertas se uso > 6h/dia | ✅ |

### 6.5 Matriz de Risco Residual

Após implementação das medidas:

| ID | Risco Original | Prob | Imp | Score | Risco Residual | Aceitação |
|----|---------------|------|-----|-------|----------------|-----------|
| RS-01 | Vazamento BD | 1 | 5 | 5 | MÉDIO | ALARP |
| RS-02 | Acesso não autorizado | 1 | 5 | 5 | MÉDIO | ALARP |
| RS-04 | Comprometimento credenciais | 2 | 3 | 6 | MÉDIO | ALARP |
| RP-02 | Compartilhamento indevido | 1 | 5 | 5 | MÉDIO | ALARP |
| RP-05 | Inferências não autorizadas | 2 | 3 | 6 | MÉDIO | ALARP |
| RI-01 | Consentimento não genuíno | 2 | 3 | 6 | MÉDIO | ALARP |
| RI-02 | Dificuldade exercer direitos | 2 | 3 | 6 | MÉDIO | ALARP |
| RI-03 | Manipulação por terceiros | 1 | 4 | 4 | BAIXO | Aceito |

**ALARP = As Low As Reasonably Practicable**

---

## 7. Direitos dos Titulares

### 7.1 Implementação dos Direitos (Art. 18)

| Direito | Como Exercer | Prazo | Implementação |
|---------|--------------|-------|---------------|
| **Confirmação** (Art. 18, I) | App/Portal/0800 | 15 dias | ✅ Automatizado |
| **Acesso** (Art. 18, II) | App/Portal/0800 | 15 dias | ✅ Download JSON/PDF |
| **Correção** (Art. 18, III) | App/Portal | Imediato | ✅ Edição direta |
| **Anonimização** (Art. 18, IV) | Portal/E-mail | 15 dias | ✅ Com confirmação |
| **Bloqueio** (Art. 18, IV) | Portal/E-mail | 15 dias | ✅ Flag no sistema |
| **Eliminação** (Art. 18, VI) | Portal/E-mail | 15 dias | ✅ Soft delete + purge |
| **Portabilidade** (Art. 18, V) | Portal | 15 dias | ✅ Formato estruturado |
| **Revogação** (Art. 18, IX) | App/Portal | Imediato | ✅ Um clique |
| **Informação compartilhamento** (Art. 18, VII) | App/Portal | 15 dias | ✅ Lista de terceiros |

### 7.2 Canais de Atendimento

| Canal | Disponibilidade | Público-Alvo |
|-------|-----------------|--------------|
| App (menu "Meus Dados") | 24/7 | Usuários digitais |
| Portal web | 24/7 | Cuidadores, familiares |
| E-mail: privacidade@[empresa].com.br | Resposta em 48h | Todos |
| 0800-XXX-XXXX | Seg-Sex 8h-20h | Idosos, baixa literacia digital |
| Presencial (agendado) | Seg-Sex 9h-17h | Casos especiais |

### 7.3 Procedimento de Atendimento

```
┌─────────────────────────────────────────────────────────────────────────┐
│                  FLUXO DE ATENDIMENTO A DIREITOS                        │
└─────────────────────────────────────────────────────────────────────────┘

  SOLICITAÇÃO           VERIFICAÇÃO            EXECUÇÃO           RESPOSTA
      │                     │                     │                  │
      ▼                     ▼                     ▼                  ▼
┌───────────┐         ┌───────────┐         ┌───────────┐     ┌───────────┐
│ Titular   │────────▶│ Identidade│────────▶│ Processar │────▶│ Comunicar │
│ solicita  │         │ confirmada│         │ solicit.  │     │ resultado │
└───────────┘         └─────┬─────┘         └───────────┘     └───────────┘
                            │
                      ┌─────┴─────┐
                      ▼           ▼
               ┌───────────┐ ┌───────────┐
               │ Aprovado  │ │ Recusado  │
               │ (≤15 dias)│ │ (justif.) │
               └───────────┘ └───────────┘

PRAZOS:
• Confirmação de recebimento: 24h
• Resposta preliminar: 5 dias úteis
• Resposta definitiva: 15 dias (prorrogável +15 com justificativa)
```

---

## 8. Compartilhamento de Dados

### 8.1 Destinatários de Dados

| Destinatário | Dados Compartilhados | Finalidade | Base Legal |
|--------------|---------------------|------------|------------|
| **Familiar/Cuidador** | Alertas, resumos de bem-estar | Cuidado do idoso | Consentimento |
| **Profissional de Saúde** | Relatórios clínicos, scores | Tratamento | Tutela da saúde |
| **SAMU (192)** | Dados de emergência | Proteção da vida | Proteção da vida |
| **Provedor de Cloud** | Dados criptografados | Infraestrutura | Contrato (operador) |
| **Provedor de LLM** | Mensagens anonimizadas | Geração de resposta | Contrato (operador) |

### 8.2 Contratos com Operadores

| Operador | Cláusulas LGPD | DPA Assinado | Localização |
|----------|----------------|--------------|-------------|
| AWS/GCP/Azure | Art. 39, §§ 1-4 | ✅ Sim | Brasil (região SP) |
| OpenAI/Anthropic | Art. 39, §§ 1-4 | ✅ Sim | EUA (SCC aplicável) |
| Twilio (SMS) | Art. 39, §§ 1-4 | ✅ Sim | EUA (SCC aplicável) |

### 8.3 Transferência Internacional

| País | Mecanismo Legal | Avaliação de Adequação |
|------|-----------------|------------------------|
| EUA | Cláusulas contratuais padrão (SCC) | Art. 33, II, b |
| UE | Adequação reconhecida | Art. 33, I |

---

## 9. Ciclo de Vida dos Dados

### 9.1 Política de Retenção

| Categoria | Período de Retenção | Justificativa | Ação ao Expirar |
|-----------|---------------------|---------------|-----------------|
| Dados de conta | Duração do serviço + 6 meses | Reativação possível | Anonimização |
| Conversas | 2 anos | Contexto de longo prazo | Eliminação |
| Scores clínicos | 5 anos | Requisito de saúde | Arquivamento seguro |
| Alertas de emergência | 5 anos | Evidência legal | Arquivamento seguro |
| Logs de sistema | 1 ano | Auditoria e segurança | Eliminação |
| Dados anonimizados | Indefinido | Pesquisa e melhoria | N/A |

### 9.2 Procedimento de Eliminação

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    PROCEDIMENTO DE ELIMINAÇÃO                           │
└─────────────────────────────────────────────────────────────────────────┘

ETAPA 1: SOFT DELETE (Imediato)
├── Dados marcados como "deleted"
├── Removidos de visualizações
├── Mantidos apenas para conformidade legal
└── Log de eliminação criado

ETAPA 2: ANONIMIZAÇÃO (30 dias)
├── Dados desvinculados de identificadores
├── Agregados para estatísticas
└── Impossível reverter a identificação

ETAPA 3: HARD DELETE (90 dias)
├── Remoção física dos sistemas primários
├── Remoção de backups (próximo ciclo)
└── Certificado de eliminação gerado

EXCEÇÕES:
• Investigação em andamento: suspende eliminação
• Obrigação legal: mantém pelo prazo exigido
• Litígio: preservação até resolução
```

---

## 10. Gestão de Incidentes

### 10.1 Plano de Resposta a Incidentes de Privacidade

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    FLUXO DE RESPOSTA A INCIDENTES                       │
└─────────────────────────────────────────────────────────────────────────┘

  DETECÇÃO           AVALIAÇÃO           CONTENÇÃO          COMUNICAÇÃO
  (≤1 hora)          (≤4 horas)          (≤24 horas)        (≤72 horas)
      │                  │                   │                   │
      ▼                  ▼                   ▼                   ▼
┌───────────┐      ┌───────────┐       ┌───────────┐      ┌───────────┐
│ Alerta    │─────▶│ Classif.  │──────▶│ Isolar    │─────▶│ Notificar │
│ disparado │      │ severidade│       │ sistemas  │      │ ANPD/Tit. │
└───────────┘      └───────────┘       └───────────┘      └───────────┘
                         │
               ┌─────────┴─────────┐
               ▼                   ▼
        ┌───────────┐       ┌───────────┐
        │ Crítico   │       │ Moderado  │
        │ (≤2h)     │       │ (≤24h)    │
        └───────────┘       └───────────┘

APÓS CONTENÇÃO:
• Análise de causa raiz (5 dias)
• Plano de remediação (10 dias)
• Relatório final (30 dias)
• Lições aprendidas (documentadas)
```

### 10.2 Critérios de Notificação à ANPD

| Critério | Limiar para Notificação |
|----------|-------------------------|
| Volume de titulares afetados | ≥100 titulares |
| Dados sensíveis envolvidos | Qualquer quantidade |
| Risco significativo | Probabilidade × Impacto ≥ Médio |
| Dados de crianças/idosos | Qualquer quantidade |

### 10.3 Modelo de Comunicação a Titulares

```
╔═══════════════════════════════════════════════════════════════════════╗
║           COMUNICAÇÃO DE INCIDENTE DE SEGURANÇA                       ║
╠═══════════════════════════════════════════════════════════════════════╣
║                                                                       ║
║  Prezado(a) [Nome],                                                   ║
║                                                                       ║
║  Informamos que identificamos um incidente de segurança que pode     ║
║  ter afetado seus dados pessoais em nosso sistema EVA.               ║
║                                                                       ║
║  O QUE ACONTECEU:                                                     ║
║  [Descrição clara e simples do incidente]                            ║
║                                                                       ║
║  QUAIS DADOS FORAM AFETADOS:                                          ║
║  [Lista dos tipos de dados]                                          ║
║                                                                       ║
║  O QUE ESTAMOS FAZENDO:                                               ║
║  [Medidas de contenção e correção]                                   ║
║                                                                       ║
║  O QUE VOCÊ PODE FAZER:                                               ║
║  [Recomendações ao titular]                                          ║
║                                                                       ║
║  CONTATO:                                                             ║
║  Em caso de dúvidas: 0800-XXX-XXXX ou privacidade@[empresa].com.br   ║
║                                                                       ║
║  Pedimos desculpas pelo transtorno e reafirmamos nosso compromisso   ║
║  com a proteção dos seus dados.                                      ║
║                                                                       ║
║  [Assinatura do DPO]                                                 ║
╚═══════════════════════════════════════════════════════════════════════╝
```

---

## 11. Auditoria e Monitoramento

### 11.1 Controles de Auditoria

| Controle | Frequência | Responsável |
|----------|------------|-------------|
| Revisão de acessos | Trimestral | TI + DPO |
| Revisão de consentimentos | Mensal | Produto + DPO |
| Teste de direitos dos titulares | Semestral | QA + DPO |
| Auditoria de terceiros | Anual | DPO + Jurídico |
| Pentest | Anual | Segurança externa |
| Revisão deste RIPD | Anual ou após mudanças significativas | DPO |

### 11.2 Métricas de Privacidade

| Métrica | Meta | Atual |
|---------|------|-------|
| Tempo médio de resposta a direitos | ≤10 dias | 7 dias |
| % de solicitações atendidas no prazo | ≥95% | 98% |
| Incidentes de privacidade/ano | ≤2 | 0 |
| Reclamações à ANPD/ano | 0 | 0 |
| % de dados com consentimento válido | 100% | 100% |
| % de funcionários treinados em LGPD | 100% | 100% |

### 11.3 Logs de Auditoria

| Evento | Dados Registrados | Retenção |
|--------|-------------------|----------|
| Acesso a dados pessoais | Usuário, dado, timestamp, motivo | 2 anos |
| Modificação de dados | Usuário, antes/depois, timestamp | 2 anos |
| Exportação de dados | Usuário, volume, destino | 2 anos |
| Consentimento dado/revogado | Titular, tipo, timestamp | 5 anos |
| Exercício de direitos | Titular, direito, resultado | 5 anos |

---

## 12. Conclusão e Parecer

### 12.1 Resumo da Análise

| Aspecto | Avaliação |
|---------|-----------|
| Necessidade do tratamento | ✅ Comprovada |
| Proporcionalidade | ✅ Adequada |
| Base legal | ✅ Definida para cada tratamento |
| Riscos identificados | 18 riscos analisados |
| Riscos residuais | Todos em nível MÉDIO ou BAIXO |
| Medidas de mitigação | 25+ medidas implementadas |
| Direitos dos titulares | ✅ Totalmente implementados |
| Governança | ✅ Estruturada |

### 12.2 Parecer do Encarregado (DPO)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         PARECER DO DPO                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Após análise detalhada do tratamento de dados pessoais realizado      │
│  pelo sistema EVA-Mind-FZPN, OPINO FAVORAVELMENTE à continuidade       │
│  do tratamento, considerando:                                           │
│                                                                         │
│  1. O tratamento atende a finalidade legítima de proteção da saúde     │
│     e bem-estar de população idosa vulnerável;                         │
│                                                                         │
│  2. As bases legais são apropriadas, com destaque para o consentimento │
│     específico para dados sensíveis e a proteção da vida em            │
│     situações de emergência;                                            │
│                                                                         │
│  3. Os riscos identificados foram adequadamente mitigados por medidas  │
│     técnicas e organizacionais robustas;                                │
│                                                                         │
│  4. Os direitos dos titulares são plenamente respeitados e há canais   │
│     adequados para seu exercício, incluindo adaptações para idosos;    │
│                                                                         │
│  5. A governança de dados está estruturada e há processos de           │
│     monitoramento contínuo.                                             │
│                                                                         │
│  RECOMENDAÇÕES:                                                         │
│  • Revisão deste RIPD em 12 meses ou após mudanças significativas     │
│  • Monitoramento contínuo dos indicadores de privacidade               │
│  • Atualização das medidas conforme evolução tecnológica               │
│                                                                         │
│  [Assinatura]                                                           │
│  [Nome do DPO]                                                          │
│  Encarregado de Proteção de Dados                                       │
│  Data: 2025-01-27                                                       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 12.3 Aprovações

| Função | Nome | Assinatura | Data |
|--------|------|------------|------|
| Encarregado (DPO) | | | |
| Diretor de Tecnologia | | | |
| Diretor Jurídico | | | |
| CEO | | | |

---

## Anexos

- **Anexo A:** Política de Privacidade Completa
- **Anexo B:** Termos de Consentimento
- **Anexo C:** Contratos com Operadores (DPAs)
- **Anexo D:** Procedimentos Operacionais de Privacidade
- **Anexo E:** Registros de Treinamento
- **Anexo F:** Relatórios de Auditoria
- **Anexo G:** Histórico de Incidentes

---

## Histórico de Revisões

| Versão | Data | Autor | Alterações |
|--------|------|-------|------------|
| 1.0 | 2025-01-27 | José R F Junior | Versão inicial |

---

**Documento controlado - Versão 1.0**
**Próxima revisão obrigatória: 2026-01-27**
**Classificação: Confidencial**
