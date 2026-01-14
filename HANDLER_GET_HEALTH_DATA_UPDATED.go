// Updated handler for get_health_data - Add this to HANDLERS_TO_ADD.go

// 9. Google Fit - get_health_data (UPDATED with database storage)
case "get_health_data":
	if !googleFeaturesWhitelist[client.CPF] {
		return map[string]interface{}{"success": false, "error": "Google Fit features not available"}
	}
	
	_, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
	if err != nil || time.Now().After(expiry) {
		return map[string]interface{}{"success": false, "error": "Google account not linked or expired"}
	}
	
	fitSvc := googlefit.NewService(context.Background())
	
	// Get all health data
	healthData, err := fitSvc.GetAllHealthData(accessToken)
	if err != nil {
		return map[string]interface{}{"success": false, "error": err.Error()}
	}
	
	// Save to database automatically
	if healthData.Steps > 0 {
		s.db.SaveVitalSign(client.IdosoID, "passos", fmt.Sprintf("%d", healthData.Steps), "steps", "google_fit", "")
	}
	if healthData.HeartRate > 0 {
		s.db.SaveVitalSign(client.IdosoID, "frequencia_cardiaca", fmt.Sprintf("%.0f", healthData.HeartRate), "bpm", "google_fit", "")
	}
	if healthData.Calories > 0 {
		s.db.SaveVitalSign(client.IdosoID, "calorias", fmt.Sprintf("%d", healthData.Calories), "kcal", "google_fit", "")
	}
	if healthData.Distance > 0 {
		s.db.SaveVitalSign(client.IdosoID, "distancia", fmt.Sprintf("%.2f", healthData.Distance), "km", "google_fit", "")
	}
	if healthData.Weight > 0 {
		s.db.SaveVitalSign(client.IdosoID, "peso", fmt.Sprintf("%.1f", healthData.Weight), "kg", "google_fit", "")
	}
	
	return map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"steps":      healthData.Steps,
			"heart_rate": healthData.HeartRate,
			"calories":   healthData.Calories,
			"distance":   healthData.Distance,
			"weight":     healthData.Weight,
		},
		"message": "Dados de saúde coletados e salvos com sucesso",
	}
