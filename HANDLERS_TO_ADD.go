// Add these handlers to the handleToolCall function in main.go after manage_calendar_event

// 1. Gmail - send_email
case "send_email":
	if !googleFeaturesWhitelist[client.CPF] {
		return map[string]interface{}{"success": false, "error": "Gmail features not available"}
	}
	
	_, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
	if err != nil || time.Now().After(expiry) {
		return map[string]interface{}{"success": false, "error": "Google account not linked or expired"}
	}
	
	to, _ := args["to"].(string)
	subject, _ := args["subject"].(string)
	body, _ := args["body"].(string)
	
	gmailSvc := gmail.NewService(context.Background())
	err = gmailSvc.SendEmail(accessToken, to, subject, body)
	if err != nil {
		return map[string]interface{}{"success": false, "error": err.Error()}
	}
	return map[string]interface{}{"success": true, "message": "Email enviado com sucesso"}

// 2. Drive - save_to_drive
case "save_to_drive":
	if !googleFeaturesWhitelist[client.CPF] {
		return map[string]interface{}{"success": false, "error": "Drive features not available"}
	}
	
	_, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
	if err != nil || time.Now().After(expiry) {
		return map[string]interface{}{"success": false, "error": "Google account not linked or expired"}
	}
	
	filename, _ := args["filename"].(string)
	content, _ := args["content"].(string)
	folder, _ := args["folder"].(string)
	
	driveSvc := drive.NewService(context.Background())
	fileID, err := driveSvc.SaveFile(accessToken, filename, content, folder)
	if err != nil {
		return map[string]interface{}{"success": false, "error": err.Error()}
	}
	return map[string]interface{}{"success": true, "message": "Arquivo salvo", "file_id": fileID}

// 3. Sheets - manage_health_sheet
case "manage_health_sheet":
	if !googleFeaturesWhitelist[client.CPF] {
		return map[string]interface{}{"success": false, "error": "Sheets features not available"}
	}
	
	_, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
	if err != nil || time.Now().After(expiry) {
		return map[string]interface{}{"success": false, "error": "Google account not linked or expired"}
	}
	
	action, _ := args["action"].(string)
	sheetsSvc := sheets.NewService(context.Background())
	
	if action == "create" {
		title, _ := args["title"].(string)
		url, err := sheetsSvc.CreateHealthSheet(accessToken, title)
		if err != nil {
			return map[string]interface{}{"success": false, "error": err.Error()}
		}
		return map[string]interface{}{"success": true, "message": "Planilha criada", "url": url}
	}
	
	// TODO: Implement append action
	return map[string]interface{}{"success": false, "error": "Action not implemented"}

// 4. Docs - create_health_doc
case "create_health_doc":
	if !googleFeaturesWhitelist[client.CPF] {
		return map[string]interface{}{"success": false, "error": "Docs features not available"}
	}
	
	_, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
	if err != nil || time.Now().After(expiry) {
		return map[string]interface{}{"success": false, "error": "Google account not linked or expired"}
	}
	
	title, _ := args["title"].(string)
	content, _ := args["content"].(string)
	
	docsSvc := docs.NewService(context.Background())
	url, err := docsSvc.CreateDocument(accessToken, title, content)
	if err != nil {
		return map[string]interface{}{"success": false, "error": err.Error()}
	}
	return map[string]interface{}{"success": true, "message": "Documento criado", "url": url}

// 5. Maps - find_nearby_places
case "find_nearby_places":
	if !googleFeaturesWhitelist[client.CPF] {
		return map[string]interface{}{"success": false, "error": "Maps features not available"}
	}
	
	placeType, _ := args["place_type"].(string)
	location, _ := args["location"].(string)
	radius := 5000
	if r, ok := args["radius"].(float64); ok {
		radius = int(r)
	}
	
	mapsSvc := maps.NewService(context.Background(), s.cfg.GoogleMapsAPIKey)
	places, err := mapsSvc.FindNearbyPlaces(placeType, location, radius)
	if err != nil {
		return map[string]interface{}{"success": false, "error": err.Error()}
	}
	return map[string]interface{}{"success": true, "places": places}

// 6. YouTube - search_videos
case "search_videos":
	if !googleFeaturesWhitelist[client.CPF] {
		return map[string]interface{}{"success": false, "error": "YouTube features not available"}
	}
	
	_, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
	if err != nil || time.Now().After(expiry) {
		return map[string]interface{}{"success": false, "error": "Google account not linked or expired"}
	}
	
	query, _ := args["query"].(string)
	maxResults := int64(5)
	if mr, ok := args["max_results"].(float64); ok {
		maxResults = int64(mr)
	}
	
	youtubeSvc := youtube.NewService(context.Background())
	videos, err := youtubeSvc.SearchVideos(accessToken, query, maxResults)
	if err != nil {
		return map[string]interface{}{"success": false, "error": err.Error()}
	}
	return map[string]interface{}{"success": true, "videos": videos}

// 7. Spotify - play_music (simplified)
case "play_music":
	if !googleFeaturesWhitelist[client.CPF] {
		return map[string]interface{}{"success": false, "error": "Spotify features not available"}
	}
	
	// TODO: Implement Spotify OAuth separately
	return map[string]interface{}{"success": false, "error": "Spotify integration pending OAuth setup"}

// 8. Uber - request_ride (simplified)
case "request_ride":
	if !googleFeaturesWhitelist[client.CPF] {
		return map[string]interface{}{"success": false, "error": "Uber features not available"}
	}
	
	// TODO: Implement Uber OAuth separately
	return map[string]interface{}{"success": false, "error": "Uber integration pending OAuth setup"}

// 9. Google Fit - get_health_data
case "get_health_data":
	if !googleFeaturesWhitelist[client.CPF] {
		return map[string]interface{}{"success": false, "error": "Google Fit features not available"}
	}
	
	_, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
	if err != nil || time.Now().After(expiry) {
		return map[string]interface{}{"success": false, "error": "Google account not linked or expired"}
	}
	
	fitSvc := googlefit.NewService(context.Background())
	steps, err := fitSvc.GetStepsToday(accessToken)
	if err != nil {
		return map[string]interface{}{"success": false, "error": err.Error()}
	}
	return map[string]interface{}{"success": true, "steps": steps}

// 10. WhatsApp - send_whatsapp
case "send_whatsapp":
	if !googleFeaturesWhitelist[client.CPF] {
		return map[string]interface{}{"success": false, "error": "WhatsApp features not available"}
	}
	
	// TODO: Implement WhatsApp Business API
	return map[string]interface{}{"success": false, "error": "WhatsApp integration pending configuration"}
