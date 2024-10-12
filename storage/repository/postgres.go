package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
	TIMEZONE string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Устанавливаем часовой пояс
	_, err = db.Exec(fmt.Sprintf("ALTER DATABASE %s SET TIMEZONE = '%s'", cfg.DBName, cfg.TIMEZONE))
	if err != nil {
		return nil, err
	}

	// Получаем текущее время
	var currentTime string
	err = db.QueryRow("SELECT NOW()").Scan(&currentTime)
	if err != nil {
		return nil, err
	}

	logrus.Printf("Created Postgres DB %s\n with time: %v", cfg.DBName, currentTime)
	return db, nil
}
