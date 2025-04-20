package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pobruno/casa360/models"
	"github.com/robfig/cron/v3"
)

func CreateTask(c *gin.Context) {
	var task models.TaskInstallment
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar expressão CRON
	if _, err := cron.ParseStandard(task.RecurrenceCron); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expressão CRON inválida"})
		return
	}

	if err := task.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func ListTasks(c *gin.Context) {
	tasks, err := models.ListTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func GetTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	task := models.TaskInstallment{ID: id}
	if err := task.Get(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarefa não encontrada"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var task models.TaskInstallment
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar expressão CRON
	if _, err := cron.ParseStandard(task.RecurrenceCron); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expressão CRON inválida"})
		return
	}

	task.ID = id
	if err := task.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	task := models.TaskInstallment{ID: id}
	if err := task.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func UpdateTaskOccurrences(c *gin.Context) {
	tasks, err := models.ListTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	totalOcorrencias := 0
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	c.Stream(func(w io.Writer) bool {
		for _, task := range tasks {
			c.SSEvent("log", "Processando tarefa: "+task.Title+" (ID: "+task.ID.String()+")")
			
			schedule, err := parser.Parse(task.RecurrenceCron)
			if err != nil {
				c.SSEvent("error", "Erro ao parsear expressão CRON para tarefa "+task.ID.String()+": "+err.Error())
				continue
			}

			// Começar da data de início da tarefa
			nextTime := task.StartDate

			// Gerar todas as ocorrências até a data atual
			for nextTime.Before(now) || nextTime.Equal(now) {
				c.SSEvent("log", "Verificando data: "+nextTime.Format("2006-01-02")+" para tarefa: "+task.Title)
				
				occurrence := models.TaskOccurrence{
					TaskID:       task.ID,
					Date:         nextTime,
					Status:      false,
					UserID:      task.UserID,
					PayerGroupID: task.PayerGroupID,
					Subtasks:    task.Subtasks,
				}

				err := occurrence.Create()
				if err != nil {
					if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
						c.SSEvent("error", "Erro ao criar ocorrência: "+err.Error())
					} else {
						c.SSEvent("log", "Ocorrência já existe para data: "+nextTime.Format("2006-01-02"))
					}
				} else {
					totalOcorrencias++
					c.SSEvent("success", "Ocorrência criada para data: "+nextTime.Format("2006-01-02"))
				}

				nextTime = schedule.Next(nextTime)
			}
		}

		c.SSEvent("complete", gin.H{
			"message": "Processamento concluído",
			"total_ocorrencias": totalOcorrencias,
		})
		return false
	})
}

// CreateTaskOccurrence cria uma nova ocorrência de tarefa
func CreateTaskOccurrence(c *gin.Context) {
	var occurrence models.TaskOccurrence
	if err := c.ShouldBindJSON(&occurrence); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := occurrence.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, occurrence)
}

// UpdateTaskOccurrence atualiza uma ocorrência de tarefa existente
func UpdateTaskOccurrence(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Primeiro, buscar a ocorrência existente
	existingOccurrence := models.TaskOccurrence{ID: id}
	if err := existingOccurrence.Get(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ocorrência não encontrada"})
		return
	}

	// Em seguida, vincular apenas os campos que foram enviados
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualizar apenas os campos que foram enviados
	if status, exists := updateData["status"]; exists {
		existingOccurrence.Status = status.(bool)
	}
	if userID, exists := updateData["user_id"]; exists && userID != nil {
		id, err := uuid.Parse(userID.(string))
		if err == nil {
			existingOccurrence.UserID = id
		}
	}
	if payerGroupID, exists := updateData["payer_group_id"]; exists && payerGroupID != nil {
		id, err := uuid.Parse(payerGroupID.(string))
		if err == nil {
			existingOccurrence.PayerGroupID = id
		}
	}
	if subtasks, exists := updateData["subtasks"]; exists && subtasks != nil {
		existingOccurrence.Subtasks = subtasks.(json.RawMessage)
	}

	// Agora atualizar a ocorrência
	if err := existingOccurrence.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingOccurrence)
}

// DeleteTaskOccurrence remove uma ocorrência de tarefa
func DeleteTaskOccurrence(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	occurrence := models.TaskOccurrence{ID: id}
	if err := occurrence.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GenerateTaskOccurrences gera ocorrências para uma tarefa baseada em seu cronograma
func GenerateTaskOccurrences(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	task := models.TaskInstallment{ID: id}
	if err := task.Get(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarefa não encontrada"})
		return
	}

	if err := task.GenerateOccurrences(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
} 