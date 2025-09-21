package main

import (
	"net/http"

	"github.com/duds-fw/xcorrelations/sdk"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(sdk.GinMiddleware())

	r.GET("/apiA", func(c *gin.Context) {
		// cid := c.GetString("correlationId")
		ctx := c.Request.Context()

		sdk.Log(ctx, "Start processing API A")

		// Simulasi hit API B
		resp, err := sdk.HttpRequest(ctx, http.MethodGet, "http://localhost:8080/apiB", nil)
		if err != nil {
			sdk.Log(ctx, "Error calling API B: "+err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		sdk.Log(ctx, "Success calling API B")
		sdk.HttpResponse(c, 200, gin.H{"data": resp.String()})
		// c.Header("X-Correlation-ID", cid)
		// c.JSON(200, gin.H{"data": resp.String()})
	})

	r.GET("/apiB", func(c *gin.Context) {
		ctx := c.Request.Context()
		sdk.Log(ctx, "API B received request")
		c.JSON(200, gin.H{"msg": "Hello from API B"})
	})

	r.Run(":8080")
}
