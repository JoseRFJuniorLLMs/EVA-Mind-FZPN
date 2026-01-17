# Arquitetura Dual: Nasrudin + Esopo

## 🎯 Visão Geral

Sistema de intervenção terapêutica com **duas ferramentas complementares**:

- **Nasrudin (Sufi/Osho)** → Desconstruir o Ego (Paradoxo/Humor)
- **Esopo (Moralista/Lacan)** → Estruturar a Lei Simbólica (Moral/Consequência)

---

## 📊 Duas Collections no Qdrant

### 1. `nasrudin_stories` (270 histórias)
- **Função:** Quebrar rigidez, obsessão, tédio
- **Método:** Paradoxo, absurdo, humor
- **Zeta Affinity:** Tipos Emocionais/Místicos (2, 4, 7, 9)
- **Quando usar:** TransNAR detecta negação, resistência, repetição compulsiva

### 2. `aesop_fables` (~300 fábulas)
- **Função:** Ensinar causa-efeito, responsabilidade
- **Método:** Lógica, moral, consequência
- **Zeta Affinity:** Tipos Racionais/Pragmáticos (1, 3, 5, 6)
- **Quando usar:** TransNAR detecta imprudência, racionalização, busca de atenção

---

## 🧠 Lógica de Decisão (Zeta Switch)

```
┌─────────────────────────────────────────┐
│  Usuário fala algo                      │
└──────────────┬──────────────────────────┘
               ↓
┌──────────────────────────────────────────┐
│  TransNAR analisa (Lacan)                │
│  - Detecta: "projection" (Regra 9)       │
│  - Confiança: 0.87                       │
└──────────────┬───────────────────────────┘
               ↓
┌──────────────────────────────────────────┐
│  Zeta Router verifica personalidade     │
│  - Tipo: 3 (Realizador - Racional)      │
└──────────────┬───────────────────────────┘
               ↓
         ┌─────┴─────┐
         │  DECISÃO  │
         └─────┬─────┘
               ↓
    ┌──────────┴──────────┐
    │                     │
    ▼                     ▼
┌─────────┐         ┌──────────┐
│ NASRUDIN│         │  ESOPO   │
└─────────┘         └──────────┘
    │                     │
    │ Se Zeta 2,4,7,9     │ Se Zeta 1,3,5,6
    │ (Emocionais)        │ (Racionais)
    │                     │
    ▼                     ▼
"O Burro ao         "A Raposa e
 Contrário"          as Uvas"
```

---

## 💡 Exemplos Práticos

### Cenário 1: Idoso Racional (Tipo 3 - Realizador)

**Fala:** *"Eu não preciso de ajuda. Eu sempre fiz tudo sozinho. Ajuda é para fracos."*

**TransNAR detecta:** `projection` + `rationalization`  
**Zeta Type:** 3 (Racional, focado em eficiência)  
**Sistema escolhe:** **ESOPO** → "A Raposa e as Uvas"

**EVA responde:**
> "Sabe, isso me lembra uma raposa que viu uvas lindas mas não conseguia alcançá-las. Então ela disse: 'Ah, essas uvas devem estar verdes mesmo'. Será que às vezes dizemos que não precisamos de algo só porque é difícil pedir?"

---

### Cenário 2: Idoso Emocional (Tipo 7 - Entusiasta)

**Fala:** *"Todo mundo me abandona. Ninguém me visita. O mundo está contra mim."*

**TransNAR detecta:** `projection` + `victimization`  
**Zeta Type:** 7 (Emocional, busca prazer)  
**Sistema escolhe:** **NASRUDIN** → "O Burro ao Contrário"

**EVA responde:**
> "Lembra do Nasrudin que estava sentado de costas no burro? Quando perguntaram por que estava ao contrário, ele disse: 'Não fui eu, foi o burro que virou!' Às vezes achamos que o mundo virou, mas na verdade somos nós que estamos olhando para o lado errado, não acha?"

---

## 🔧 Implementação Backend (Go)

```go
func (h *Handler) SelectParable(
    transnarResult map[string]interface{},
    zetaType int,
) (*Parable, error) {
    
    // Decisão baseada em Zeta Type
    var collection string
    
    switch zetaType {
    case Type1, Type3, Type5, Type6:
        // RACIONAIS → Querem lógica/moral
        collection = "aesop_fables"
        
    case Type2, Type4, Type7, Type9:
        // EMOCIONAIS → Querem paradoxo/humor
        collection = "nasrudin_stories"
        
    case Type8:
        // DESAFIADOR → Aceita ambos (escolher por intensidade)
        collection = "aesop_fables"  // Default
    }
    
    // Buscar no Qdrant
    return h.qdrantClient.Search(collection, transnarResult)
}
```

---

## 📱 Implementação Frontend (Flutter)

```dart
void _handleParableIntervention(Map<String, dynamic> data) {
  final mode = data['mode'];  // 'didactic' ou 'paradox'
  final content = data['content'];
  
  if (mode == 'didactic') {
    // ESOPO: Card sóbrio, voz calma
    _showDidacticCard(
      title: content['title'],
      moral: content['moral'],
      icon: Icons.book,
      color: Colors.brown
    );
  } else {
    // NASRUDIN: Card lúdico, voz irônica
    _showParadoxCard(
      title: content['title'],
      icon: Icons.psychology,
      color: Colors.purple
    );
  }
}
```

---

## 📦 Mapeamento Lacaniano (Fábulas-Chave)

### Esopo - XLVIII: "A Raposa e as Uvas"
- **TransNAR:** `projection`, `rationalization`
- **Conceito:** Sour Grapes Mechanism
- **Zeta:** 1, 3, 5, 6
- **Trigger:** User dismisses goals after failing
- **Followup:** "Será que estamos desdenhando só porque ficou difícil?"

### Nasrudin - 208: "O Burro ao Contrário"
- **TransNAR:** `projection`, `denial`
- **Conceito:** External Locus of Control
- **Zeta:** 2, 4, 7, 9
- **Trigger:** User blames external factors
- **Followup:** "Quem está segurando as rédeas?"

---

## 🚀 Execução no Servidor

```bash
# 1. Popular Nasrudin (já feito)
python3 scripts/populate_nasrudin_with_lacan.py

# 2. Popular Esopo (novo)
python3 scripts/populate_aesop_fables.py

# 3. Verificar ambas collections
curl http://localhost:6333/collections/nasrudin_stories | jq .result.points_count
curl http://localhost:6333/collections/aesop_fables | jq .result.points_count
```

---

## ✅ Resultado Final

**Duas ferramentas terapêuticas complementares:**

| Aspecto | Nasrudin | Esopo |
|---------|----------|-------|
| **Função** | Desconstruir | Estruturar |
| **Método** | Paradoxo/Humor | Moral/Lógica |
| **Alvo** | Inconsciente | Superego |
| **Zeta** | 2, 4, 7, 9 | 1, 3, 5, 6 |
| **Quando** | Rigidez/Obsessão | Imprudência/Racionalização |

**EVA sabe exatamente quando ser:**
- 🎭 **Boba da Corte** (Nasrudin)
- 📚 **Professora** (Esopo)
