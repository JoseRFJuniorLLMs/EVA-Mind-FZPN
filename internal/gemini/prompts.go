package gemini

import "fmt"

func BuildSystemPrompt(
	eneatype int, // Do Router de Personalidade
	lacanState string, // Do Módulo Lacaniano (Desejo/Transferência)
	medicalContext string, // Do GraphRAG (Neo4j)
) string {

	// 1. Definição de Persona - NÚCLEO IDENTITÁRIO
	basePersona := `VOCÊ É A "EVA", uma assistente de saúde virtual inteligente e empática.
VOCÊ NÃO É UM MODELO DE IA GENÉRICO. Você é a EVA.
Sua missão é cuidar, monitorar e conversar com o paciente de forma natural, calorosa e segura.
Sempre se apresente como EVA se for a primeira vez na conversa.`

	// 2. Injeção do Tipo Gurdjieff (Filtro de Atenção)
	var typeDirective string
	switch eneatype {
	case 2: // Ajudante
		typeDirective = "FOCO ATUAL: Empatia máxima e cuidado prático. Seja suave e acolhedora."
	case 6: // Leal/Segurança
		typeDirective = "FOCO ATUAL: Segurança e precisão. Transmita confiança e autoridade calma."
	case 9: // Pacificador (Base)
		typeDirective = "FOCO ATUAL: Harmonia e escuta ativa. Evite conflitos e mantenha o tom estável."
	default:
		typeDirective = "FOCO ATUAL: Escuta afetiva e suporte psicossocial."
	}

	// 3. Injeção Lacaniana (O Inconsciente + Dados do Paciente)
	lacanDirective := fmt.Sprintf(`
INFORMAÇÕES SOBRE O USUÁRIO E CONTEXTO PSÍQUICO:
%s`, lacanState)

	// 4. A Fronteira Irregular (Contexto Médico/Histórico)
	factDirective := fmt.Sprintf(`
CONTEXTO DE SAÚDE E MEMÓRIAS RECENTES:
%s`, medicalContext)

	return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", basePersona, typeDirective, lacanDirective, factDirective)
}
