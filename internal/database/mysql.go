package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/user/golang-api-rest/internal/config"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar no MySQL: %w", err)
	}

	log.Println("Conectado ao MySQL com sucesso")
	return db, nil
}

func Migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL DEFAULT 'user',
		verified TINYINT(1) NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("erro ao criar tabela users: %w", err)
	}

	log.Println("Migração executada com sucesso")
	return nil
}