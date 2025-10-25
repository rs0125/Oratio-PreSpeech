package services

import (
	"Oratio/models"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

func SaveSession(speech string, questions []models.Question) (*models.Session, error) {
	// Marshal []Question into JSON
	questionsJSON, err := json.Marshal(questions)
	if err != nil {
		return nil, err
	}

	session := models.Session{
		Speech:      speech,
		Questions:   datatypes.JSON(questionsJSON),
		GeneratedBy: "GPT-4o-mini",
		CreatedAt:   time.Now(),
	}

	if err := DB.Create(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func GetSessionByID(id uint) (*models.Session, error) {
	var session models.Session
	if err := DB.First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
