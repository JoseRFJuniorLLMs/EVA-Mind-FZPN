// Add these 6 tool definitions after manage_calendar_event in tools.go (around line 151)

// 1. Gmail - send_email
map[string]interface{}{
	"name":        "send_email",
	"description": "Envia um email usando Gmail do usuário",
	"parameters": map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"to": map[string]interface{}{
				"type":        "string",
				"description": "Email do destinatário",
			},
			"subject": map[string]interface{}{
				"type":        "string",
				"description": "Assunto do email",
			},
			"body": map[string]interface{}{
				"type":        "string",
				"description": "Corpo do email",
			},
		},
		"required": []string{"to", "subject", "body"},
	},
},

// 2. Drive - save_to_drive
map[string]interface{}{
	"name":        "save_to_drive",
	"description": "Salva um documento no Google Drive do usuário",
	"parameters": map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"filename": map[string]interface{}{
				"type":        "string",
				"description": "Nome do arquivo",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "Conteúdo do documento",
			},
			"folder": map[string]interface{}{
				"type":        "string",
				"description": "Nome da pasta (opcional)",
			},
		},
		"required": []string{"filename", "content"},
	},
},

// 3. Sheets - manage_health_sheet
map[string]interface{}{
	"name":        "manage_health_sheet",
	"description": "Gerencia planilha de saúde no Google Sheets",
	"parameters": map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Ação: 'create' ou 'append'",
				"enum":        []string{"create", "append"},
			},
			"title": map[string]interface{}{
				"type":        "string",
				"description": "Título da planilha (para create)",
			},
			"data": map[string]interface{}{
				"type":        "object",
				"description": "Dados de saúde (date, time, blood_pressure, glucose, medication, notes)",
			},
		},
		"required": []string{"action"},
	},
},

// 4. Docs - create_health_doc
map[string]interface{}{
	"name":        "create_health_doc",
	"description": "Cria um documento de saúde no Google Docs",
	"parameters": map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"title": map[string]interface{}{
				"type":        "string",
				"description": "Título do documento",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "Conteúdo do documento",
			},
		},
		"required": []string{"title", "content"},
	},
},

// 5. Maps - find_nearby_places
map[string]interface{}{
	"name":        "find_nearby_places",
	"description": "Busca lugares próximos (farmácias, hospitais, restaurantes)",
	"parameters": map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"place_type": map[string]interface{}{
				"type":        "string",
				"description": "Tipo: pharmacy, hospital, restaurant, etc",
			},
			"location": map[string]interface{}{
				"type":        "string",
				"description": "Localização (lat,lng)",
			},
			"radius": map[string]interface{}{
				"type":        "integer",
				"description": "Raio em metros (padrão: 5000)",
			},
		},
		"required": []string{"place_type", "location"},
	},
},

// 6. YouTube - search_videos
map[string]interface{}{
	"name":        "search_videos",
	"description": "Busca vídeos no YouTube",
	"parameters": map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "Termo de busca",
			},
			"max_results": map[string]interface{}{
				"type":        "integer",
				"description": "Número máximo de resultados (padrão: 5)",
			},
		},
		"required": []string{"query"},
	},
},
