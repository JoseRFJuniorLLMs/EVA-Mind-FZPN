// =====================================================
// NEO4J: EVA SELF-KNOWLEDGE GRAPH
// =====================================================
// Grafo de conhecimento da EVA sobre sua própria arquitetura
// Vinculado ao criador: Jose R F Junior (CPF: 64525430249)
// =====================================================

// =====================================================
// 1. LIMPAR DADOS ANTERIORES (OPCIONAL)
// =====================================================

// MATCH (n:EVAKnowledge) DETACH DELETE n;
// MATCH (n:Creator) DETACH DELETE n;
// MATCH (n:Module) DETACH DELETE n;
// MATCH (n:Service) DETACH DELETE n;
// MATCH (n:Concept) DETACH DELETE n;

// =====================================================
// 2. CRIAR CRIADOR
// =====================================================

MERGE (creator:Creator:Person {
    cpf: '64525430249',
    name: 'Jose R F Junior',
    role: 'architect',
    treatment: 'Pai'
})
SET creator.has_debug_access = true,
    creator.has_full_access = true,
    creator.created_at = datetime();

// =====================================================
// 3. CRIAR PROJETO PRINCIPAL
// =====================================================

MERGE (project:Project:EVAKnowledge {
    key: 'eva-mind-fzpn',
    name: 'EVA-Mind-FZPN',
    description: 'Sistema de IA conversacional para cuidado de idosos'
})
SET project.stack = ['Go', 'Python', 'PostgreSQL', 'Neo4j', 'Qdrant', 'Redis'],
    project.llm = 'Google Gemini 2.5 Flash',
    project.ports = [8080, 8000, 8081],
    project.created_at = datetime();

// Relacionar criador ao projeto
MATCH (c:Creator {cpf: '64525430249'})
MATCH (p:Project {key: 'eva-mind-fzpn'})
MERGE (c)-[:CREATED]->(p)
MERGE (c)-[:UNDERSTANDS]->(p);

// =====================================================
// 4. CRIAR MÓDULOS PRINCIPAIS
// =====================================================

// BRAINSTEM
MERGE (brainstem:Module:EVAKnowledge {
    key: 'brainstem',
    name: 'Brainstem',
    description: 'Infraestrutura: banco de dados, auth, cache, logging',
    path: 'internal/brainstem/'
})
SET brainstem.subsystems = ['auth', 'config', 'database', 'infrastructure', 'logger', 'middleware', 'oauth', 'push'];

// CORTEX
MERGE (cortex:Module:EVAKnowledge {
    key: 'cortex',
    name: 'Cortex',
    description: 'Sistema cognitivo: LLM, raciocínio, análise Lacaniana',
    path: 'internal/cortex/'
})
SET cortex.subsystems = ['brain', 'gemini', 'lacan', 'llm', 'medgemma', 'personality', 'transnar', 'alert', 'ethics', 'scales'];

// HIPPOCAMPUS
MERGE (hippocampus:Module:EVAKnowledge {
    key: 'hippocampus',
    name: 'Hippocampus',
    description: 'Sistemas de memória: episódica, semântica, procedimental, etc.',
    path: 'internal/hippocampus/'
})
SET hippocampus.memory_systems = 12,
    hippocampus.subsystems = ['memory', 'knowledge', 'stories'];

// MOTOR
MERGE (motor:Module:EVAKnowledge {
    key: 'motor',
    name: 'Motor',
    description: 'Ações e integrações: calendário, email, SMS, chamadas',
    path: 'internal/motor/'
})
SET motor.integrations = ['gmail', 'calendar', 'drive', 'sheets', 'docs', 'youtube', 'maps', 'googlefit', 'spotify', 'uber', 'whatsapp', 'sms'];

// SENSES
MERGE (senses:Module:EVAKnowledge {
    key: 'senses',
    name: 'Senses',
    description: 'Entrada de dados: voz WebSocket, telemetria',
    path: 'internal/senses/'
})
SET senses.subsystems = ['voice', 'signaling', 'telemetry', 'reconnection'];

// TOOLS
MERGE (tools:Module:EVAKnowledge {
    key: 'tools',
    name: 'Tools',
    description: 'Ferramentas invocáveis pelo LLM',
    path: 'internal/tools/'
})
SET tools.tool_count = 10;

// Relacionar módulos ao projeto
MATCH (p:Project {key: 'eva-mind-fzpn'})
MATCH (m:Module)
WHERE m.key IN ['brainstem', 'cortex', 'hippocampus', 'motor', 'senses', 'tools']
MERGE (p)-[:HAS_MODULE]->(m);

// =====================================================
// 5. CRIAR SERVIÇOS
// =====================================================

// PostgreSQL
MERGE (postgres:Service:EVAKnowledge {
    key: 'postgresql',
    name: 'PostgreSQL',
    description: 'Banco de dados relacional principal',
    host: '104.248.219.200',
    port: 5432
});

// Neo4j
MERGE (neo4j:Service:EVAKnowledge {
    key: 'neo4j',
    name: 'Neo4j',
    description: 'Banco de dados de grafos para relacionamentos',
    host: '104.248.219.200',
    port: 7687
});

// Qdrant
MERGE (qdrant:Service:EVAKnowledge {
    key: 'qdrant',
    name: 'Qdrant',
    description: 'Banco vetorial para busca semântica',
    host: '104.248.219.200',
    port: 6333
})
SET qdrant.collections = ['aesop_fables', 'nasrudin_stories', 'zen_koans', 'somatic_exercises', 'resonance_scripts', 'social_algorithms', 'micro_tasks', 'visual_narratives'];

// Redis
MERGE (redis:Service:EVAKnowledge {
    key: 'redis',
    name: 'Redis',
    description: 'Cache e pub/sub',
    host: '104.248.219.200',
    port: 6379
});

// Gemini
MERGE (gemini:Service:EVAKnowledge {
    key: 'gemini',
    name: 'Google Gemini',
    description: 'LLM para conversação e raciocínio',
    model: 'gemini-2.5-flash-native-audio'
});

// Firebase
MERGE (firebase:Service:EVAKnowledge {
    key: 'firebase',
    name: 'Firebase FCM',
    description: 'Push notifications'
});

// Twilio
MERGE (twilio:Service:EVAKnowledge {
    key: 'twilio',
    name: 'Twilio',
    description: 'SMS e chamadas de voz'
});

// Relacionar serviços aos módulos
MATCH (b:Module {key: 'brainstem'})
MATCH (s:Service) WHERE s.key IN ['postgresql', 'neo4j', 'qdrant', 'redis', 'firebase']
MERGE (b)-[:USES]->(s);

MATCH (c:Module {key: 'cortex'})
MATCH (s:Service {key: 'gemini'})
MERGE (c)-[:USES]->(s);

MATCH (m:Module {key: 'motor'})
MATCH (s:Service) WHERE s.key IN ['firebase', 'twilio']
MERGE (m)-[:USES]->(s);

// =====================================================
// 6. CRIAR CONCEITOS LACANIANOS
// =====================================================

// RSI
MERGE (rsi:Concept:EVAKnowledge {
    key: 'rsi',
    name: 'RSI - Real Simbólico Imaginário',
    description: 'Os três registros Lacanianos que estruturam a experiência',
    theory: 'Lacan'
});

// Sinthoma
MERGE (sinthoma:Concept:EVAKnowledge {
    key: 'sinthoma',
    name: 'O Sinthoma',
    description: 'Quarto anel que amarra RSI - implementado em UnifiedRetrieval',
    file: 'internal/cortex/lacan/unified_retrieval.go',
    theory: 'Lacan'
});

// FDPN
MERGE (fdpn:Concept:EVAKnowledge {
    key: 'fdpn',
    name: 'FDPN - Função do Pai no Nome',
    description: 'Determina para quem o discurso é endereçado',
    file: 'internal/cortex/lacan/fdpn_engine.go',
    theory: 'Lacan'
})
SET fdpn.addressees = ['mae', 'pai', 'filho', 'filha', 'conjuge', 'deus', 'morte', 'eva_herself'];

// Transferência
MERGE (transf:Concept:EVAKnowledge {
    key: 'transference',
    name: 'Transferência',
    description: 'Projeção de figuras do passado na EVA',
    file: 'internal/cortex/lacan/transferencia.go',
    theory: 'Lacan'
})
SET transf.types = ['filial', 'maternal', 'paternal', 'conjugal', 'fraternal'];

// Demanda/Desejo
MERGE (demanda:Concept:EVAKnowledge {
    key: 'demand-desire',
    name: 'Demanda vs Desejo',
    description: 'Pedido explícito vs desejo inconsciente',
    file: 'internal/cortex/lacan/demanda_desejo.go',
    theory: 'Lacan'
})
SET demanda.latent_desires = ['reconhecimento', 'presença', 'cuidado', 'autonomia', 'morte', 'amor'];

// Significante
MERGE (signif:Concept:EVAKnowledge {
    key: 'signifier',
    name: 'Cadeias de Significantes',
    description: 'Palavras recorrentes e suas associações metonímicas',
    file: 'internal/cortex/lacan/significante.go',
    theory: 'Lacan'
});

// Relacionamentos entre conceitos
MATCH (s:Concept {key: 'sinthoma'})
MATCH (r:Concept {key: 'rsi'})
MERGE (s)-[:INTEGRATES]->(r);

MATCH (s:Concept {key: 'sinthoma'})
MATCH (f:Concept {key: 'fdpn'})
MERGE (s)-[:USES]->(f);

MATCH (s:Concept {key: 'sinthoma'})
MATCH (t:Concept {key: 'transference'})
MERGE (s)-[:USES]->(t);

MATCH (s:Concept {key: 'sinthoma'})
MATCH (d:Concept {key: 'demand-desire'})
MERGE (s)-[:USES]->(d);

MATCH (f:Concept {key: 'fdpn'})
MATCH (sig:Concept {key: 'signifier'})
MERGE (f)-[:ANALYZES]->(sig);

// Relacionar conceitos ao módulo cortex
MATCH (c:Module {key: 'cortex'})
MATCH (con:Concept) WHERE con.theory = 'Lacan'
MERGE (c)-[:IMPLEMENTS]->(con);

// =====================================================
// 7. CRIAR SISTEMA DE MEMÓRIA
// =====================================================

MERGE (mem:MemorySystem:EVAKnowledge {
    key: 'memory-system',
    name: 'Sistema de Memória 12-Camadas',
    description: 'Arquitetura de memória baseada em Schacter, van der Kolk e Gurdjieff'
})
SET mem.layers = [
    'episodic', 'semantic', 'procedural', 'perceptual', 'working',
    'implicit', 'explicit', 'state-dependent',
    'enneagram', 'self-core', 'deep-archetypal', 'consciousness'
];

MATCH (h:Module {key: 'hippocampus'})
MATCH (m:MemorySystem {key: 'memory-system'})
MERGE (h)-[:IMPLEMENTS]->(m);

// =====================================================
// 8. CRIAR FERRAMENTAS
// =====================================================

MERGE (t1:Tool:EVAKnowledge {key: 'get_vitals', name: 'Obter Sinais Vitais', category: 'health'});
MERGE (t2:Tool:EVAKnowledge {key: 'get_agendamentos', name: 'Obter Agendamentos', category: 'calendar'});
MERGE (t3:Tool:EVAKnowledge {key: 'scan_medication_visual', name: 'Escanear Medicamento', category: 'vision'});
MERGE (t4:Tool:EVAKnowledge {key: 'apply_phq9', name: 'Aplicar PHQ-9', category: 'assessment'});
MERGE (t5:Tool:EVAKnowledge {key: 'apply_gad7', name: 'Aplicar GAD-7', category: 'assessment'});
MERGE (t6:Tool:EVAKnowledge {key: 'apply_cssrs', name: 'Aplicar C-SSRS', category: 'crisis', critical: true});
MERGE (t7:Tool:EVAKnowledge {key: 'analyze_voice_prosody', name: 'Analisar Prosódia', category: 'voice'});

MATCH (tm:Module {key: 'tools'})
MATCH (t:Tool)
MERGE (tm)-[:PROVIDES]->(t);

// =====================================================
// 9. CRIAR FLUXO DE DADOS
// =====================================================

// Nós de processo
MERGE (p1:Process:EVAKnowledge {key: 'audio-input', name: 'Entrada de Áudio', order: 1});
MERGE (p2:Process:EVAKnowledge {key: 'websocket', name: 'WebSocket Signaling', order: 2});
MERGE (p3:Process:EVAKnowledge {key: 'context-build', name: 'Construção de Contexto', order: 3});
MERGE (p4:Process:EVAKnowledge {key: 'unified-retrieval', name: 'UnifiedRetrieval (RSI)', order: 4});
MERGE (p5:Process:EVAKnowledge {key: 'gemini-call', name: 'Chamada Gemini LLM', order: 5});
MERGE (p6:Process:EVAKnowledge {key: 'tool-exec', name: 'Execução de Tools', order: 6});
MERGE (p7:Process:EVAKnowledge {key: 'memory-save', name: 'Salvar Memória', order: 7});
MERGE (p8:Process:EVAKnowledge {key: 'response', name: 'Resposta ao Cliente', order: 8});

// Encadear processos
MATCH (p1:Process {key: 'audio-input'})
MATCH (p2:Process {key: 'websocket'})
MERGE (p1)-[:FLOWS_TO]->(p2);

MATCH (p2:Process {key: 'websocket'})
MATCH (p3:Process {key: 'context-build'})
MERGE (p2)-[:FLOWS_TO]->(p3);

MATCH (p3:Process {key: 'context-build'})
MATCH (p4:Process {key: 'unified-retrieval'})
MERGE (p3)-[:FLOWS_TO]->(p4);

MATCH (p4:Process {key: 'unified-retrieval'})
MATCH (p5:Process {key: 'gemini-call'})
MERGE (p4)-[:FLOWS_TO]->(p5);

MATCH (p5:Process {key: 'gemini-call'})
MATCH (p6:Process {key: 'tool-exec'})
MERGE (p5)-[:MAY_TRIGGER]->(p6);

MATCH (p5:Process {key: 'gemini-call'})
MATCH (p7:Process {key: 'memory-save'})
MERGE (p5)-[:FLOWS_TO]->(p7);

MATCH (p7:Process {key: 'memory-save'})
MATCH (p8:Process {key: 'response'})
MERGE (p7)-[:FLOWS_TO]->(p8);

// =====================================================
// 10. ÍNDICES E CONSTRAINTS
// =====================================================

CREATE INDEX eva_knowledge_key IF NOT EXISTS FOR (n:EVAKnowledge) ON (n.key);
CREATE INDEX module_key IF NOT EXISTS FOR (n:Module) ON (n.key);
CREATE INDEX service_key IF NOT EXISTS FOR (n:Service) ON (n.key);
CREATE INDEX concept_key IF NOT EXISTS FOR (n:Concept) ON (n.key);
CREATE INDEX tool_key IF NOT EXISTS FOR (n:Tool) ON (n.key);
CREATE INDEX creator_cpf IF NOT EXISTS FOR (n:Creator) ON (n.cpf);

// =====================================================
// 11. QUERY DE VERIFICAÇÃO
// =====================================================

// Contar nós criados
MATCH (n:EVAKnowledge) RETURN labels(n)[0] as type, count(*) as count ORDER BY count DESC;
