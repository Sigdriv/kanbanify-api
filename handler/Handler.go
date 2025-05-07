package handler

import "github.com/gin-gonic/gin"

func Handler() {
	router := gin.Default()
	router.Use(gin.Logger())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.GET("/api/issues", getIssues)
	router.POST("/api/issue", createIssue)
	router.PATCH("/api/issue", updateIssue)
	router.DELETE("/api/issues/:id", deleteIssue)

	router.Run(":8080")
}
