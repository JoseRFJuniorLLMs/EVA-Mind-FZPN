# 🌊 Deep Resonance Engine (DRE)

**EVA-Mind-FZPN - Módulo de Ressonância Profunda**  
**Technical Specification & Implementation Guide v1.0**  
**Classificação:** Módulo de Intervenção Nível 4 (Subconsciente)

---

## 📋 Visão Geral

O **Deep Resonance Engine** é um sistema de sugestão indireta baseado em **Hipnose Ericksoniana** e **PNL**. Seu objetivo é bypassar a resistência crítica do usuário (lógica/teimosia) e entregar comandos terapêuticos de relaxamento, alívio de dor e sono diretamente ao inconsciente.

### **Diferença Chave:**

| Aspecto | EVA Normal | EVA Ressonância |
|---------|------------|-----------------|
| **Velocidade** | 1.0x (normal) | 0.75-0.85x (mais lento) |
| **Tom** | Empático, neutro | Grave, hipnótico (-2 a -3 semitones) |
| **Pausas** | Naturais (500ms) | Longas (2-6 segundos) |
| **Foco** | Conteúdo lógico | Prosódia e metáfora |
| **Visual** | Avatar falante | Gradiente abstrato lento |
| **Áudio de fundo** | Nenhum | Binaural/Ambiente (30% volume) |

---

## 🎯 Quando Usar

### **Gatilhos (Triggers):**

O sistema não ativa aleatoriamente. Requer detecção de **"Travamento Lógico"**:

#### **Cenário A: Via Texto/Voz**
- Usuário repete: *"Dói muito, nada funciona, o remédio não faz efeito"*
- **TransNAR detecta:** Loop de dor psicossomática

#### **Cenário B: Via Sensores**
- Frequência cardíaca alta às 23:00
- Acelerômetro: Usuário deitado/parado
- **Sistema detecta:** Insônia

#### **Cenário C: Via Pedido Direto**
- Usuário: *"EVA, me ajuda a dormir"*
- **Sistema ativa:** Script de sono profundo

---

## 🗄️ Arquitetura de Dados

### **1. Qdrant (Vector DB)**

**Collection:** `resonance_scripts`

#### **Schema do Objeto:**

```json
{
  "id": "dre_005_pain_melting_lake",
  "vector": [0.15, -0.88, ...],
  "payload": {
    "title": "O Lago de Cera Derretida",
    "category": "pain_management",
    "target_symptom": ["chronic_pain", "muscle_tension", "joint_stiffness"],
    "technique": "ericksonian_metaphor",
    "duration_seconds": 300,
    
    "voice_settings": {
      "speaking_rate": 0.85,
      "pitch": -2.0,
      "volume_gain_db": 0.0,
      "post_processing": "reverb_light"
    },
    
    "background_track_id": "bg_theta_waves_rain",
    "background_volume": 0.3,
    
    "script_ssml": "<speak>José... <break time='1500ms'/> eu não sei se você consegue imaginar... agora... <break time='1000ms'/> como seria sentir o calor do sol... <break time='2000ms'/> tocando suavemente essa tensão... derretendo como cera morna... <break time='3000ms'/> escorrendo para longe...</speak>",
    
    "safety_level": "green",
    "zeta_affinity": [1, 3, 5]
  }
}
```

---

### **2. PostgreSQL (Controle e Logs)**

```sql
CREATE TABLE resonance_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    script_id VARCHAR(50),
    started_at TIMESTAMP DEFAULT NOW(),
    completed BOOLEAN DEFAULT FALSE,
    user_feedback_score INT,
    biometrics_start JSONB,
    biometrics_end JSONB,
    interruption_reason VARCHAR(100)
);

CREATE INDEX idx_resonance_user ON resonance_sessions(user_id);
CREATE INDEX idx_resonance_script ON resonance_sessions(script_id);
```

---

## 🔄 Fluxo de Funcionamento

### **Passo 1: Detecção (Trigger)**
```
TransNAR detecta → "Loop de dor psicossomática"
        ↓
Sensores confirmam → Usuário deitado, FC alta
        ↓
Sistema decide → Ativar DRE
```

---

### **Passo 2: Safety Check**

```go
func (d *DRE) CanActivate(userID string) (bool, string) {
    // 1. Verificar acelerômetro
    if user.IsMoving() {
        return false, "Usuário em movimento - BLOQUEADO"
    }
    
    // 2. Verificar histórico
    lastSession := db.GetLastResonanceSession(userID)
    if time.Since(lastSession) < 2*time.Hour {
        return false, "Intervalo mínimo não respeitado"
    }
    
    // 3. Verificar contraindicações
    if user.HasEpilepsy {
        return false, "Epilepsia - usar método alternativo"
    }
    
    return true, "Seguro para ativar"
}
```

---

### **Passo 3: Seleção do Script**

```go
func (d *DRE) SelectScript(symptom string, zetaType int) (*ResonanceScript, error) {
    // 1. Gerar embedding do sintoma
    vector := d.embedder.Generate(symptom)
    
    // 2. Buscar no Qdrant
    results := d.qdrant.Search("resonance_scripts", vector, filters: {
        "zeta_affinity": zetaType,
        "safety_level": "green"
    })
    
    // 3. Retornar melhor match
    return results[0], nil
}
```

---

### **Passo 4: Geração de Áudio (Mixagem)**

#### **Backend (Go):**
```go
func (d *DRE) GenerateAudio(script *ResonanceScript) (*AudioBundle, error) {
    // 1. Gerar TTS com SSML
    voiceURL := d.tts.GenerateWithSSML(
        text: script.ScriptSSML,
        rate: script.VoiceSettings.SpeakingRate,
        pitch: script.VoiceSettings.Pitch,
    )
    
    // 2. Buscar áudio de fundo
    bgURL := d.audioLibrary.Get(script.BackgroundTrackID)
    
    // 3. Retornar bundle
    return &AudioBundle{
        VoiceURL: voiceURL,
        BackgroundURL: bgURL,
        BackgroundVolume: script.BackgroundVolume,
    }, nil
}
```

---

#### **Frontend (Flutter):**
```dart
class DualAudioPlayer {
  final AudioPlayer voicePlayer = AudioPlayer();
  final AudioPlayer bgPlayer = AudioPlayer();
  
  Future<void> playResonance(AudioBundle bundle) async {
    // 1. Iniciar áudio de fundo (loop)
    await bgPlayer.setUrl(bundle.backgroundURL);
    await bgPlayer.setVolume(bundle.backgroundVolume);
    await bgPlayer.setLoopMode(LoopMode.one);
    await bgPlayer.play();
    
    // 2. Tocar voz da EVA
    await voicePlayer.setUrl(bundle.voiceURL);
    await voicePlayer.play();
    
    // 3. Quando voz terminar, fade out do fundo
    voicePlayer.playerStateStream.listen((state) {
      if (state.processingState == ProcessingState.completed) {
        _fadeOutBackground();
      }
    });
  }
  
  Future<void> _fadeOutBackground() async {
    for (double vol = 0.3; vol >= 0; vol -= 0.01) {
      await bgPlayer.setVolume(vol);
      await Future.delayed(Duration(milliseconds: 100));
    }
    await bgPlayer.stop();
  }
}
```

---

### **Passo 5: Visual Hipnótico**

```dart
class ResonanceVisualizer extends StatefulWidget {
  @override
  _ResonanceVisualizerState createState() => _ResonanceVisualizerState();
}

class _ResonanceVisualizerState extends State<ResonanceVisualizer>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  
  @override
  void initState() {
    super.initState();
    
    // Animação MUITO lenta (30 segundos por ciclo)
    _controller = AnimationController(
      duration: Duration(seconds: 30),
      vsync: this,
    )..repeat(reverse: true);
    
    // Diminuir brilho da tela
    ScreenBrightness().setScreenBrightness(0.2);
  }
  
  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _controller,
      builder: (context, child) {
        return Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [
                Color.lerp(Colors.deepPurple[900], Colors.indigo[900], _controller.value)!,
                Color.lerp(Colors.indigo[900], Colors.blue[900], _controller.value)!,
              ],
            ),
          ),
          child: Center(
            child: Text(
              'Apenas respire...',
              style: TextStyle(
                color: Colors.white.withOpacity(0.3),
                fontSize: 18,
                fontWeight: FontWeight.w200,
              ),
            ),
          ),
        );
      },
    );
  }
  
  @override
  void dispose() {
    _controller.dispose();
    ScreenBrightness().resetScreenBrightness();
    super.dispose();
  }
}
```

---

## 🔒 Protocolos de Segurança

### **1. Bloqueio por Movimento**

```dart
class MovementMonitor {
  final AccelerometerEvents accelerometer;
  
  Stream<bool> get isMoving {
    return accelerometer.listen((event) {
      double magnitude = sqrt(
        event.x * event.x + 
        event.y * event.y + 
        event.z * event.z
      );
      
      // Se magnitude > 12 m/s², usuário está andando
      return magnitude > 12.0;
    });
  }
}

// No CallScreen
_movementMonitor.isMoving.listen((moving) {
  if (moving && _resonanceActive) {
    _stopResonance();
    _showAlert("José, este exercício exige que você pare. Vamos continuar quando sentar?");
  }
});
```

---

### **2. Wake-Up Protocols**

#### **Standard Wake Up:**
```
"Contando de 1 a 5... 
1... sentindo energia voltando... 
2... movendo os dedos... 
3... respirando fundo... 
4... abrindo os olhos devagar... 
5... totalmente desperto e alerta."
```

#### **Sleep Mode (Sem Wake Up):**
```
"... e você pode continuar dormindo... 
profundamente... 
até a manhã..."
[fade out completo]
```

---

## 📦 Scripts Prontos (Seed Data)

### **Arquivo 1: Core Scripts** (`resonance_scripts_core.json`)

10 scripts essenciais:
1. **O Trem Antigo** (Insônia)
2. **A Luva de Anestesia** (Dor crônica)
3. **O Museu das Nuvens** (Ansiedade)
4. **A Árvore Solar** (Letargia/Depressão)
5. **A Fogueira Ancestral** (Solidão)
6. **O Espelho da Alma** (Autoestima)
7. **A Biblioteca da Vida** (Culpa/Arrependimento)
8. **O Mel Dourado** (Tensão muscular)
9. **O Jardineiro** (Medo do futuro)
10. **O Barco Sem Cordas** (Insônia profunda)

---

## 🚀 Implementação

### **Checklist Backend (Go):**

- [ ] Criar `pkg/resonance/engine.go`
- [ ] Implementar `SafetyChecker`
- [ ] Implementar `ScriptSelector` (Qdrant search)
- [ ] Criar endpoint `/intervention/resonance`
- [ ] Integrar com TTS (SSML support)
- [ ] Criar serviço de mixagem de áudio

---

### **Checklist Frontend (Flutter):**

- [ ] Criar `lib/widgets/resonance_visualizer.dart`
- [ ] Implementar `DualAudioPlayer`
- [ ] Adicionar `MovementMonitor` (acelerômetro)
- [ ] Integrar em `call_screen.dart`
- [ ] Implementar controle de brilho de tela
- [ ] Adicionar wake-up alerts

---

### **Checklist Database:**

- [ ] Criar collection `resonance_scripts` no Qdrant
- [ ] Criar tabela `resonance_sessions` no PostgreSQL
- [ ] Popular com 10 scripts core
- [ ] Gerar embeddings para cada script
- [ ] Testar busca semântica

---

## 📊 Métricas de Sucesso

### **Como Medir Eficácia:**

```sql
-- Taxa de conclusão
SELECT 
    script_id,
    COUNT(*) as total_sessions,
    SUM(CASE WHEN completed THEN 1 ELSE 0 END) as completed,
    AVG(user_feedback_score) as avg_score
FROM resonance_sessions
GROUP BY script_id
ORDER BY avg_score DESC;
```

### **Biometria (Se disponível):**

```sql
-- Redução de FC média
SELECT 
    AVG((biometrics_start->>'heart_rate')::int - 
        (biometrics_end->>'heart_rate')::int) as avg_hr_reduction
FROM resonance_sessions
WHERE completed = true;
```

---

## 🎓 Fundamentos Teóricos

### **Técnicas Ericksonianas Usadas:**

1. **Pacing & Leading:** Validar realidade atual → Guiar para novo estado
2. **Embedded Commands:** Comandos escondidos em frases complexas
3. **Metaphor:** Histórias que o inconsciente decodifica
4. **Confusion:** Sobrecarregar a mente consciente para acessar o inconsciente
5. **Time Distortion:** Fazer 5 minutos parecerem 30 (ou vice-versa)

### **Referências:**

- **"Patterns of the Hypnotic Techniques of Milton H. Erickson"** (Bandler & Grinder)
- **"Trance-formations"** (Bandler & Grinder)
- **"My Voice Will Go With You"** (Sidney Rosen)

---

## ⚠️ Contraindicações

**NÃO usar DRE se:**
- ❌ Usuário tem epilepsia
- ❌ Usuário está dirigindo/andando
- ❌ Usuário tem histórico de psicose
- ❌ Menos de 2h desde última sessão

**Fallback:** Usar Zen Breathing (mais seguro)

---

## ✅ Resultado Final

**EVA agora tem 5 camadas de intervenção:**

1. 🎭 **Nasrudin** → Paradoxo/Humor (Consciente)
2. 📚 **Esopo** → Moral/Lógica (Superego)
3. 🧘 **Zen** → Insight/Silêncio (Mente)
4. 🫁 **Somático** → Aterramento (Corpo)
5. 🌊 **Ressonância** → Hipnose (Subconsciente)

**Sistema completo de intervenção psicofisiológica multinível!** ✨
