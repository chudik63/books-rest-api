package routes

import (
	"go-books-api/internal/delivery/http/handler"

	_ "go-books-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegistrateRoutes(app *gin.Engine, h *handler.Handler) {
	books := app.Group("/books")

	books.POST("/", h.AddBook)
	books.GET("/", h.ListBooks)
	books.GET("/:id", h.GetBook)
	books.DELETE("/:id", h.DeleteBook)
	books.PUT("/:id", h.UpdateBook)

	app.GET("/docs/*any", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
