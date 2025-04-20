package handlers

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pobruno/casa360/models"
)

// Handlers para Centro de Custo
func CreateFinanceCC(c *gin.Context) {
	var cc models.FinanceCC
	if err := c.ShouldBindJSON(&cc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cc)
}

func ListFinanceCCs(c *gin.Context) {
	ccs, err := models.ListFinanceCCs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ccs)
}

// Handlers para Moedas
func CreateFinanceCurrency(c *gin.Context) {
	var currency models.FinanceCurrency
	if err := c.ShouldBindJSON(&currency); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := currency.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, currency)
}

func ListFinanceCurrencies(c *gin.Context) {
	currencies, err := models.ListFinanceCurrencies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, currencies)
}

// Handlers para Finanças
func CreateFinance(c *gin.Context) {
	var finance models.FinanceInstallment
	if err := c.ShouldBindJSON(&finance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := finance.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, finance)
}

func ListFinances(c *gin.Context) {
	finances, err := models.ListFinanceInstallments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, finances)
}

func GetFinance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	finance := models.FinanceInstallment{ID: id}
	if err := finance.Get(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Financeiro não encontrado"})
		return
	}

	c.JSON(http.StatusOK, finance)
}

func UpdateFinance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var finance models.FinanceInstallment
	if err := c.ShouldBindJSON(&finance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	finance.ID = id
	if err := finance.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, finance)
}

func DeleteFinance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	finance := models.FinanceInstallment{ID: id}
	if err := finance.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func UpdateFinanceOccurrences(c *gin.Context) {
	finances, err := models.ListFinanceInstallments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	totalOcorrencias := 0

	c.Stream(func(w io.Writer) bool {
		for _, finance := range finances {
			c.SSEvent("log", "Processando finança: "+finance.Title+" (ID: "+finance.ID.String()+")")

			// Validar recurrence_days
			if finance.RecurrenceDays <= 0 {
				c.SSEvent("error", "Dias de recorrência inválidos para finança "+finance.Title+": "+string(finance.RecurrenceDays))
				continue
			}

			// Definir data final como:
			// 1. Data final da finança (se existir)
			// 2. Data atual (se não existir data final)
			endDate := now
			if finance.EndDate != nil && finance.EndDate.Before(now) {
				endDate = *finance.EndDate
			}

			// Começar da data de início
			nextDate := finance.StartDate

			// Gerar todas as ocorrências até a data final
			for nextDate.Before(endDate) || nextDate.Equal(endDate) {
				c.SSEvent("log", "Verificando data: "+nextDate.Format("2006-01-02")+" para finança: "+finance.Title)

				occurrence := models.FinanceOccurrence{
					FinanceID: finance.ID,
					Date:      nextDate,
					Amount:    finance.Amount,
					Status:    false,
				}

				err := occurrence.Create()
				if err != nil {
					if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
						c.SSEvent("error", "Erro ao criar ocorrência: "+err.Error())
					} else {
						c.SSEvent("log", "Ocorrência já existe para data: "+nextDate.Format("2006-01-02"))
					}
				} else {
					totalOcorrencias++
					c.SSEvent("success", "Ocorrência criada para data: "+nextDate.Format("2006-01-02"))
				}

				nextDate = nextDate.AddDate(0, 0, finance.RecurrenceDays)
			}
		}

		c.SSEvent("complete", gin.H{
			"message": "Processamento concluído",
			"total_ocorrencias": totalOcorrencias,
		})
		return false
	})
}

// CreateFinanceOccurrence cria uma nova ocorrência financeira
func CreateFinanceOccurrence(c *gin.Context) {
	var occurrence models.FinanceOccurrence
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

// UpdateFinanceOccurrence atualiza uma ocorrência financeira existente
func UpdateFinanceOccurrence(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var occurrence models.FinanceOccurrence
	if err := c.ShouldBindJSON(&occurrence); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	occurrence.ID = id
	if err := occurrence.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, occurrence)
}

// DeleteFinanceOccurrence remove uma ocorrência financeira
func DeleteFinanceOccurrence(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	occurrence := models.FinanceOccurrence{ID: id}
	if err := occurrence.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GenerateFinanceOccurrences gera ocorrências para uma finança baseada em sua recorrência
func GenerateFinanceOccurrences(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	finance := models.FinanceInstallment{ID: id}
	if err := finance.Get(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Finança não encontrada"})
		return
	}

	if err := finance.GenerateOccurrences(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
} 