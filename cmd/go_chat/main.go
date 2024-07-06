package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type PostgresCredentials struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

var Conn *sql.DB

func Connect(cred PostgresCredentials) error {
	credentials := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cred.Host, cred.Port, cred.Username, cred.Password, cred.DatabaseName)

	var err error

	Conn, err = sql.Open("postgres", credentials)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Connect(PostgresCredentials{
		Host:         "localhost",
		Port:         "5432",
		Username:     "postgres",
		Password:     "password",
		DatabaseName: "gochat",
	}); err != nil {
		log.Panic(err)
		os.Exit(1)
	}

	if _, err := Conn.Exec(`CREATE TABLE IF NOT EXISTS Messages (
        id         SERIAL,
        username   VARCHAR(255) NOT NULL CHECK (username <> ''),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`); err != nil {
		log.Panic(err)
		os.Exit(1)
	}

	r := gin.Default()

	r.LoadHTMLGlob("views/*")

	r.GET("/:name", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "chat.html", struct{ Username string }{
			Username: ctx.Param("name"),
		})
	})

	// TODO: Adicionar uma rota que retorna todos as mensagens em HTML.

	r.Run()
}
