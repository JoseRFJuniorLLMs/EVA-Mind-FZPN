package tools

// Schema definition for Gemini Function Calling
type FunctionDeclaration struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Parameters  *FunctionParameters `json:"parameters"`
}

type FunctionParameters struct {
	Type       string               `json:"type"` // "OBJECT"
	Properties map[string]*Property `json:"properties"`
	Required   []string             `json:"required"`
}

type Property struct {
	Type        string   `json:"type"` // "STRING", "INTEGER", "BOOLEAN"
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// GetVitalsDefinition returns the schema for the GetVitals tool
func GetVitalsDefinition() FunctionDeclaration {
	return FunctionDeclaration{
		Name:        "get_vitals",
		Description: "Recupera os sinais vitais mais recentes do idoso (pressão arterial, glicose, batimentos cardíacos, peso, saturação). Use para verificar o estado de saúde física atual ou histórico recente.",
		Parameters: &FunctionParameters{
			Type: "OBJECT",
			Properties: map[string]*Property{
				"vitals_type": {
					Type:        "STRING",
					Description: "O tipo de sinal vital a ser buscado. Exemplos: 'pressao_arterial', 'glicemia', 'batimentos', 'saturacao_o2', 'peso', 'temperatura'. Se vazio, tenta buscar um resumo geral.",
					Enum:        []string{"pressao_arterial", "glicemia", "batimentos", "saturacao_o2", "peso", "temperatura"},
				},
				"limit": {
					Type:        "INTEGER",
					Description: "Número máximo de registros a retornar (padrão: 3).",
				},
			},
			Required: []string{"vitals_type"},
		},
	}
}

// GetAgendamentosDefinition returns the schema for GetAgendamentos tool
func GetAgendamentosDefinition() FunctionDeclaration {
	return FunctionDeclaration{
		Name:        "get_agendamentos",
		Description: "Recupera a lista de próximos agendamentos, compromissos médicos ou lembretes de medicação do idoso.",
		Parameters: &FunctionParameters{
			Type: "OBJECT",
			Properties: map[string]*Property{
				"limit": {
					Type:        "INTEGER",
					Description: "Número de agendamentos futuros a retornar (padrão: 5).",
				},
			},
			Required: []string{},
		},
	}
}

// GetToolDefinitions returns all available tool definitions
func GetToolDefinitions() []FunctionDeclaration {
	return []FunctionDeclaration{
		GetVitalsDefinition(),
		GetAgendamentosDefinition(),
	}
}
