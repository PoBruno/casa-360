package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pobruno/casa360/config"
	"github.com/pobruno/casa360/handlers"
)

func main() {
	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado")
	}

	// Inicializa o banco de dados
	config.InitDB()

	// Inicializa o router
	r := gin.Default()

	// Configura as rotas
	setupRoutes(r)

	// Inicia o servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado na porta %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// Grupo de rotas para usuários
	setupUserRoutes(r)

	// Grupo de rotas para grupos de pagadores
	setupPayerGroupRoutes(r)

	// Grupo de rotas para centro de custo
	setupFinanceCCRoutes(r)

	// Grupo de rotas para moedas
	setupCurrencyRoutes(r)

	// Grupo de rotas para tarefas
	setupTaskRoutes(r)

	// Grupo de rotas para finanças
	setupFinanceRoutes(r)

	// Grupo de rotas para dashboard e carteiras
	setupDashboardRoutes(r)
}

func setupUserRoutes(r *gin.Engine) {
	// Rotas sem barra final
	r.POST("/users", handlers.CreateUser)
	r.GET("/users", handlers.ListUsers)
	r.GET("/users/:id", handlers.GetUser)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)

	// Rotas com barra final
	r.POST("/users/", handlers.CreateUser)
	r.GET("/users/", handlers.ListUsers)
	r.GET("/users/:id/", handlers.GetUser)
	r.PUT("/users/:id/", handlers.UpdateUser)
	r.DELETE("/users/:id/", handlers.DeleteUser)
}

func setupPayerGroupRoutes(r *gin.Engine) {
	// Rotas sem barra final
	r.POST("/payer-groups", handlers.CreatePayerGroup)
	r.GET("/payer-groups", handlers.ListPayerGroups)
	r.GET("/payer-groups/:id", handlers.GetPayerGroup)
	r.PUT("/payer-groups/:id", handlers.UpdatePayerGroup)
	r.DELETE("/payer-groups/:id", handlers.DeletePayerGroup)
	r.POST("/payer-groups/:id/members", handlers.CreatePayerGroupMember)
	r.GET("/payer-groups/:id/members", handlers.ListPayerGroupMembers)
	r.DELETE("/payer-groups/:id/members/:member_id", handlers.DeletePayerGroupMember)

	// Rotas com barra final
	r.POST("/payer-groups/", handlers.CreatePayerGroup)
	r.GET("/payer-groups/", handlers.ListPayerGroups)
	r.GET("/payer-groups/:id/", handlers.GetPayerGroup)
	r.PUT("/payer-groups/:id/", handlers.UpdatePayerGroup)
	r.DELETE("/payer-groups/:id/", handlers.DeletePayerGroup)
	r.POST("/payer-groups/:id/members/", handlers.CreatePayerGroupMember)
	r.GET("/payer-groups/:id/members/", handlers.ListPayerGroupMembers)
	r.DELETE("/payer-groups/:id/members/:member_id/", handlers.DeletePayerGroupMember)
}

func setupFinanceCCRoutes(r *gin.Engine) {
	// Rotas sem barra final
	r.POST("/finance-cc", handlers.CreateFinanceCC)
	r.GET("/finance-cc", handlers.ListFinanceCCs)

	// Rotas com barra final
	r.POST("/finance-cc/", handlers.CreateFinanceCC)
	r.GET("/finance-cc/", handlers.ListFinanceCCs)
}

func setupCurrencyRoutes(r *gin.Engine) {
	// Rotas sem barra final
	r.POST("/currencies", handlers.CreateFinanceCurrency)
	r.GET("/currencies", handlers.ListFinanceCurrencies)

	// Rotas com barra final
	r.POST("/currencies/", handlers.CreateFinanceCurrency)
	r.GET("/currencies/", handlers.ListFinanceCurrencies)
}

func setupTaskRoutes(r *gin.Engine) {
	// Rotas sem barra final
	r.POST("/tasks", handlers.CreateTask)
	r.GET("/tasks", handlers.ListTasks)
	r.GET("/tasks/:id", handlers.GetTask)
	r.PUT("/tasks/:id", handlers.UpdateTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)
	r.POST("/tasks/update-occurrences", handlers.UpdateTaskOccurrences)

	// Ocorrências de tarefas
	r.POST("/tasks/:id/occurrences", handlers.GenerateTaskOccurrences)
	r.POST("/task-occurrences", handlers.CreateTaskOccurrence)
	r.GET("/task-occurrences", handlers.ListTaskOccurrences)
	r.PUT("/task-occurrences/:id", handlers.UpdateTaskOccurrence)
	r.DELETE("/task-occurrences/:id", handlers.DeleteTaskOccurrence)

	// Rotas com barra final
	r.POST("/tasks/", handlers.CreateTask)
	r.GET("/tasks/", handlers.ListTasks)
	r.GET("/tasks/:id/", handlers.GetTask)
	r.PUT("/tasks/:id/", handlers.UpdateTask)
	r.DELETE("/tasks/:id/", handlers.DeleteTask)
	r.POST("/tasks/update-occurrences/", handlers.UpdateTaskOccurrences)

	// Ocorrências de tarefas com barra final
	r.POST("/tasks/:id/occurrences/", handlers.GenerateTaskOccurrences)
	r.POST("/task-occurrences/", handlers.CreateTaskOccurrence)
	r.GET("/task-occurrences/", handlers.ListTaskOccurrences)
	r.PUT("/task-occurrences/:id/", handlers.UpdateTaskOccurrence)
	r.DELETE("/task-occurrences/:id/", handlers.DeleteTaskOccurrence)
}

func setupFinanceRoutes(r *gin.Engine) {
	// Rotas sem barra final
	r.POST("/finances", handlers.CreateFinance)
	r.GET("/finances", handlers.ListFinances)
	r.GET("/finances/:id", handlers.GetFinance)
	r.PUT("/finances/:id", handlers.UpdateFinance)
	r.DELETE("/finances/:id", handlers.DeleteFinance)
	r.POST("/finances/update-occurrences", handlers.UpdateFinanceOccurrences)

	// Ocorrências financeiras
	r.POST("/finances/:id/occurrences", handlers.GenerateFinanceOccurrences)
	r.POST("/finance-occurrences", handlers.CreateFinanceOccurrence)
	r.GET("/finance-occurrences", handlers.ListFinanceOccurrences)
	r.PUT("/finance-occurrences/:id", handlers.UpdateFinanceOccurrence)
	r.DELETE("/finance-occurrences/:id", handlers.DeleteFinanceOccurrence)

	// Rotas com barra final
	r.POST("/finances/", handlers.CreateFinance)
	r.GET("/finances/", handlers.ListFinances)
	r.GET("/finances/:id/", handlers.GetFinance)
	r.PUT("/finances/:id/", handlers.UpdateFinance)
	r.DELETE("/finances/:id/", handlers.DeleteFinance)
	r.POST("/finances/update-occurrences/", handlers.UpdateFinanceOccurrences)

	// Ocorrências financeiras com barra final
	r.POST("/finances/:id/occurrences/", handlers.GenerateFinanceOccurrences)
	r.POST("/finance-occurrences/", handlers.CreateFinanceOccurrence)
	r.GET("/finance-occurrences/", handlers.ListFinanceOccurrences)
	r.PUT("/finance-occurrences/:id/", handlers.UpdateFinanceOccurrence)
	r.DELETE("/finance-occurrences/:id/", handlers.DeleteFinanceOccurrence)
}

func setupDashboardRoutes(r *gin.Engine) {
	// Dashboard de ocorrências
	r.GET("/occurrences/dashboard", handlers.ListOccurrencesDashboard)
	r.GET("/occurrences/dashboard/", handlers.ListOccurrencesDashboard)

	// Carteiras
	r.GET("/wallets/:user_id", handlers.GetLastWallet)
	r.GET("/wallets/:user_id/", handlers.GetLastWallet)

	// Transações
	r.GET("/transactions/:occurrence_id", handlers.ListTransactions)
	r.GET("/transactions/:occurrence_id/", handlers.ListTransactions)
} 