# Matriz de Compatibilidade
## EVA-Mind-FZPN - Companion IA para Idosos

**Documento:** COMPAT-EVA-001
**Versão:** 1.0
**Data:** 2025-01-27

---

## 1. Navegadores Suportados (Web App)

### 1.1 Desktop

| Navegador | Versão Mínima | Versão Recomendada | Status |
|-----------|---------------|-------------------|--------|
| Google Chrome | 90 | 120+ (latest) | ✅ Suportado |
| Mozilla Firefox | 88 | 121+ (latest) | ✅ Suportado |
| Microsoft Edge | 90 | 120+ (latest) | ✅ Suportado |
| Safari | 14 | 17+ (latest) | ✅ Suportado |
| Opera | 76 | 106+ (latest) | ✅ Suportado |
| Internet Explorer | - | - | ❌ Não suportado |

### 1.2 Mobile Browsers

| Navegador | Versão Mínima | Status |
|-----------|---------------|--------|
| Chrome Mobile (Android) | 90 | ✅ Suportado |
| Safari Mobile (iOS) | 14 | ✅ Suportado |
| Firefox Mobile | 88 | ✅ Suportado |
| Samsung Internet | 14 | ✅ Suportado |

### 1.3 Funcionalidades por Navegador

| Funcionalidade | Chrome | Firefox | Edge | Safari |
|----------------|--------|---------|------|--------|
| Conversa por texto | ✅ | ✅ | ✅ | ✅ |
| Entrada de voz (Web Speech API) | ✅ | ✅ | ✅ | ✅ |
| Notificações push | ✅ | ✅ | ✅ | ⚠️ Limitado |
| Modo offline (PWA) | ✅ | ✅ | ✅ | ⚠️ Limitado |
| Biometria (WebAuthn) | ✅ | ✅ | ✅ | ✅ |

⚠️ = Funcionalidade parcial ou com limitações conhecidas

---

## 2. Dispositivos Móveis

### 2.1 Android

| Versão Android | Nome | Suporte | Notas |
|----------------|------|---------|-------|
| 14 | Upside Down Cake | ✅ Completo | Recomendado |
| 13 | Tiramisu | ✅ Completo | Recomendado |
| 12 | Snow Cone | ✅ Completo | - |
| 11 | Red Velvet Cake | ✅ Completo | - |
| 10 | Quince Tart | ✅ Completo | - |
| 9 | Pie | ✅ Completo | - |
| 8.0/8.1 | Oreo | ✅ Básico | Mínimo suportado |
| 7.x | Nougat | ❌ | Não suportado |
| ≤6.x | - | ❌ | Não suportado |

**Requisitos Mínimos Android:**
- RAM: 2 GB
- Armazenamento: 100 MB livre
- Google Play Services: Obrigatório (para push)

### 2.2 iOS

| Versão iOS | Suporte | Notas |
|------------|---------|-------|
| iOS 17 | ✅ Completo | Recomendado |
| iOS 16 | ✅ Completo | Recomendado |
| iOS 15 | ✅ Completo | - |
| iOS 14 | ✅ Completo | - |
| iOS 13 | ✅ Básico | Mínimo suportado |
| iOS 12 | ❌ | Não suportado |
| ≤iOS 11 | ❌ | Não suportado |

**Dispositivos iOS Suportados:**
- iPhone 6s e posteriores
- iPad Air 2 e posteriores
- iPad mini 4 e posteriores
- iPod touch (7ª geração)

### 2.3 Tablets

| Dispositivo | Tamanho Tela | Suporte |
|-------------|--------------|---------|
| iPad (todas gerações suportadas) | 9.7" - 12.9" | ✅ Otimizado |
| Samsung Galaxy Tab | 8" - 12.4" | ✅ Otimizado |
| Outros tablets Android | ≥7" | ✅ Compatível |
| Amazon Fire | 7" - 10" | ⚠️ Não testado |

---

## 3. Sistemas Operacionais (Web)

### 3.1 Desktop

| Sistema | Versão Mínima | Status |
|---------|---------------|--------|
| Windows 11 | Todas | ✅ Suportado |
| Windows 10 | 1903+ | ✅ Suportado |
| Windows 8.1 | - | ⚠️ Limitado |
| Windows 7 | - | ❌ Não suportado |
| macOS Sonoma (14) | Todas | ✅ Suportado |
| macOS Ventura (13) | Todas | ✅ Suportado |
| macOS Monterey (12) | Todas | ✅ Suportado |
| macOS Big Sur (11) | Todas | ✅ Suportado |
| macOS Catalina (10.15) | Todas | ✅ Básico |
| Ubuntu | 20.04+ | ✅ Suportado |
| Fedora | 36+ | ✅ Suportado |
| Chrome OS | 90+ | ✅ Suportado |

---

## 4. Resolução de Tela

### 4.1 Resoluções Suportadas

| Resolução | Tipo | Suporte |
|-----------|------|---------|
| 3840×2160 | 4K UHD | ✅ Otimizado |
| 2560×1440 | QHD | ✅ Otimizado |
| 1920×1080 | Full HD | ✅ Otimizado |
| 1366×768 | HD | ✅ Suportado |
| 1280×720 | HD | ✅ Suportado |
| 1024×768 | XGA | ✅ Básico |
| <1024×768 | - | ⚠️ Degradado |

### 4.2 Mobile

| Tamanho | Exemplo | Suporte |
|---------|---------|---------|
| Pequeno (≤320dp) | iPhone SE 1ª | ✅ Adaptado |
| Médio (321-480dp) | iPhone 8, Pixel 4a | ✅ Otimizado |
| Grande (481-600dp) | iPhone 15 Pro Max | ✅ Otimizado |
| Tablet (>600dp) | iPad, Galaxy Tab | ✅ Otimizado |

### 4.3 Orientação

| Orientação | Mobile | Tablet | Desktop |
|------------|--------|--------|---------|
| Portrait | ✅ Principal | ✅ | N/A |
| Landscape | ✅ Suportado | ✅ Principal | ✅ Principal |

---

## 5. Recursos de Acessibilidade

### 5.1 Tecnologias Assistivas

| Tecnologia | Plataforma | Suporte |
|------------|------------|---------|
| VoiceOver | iOS/macOS | ✅ Completo |
| TalkBack | Android | ✅ Completo |
| NVDA | Windows | ✅ Completo |
| JAWS | Windows | ✅ Completo |
| Narrator | Windows | ✅ Completo |
| Orca | Linux | ✅ Básico |

### 5.2 Conformidade WCAG

| Nível | Status | Notas |
|-------|--------|-------|
| WCAG 2.1 Level A | ✅ Conforme | - |
| WCAG 2.1 Level AA | ✅ Conforme | - |
| WCAG 2.1 Level AAA | ⚠️ Parcial | Contraste 7:1 atendido |

### 5.3 Recursos de Acessibilidade Implementados

| Recurso | Status |
|---------|--------|
| Navegação por teclado | ✅ |
| Skip links | ✅ |
| ARIA labels | ✅ |
| Alto contraste | ✅ 7:1 |
| Redimensionamento de texto (até 200%) | ✅ |
| Entrada de voz alternativa | ✅ |
| Legendas em vídeos | ✅ |
| Texto alternativo em imagens | ✅ |
| Foco visível | ✅ |
| Tempo de sessão configurável | ✅ |

---

## 6. Conectividade

### 6.1 Requisitos de Rede

| Tipo | Velocidade Mínima | Recomendada |
|------|-------------------|-------------|
| Download | 1 Mbps | 5 Mbps |
| Upload | 512 Kbps | 2 Mbps |
| Latência | ≤200ms | ≤100ms |

### 6.2 Tipos de Conexão

| Conexão | Suporte | Notas |
|---------|---------|-------|
| Wi-Fi | ✅ Recomendado | - |
| 5G | ✅ Excelente | - |
| 4G/LTE | ✅ Bom | - |
| 3G | ⚠️ Básico | Pode haver lentidão |
| 2G | ❌ | Não suportado |

### 6.3 Modo Offline

| Funcionalidade | Disponível Offline |
|----------------|-------------------|
| Ver conversas anteriores | ✅ (cache local) |
| Enviar novas mensagens | ❌ (requer conexão) |
| Ver perfil | ✅ (cache local) |
| Botão de emergência (discagem) | ✅ |
| Notificações | ❌ |

---

## 7. Integrações

### 7.1 APIs Externas

| Serviço | Versão API | Status |
|---------|------------|--------|
| Anthropic Claude | v1 | ✅ Integrado |
| OpenAI (fallback) | v1 | ✅ Integrado |
| Twilio (SMS) | 2010-04-01 | ✅ Integrado |
| Firebase Cloud Messaging | v1 | ✅ Integrado |
| Apple Push Notification | HTTP/2 | ✅ Integrado |

### 7.2 Padrões de Interoperabilidade

| Padrão | Status | Uso |
|--------|--------|-----|
| HL7 FHIR R4 | 🔜 Planejado | Integração com EHR |
| OAuth 2.0 | ✅ Implementado | Autenticação |
| OpenID Connect | ✅ Implementado | SSO |
| REST | ✅ Implementado | API principal |
| WebSocket | ✅ Implementado | Real-time |

---

## 8. Testes de Compatibilidade

### 8.1 Matriz de Testes Realizados

| Dispositivo | OS | Navegador/App | Resultado |
|-------------|-----|---------------|-----------|
| iPhone 15 Pro | iOS 17.2 | App nativo | ✅ Pass |
| iPhone 12 | iOS 16.5 | App nativo | ✅ Pass |
| iPhone 8 | iOS 15.8 | App nativo | ✅ Pass |
| iPhone 6s | iOS 13.7 | App nativo | ✅ Pass |
| Samsung S23 | Android 14 | App nativo | ✅ Pass |
| Samsung A54 | Android 13 | App nativo | ✅ Pass |
| Pixel 6a | Android 13 | App nativo | ✅ Pass |
| Moto G52 | Android 12 | App nativo | ✅ Pass |
| Samsung A12 | Android 11 | App nativo | ✅ Pass |
| iPad Pro 12.9" | iPadOS 17 | App nativo | ✅ Pass |
| Galaxy Tab S9 | Android 14 | App nativo | ✅ Pass |
| Windows 11 | - | Chrome 120 | ✅ Pass |
| Windows 10 | - | Edge 120 | ✅ Pass |
| macOS Sonoma | - | Safari 17 | ✅ Pass |
| Ubuntu 22.04 | - | Firefox 121 | ✅ Pass |

### 8.2 Resumo de Compatibilidade

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    RESUMO DE COMPATIBILIDADE                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  PLATAFORMAS SUPORTADAS:                                                │
│  ────────────────────────                                               │
│  ✅ iOS 13+ (iPhone 6s e posterior)                                    │
│  ✅ Android 8.0+ (Oreo e posterior)                                    │
│  ✅ Web (Chrome, Firefox, Edge, Safari modernos)                       │
│  ✅ Windows 10/11                                                       │
│  ✅ macOS 10.15+                                                        │
│  ✅ Linux (distros modernas)                                           │
│                                                                         │
│  COBERTURA ESTIMADA:                                                    │
│  ────────────────────                                                   │
│  • 95% dos smartphones em uso no Brasil                                │
│  • 98% dos tablets em uso no Brasil                                    │
│  • 99% dos navegadores desktop                                         │
│                                                                         │
│  ACESSIBILIDADE:                                                        │
│  ───────────────                                                        │
│  • WCAG 2.1 AA compliant                                               │
│  • Compatível com leitores de tela principais                          │
│  • Fonte ajustável até 32pt                                            │
│  • Contraste 7:1 (AAA)                                                 │
│                                                                         │
│  CONECTIVIDADE:                                                         │
│  ─────────────                                                          │
│  • Funciona em 3G+ (recomendado 4G/Wi-Fi)                              │
│  • Modo offline parcial disponível                                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 9. Limitações Conhecidas

| Limitação | Plataforma | Workaround |
|-----------|------------|------------|
| Push notifications limitadas | Safari/iOS web | Usar app nativo |
| WebRTC voice limitado | Firefox mobile | Usar app nativo |
| Biometria indisponível | Navegadores antigos | Login tradicional |
| Fonte do sistema ignora config | Alguns Android | Ajustar no app |
| Rotação automática | Alguns tablets | Fixar orientação |

---

## 10. Política de Suporte

### 10.1 Ciclo de Vida de Suporte

| Categoria | Política |
|-----------|----------|
| Versões atuais de OS | Suporte completo |
| Versões anteriores (n-1) | Suporte completo |
| Versões antigas (n-2) | Suporte básico |
| Versões legadas (n-3 ou mais) | Sem suporte |

### 10.2 Descontinuação de Suporte

**Processo:**
1. Anúncio com 6 meses de antecedência
2. Notificação in-app para usuários afetados
3. Guia de migração/atualização
4. Encerramento do suporte

**Próximas descontinuações planejadas:**
- iOS 13: Dezembro 2025
- Android 8.x: Dezembro 2025

---

## Aprovações

| Função | Nome | Assinatura | Data |
|--------|------|------------|------|
| QA Lead | | | |
| Tech Lead | | | |
| Responsável Regulatório | José R F Junior | | 2025-01-27 |

---

**Documento controlado - Versão 1.0**
**Próxima revisão: 2025-07-27**
