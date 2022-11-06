package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/api/v1/issueClaim", issueClaim)

	router.Run("localhost:8090")
}

func issueClaim(c *gin.Context) {
	claims := IssueClaim("holderID_abc")
	c.IndentedJSON(http.StatusOK, claims)
}
