package tasks

import (
	"log"
	"net/http"

	db "github.com/AnthuanGarcia/RestApiGo/db"
	model "github.com/AnthuanGarcia/RestApiGo/src/models"
	"github.com/gin-gonic/gin"
)

// HandleGetTasks - EndPoint Todas las tareas
func HandleGetTasks(c *gin.Context) {
	loadedTasks, err := db.GetAllTasks()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": loadedTasks})
}

// HandleGetTask - Endpoint una sola Tarea
func HandleGetTask(c *gin.Context) {
	var task model.Task

	if err := c.BindUri(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
	}

	loadedTask, err := db.GetTaskID(task.ID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ID":   loadedTask.ID,
		"Body": loadedTask.Body,
	})
}

// HandleCreateTask - Crea una tarea en un documento
func HandleCreateTask(c *gin.Context) {
	var task model.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}

	id, err := db.Create(&task)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// HandleUpdateTask - Actualiza una tarea en un documento
func HandleUpdateTask(c *gin.Context) {
	var task model.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}

	savedTask, err := db.Update(&task)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": savedTask})
}
