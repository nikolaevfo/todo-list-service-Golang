package handler

import (
	"to-do-list/pkg/service"

	"github.com/gin-gonic/gin"

	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "to-do-list/docs"
)

// структура handlers использует
type Handler struct {
	services *service.Service
}

// метод для инициализации из main.go
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// Метод запускается из main.go и инициализирует все endpoints.
// Для разработки API применяется фреймворк для golang - gin.
func (h *Handler) InitRoutes() *gin.Engine {
	// инициализация роутера
	router := gin.New()

	// подключаем swagger к роутеру
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Прописываем endpoints к обработчикам текущего модуля handler:
	// в файлах с соответствующими именами в текущей папке:
	// auth.go, item.go, list.go.
	// Обработчики во фреймворке gin приниают в качестве параметра
	// указатель - *gin.Context.
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}
		items := api.Group("items")
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}
	}

	return router
}
