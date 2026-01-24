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

// ScanMedicationVisualDefinition returns the schema for visual medication scanning
func ScanMedicationVisualDefinition() FunctionDeclaration {
	return FunctionDeclaration{
		Name:        "scan_medication_visual",
		Description: "Abre a câmera do celular para identificar medicamentos visualmente via Gemini Vision. Use quando o paciente expressar confusão sobre qual remédio tomar ou pedir ajuda para identificar medicação.",
		Parameters: &FunctionParameters{
			Type: "OBJECT",
			Properties: map[string]*Property{
				"reason": {
					Type:        "STRING",
					Description: "Motivo da solicitação de scan (ex: 'paciente confuso sobre medicação matinal', 'não sabe qual tomar agora')",
				},
				"time_of_day": {
					Type:        "STRING",
					Description: "Período do dia para filtrar medicamentos candidatos",
					Enum:        []string{"morning", "afternoon", "evening", "night"},
				},
			},
			Required: []string{"reason", "time_of_day"},
		},
	}
}

// AnalyzeVoiceProsodyDefinition returns the schema for voice prosody analysis
func AnalyzeVoiceProsodyDefinition() FunctionDeclaration {
	return FunctionDeclaration{
		Name:        "analyze_voice_prosody",
		Description: "Analisa biomarcadores vocais (pitch, ritmo, pausas, tremor) para detectar sinais de depressão, ansiedade, Parkinson ou desidratação. Use quando perceber mudanças significativas no padrão de fala do paciente.",
		Parameters: &FunctionParameters{
			Type: "OBJECT",
			Properties: map[string]*Property{
				"analysis_type": {
					Type:        "STRING",
					Description: "Tipo de análise específica a realizar",
					Enum:        []string{"depression", "anxiety", "parkinson", "hydration", "full"},
				},
				"audio_segment_seconds": {
					Type:        "INTEGER",
					Description: "Duração do segmento de áudio a analisar em segundos (padrão: 30)",
				},
			},
			Required: []string{},
		},
	}
}

// ApplyPHQ9Definition returns the schema for PHQ-9 depression assessment
func ApplyPHQ9Definition() FunctionDeclaration {
	return FunctionDeclaration{
		Name:        "apply_phq9",
		Description: "Aplica a escala PHQ-9 (Patient Health Questionnaire) conversacionalmente para avaliar depressão. Faça as 9 perguntas de forma natural e empática, uma por vez.",
		Parameters: &FunctionParameters{
			Type: "OBJECT",
			Properties: map[string]*Property{
				"start_assessment": {
					Type:        "BOOLEAN",
					Description: "Iniciar aplicação da escala PHQ-9",
				},
			},
			Required: []string{"start_assessment"},
		},
	}
}

// ApplyGAD7Definition returns the schema for GAD-7 anxiety assessment
func ApplyGAD7Definition() FunctionDeclaration {
	return FunctionDeclaration{
		Name:        "apply_gad7",
		Description: "Aplica a escala GAD-7 (Generalized Anxiety Disorder) conversacionalmente para avaliar ansiedade. Faça as 7 perguntas de forma natural e empática.",
		Parameters: &FunctionParameters{
			Type: "OBJECT",
			Properties: map[string]*Property{
				"start_assessment": {
					Type:        "BOOLEAN",
					Description: "Iniciar aplicação da escala GAD-7",
				},
			},
			Required: []string{"start_assessment"},
		},
	}
}

// ApplyCSSRSDefinition returns the schema for C-SSRS suicide risk assessment
func ApplyCSSRSDefinition() FunctionDeclaration {
	return FunctionDeclaration{
		Name:        "apply_cssrs",
		Description: "🚨 CRÍTICO: Aplica a Columbia Suicide Severity Rating Scale (C-SSRS) para avaliar risco suicida. Use APENAS se o paciente mencionar suicídio, autolesão ou desejo de morrer. Faça as perguntas com extremo cuidado e empatia.",
		Parameters: &FunctionParameters{
			Type: "OBJECT",
			Properties: map[string]*Property{
				"trigger_phrase": {
					Type:        "STRING",
					Description: "Frase que disparou a necessidade da avaliação (ex: 'não quero mais viver')",
				},
				"start_assessment": {
					Type:        "BOOLEAN",
					Description: "Iniciar aplicação da escala C-SSRS",
				},
			},
			Required: []string{"trigger_phrase", "start_assessment"},
		},
	}
}

// GetToolDefinitions returns all available tool definitions
func GetToolDefinitions() []FunctionDeclaration {
	return []FunctionDeclaration{
		GetVitalsDefinition(),
		GetAgendamentosDefinition(),
		ScanMedicationVisualDefinition(),
		AnalyzeVoiceProsodyDefinition(),
		ApplyPHQ9Definition(),
		ApplyGAD7Definition(),
		ApplyCSSRSDefinition(),
	}
}
