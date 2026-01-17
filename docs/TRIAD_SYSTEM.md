# Arquitetura Tríade: Esopo + Nasrudin + Zen

## 🎯 A Tríade Completa da Psique

Sistema de intervenção terapêutica com **três ferramentas complementares**:

1. **Esopo** → Superego (Moral/Estrutura)
2. **Nasrudin** → Inconsciente (Paradoxo/Desconstrução)  
3. **Zen** → Self/Vazio (Presença/Centramento)

---

## 📊 Quatro Collections no Qdrant

### 1. `nasrudin_stories` (270 histórias)
- **Função:** Quebrar rigidez, obsessão, tédio
- **Método:** Paradoxo, absurdo, humor
- **Zeta Affinity:** Tipos Emocionais (2, 4, 7, 9)
- **Quando:** Negação, resistência, repetição compulsiva

### 2. `aesop_fables` (~300 fábulas)
- **Função:** Ensinar causa-efeito, responsabilidade
- **Método:** Lógica, moral, consequência
- **Zeta Affinity:** Tipos Racionais (1, 3, 5, 6)
- **Quando:** Imprudência, racionalização, busca de atenção

### 3. `zen_koans` (~50 histórias)
- **Função:** Esvaziar a mente, parar o pensamento
- **Método:** Choque de insight, silêncio
- **Zeta Affinity:** Tipos Analíticos/Introspectivos (1, 4, 5, 9)
- **Quando:** Overthinking, ansiedade mental, saturação cognitiva

### 4. `somatic_exercises` (~20 exercícios)
- **Função:** Aterramento, presença corporal
- **Método:** Comandos somáticos diretos
- **Sintomas:** Pânico, dissociação, hiperventilação
- **Quando:** Crise aguda, ansiedade física

---

## 🧠 Lógica de Decisão Expandida

```
┌─────────────────────────────────────────┐
│  Usuário fala algo                      │
└──────────────┬──────────────────────────┘
               ↓
┌──────────────────────────────────────────┐
│  TransNAR analisa (Lacan)                │
│  - Detecta padrão clínico                │
│  - Mede intensidade (0-1)                │
└──────────────┬───────────────────────────┘
               ↓
┌──────────────────────────────────────────┐
│  Zeta Router + Estado Fisiológico       │
└──────────────┬───────────────────────────┘
               ↓
         ┌─────┴─────┐
         │  DECISÃO  │
         └─────┬─────┘
               ↓
    ┌──────────┼──────────┬──────────┐
    │          │          │          │
    ▼          ▼          ▼          ▼
┌────────┐ ┌────────┐ ┌──────┐ ┌─────────┐
│ ESOPO  │ │NASRUDIN│ │ ZEN  │ │SOMATIC  │
│        │ │        │ │KOAN  │ │EXERCISE │
└────────┘ └────────┘ └──────┘ └─────────┘
    │          │          │          │
Racional   Emocional  Mental    Físico
1,3,5,6    2,4,7,9    1,4,5,9   Crise
```

---

## 💡 Exemplos Práticos

### Cenário 1: Overthinking (Tipo 5 - Investigador)

**Fala:** *"Eu li tudo sobre minha doença, analisei os dados, mas não entendo por que não melhoro."*

**TransNAR detecta:** `intellectualization` (defesa excessiva)  
**Zeta Type:** 5 (Investigador - Analítico)  
**Estado:** Mental saturation  
**Sistema escolhe:** **ZEN KOAN** → "Uma Xícara de Chá"

**EVA responde:**
> "Sua mente parece a xícara de chá do mestre Nan-in. Ele tentava servir chá para um professor, mas a xícara já estava cheia e transbordou. 'Como posso lhe mostrar o Zen', disse o mestre, 'se sua mente já está cheia de opiniões?' Que tal pararmos de pensar por um minuto?"

---

### Cenário 2: Crise de Pânico (Qualquer Tipo)

**Fala:** *"Não consigo respirar! Meu coração está acelerado! Vou morrer!"*

**TransNAR detecta:** `panic_attack` (fisiológico)  
**Sensores:** Frequência cardíaca 140 bpm  
**Sistema escolhe:** **SOMATIC EXERCISE** → "Atenção na Respiração"

**EVA responde (voz calma, comando direto):**
> "Pare tudo agora. Feche os olhos. Sinta o ar entrando pelo nariz. Você está entrando com ele. Sinta o ar saindo. Você está saindo com ele. Só isso. Nada mais."

*(EVA monitora frequência cardíaca em tempo real e ajusta duração do exercício)*

---

### Cenário 3: Racionalização (Tipo 3 - Realizador)

**Fala:** *"Eu não preciso de ajuda. Sempre fiz tudo sozinho."*

**TransNAR detecta:** `rationalization` + `projection`  
**Zeta Type:** 3 (Racional)  
**Sistema escolhe:** **ESOPO** → "A Raposa e as Uvas"

*(Já documentado anteriormente)*

---

### Cenário 4: Depressão/Tédio (Tipo 9 - Pacificador)

**Fala:** *"Nada faz sentido. Tudo é sempre igual."*

**TransNAR detecta:** `learned_helplessness` + `apathy`  
**Zeta Type:** 9 (Emocional/Pacificador)  
**Sistema escolhe:** **NASRUDIN** → "A Chave e a Luz"

*(Já documentado anteriormente)*

---

## 🔧 Implementação Backend (Go)

```go
func (h *Handler) SelectIntervention(
    transnarResult map[string]interface{},
    zetaType int,
    vitalSigns VitalSigns,
) (*Intervention, error) {
    
    // PRIORIDADE 1: Crise Física (sobrepõe tudo)
    if vitalSigns.HeartRate > 120 || vitalSigns.Panic {
        return h.somaticMatcher.FindExercise("panic_attack")
    }
    
    // PRIORIDADE 2: Saturação Mental
    if transnarResult["defense"] == "intellectualization" {
        return h.zenMatcher.FindKoan("mental_saturation")
    }
    
    // PRIORIDADE 3: Roteamento Zeta Normal
    switch zetaType {
    case Type1, Type3, Type5, Type6:
        return h.aesopMatcher.FindFable(transnarResult)
        
    case Type2, Type4, Type7, Type9:
        return h.nasrudinMatcher.FindStory(transnarResult)
        
    default:
        return h.aesopMatcher.FindFable(transnarResult)
    }
}
```

---

## 📱 Implementação Frontend (Flutter)

```dart
void _handleIntervention(Map<String, dynamic> data) {
  final type = data['intervention_type'];
  
  switch (type) {
    case 'didactic':  // Esopo
      _showMoralCard(data);
      break;
      
    case 'paradox':  // Nasrudin
      _showHumorCard(data);
      break;
      
    case 'zen_koan':  // Zen Narrativo
      _showZenCard(data);
      _playGongSound();  // Som de sino zen
      break;
      
    case 'somatic':  // Exercício Físico
      _startGuidedExercise(data);
      _showBreathingAnimation();  // Animação de respiração
      _monitorVitals();  // Monitorar sinais vitais
      break;
  }
}
```

---

## 📦 Schema Lacaniano (Exemplos)

### Zen Koan: "Uma Xícara de Chá"
```json
{
  "koan_id": "zen_001_xicara_cha",
  "title": "Uma Xícara de Chá",
  "text": "Nan-in serviu chá...",
  "clinical_tags": {
    "transnar_rule": "intellectualization",
    "target_state": "mental_saturation",
    "intervention_type": "shock_insight",
    "zeta_affinity": [1, 5, 6]
  },
  "trigger_condition": "User overthinks, cannot stop mental chatter",
  "eva_followup": "Sua mente está cheia como a xícara..."
}
```

### Somatic Exercise: "Equilíbrio"
```json
{
  "exercise_id": "somatic_001_equilibrio",
  "title": "Centralização pelo Equilíbrio",
  "instruction": "Fique em pé. Sinta os dois pés...",
  "symptoms": ["panic_attack", "dizziness", "disassociation"],
  "action": "grounding",
  "duration_seconds": 60,
  "eva_voice_command": "Pare tudo. Fique em pé..."
}
```

---

## 🚀 Execução no Servidor

```bash
# 1. Popular Nasrudin (270 histórias)
python3 scripts/populate_nasrudin_with_lacan.py

# 2. Popular Esopo (~300 fábulas)
python3 scripts/populate_aesop_fables.py

# 3. Popular Zen (koans + exercícios)
python3 scripts/populate_zen_content.py

# 4. Verificar todas collections
curl http://localhost:6333/collections | jq '.result.collections[].name'
```

---

## ✅ Resultado Final: A Tríade Completa

| Ferramenta | Função | Alvo | Zeta | Quando |
|------------|--------|------|------|--------|
| **Esopo** | Estruturar | Superego | 1,3,5,6 | Imprudência |
| **Nasrudin** | Desconstruir | Inconsciente | 2,4,7,9 | Rigidez |
| **Zen Koan** | Esvaziar | Mente | 1,4,5,9 | Overthinking |
| **Somatic** | Aterrar | Corpo | Todos | Crise Física |

**EVA agora é um Mestre Zen Digital completo!** 🧘‍♀️✨
