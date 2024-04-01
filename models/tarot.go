package models

import (
	"errors"
	"math/rand"
)

// @Description Tarot Model
type TarotModel struct {
	Number   string `json:"number" binding:"required" example:"0" maxLength:"15"`
	Type     string `json:"type" binding:"required" example:"major" maxLength:"63"`
	CardName string `json:"cardname" binding:"required" example:"The Fool" maxLength:"255"`
}

var tarots = []TarotModel{
	{Number: "0", Type: "major", CardName: "The Fool"},
	{Number: "1", Type: "major", CardName: "The Magician"},
	{Number: "2", Type: "major", CardName: "The High Priestess"},
	{Number: "3", Type: "major", CardName: "The Empress"},
	{Number: "4", Type: "major", CardName: "The Emperor"},
	{Number: "5", Type: "major", CardName: "The Hierophant"},
	{Number: "6", Type: "major", CardName: "The Lovers"},
	{Number: "7", Type: "major", CardName: "The Chariot"},
	{Number: "8", Type: "major", CardName: "Strength"},
	{Number: "9", Type: "major", CardName: "The Hermit"},
	{Number: "10", Type: "major", CardName: "Wheel of Fortune"},
	{Number: "11", Type: "major", CardName: "Justice"},
	{Number: "12", Type: "major", CardName: "The Hanged Man"},
	{Number: "13", Type: "major", CardName: "Death"},
	{Number: "14", Type: "major", CardName: "Temperance"},
	{Number: "15", Type: "major", CardName: "The Devil"},
	{Number: "16", Type: "major", CardName: "The Tower"},
	{Number: "17", Type: "major", CardName: "The Star"},
	{Number: "18", Type: "major", CardName: "The Moon"},
	{Number: "19", Type: "major", CardName: "The Sun"},
	{Number: "20", Type: "major", CardName: "Judgement"},
	{Number: "21", Type: "major", CardName: "The World"},
}

func GetTarots() ([]TarotModel, error) {
	return tarots, nil
}

func GetTarotByCardNumber(cardnumber string) (TarotModel, error) {
	// Loop over the list of tarot deck, looking for a tarot card whose number value matches the parameter.
	var tarot TarotModel
	for _, a := range tarots {
		if a.Number == cardnumber {
			tarot = a
			return tarot, nil
		}
	}
	return tarot, errors.New("tarot card not found, is card number correct?")
}

func GetRandomTarot() (TarotModel, error) {
	return tarots[rand.Intn(len(tarots))], nil
}
