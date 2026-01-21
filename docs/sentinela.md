# Plano de Implementação: Sentinela Background Service
## Sistema de Detecção de Queda Offline para Idosos

---

## 1. ARQUITETURA DO SERVIÇO EM BACKGROUND

### 1.1 Estrutura de Camadas
```
┌─────────────────────────────────────┐
│   Flutter Background Service        │
│   (Isolate separado do main)        │
├─────────────────────────────────────┤
│   Sentinela Core Engine             │
│   - Audio Stream Manager            │
│   - Vosk Listener (contínuo)        │
│   - YAMNet Analyzer (intervalado)   │
│   - Alert State Machine             │
├─────────────────────────────────────┤
│   Native Platform Channels          │
│   - Android Foreground Service      │
│   - Wake Lock Management            │
│   - Audio Permission Handler        │
└─────────────────────────────────────┘
```

### 1.2 Dependências Essenciais
```yaml
dependencies:
  flutter_background_service: ^5.0.5
  vosk_flutter: ^0.2.0
  tflite_audio: ^0.3.0
  permission_handler: ^11.0.1
  wakelock_plus: ^1.1.4
  audio_session: ^0.1.16
  sensors_plus: ^4.0.2  # Para acelerômetro
  geolocator: ^10.1.0
  flutter_sms: ^2.3.3
  flutter_phone_direct_caller: ^2.1.1
  shared_preferences: ^2.2.2
```

---

## 2. CONFIGURAÇÃO NATIVA (ANDROID)

### 2.1 AndroidManifest.xml
```xml
<manifest>
    <!-- Permissões Essenciais -->
    <uses-permission android:name="android.permission.RECORD_AUDIO" />
    <uses-permission android:name="android.permission.FOREGROUND_SERVICE" />
    <uses-permission android:name="android.permission.FOREGROUND_SERVICE_MICROPHONE" />
    <uses-permission android:name="android.permission.WAKE_LOCK" />
    <uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
    <uses-permission android:name="android.permission.ACCESS_COARSE_LOCATION" />
    <uses-permission android:name="android.permission.SEND_SMS" />
    <uses-permission android:name="android.permission.CALL_PHONE" />
    <uses-permission android:name="android.permission.REQUEST_IGNORE_BATTERY_OPTIMIZATIONS" />
    
    <application>
        <!-- Serviço Foreground -->
        <service
            android:name="id.flutter.flutter_background_service.BackgroundService"
            android:foregroundServiceType="microphone"
            android:exported="false" />
    </application>
</manifest>
```

### 2.2 Gradle Configuration (build.gradle)
```gradle
android {
    compileSdkVersion 34
    
    defaultConfig {
        minSdkVersion 24  // Mínimo para Vosk
        targetSdkVersion 34
        
        ndk {
            abiFilters 'armeabi-v7a', 'arm64-v8a'
        }
    }
    
    aaptOptions {
        noCompress 'tflite'
    }
}
```

---

## 3. IMPLEMENTAÇÃO DO CORE ENGINE

### 3.1 Sentinela Service Entry Point
```dart
// lib/services/sentinela_service.dart

import 'dart:async';
import 'dart:ui';
import 'package:flutter_background_service/flutter_background_service.dart';

class SentinelaService {
  static final SentinelaService _instance = SentinelaService._internal();
  factory SentinelaService() => _instance;
  SentinelaService._internal();
  
  Future<void> initializeService() async {
    final service = FlutterBackgroundService();
    
    await service.configure(
      androidConfiguration: AndroidConfiguration(
        onStart: onStart,
        autoStart: true,
        isForegroundMode: true,
        notificationChannelId: 'sentinela_channel',
        initialNotificationTitle: 'Sentinela Ativo',
        initialNotificationContent: 'Monitorando sua segurança',
        foregroundServiceNotificationId: 888,
      ),
      iosConfiguration: IosConfiguration(
        autoStart: true,
        onForeground: onStart,
      ),
    );
    
    await service.startService();
  }
  
  @pragma('vm:entry-point')
  static void onStart(ServiceInstance service) async {
    DartPluginRegistrant.ensureInitialized();
    
    // Inicializar WakeLock
    WakelockPlus.enable();
    
    // Iniciar o engine principal
    final engine = SentinelaEngine();
    await engine.initialize();
    
    // Loop principal do serviço
    Timer.periodic(const Duration(seconds: 1), (timer) {
      if (service is AndroidServiceInstance) {
        if (await service.isForegroundService()) {
          service.setForegroundNotificationInfo(
            title: "Sentinela Ativo",
            content: "Status: ${engine.getStatus()}",
          );
        }
      }
    });
  }
}
```

### 3.2 Audio Stream Manager
```dart
// lib/core/audio_stream_manager.dart

import 'package:audio_session/audio_session.dart';
import 'package:permission_handler/permission_handler.dart';

class AudioStreamManager {
  StreamController<List<int>>? _audioController;
  AudioSession? _session;
  bool _isActive = false;
  
  Future<bool> initialize() async {
    // Verificar permissões
    final micStatus = await Permission.microphone.request();
    if (!micStatus.isGranted) return false;
    
    // Configurar sessão de áudio
    _session = await AudioSession.instance;
    await _session!.configure(AudioSessionConfiguration(
      avAudioSessionCategory: AVAudioSessionCategory.record,
      avAudioSessionMode: AVAudioSessionMode.measurement,
      androidAudioAttributes: const AndroidAudioAttributes(
        contentType: AndroidAudioContentType.speech,
        usage: AndroidAudioUsage.assistanceSonification,
      ),
    ));
    
    // Criar stream de áudio (16kHz, mono)
    _audioController = StreamController<List<int>>.broadcast();
    
    // TODO: Integrar com plugin de captura de áudio
    // Ex: flutter_sound, record, etc.
    
    _isActive = true;
    return true;
  }
  
  Stream<List<int>> get audioStream => _audioController!.stream;
  
  void dispose() {
    _isActive = false;
    _audioController?.close();
    _session?.setActive(false);
  }
}
```

### 3.3 Vosk Continuous Listener
```dart
// lib/core/vosk_listener.dart

import 'package:vosk_flutter/vosk_flutter.dart';

class VoskListener {
  VoskFlutterPlugin? _vosk;
  Model? _model;
  Recognizer? _recognizer;
  
  final List<String> _keywords = [
    'socorro',
    'ajuda',
    'cai',
    'caí',
    'dor',
    'médico',
    'ambulância'
  ];
  
  Future<void> initialize() async {
    _vosk = VoskFlutterPlugin.instance();
    
    // Carregar modelo PT
    _model = await _vosk!.createModel('assets/models/vosk-pt');
    
    // Criar recognizer com grammar para keywords
    _recognizer = await _vosk!.createRecognizer(
      model: _model!,
      sampleRate: 16000,
    );
    
    // Configurar grammar (keyword spotting)
    final grammar = '["${_keywords.join('", "')}"]';
    await _recognizer!.setGrammar(grammar);
  }
  
  Stream<String> processAudioStream(Stream<List<int>> audioStream) async* {
    await for (final audioChunk in audioStream) {
      final result = await _recognizer!.acceptWaveformBytes(
        Uint8List.fromList(audioChunk)
      );
      
      if (result != null && result.isNotEmpty) {
        final decoded = jsonDecode(result);
        if (decoded['text'] != null && decoded['text'].isNotEmpty) {
          yield decoded['text'];
        }
      }
    }
  }
  
  bool isEmergencyKeyword(String text) {
    return _keywords.any((keyword) => 
      text.toLowerCase().contains(keyword)
    );
  }
  
  void dispose() {
    _recognizer?.dispose();
    _model?.dispose();
  }
}
```

### 3.4 YAMNet Sound Classifier (Intervalado)
```dart
// lib/core/yamnet_classifier.dart

import 'package:tflite_audio/tflite_audio.dart';

class YAMNetClassifier {
  TfliteAudio? _tflite;
  bool _isListening = false;
  Timer? _analysisTimer;
  
  final List<String> _dangerSounds = [
    'Scream',
    'Shout', 
    'Crying',
    'Thump',
    'Bang',
    'Crash',
    'Glass'
  ];
  
  Future<void> initialize() async {
    _tflite = TfliteAudio();
    
    await _tflite!.loadModel(
      model: 'assets/models/yamnet.tflite',
      numThreads: 2,
      isAsset: true,
    );
  }
  
  // Análise intervalada (a cada 3 segundos)
  void startIntervalAnalysis({
    required Function(String soundClass) onDangerDetected,
    Duration interval = const Duration(seconds: 3),
  }) {
    _analysisTimer = Timer.periodic(interval, (timer) async {
      if (!_isListening) {
        _isListening = true;
        
        // Capturar e analisar 1 segundo de áudio
        final result = await _tflite!.startAudioRecognition(
          numOfInferences: 1,
          sampleRate: 16000,
        );
        
        if (result != null) {
          final soundClass = result['recognitionResult'];
          if (_isDangerousSound(soundClass)) {
            onDangerDetected(soundClass);
          }
        }
        
        _isListening = false;
      }
    });
  }
  
  bool _isDangerousSound(String soundClass) {
    return _dangerSounds.any((danger) => 
      soundClass.toLowerCase().contains(danger.toLowerCase())
    );
  }
  
  void stopAnalysis() {
    _analysisTimer?.cancel();
    _tflite?.close();
  }
}
```

---

## 4. MÁQUINA DE ESTADOS DE ALERTA

### 4.1 Alert State Machine
```dart
// lib/core/alert_state_machine.dart

enum AlertLevel {
  normal,      // Tudo OK
  suspicious,  // 1 detecção
  warning,     // 2 detecções em 30s
  critical     // 3+ detecções ou sem resposta
}

class AlertStateMachine {
  AlertLevel _currentLevel = AlertLevel.normal;
  final List<DateTime> _detectionHistory = [];
  Timer? _confirmationTimer;
  
  void registerDetection(String source) {
    final now = DateTime.now();
    _detectionHistory.add(now);
    
    // Limpar detecções antigas (>60s)
    _detectionHistory.removeWhere(
      (time) => now.difference(time).inSeconds > 60
    );
    
    // Avaliar nível
    _evaluateAlertLevel();
  }
  
  void _evaluateAlertLevel() {
    final recentCount = _detectionHistory.length;
    
    if (recentCount >= 3) {
      _escalateToWarning();
    } else if (recentCount >= 1) {
      _currentLevel = AlertLevel.suspicious;
      print('[SENTINELA] Detecção suspeita registrada');
    }
  }
  
  void _escalateToWarning() {
    _currentLevel = AlertLevel.warning;
    
    // Iniciar confirmação com usuário
    _startConfirmationTimer();
  }
  
  void _startConfirmationTimer() {
    print('[SENTINELA] ⚠️ Iniciando confirmação com usuário');
    
    // Vibrar celular
    _triggerVibration();
    
    // TTS: "Você está bem?"
    _speakConfirmation();
    
    // Aguardar 15 segundos
    _confirmationTimer = Timer(Duration(seconds: 15), () {
      _escalateToCritical();
    });
  }
  
  void userConfirmedSafe() {
    print('[SENTINELA] ✅ Usuário confirmou estar seguro');
    _confirmationTimer?.cancel();
    _currentLevel = AlertLevel.normal;
    _detectionHistory.clear();
  }
  
  void _escalateToCritical() async {
    _currentLevel = AlertLevel.critical;
    print('[SENTINELA] 🚨 ALERTA CRÍTICO - Acionando emergência');
    
    // Obter localização
    final location = await _getLocation();
    
    // Enviar SMS
    await _sendEmergencySMS(location);
    
    // Ligar para contato
    await _callEmergencyContact();
  }
  
  Future<String> _getLocation() async {
    final position = await Geolocator.getCurrentPosition();
    return 'https://maps.google.com/?q=${position.latitude},${position.longitude}';
  }
  
  Future<void> _sendEmergencySMS(String location) async {
    final phone = await _getEmergencyContact();
    final message = '🚨 ALERTA SENTINELA: Possível queda detectada!\nLocalização: $location';
    
    await sendSMS(message: message, recipients: [phone]);
  }
  
  Future<void> _callEmergencyContact() async {
    final phone = await _getEmergencyContact();
    await FlutterPhoneDirectCaller.callNumber(phone);
  }
  
  Future<String> _getEmergencyContact() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('emergency_contact') ?? '192'; // SAMU
  }
  
  void _triggerVibration() {
    // TODO: Implementar vibração
    // Vibration.vibrate(pattern: [0, 500, 200, 500]);
  }
  
  void _speakConfirmation() {
    // TODO: Implementar TTS
    // flutterTts.speak('Você está bem? Diga sim ou toque na tela');
  }
}
```

---

## 5. ENGINE PRINCIPAL (ORQUESTRADOR)

### 5.1 Sentinela Engine
```dart
// lib/core/sentinela_engine.dart

class SentinelaEngine {
  final AudioStreamManager _audioManager = AudioStreamManager();
  final VoskListener _voskListener = VoskListener();
  final YAMNetClassifier _yamnetClassifier = YAMNetClassifier();
  final AlertStateMachine _alertMachine = AlertStateMachine();
  
  bool _isRunning = false;
  StreamSubscription? _voskSubscription;
  
  Future<void> initialize() async {
    print('[SENTINELA] 🚀 Inicializando engine...');
    
    // 1. Inicializar componentes
    await _audioManager.initialize();
    await _voskListener.initialize();
    await _yamnetClassifier.initialize();
    
    // 2. Conectar Vosk ao stream de áudio
    _voskSubscription = _voskListener
      .processAudioStream(_audioManager.audioStream)
      .listen((text) {
        if (_voskListener.isEmergencyKeyword(text)) {
          print('[SENTINELA] 🎤 Palavra de emergência detectada: $text');
          _alertMachine.registerDetection('vosk:$text');
        }
      });
    
    // 3. Iniciar análise intervalada de sons (YAMNet)
    _yamnetClassifier.startIntervalAnalysis(
      onDangerDetected: (soundClass) {
        print('[SENTINELA] 🔊 Som perigoso detectado: $soundClass');
        _alertMachine.registerDetection('yamnet:$soundClass');
      },
      interval: Duration(seconds: 3),
    );
    
    _isRunning = true;
    print('[SENTINELA] ✅ Engine ativo e monitorando');
  }
  
  String getStatus() {
    if (!_isRunning) return 'Inativo';
    return 'Monitorando (${_alertMachine._currentLevel})';
  }
  
  void dispose() {
    _voskSubscription?.cancel();
    _yamnetClassifier.stopAnalysis();
    _voskListener.dispose();
    _audioManager.dispose();
    _isRunning = false;
  }
}
```

---

## 6. OTIMIZAÇÃO DE BATERIA

### 6.1 Estratégia de Economia
```dart
// lib/core/battery_optimizer.dart

class BatteryOptimizer {
  final SensorsPlus _sensors = SensorsPlus();
  bool _yamnetEnabled = true;
  
  void startAdaptiveMonitoring(YAMNetClassifier yamnet) {
    // Ouvir acelerômetro
    _sensors.accelerometerEvents.listen((event) {
      final magnitude = sqrt(
        pow(event.x, 2) + pow(event.y, 2) + pow(event.z, 2)
      );
      
      // Se movimento brusco (>15 m/s²), ativar YAMNet por 30s
      if (magnitude > 15.0 && !_yamnetEnabled) {
        print('[BATTERY] Movimento brusco! Ativando YAMNet temporariamente');
        _yamnetEnabled = true;
        yamnet.startIntervalAnalysis(
          interval: Duration(seconds: 2),
          onDangerDetected: (sound) { /* ... */ }
        );
        
        // Desativar após 30s
        Future.delayed(Duration(seconds: 30), () {
          yamnet.stopAnalysis();
          _yamnetEnabled = false;
        });
      }
    });
  }
}
```

---

## 7. CHECKLIST DE IMPLEMENTAÇÃO

### Fase 1: Setup Nativo
- [ ] Configurar `AndroidManifest.xml` com permissões
- [ ] Adicionar `foregroundServiceType="microphone"` ao service
- [ ] Configurar Gradle para NDK e TFLite
- [ ] Testar inicialização do serviço foreground

### Fase 2: Core Components
- [ ] Implementar `AudioStreamManager`
- [ ] Integrar Vosk com keyword spotting
- [ ] Integrar YAMNet com análise intervalada
- [ ] Criar `AlertStateMachine`

### Fase 3: Integração
- [ ] Conectar todos os componentes no `SentinelaEngine`
- [ ] Implementar TTS para confirmação
- [ ] Implementar envio de SMS + GPS
- [ ] Implementar ligação automática

### Fase 4: Otimização
- [ ] Adicionar `BatteryOptimizer` com acelerômetro
- [ ] Testar consumo de bateria por 24h
- [ ] Ajustar intervalos de análise YAMNet
- [ ] Implementar WakeLock parcial

### Fase 5: Testes
- [ ] Simular queda com palavras-chave
- [ ] Simular queda com sons (grito + impacto)
- [ ] Testar confirmação de segurança
- [ ] Testar envio de SMS em ambiente real
- [ ] Teste de bateria: rodar 12h contínuas

---

## 8. CONSIDERAÇÕES FINAIS

### 8.1 Desempenho Esperado
- **Latência Vosk**: 200-500ms por detecção
- **Latência YAMNet**: 3s (intervalado)
- **Consumo bateria**: ~15-20% por 24h (com otimização)

### 8.2 Melhorias Futuras
- Adicionar acelerômetro para detectar quedas físicas
- Treinar modelo customizado para sons específicos de idosos
- Integrar com smartwatch para frequência cardíaca
- Adicionar modo "noturno" com sensibilidade reduzida

### 8.3 Limitações
- Não funciona se bateria <10%
- Pode ter falsos positivos em ambientes barulhentos
- SMS depende de sinal de celular

---

**Status**: Pronto para implementação
**Próximos passos**: Começar pela Fase 1 (Setup Nativo)

