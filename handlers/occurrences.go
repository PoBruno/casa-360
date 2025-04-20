package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pobruno/casa360/models"
)

// ListTaskOccurrences lista todas as ocorrências de tarefas
func ListTaskOccurrences(c *gin.Context) {
	occurrences, err := models.ListTaskOccurrences()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, occurrences)
}

// ListFinanceOccurrences lista todas as ocorrências financeiras
func ListFinanceOccurrences(c *gin.Context) {
	occurrences, err := models.ListFinanceOccurrences()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, occurrences)
}

// ListOccurrencesDashboard retorna todas as ocorrências do dashboard
func ListOccurrencesDashboard(c *gin.Context) {
	occurrences, err := models.ListOccurrencesDashboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, occurrences)
}

// GetLastWallet retorna o último registro da carteira de um usuário
func GetLastWallet(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	wallet, err := models.GetLastWalletByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if wallet == nil {
		c.JSON(http.StatusOK, gin.H{"amount": 0})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// ListTransactions retorna todas as transações de uma ocorrência
func ListTransactions(c *gin.Context) {
	occurrenceID, err := uuid.Parse(c.Param("occurrence_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ocorrência inválido"})
		return
	}

	transactions, err := models.ListTransactionsByOccurrenceID(occurrenceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
} 