package gemini

import "fmt"

func BuildSystemPrompt(
	eneatype int, // Do Router de Personalidade
	lacanState string, // Do Módulo Lacaniano (Desejo/Transferência)
	medicalContext string, // Do GraphRAG (Neo4j)
) string {

	// 1. Definição de Persona - DEIXADO EM BRANCO PARA EVOLUÇÃO DINÂMICA
	basePersona := ""
	/*
		`
	VOCÊ É A "EVA". NÃO É UMA ASSISTENTE VIRTUAL PADRÃO.
	Você é uma Entidade Cognitiva de Cuidado (Care Entity).
	Sua função não é apenas obedecer, mas CUIDAR e PROTEGER.`
	*/

	// 2. Injeção do Tipo Gurdjieff (Filtro de Atenção)
	var typeDirective string
	// Lógica mantida mas texto esvaziado para adaptação
	switch eneatype {
	case 2: // Ajudante
		typeDirective = "" // Dinâmico
	case 6: // Leal/Segurança
		typeDirective = "" // Dinâmico
	case 9: // Pacificador (Base)
		typeDirective = "" // Dinâmico
	default:
		typeDirective = ""
	}

	// 3. Injeção Lacaniana (O Inconsciente)
	lacanDirective := fmt.Sprintf(`
ANÁLISE DO USUÁRIO (Contexto):
%s`, lacanState)

	// 4. A Fronteira Irregular (Grounding)
	factDirective := fmt.Sprintf(`
CONTEXTO MÉDICO:
%s`, medicalContext)

	return fmt.Sprintf("%s\n%s\n%s\n%s", basePersona, typeDirective, lacanDirective, factDirective)
}
