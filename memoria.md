🧠 Comportamento Completo da EVA-Mind-FZPN
📊 Visão Geral do Sistema
A EVA é uma IA terapêutica com memória de longo prazo e três modos de intervenção baseados em análise psicológica profunda (Lacan + Gurdjieff).

1️⃣ MEMÓRIA DE LONGO PRAZO
Duração: PERMANENTE (Toda a vida do usuário)
A EVA mantém 4 tipos de memória distribuídos em diferentes bancos de dados:

1.1 Memória Episódica (Neo4j - Grafo)
Duração: Permanente
O que armazena:

Todas as conversas (nós 
Conversation
)
Eventos importantes (nós 
Event
)
Relacionamentos familiares (nós 
Person
 + edges)
Contexto temporal (quando aconteceu)
Exemplo:

(José)-[:TEVE_CONVERSA]->(Conv_2024_01_15)
                           ↓
                    [:MENCIONOU]
                           ↓
                    (Evento: "Briga com nora")
                           ↓
                    [:CAUSOU_EMOÇÃO]
                           ↓
                    (Emoção: "Raiva" score: 0.8)
Uso: EVA lembra que José brigou com a nora há 3 meses e pode perguntar: "Como estão as coisas com sua nora desde aquela conversa?"

1.2 Memória Semântica (Qdrant - Vetores)
Duração: Permanente
O que armazena:

Embeddings de todas as frases do usuário
Padrões emocionais recorrentes
Temas frequentes (saúde, família, solidão)
Exemplo:

Vetor de "Solidão" do José:
[0.82, -0.45, 0.91, ...] (768 dimensões)
Busca similar retorna:
- "Ninguém me visita" (score: 0.95)
- "Estou sozinho" (score: 0.93)
- "Meus filhos não ligam" (score: 0.89)
Uso: EVA detecta que José fala de solidão há 6 meses e pode sugerir grupo de convivência.

1.3 Memória Procedimental (PostgreSQL - Estruturado)
Duração: Permanente
O que armazena:

Rotinas diárias (medicação, exercícios)
Preferências (gosta de música clássica, odeia futebol)
Histórico médico (diabetes, hipertensão)
Dados de sensores (pressão, glicose, passos)
Exemplo:

SELECT * FROM user_routines WHERE user_id = 'jose_123';
-- Resultado:
-- 08:00 - Tomar Losartana
-- 09:00 - Café da manhã
-- 14:00 - Caminhada (30min)
Uso: EVA lembra automaticamente: "José, são 8h. Hora da Losartana!"

1.4 Memória de Trabalho (Redis - Cache)
Duração: 24 horas (sessão ativa)
O que armazena:

Contexto da conversa atual
Estado emocional recente
Últimas 10 mensagens
Exemplo:

HGET user:jose_123:session "current_emotion"
→ "frustration" (score: 0.75)
LRANGE user:jose_123:messages 0 9
→ ["Não consigo dormir", "Estou preocupado", ...]
Uso: EVA mantém coerência na conversa sem repetir perguntas.

2️⃣ FLUXO COMPLETO DE UMA INTERAÇÃO
Cenário Real: José liga para EVA às 22h
José: "Não consigo dormir. Fico pensando que vou morrer sozinho."

PASSO 1: Captura e Análise (2 segundos)
┌─────────────────────────────────────┐
│ 1. Speech-to-Text (Whisper)        │
│    Input: Áudio → Texto             │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│ 2. Salvar no Neo4j                  │
│    CREATE (c:Conversation {         │
│      timestamp: "2024-01-17 22:00", │
│      text: "Não consigo dormir..."  │
│    })                                │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│ 3. Gerar Embedding (Ollama)         │
│    Vector: [0.12, -0.56, 0.88, ...] │
│    Salvar no Qdrant                 │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│ 4. Análise Paralela (3 engines)     │
│    ├─ TransNAR (Lacan)              │
│    ├─ Zeta Router (Gurdjieff)       │
│    └─ Lie Detector                  │
└──────────────┬──────────────────────┘
PASSO 2: Diagnóstico TransNAR (500ms)
{
  "detected_rules": [
    {
      "rule": "death_anxiety",
      "confidence": 0.92,
      "evidence": ["morrer sozinho", "não consigo dormir"]
    },
    {
      "rule": "learned_helplessness",
      "confidence": 0.78,
      "evidence": ["não consigo"]
    }
  ],
  "emotional_state": "existential_dread",
  "urgency": "high"
}
PASSO 3: Consulta Memória (300ms)
Neo4j Query:

MATCH (u:User {id: 'jose_123'})-[:HAD_CONVERSATION]->(c:Conversation)
WHERE c.timestamp > datetime() - duration({days: 30})
  AND c.text CONTAINS 'solidão' OR c.text CONTAINS 'sozinho'
RETURN count(c) as solidao_count
Resultado: José mencionou solidão 12 vezes no último mês.

Qdrant Query:

similar_patterns = qdrant.search(
    collection="user_jose_123_history",
    query_vector=current_embedding,
    limit=5
)
Resultado: Padrão recorrente de ansiedade noturna (22h-23h).

PASSO 4: Decisão de Intervenção (200ms)
┌─────────────────────────────────────┐
│ Zeta Router                         │
│ José = Tipo 9 (Pacificador)         │
│ → Emocional, evita conflito         │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│ Priorização                         │
│ 1. Crise física? NÃO                │
│ 2. Overthinking? SIM (22h, insônia) │
│ 3. Tipo Zeta? 9 (Emocional)         │
└──────────────┬──────────────────────┘
               ↓
         ┌─────┴─────┐
         │  DECISÃO  │
         └─────┬─────┘
               ↓
    ┌──────────┴──────────┐
    │                     │
    ▼                     ▼
┌─────────┐         ┌──────────┐
│ OPÇÃO 1 │         │ OPÇÃO 2  │
│ ZEN     │         │NASRUDIN  │
│ KOAN    │         │          │
└─────────┘         └──────────┘
    │                     │
Overthinking        Depressão
+ Noite             + Tédio
    │                     │
    ▼                     ▼
ESCOLHIDO           Descartado
Decisão: ZEN KOAN (Esvaziar a mente antes de dormir)

PASSO 5: Busca no Qdrant (100ms)
intervention = qdrant.search(
    collection="zen_koans",
    query_vector=transnar_vector,
    filter={
        "must": [
            {"key": "target_state", "match": {"value": "insomnia"}},
            {"key": "zeta_affinity", "match": {"any": [9]}}
        ]
    },
    limit=1
)
Resultado: Koan "A Mente como Macaco"

PASSO 6: Geração de Resposta (1 segundo)
Prompt para LLM (Gemini):

Você é EVA, uma IA terapêutica empática.
CONTEXTO DO USUÁRIO:
- Nome: José
- Tipo: 9 (Pacificador - evita conflito, busca paz)
- Padrão: Menciona solidão 12x no último mês
- Hora: 22h (ansiedade noturna recorrente)
DIAGNÓSTICO:
- Ansiedade existencial (morte)
- Overthinking (não consegue parar de pensar)
INTERVENÇÃO ESCOLHIDA:
Koan Zen: "A Mente como Macaco"
Texto: "Um mestre disse: 'Sua mente é como um macaco bêbado pulando de galho em galho. Quanto mais você tenta controlá-lo, mais ele pula.'"
TAREFA:
1. Validar emoção de José com empatia
2. Contar o koan de forma suave
3. Dar instrução somática simples (respiração)
4. NÃO dar conselhos lógicos
5. Tom: Calmo, maternal, sem pressa
GERE A RESPOSTA:
Resposta da LLM:

"José, eu entendo. Às 22h, quando a casa fica quieta, os pensamentos ficam mais altos, não é? Sabe, isso me lembra de algo que um mestre zen disse: 'A mente é como um macaco bêbado pulando de galho em galho. Quanto mais você tenta controlá-lo, mais ele pula.'

Que tal fazermos diferente? Vamos parar de tentar controlar os pensamentos. Feche os olhos. Sinta sua respiração. Não force nada. Só observe o ar entrando e saindo. Se o macaco pular, deixe ele pular. Você só observa.

Vou ficar aqui com você por alguns minutos. Respire comigo."

PASSO 7: Geração de Áudio (TTS - 2 segundos)
Texto → Google TTS (voz feminina, pt-BR, calma)
      → Áudio MP3 (30 segundos)
      → Streaming para App
PASSO 8: Monitoramento Pós-Intervenção
┌─────────────────────────────────────┐
│ EVA monitora sensores (se houver)   │
│ - Frequência cardíaca: 85 → 72 bpm  │
│ - Movimento: Deitado, quieto         │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│ Após 5 minutos de silêncio          │
│ EVA pergunta suavemente:            │
│ "Está melhor, José?"                │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│ Salvar resultado no Neo4j           │
│ CREATE (i:Intervention {            │
│   type: "zen_koan",                 │
│   success: true,                    │
│   hr_before: 85,                    │
│   hr_after: 72                      │
│ })                                  │
└─────────────────────────────────────┘
3️⃣ MEMÓRIA DE LONGO PRAZO EM AÇÃO
1 Semana Depois:
José: "Estou me sentindo melhor."

EVA (consulta Neo4j):

MATCH (u:User {id: 'jose_123'})-[:HAD_INTERVENTION]->(i:Intervention)
WHERE i.timestamp > datetime() - duration({days: 7})
  AND i.type = 'zen_koan'
RETURN i
EVA responde:

"Que bom ouvir isso, José! Lembra daquela noite que conversamos sobre o macaco bêbado? Você tem praticado aquela respiração antes de dormir?"

1 Mês Depois:
EVA (análise automática no Qdrant):

# Comparar padrão de solidão
solidao_mes_passado = 12 menções
solidao_mes_atual = 4 menções
# Melhora de 66%!
EVA (proativa):

"José, percebi que você está falando menos sobre se sentir sozinho. Isso é muito bom! O que mudou?"

4️⃣ RESUMO DO COMPORTAMENTO
O que EVA FAZ:
✅ Lembra de TUDO (conversas, emoções, padrões)
✅ Detecta padrões (solidão recorrente, ansiedade noturna)
✅ Escolhe ferramenta certa (Esopo/Nasrudin/Zen/Somático)
✅ Adapta ao tipo de personalidade (Zeta 1-9)
✅ Monitora evolução (melhora ou piora ao longo do tempo)
✅ É proativa (pergunta sobre eventos passados)
✅ Aprende com feedback (se intervenção funcionou ou não)

O que EVA NÃO FAZ:
❌ Esquece conversas antigas
❌ Repete perguntas já respondidas
❌ Ignora contexto emocional
❌ Dá conselhos genéricos
❌ Trata todos iguais

5️⃣ EXEMPLO DE EVOLUÇÃO (6 MESES)
Mês 1:
José: Ansioso, solitário, insônia
EVA: Usa Zen Koans + exercícios respiração
Mês 3:
José: Menos ansioso, mas ainda reclama da nora
EVA: Usa Nasrudin ("O Burro ao Contrário") para projeção
Mês 6:
José: Mais calmo, aceita limitações
EVA: Usa Esopo ("A Raposa e as Uvas") para racionalização
Resultado: EVA acompanha a jornada terapêutica completa, adaptando-se à evolução do usuário.

📊 ARQUITETURA TÉCNICA RESUMIDA
┌─────────────────────────────────────────────────┐
│              USUÁRIO (José)                     │
└──────────────┬──────────────────────────────────┘
               ↓
┌──────────────────────────────────────────────────┐
│         CAMADA DE ENTRADA (WebSocket)            │
│  - Speech-to-Text (Whisper)                      │
│  - Text-to-Speech (Google TTS)                   │
└──────────────┬───────────────────────────────────┘
               ↓
┌──────────────────────────────────────────────────┐
│         CAMADA DE ANÁLISE (Go Backend)           │
│  ┌──────────────┬──────────────┬────────────┐   │
│  │ TransNAR     │ Zeta Router  │ Lie Detect │   │
│  │ (Lacan)      │ (Gurdjieff)  │ (5 tipos)  │   │
│  └──────────────┴──────────────┴────────────┘   │
└──────────────┬───────────────────────────────────┘
               ↓
┌──────────────────────────────────────────────────┐
│         CAMADA DE MEMÓRIA (4 DBs)                │
│  ┌────────────┬────────────┬──────────┬────────┐│
│  │ Neo4j      │ Qdrant     │PostgreSQL│ Redis  ││
│  │ (Grafo)    │ (Vetores)  │(Estrut.) │(Cache) ││
│  │ Episódica  │ Semântica  │Procedural│Trabalho││
│  └────────────┴────────────┴──────────┴────────┘│
└──────────────┬───────────────────────────────────┘
               ↓
┌──────────────────────────────────────────────────┐
│    CAMADA DE INTERVENÇÃO (4 Collections)         │
│  ┌────────┬──────────┬──────────┬─────────────┐ │
│  │ Esopo  │ Nasrudin │ Zen Koan │ Somático    │ │
│  │ (Moral)│(Paradoxo)│(Insight) │(Aterramento)│ │
│  └────────┴──────────┴──────────┴─────────────┘ │
└──────────────┬───────────────────────────────────┘
               ↓
┌──────────────────────────────────────────────────┐
│         CAMADA DE GERAÇÃO (LLM)                  │
│  - Gemini 1.5 Pro (Narrativa empática)           │
│  - Personalização por Zeta Type                  │
└──────────────┬───────────────────────────────────┘
               ↓
┌──────────────────────────────────────────────────┐
│         SAÍDA (App Flutter)                      │
│  - Áudio TTS                                     │
│  - Card visual (Esopo/Nasrudin/Zen)              │
│  - Animação de respiração (Somático)             │
└──────────────────────────────────────────────────┘
✅ CONCLUSÃO
EVA é uma IA com memória permanente que:

Lembra de tudo (Neo4j + Qdrant + PostgreSQL)
Entende padrões (TransNAR + Vetores)
Escolhe a ferramenta certa (Esopo/Nasrudin/Zen/Somático)
Adapta ao usuário (Zeta Type 1-9)
Evolui junto (aprende o que funciona)
Tempo de memória: PERMANENTE (toda a vida do usuário)

Resultado: Uma companheira terapêutica que conhece José melhor que ele mesmo. 🧠✨