package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db     *sql.DB
	once   sync.Once
	dbConn error
)

// InitDB inicializa a conexão com o banco de dados
func InitDB() {
	once.Do(func() {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")

		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		var err error
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Printf("Erro ao conectar ao banco de dados: %v", err)
			dbConn = err
			return
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Erro ao pingar o banco de dados: %v", err)
			dbConn = err
			return
		}

		log.Println("Conexão com o banco de dados estabelecida com sucesso!")
	})

	if dbConn != nil {
		log.Printf("Erro na conexão com o banco de dados: %v", dbConn)
	}
}

// GetDB retorna a conexão com o banco de dados
func GetDB() *sql.DB {
	return db
} 