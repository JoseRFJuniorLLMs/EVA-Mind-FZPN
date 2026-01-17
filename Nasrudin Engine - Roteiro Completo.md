# Nasrudin Engine - Implementação Passo a Passo
## EVA-Mind-FZPN - Fase 2

---

## 📋 Índice

1. [Visão Geral da Arquitetura](#visao-geral)
2. [Fase 1: Banco de Dados (Qdrant)](#fase-1-qdrant)
3. [Fase 2: Backend - Lógica de Detecção](#fase-2-backend)
4. [Fase 3: Backend - Geração de Resposta](#fase-3-geracao)
5. [Fase 4: Frontend - Apresentação Visual](#fase-4-frontend)
6. [Fase 5: Testes e Validação](#fase-5-testes)
7. [Checklist de Implementação](#checklist)

---

## 1. Visão Geral da Arquitetura {#visao-geral}

### Fluxo de Dados

```
[Usuário] → [Mensagem] 
    ↓
[TransNAR/Zeta] → Detecta padrão (ex: Learned Helplessness)
    ↓
[Nasrudin Matcher] → Busca história no Qdrant
    ↓
[LLM Gemini] → Narra história de forma empática
    ↓
[WebSocket] → Envia áudio + visual para Flutter
    ↓
[EVA App] → Exibe card "💡 Momento de Sabedoria"
```

### Componentes Afetados

| Componente | Ação | Arquivo |
|------------|------|---------|
| **Qdrant** | Nova collection `nasrudin_stories` | - |
| **Backend Go** | Nasrudin Detector + Matcher | `pkg/nasrudin/` |
| **LLM Service** | Prompt especial para narrativa | `pkg/llm/gemini.go` |
| **WebSocket** | Novo tipo de mensagem | `pkg/websocket/handler.go` |
| **Flutter** | Card visual + animação | `lib/services/eva_webview_service.dart` |

---

## 2. Fase 1: Banco de Dados (Qdrant) {#fase-1-qdrant}

### 2.1 Criar Collection no Qdrant

```bash
# SSH no servidor
ssh root@104.248.219.200

# Criar collection
curl -X PUT http://localhost:6333/collections/nasrudin_stories \
  -H 'Content-Type: application/json' \
  -d '{
    "vectors": {
      "size": 768,
      "distance": "Cosine"
    },
    "on_disk_payload": true
  }'
```

**Verificar criação:**
```bash
curl http://localhost:6333/collections/nasrudin_stories | jq
```

### 2.2 Estrutura de Dados das Histórias

Cada história terá este formato:

```json
{
  "id": 1,
  "vector": [0.1, 0.2, ..., 0.768],
  "payload": {
    "story_id": "nasrudin_001",
    "title": "A Chave Perdida",
    "text": "Nasrudin procurava algo no chão, sob um poste de luz. 'O que você perdeu?', perguntou um amigo. 'Minha chave', respondeu. 'Onde você a perdeu?' 'Lá dentro de casa.' 'Então por que procura aqui fora?' 'Porque aqui tem mais luz.'",
    "moral": "Às vezes procuramos soluções nos lugares mais fáceis, não nos lugares certos.",
    "tags": ["procurar_fora", "solucao_facil", "negacao", "resistencia"],
    "transnar_triggers": ["negation_as_desire", "learned_helplessness", "resistance"],
    "personality_types": [1, 3, 5, 6],
    "emotional_tone": "gentle_humor",
    "language": "pt-BR"
  }
}
```

### 2.3 Script para Popular Histórias

Crie o arquivo: `/root/EVA-Mind-FZPN/scripts/populate_nasrudin.py`

```python
#!/usr/bin/env python3
import requests
import json
from openai import OpenAI

# Configuração
QDRANT_URL = "http://localhost:6333"
GEMINI_API_KEY = "AIzaSyBnSKHtKNKJVxO-qvABiWPPJVWJzLlOhYo"

# Cliente OpenAI (compatível com Gemini via proxy)
client = OpenAI(
    api_key=GEMINI_API_KEY,
    base_url="https://generativelanguage.googleapis.com/v1beta/openai/"
)

# Histórias de Nasrudin (10 iniciais)
STORIES = [
    {
        "story_id": "nasrudin_001",
        "title": "A Chave Perdida",
        "text": "Nasrudin procurava algo no chão, sob um poste de luz. 'O que você perdeu?', perguntou um amigo. 'Minha chave', respondeu. 'Onde você a perdeu?' 'Lá dentro de casa.' 'Então por que procura aqui fora?' 'Porque aqui tem mais luz.'",
        "moral": "Às vezes procuramos soluções nos lugares mais fáceis, não nos lugares certos.",
        "tags": ["procurar_fora", "solucao_facil", "negacao"],
        "transnar_triggers": ["negation_as_desire", "learned_helplessness"],
        "personality_types": [1, 3, 5]
    },
    {
        "story_id": "nasrudin_002",
        "title": "O Vinagre",
        "text": "Nasrudin comprou um barril de vinagre. No caminho de casa, começou a chover. 'Corra!', gritaram. 'Seu vinagre vai se molhar!' Nasrudin parou e respondeu: 'E o que isso importa? O vinagre já é molhado.'",
        "moral": "Às vezes nos preocupamos com problemas que já existem dentro de nós.",
        "tags": ["preocupacao_inutil", "aceitacao", "perspectiva"],
        "transnar_triggers": ["resistance", "internal_contradiction"],
        "personality_types": [6, 7, 9]
    },
    {
        "story_id": "nasrudin_003",
        "title": "A Porta",
        "text": "Ladrões invadiram a casa de Nasrudin e roubaram tudo. Quando foram embora, Nasrudin os seguiu. 'Por que está nos seguindo?', perguntaram assustados. 'Estou me mudando também', respondeu. 'Vocês levaram tudo, inclusive minha casa.'",
        "moral": "Quando perdemos tudo, ainda temos a escolha de como reagir.",
        "tags": ["perda", "adaptacao", "humor_negro"],
        "transnar_triggers": ["death_drive", "learned_helplessness"],
        "personality_types": [4, 8, 9]
    },
    {
        "story_id": "nasrudin_004",
        "title": "O Espelho",
        "text": "Nasrudin viu seu reflexo na água pela primeira vez. 'Finalmente!', exclamou. 'Encontrei alguém mais feio que eu!'",
        "moral": "O que vemos nos outros muitas vezes é nosso próprio reflexo.",
        "tags": ["projecao", "autoconhecimento", "julgamento"],
        "transnar_triggers": ["projection", "reactive_formation"],
        "personality_types": [1, 3, 4]
    },
    {
        "story_id": "nasrudin_005",
        "title": "A Sopa de Sopa",
        "text": "Um amigo trouxe um pato para Nasrudin. Ele fez uma sopa deliciosa. No dia seguinte, outro homem apareceu dizendo 'Sou amigo do amigo que trouxe o pato'. Nasrudin serviu-lhe água quente. 'Que sopa é esta?' 'É a sopa da sopa do pato.'",
        "moral": "Quanto mais longe da fonte, mais diluída fica a verdade.",
        "tags": ["aproveitadores", "limites", "diluicao"],
        "transnar_triggers": ["transference", "compulsive_repetition"],
        "personality_types": [2, 6, 9]
    },
    {
        "story_id": "nasrudin_006",
        "title": "O Burro Morto",
        "text": "Nasrudin vendeu rifas de um burro morto. Quando descobriram, perguntaram: 'E agora?' 'Devolvo o dinheiro de quem reclamar', disse. 'E quantos reclamaram?' 'Só o dono original. Mas já dei outro burro morto para ele.'",
        "moral": "Às vezes a mentira funciona porque as pessoas preferem acreditar.",
        "tags": ["ilusao", "expectativa", "conformismo"],
        "transnar_triggers": ["denial_as_desire", "compulsive_repetition"],
        "personality_types": [3, 7, 8]
    },
    {
        "story_id": "nasrudin_007",
        "title": "A Lua no Poço",
        "text": "Nasrudin viu o reflexo da lua num poço. 'A lua caiu!', gritou. Jogou uma corda, ela prendeu numa pedra. Ao puxar com força, caiu de costas. Olhou para o céu e viu a lua. 'Que bom! Consegui colocá-la de volta.'",
        "moral": "Às vezes nos damos crédito por mudanças que aconteceriam de qualquer forma.",
        "tags": ["ilusao_controle", "falsa_causalidade", "ego"],
        "transnar_triggers": ["reactive_formation", "sublimation"],
        "personality_types": [1, 3, 8]
    },
    {
        "story_id": "nasrudin_008",
        "title": "A Morte do Burro",
        "text": "Nasrudin treinou seu burro a não comer. Ia muito bem, até que o burro morreu. 'Que pena', disse Nasrudin. 'Bem quando ele estava aprendendo.'",
        "moral": "Podemos persistir em estratégias que matam aquilo que amamos.",
        "tags": ["persistencia_errada", "controle", "morte"],
        "transnar_triggers": ["death_drive", "compulsive_repetition"],
        "personality_types": [1, 5, 8]
    },
    {
        "story_id": "nasrudin_009",
        "title": "O Juiz",
        "text": "Dois homens brigavam. O primeiro falou. 'Você tem razão', disse Nasrudin. O segundo falou. 'Você também tem razão', disse Nasrudin. Sua esposa protestou: 'Eles dizem coisas opostas!' 'Você também tem razão', respondeu Nasrudin.",
        "moral": "Concordar com todos é não ter posição própria.",
        "tags": ["evitacao_conflito", "indecisao", "complacencia"],
        "transnar_triggers": ["resistance", "transference"],
        "personality_types": [2, 6, 9]
    },
    {
        "story_id": "nasrudin_010",
        "title": "O Anel Mágico",
        "text": "Nasrudin vendia um 'anel mágico que torna invisível'. Um homem comprou. Colocou o anel e perguntou: 'Estou invisível?' 'Não consigo te ver', mentiu Nasrudin. O homem saiu feliz roubando no mercado. Foi preso imediatamente.",
        "moral": "Acreditar em soluções mágicas pode nos cegar para a realidade.",
        "tags": ["solucao_magica", "negacao_realidade", "consequencias"],
        "transnar_triggers": ["denial_as_desire", "learned_helplessness"],
        "personality_types": [7, 3, 4]
    }
]

def generate_embedding(text):
    """Gera embedding usando Gemini via OpenAI SDK"""
    response = client.embeddings.create(
        model="text-embedding-004",
        input=text
    )
    return response.data[0].embedding

def insert_story(story, point_id):
    """Insere uma história no Qdrant"""
    # Gerar embedding do texto + moral (contexto semântico completo)
    full_text = f"{story['title']}. {story['text']} Moral: {story['moral']}"
    vector = generate_embedding(full_text)
    
    # Estrutura do ponto
    point = {
        "id": point_id,
        "vector": vector,
        "payload": story
    }
    
    # Inserir no Qdrant
    response = requests.put(
        f"{QDRANT_URL}/collections/nasrudin_stories/points",
        json={"points": [point]}
    )
    
    if response.status_code == 200:
        print(f"✅ Inserido: {story['title']} (ID: {point_id})")
    else:
        print(f"❌ Erro ao inserir {story['title']}: {response.text}")

def main():
    print("🚀 Populando Qdrant com histórias de Nasrudin...\n")
    
    for idx, story in enumerate(STORIES, start=1):
        insert_story(story, idx)
    
    print(f"\n✅ {len(STORIES)} histórias inseridas com sucesso!")
    
    # Verificar
    response = requests.get(f"{QDRANT_URL}/collections/nasrudin_stories")
    if response.status_code == 200:
        data = response.json()
        print(f"\n📊 Collection info:")
        print(f"   Points: {data['result']['points_count']}")
        print(f"   Vectors: {data['result']['vectors_count']}")

if __name__ == "__main__":
    main()
```

**Executar:**
```bash
cd /root/EVA-Mind-FZPN
chmod +x scripts/populate_nasrudin.py
python3 scripts/populate_nasrudin.py
```

---

## 3. Fase 2: Backend - Lógica de Detecção {#fase-2-backend}

### 3.1 Estrutura de Arquivos

```
pkg/nasrudin/
├── detector.go        # Detecta quando usar Nasrudin
├── matcher.go         # Busca história apropriada no Qdrant
├── narrator.go        # Gera narrativa empática via LLM
└── models.go          # Structs e tipos
```

### 3.2 Código: `pkg/nasrudin/models.go`

```go
package nasrudin

import "time"

// NasrudinStory representa uma história de Nasrudin
type NasrudinStory struct {
    StoryID          string   `json:"story_id"`
    Title            string   `json:"title"`
    Text             string   `json:"text"`
    Moral            string   `json:"moral"`
    Tags             []string `json:"tags"`
    TransnarTriggers []string `json:"transnar_triggers"`
    PersonalityTypes []int    `json:"personality_types"`
    EmotionalTone    string   `json:"emotional_tone"`
    Language         string   `json:"language"`
}

// NasrudinIntervention representa uma intervenção completa
type NasrudinIntervention struct {
    Story         NasrudinStory `json:"story"`
    NarrativeText string        `json:"narrative_text"` // História narrada pela LLM
    Trigger       string        `json:"trigger"`        // Qual padrão ativou
    Confidence    float64       `json:"confidence"`
    Timestamp     time.Time     `json:"timestamp"`
}

// DetectionContext é o contexto para decisão de intervenção
type DetectionContext struct {
    TransnarInference map[string]interface{} // Resultado do TransNAR
    UserPersonality   int                    // Tipo Eneagrama
    ConversationTone  string                 // "resistant", "open", "desperate"
    RecentPatterns    []string               // Padrões detectados recentemente
}
```

### 3.3 Código: `pkg/nasrudin/detector.go`

```go
package nasrudin

import (
    "fmt"
)

// Detector decide SE e QUANDO usar Nasrudin
type Detector struct {
    minConfidence float64
}

// NewDetector cria um novo detector
func NewDetector() *Detector {
    return &Detector{
        minConfidence: 0.75, // Só atua com alta confiança
    }
}

// ShouldIntervene decide se deve usar Nasrudin
func (d *Detector) ShouldIntervene(ctx DetectionContext) (bool, string, float64) {
    // Regra 1: TransNAR detectou padrão com confiança suficiente?
    if transnarConfidence, ok := ctx.TransnarInference["confidence"].(float64); ok {
        if transnarConfidence < d.minConfidence {
            return false, "", 0.0
        }
    }

    // Regra 2: Há um trigger compatível com Nasrudin?
    trigger, confidence := d.findBestTrigger(ctx)
    if trigger == "" {
        return false, "", 0.0
    }

    // Regra 3: O tom da conversa permite intervenção paradoxal?
    // Não usar se usuário está em crise aguda
    if ctx.ConversationTone == "crisis" || ctx.ConversationTone == "suicidal" {
        return false, "", 0.0
    }

    // Regra 4: Não repetir intervenção (checar padrões recentes)
    for _, pattern := range ctx.RecentPatterns {
        if pattern == "nasrudin_used" {
            return false, "", 0.0 // Já usou recentemente
        }
    }

    return true, trigger, confidence
}

// findBestTrigger identifica qual padrão TransNAR é melhor para Nasrudin
func (d *Detector) findBestTrigger(ctx DetectionContext) (string, float64) {
    triggers := map[string]float64{
        "negation_as_desire":      0.0,
        "learned_helplessness":    0.0,
        "resistance":              0.0,
        "projection":              0.0,
        "reactive_formation":      0.0,
        "compulsive_repetition":   0.0,
    }

    // Extrair scores do TransNAR
    if rules, ok := ctx.TransnarInference["rules"].(map[string]interface{}); ok {
        for rule, data := range rules {
            if dataMap, ok := data.(map[string]interface{}); ok {
                if confidence, ok := dataMap["confidence"].(float64); ok {
                    triggers[rule] = confidence
                }
            }
        }
    }

    // Encontrar o maior score
    bestTrigger := ""
    bestScore := 0.0
    for trigger, score := range triggers {
        if score > bestScore {
            bestScore = score
            bestTrigger = trigger
        }
    }

    if bestScore < d.minConfidence {
        return "", 0.0
    }

    return bestTrigger, bestScore
}
```

### 3.4 Código: `pkg/nasrudin/matcher.go`

```go
package nasrudin

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
)

// Matcher busca a história mais apropriada no Qdrant
type Matcher struct {
    qdrantURL string
    client    *http.Client
}

// NewMatcher cria um novo matcher
func NewMatcher(qdrantURL string) *Matcher {
    return &Matcher{
        qdrantURL: qdrantURL,
        client:    &http.Client{},
    }
}

// FindStory busca história baseada no trigger e personalidade
func (m *Matcher) FindStory(ctx context.Context, trigger string, personality int, queryVector []float64) (*NasrudinStory, error) {
    // Construir query para Qdrant
    searchQuery := map[string]interface{}{
        "vector": queryVector,
        "limit":  3,
        "score_threshold": 0.7,
        "filter": map[string]interface{}{
            "must": []map[string]interface{}{
                {
                    "key": "transnar_triggers",
                    "match": map[string]interface{}{
                        "any": []string{trigger},
                    },
                },
            },
        },
    }

    // Se personalidade foi informada, filtrar também
    if personality > 0 {
        searchQuery["filter"].(map[string]interface{})["should"] = []map[string]interface{}{
            {
                "key": "personality_types",
                "match": map[string]interface{}{
                    "any": []int{personality},
                },
            },
        }
    }

    payload, _ := json.Marshal(searchQuery)
    
    // Fazer requisição ao Qdrant
    url := fmt.Sprintf("%s/collections/nasrudin_stories/points/search", m.qdrantURL)
    req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(payload)))
    if err != nil {
        return nil, fmt.Errorf("erro ao criar request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := m.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar no Qdrant: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("Qdrant retornou erro %d: %s", resp.StatusCode, string(body))
    }

    // Parse da resposta
    var qdrantResp struct {
        Result []struct {
            Score   float64         `json:"score"`
            Payload NasrudinStory   `json:"payload"`
        } `json:"result"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&qdrantResp); err != nil {
        return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
    }

    if len(qdrantResp.Result) == 0 {
        return nil, fmt.Errorf("nenhuma história encontrada para trigger: %s", trigger)
    }

    // Retornar a história com maior score
    return &qdrantResp.Result[0].Payload, nil
}
```

### 3.5 Código: `pkg/nasrudin/narrator.go`

```go
package nasrudin

import (
    "context"
    "fmt"
)

// Narrator gera a narrativa empática da história via LLM
type Narrator struct {
    llmClient LLMClient // Interface para Gemini
}

// LLMClient interface (já deve existir no seu código)
type LLMClient interface {
    Generate(ctx context.Context, prompt string) (string, error)
}

// NewNarrator cria um novo narrador
func NewNarrator(llmClient LLMClient) *Narrator {
    return &Narrator{llmClient: llmClient}
}

// Narrate gera a versão narrada da história
func (n *Narrator) Narrate(ctx context.Context, story NasrudinStory, userMessage string) (string, error) {
    prompt := fmt.Sprintf(`Você é EVA, uma assistente empática. O usuário disse: "%s"

Você detectou um padrão de pensamento que pode ser explorado através desta história:

**%s**

%s

**Moral**: %s

INSTRUÇÕES:
1. NÃO dê conselhos diretos
2. NÃO explique a moral explicitamente
3. Conte a história de forma natural e envolvente
4. Termine com uma pergunta aberta tipo: "O que você acha dessa história?"
5. Mantenha tom leve e gentil
6. Máximo 4 frases antes da história, máximo 2 depois

Responda APENAS com a narrativa, sem formatação markdown.`, 
        userMessage,
        story.Title,
        story.Text,
        story.Moral,
    )

    narrative, err := n.llmClient.Generate(ctx, prompt)
    if err != nil {
        return "", fmt.Errorf("erro ao gerar narrativa: %w", err)
    }

    return narrative, nil
}
```

---

## 4. Fase 3: Backend - Integração no Handler {#fase-3-geracao}

### 4.1 Modificar `pkg/websocket/handler.go`

Adicione ao handler principal:

```go
package websocket

import (
    "github.com/seu-repo/eva/pkg/nasrudin"
    "github.com/seu-repo/eva/pkg/transnar"
    // ... outros imports
)

type Handler struct {
    // ... campos existentes
    nasrudinDetector *nasrudin.Detector
    nasrudinMatcher  *nasrudin.Matcher
    nasrudinNarrator *nasrudin.Narrator
}

func NewHandler(/* params existentes */, qdrantURL string, llmClient nasrudin.LLMClient) *Handler {
    return &Handler{
        // ... inicializações existentes
        nasrudinDetector: nasrudin.NewDetector(),
        nasrudinMatcher:  nasrudin.NewMatcher(qdrantURL),
        nasrudinNarrator: nasrudin.NewNarrator(llmClient),
    }
}

func (h *Handler) ProcessMessage(ctx context.Context, userID int, message string) (*Response, error) {
    // 1. Processar TransNAR (já existe)
    transnarResult, err := h.transnarEngine.Analyze(ctx, userID, message)
    if err != nil {
        return nil, err
    }

    // 2. Verificar se deve usar Nasrudin
    detectionCtx := nasrudin.DetectionContext{
        TransnarInference: transnarResult,
        UserPersonality:   h.getUserPersonality(userID), // Implementar
        ConversationTone:  h.analyzeTone(message),       // Implementar
        RecentPatterns:    h.getRecentPatterns(userID),  // Implementar
    }

    shouldIntervene, trigger, confidence := h.nasrudinDetector.ShouldIntervene(detectionCtx)

    if shouldIntervene {
        return h.handleNasrudinIntervention(ctx, userID, message, trigger, confidence)
    }

    // 3. Processar normalmente se não usar Nasrudin
    return h.processNormalResponse(ctx, userID, message, transnarResult)
}

func (h *Handler) handleNasrudinIntervention(ctx context.Context, userID int, message, trigger string, confidence float64) (*Response, error) {
    // 1. Gerar embedding da mensagem do usuário
    queryVector, err := h.embeddingService.Generate(ctx, message)
    if err != nil {
        return nil, err
    }

    // 2. Buscar história apropriada
    personality := h.getUserPersonality(userID)
    story, err := h.nasrudinMatcher.FindStory(ctx, trigger, personality, queryVector)
    if err != nil {
        // Fallback: resposta normal se não encontrar história
        return h.processNormalResponse(ctx, userID, message, nil)
    }

    // 3. Gerar narrativa empática
    narrative, err := h.nasrudinNarrator.Narrate(ctx, *story, message)
    if err != nil {
        return nil, err
    }

    // 4. Gerar áudio (TTS)
    audioURL, err := h.ttsService.Generate(ctx, narrative)
    if err != nil {
        return nil, err
    }

    // 5. Construir resposta especial
    response := &Response{
        Type: "intervention_nasrudin",
        Data: map[string]interface{}{
            "story": map[string]interface{}{
                "title":     story.Title,
                "text":      story.Text,
                "narrative": narrative,
            },
            "audio_url":  audioURL,
            "trigger":    trigger,
            "confidence": confidence,
        },
    }

    // 6. Registrar uso (para evitar repetição)
    h.recordPattern(userID, "nasrudin_used")

    return response, nil
}
```

---

## 5. Fase 4: Frontend - Apresentação Visual {#fase-4-frontend}

### 5.1 Modificar `lib/services/eva_webview_service.dart`

```dart
// Adicionar em EVAWebViewService

void _handleMessageFromJS(Map<String, dynamic> message) {
  final type = message['type'];
  final data = message['data'];

  switch (type) {
    case 'user_message':
      _handleUserMessage(data);
      break;
      
    case 'intervention_nasrudin':  // NOVO
      _handleNasrudinIntervention(data);
      break;
      
    case 'assistant_response':
      _handleAssistantResponse(data);
      break;
      
    // ... outros casos
  }
}

void _handleNasrudinIntervention(Map<String, dynamic> data) async {
  final story = data['story'];
  final audioUrl = data['audio_url'];
  
  // 1. Mostrar notificação visual
  _showNasrudinCard(story);
  
  // 2. Tocar áudio da narrativa
  if (audioUrl != null) {
    await _playAudio(audioUrl);
  }
  
  // 3. Notificar UI para animar card
  onInterventionReceived?.call('💡 Momento de Sabedoria', story['narrative']);
}

void _showNasrudin