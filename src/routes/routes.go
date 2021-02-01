package routes

import (
	handleTask "github.com/AnthuanGarcia/IntegradoraII/src/handlers/tasks"
	"github.com/gin-gonic/gin"
)

// Routes - Generacion de Rutas
type Routes struct{}

// StartGin - Deploy Api
func (c Routes) StartGin() {
	r := gin.Default()

	api := r.Group("/prueba")
	{
		api.GET("/tasks/:id", handleTask.HandleGetTask)
		api.GET("/tasks/", handleTask.HandleGetTasks)
		api.PUT("/tasks/", handleTask.HandleCreateTask)
		api.POST("/tasks/", handleTask.HandleUpdateTask)
	}

	r.Run(":8000")
}
