package models

import (
	"errors"
	"math/rand"
)

// @Description Tarot Model
type TarotModel struct {
	ID       string `json:"id" binding:"required" example:"1" maxLength:"15"`
	Type     string `json:"type" binding:"required" example:"major" maxLength:"63"`
	Number   string `json:"number" binding:"required" example:"0" maxLength:"15"`
	CardName string `json:"cardname" binding:"required" example:"The Fool" maxLength:"255"`
	ImageURL string `json:"imageurl" binding:"required" example:"https://upload.wikimedia.org/wikipedia/commons/9/90/RWS_Tarot_00_Fool.jpg" maxLength:"255"`
}

var tarots = []TarotModel{
	{ID: "1", Type: "major", Number: "0", CardName: "The Fool", ImageURL: "https://upload.wikimedia.org/wikipedia/commons/9/90/RWS_Tarot_00_Fool.jpg"},
	{ID: "2", Type: "major", Number: "1", CardName: "The Magician", ImageURL: "https://upload.wikimedia.org/wikipedia/commons/d/de/RWS_Tarot_01_Magician.jpg"},
	{ID: "3", Type: "major", Number: "2", CardName: "The High Priestess", ImageURL: "https://upload.wikimedia.org/wikipedia/commons/8/88/RWS_Tarot_02_High_Priestess.jpg"},
	{ID: "4", Type: "minor", Number: "1c", CardName: "Ace of Cups", ImageURL: "https://upload.wikimedia.org/wikipedia/commons/3/36/Cups01.jpg"},
	{ID: "5", Type: "minor", Number: "2c", CardName: "Two of Cups", ImageURL: "https://upload.wikimedia.org/wikipedia/commons/f/f8/Cups02.jpg"},
}

func GetTarots() ([]TarotModel, error) {
	return tarots, nil
}

func GetTarotByID(id string) (TarotModel, error) {
	// Loop over the list of tarot deck, looking for
	// a tarot card whose ID value matches the parameter.
	var tarot TarotModel
	for _, a := range tarots {
		if a.ID == id {
			tarot = a
			return tarot, nil
		}
	}
	return tarot, errors.New("tarot not found, is ID correct?")
}

func GetRandomTarot() (TarotModel, error) {
	return tarots[rand.Intn(len(tarots))], nil
}
