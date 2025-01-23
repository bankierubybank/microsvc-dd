package routes

import (
	"net/http"

	"github.com/bankierubybank/microsvc-dd/tarot/models"
	"github.com/gin-gonic/gin"
)

func Tarots(g *gin.RouterGroup) {
	g.GET("", GetTarots)
	g.GET(":cardnumber", GetTarotByCardNumber)
	g.GET("/random", GetRandomTarot)
}

// GetTarots godoc
//	@Summary		Get all tarot cards
//	@Description	Get all tarot cards
//	@Tags			tarot
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	models.TarotModel
//	@Router			/tarots/ [get]
func GetTarots(c *gin.Context) {
	tarots, _ := models.GetTarots()
	c.JSON(http.StatusOK, tarots)
}

//	@Summary		Get a tarot card by ID
//	@Description	Get a tarot card by ID
//	@Tags			tarot
//	@Accept			json
//	@Param			cardnumber	path	int	true	"Tarot ID"
//	@Produce		json
//	@Success		200	{object}	models.TarotModel
//	@failure		404	{string}	string
//	@Router			/tarots/{cardnumber} [get]
func GetTarotByCardNumber(c *gin.Context) {
	cardnumber := c.Param("cardnumber")

	tarot, err := models.GetTarotByCardNumber(cardnumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "tarot card not found, is card number correct?"})
	}
	c.JSON(http.StatusOK, tarot)
}

//	@Summary		Get a random tarot card
//	@Description	Get a random tarot card
//	@Tags			tarot
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.TarotModel
//	@Router			/tarots/random [get]
func GetRandomTarot(c *gin.Context) {
	tarot, _ := models.GetRandomTarot()
	c.JSON(http.StatusOK, tarot)
}
