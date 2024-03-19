package routes

import (
	"net/http"

	"github.com/bankierubybank/microsvc-dd/models"
	"github.com/gin-gonic/gin"
)

func Tarots(g *gin.RouterGroup) {
	g.GET("", GetTarots)
	g.GET(":id", GetTarotByID)
	g.GET("/random", GetRandomTarot)
}

// GetTarots godoc
// @Summary	Get all tarot cards
// @Description	Get all tarot cards
// @Tags			tarot
// @Accept			json
// @Produce		json
// @Success		200 {array}  models.TarotModel
// @Router			/tarots/ [get]
func GetTarots(c *gin.Context) {
	us, _ := models.GetTarots()
	c.JSON(http.StatusOK, us)
}

// @Summary	Get a tarot card by ID
// @Description	Get a tarot card by ID
// @Tags			tarot
// @Accept			json
// @Param			id	path	int	true	"Tarot ID"
// @Produce		json
// @Success		200 {object}  models.TarotModel
// @Router			/tarots/{id} [get]
func GetTarotByID(c *gin.Context) {
	id := c.Param("id")

	u, err := models.GetTarotByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "tarot not found"})
	}
	c.JSON(http.StatusOK, u)
}

// @Summary	Get a random tarot card
// @Description	Get a random tarot card
// @Tags			tarot
// @Accept			json
// @Produce		json
// @Success		200 {object}  models.TarotModel
// @Router			/tarots/random [get]
func GetRandomTarot(c *gin.Context) {
	u, _ := models.GetRandomTarot()
	c.JSON(http.StatusOK, u)
}
