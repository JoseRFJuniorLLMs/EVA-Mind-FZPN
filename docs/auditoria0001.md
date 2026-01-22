Auditoria Profunda Backend: EVA-Mind-FZPN
📊 Ciclo 1: Diagnóstico de Instabilidade (Logs)
Os logs de produção revelaram um padrão de desconexão crítica:

Erro: websocket: close 1011 (internal server error): Deadline expired before operation could complete.
Causa Raiz: O Google Gemini Live API possui um timeout de ociosidade de aproximadamente 60 segundos. Se o idoso ficar em silêncio e o sistema não enviar "heartbeats" ou áudio, a Google encerra a sessão.
Consequência: A IA para de responder subitamente durante a chamada.
📊 Ciclo 2: Auditoria Estrutural e Redundância
Encontramos uma inconsistência grave na base de código:

Código Ativo: 
main.go
 (1800+ linhas) contém toda a lógica de sinalização, gerenciamento de clientes e integração com Gemini.
Código Fantasma: O diretório 
internal/senses/signaling/websocket.go
 contém uma implementação quase idêntica (1500 linhas) que não está sendo usada pelo 
main.go
, mas que possui parâmetros diferentes (ex: ReadDeadline de 60s em vez de 5min).
Dívida Técnica: A lógica está duplicada e o 
main.go
 tornou-se um "God-file" difícil de manter.
🚀 Plano de Ação: Estabilização e Limpeza
Fase 1: Estabilização de Conexão (P0)
Heartbeat Silencioso: Implementar o envio de frames de áudio "silenciosos" para o Gemini a cada 30 segundos de inatividade para evitar o erro 1011.
Sincronia de Timeouts: Alinhar o 
monitorClientActivity
 do 
main.go
 (atualmente 5min) com a sensibilidade do WebSocket (60s).
Fase 2: Unificação Arquitetural (P1)
Refatoração: Mover a lógica de sinalização do 
main.go
 para um pacote dedicado (internal/brainstem/signaling).
Deleção: Remover o diretório redundante internal/senses/signaling para evitar confusão.
Fase 3: Observabilidade
Logs de Contexto: Adicionar logs que diferenciem claramente entre "Desconexão do Usuário", "Timeout de Inatividade" e "Erro da Google".
