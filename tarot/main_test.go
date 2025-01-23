package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"

	"github.com/bankierubybank/microsvc-dd/tarot/models"
	"github.com/bankierubybank/microsvc-dd/tarot/routes"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestGetTarots(t *testing.T) {
	r := SetUpRouter()
	r.GET("/api/v1/tarots", routes.GetTarots)
	req, _ := http.NewRequest("GET", "/api/v1/tarots", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var tarots []models.TarotModel
	json.Unmarshal(w.Body.Bytes(), &tarots)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTarotByCardNumber(t *testing.T) {
	r := SetUpRouter()
	r.GET("/api/v1/tarots/1", routes.GetTarots)
	req, _ := http.NewRequest("GET", "/api/v1/tarots/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var tarot models.TarotModel
	json.Unmarshal(w.Body.Bytes(), &tarot)

	assert.Equal(t, http.StatusOK, w.Code)
}
