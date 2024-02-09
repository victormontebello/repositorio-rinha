package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"database/sql"

	_ "github.com/lib/pq"
)

/*
POST /clientes/[id]/transacoes
GET /clientes/[id]/extrato
*/

type TransactionPayload struct {
	Value       int32  `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

func main() {
	env := os.Getenv("DB_DSN")
	if env == "" {
		env = "host=localhost user=postgres password=postgres dbname=rinha sslmode=disable"
	}

	db, err := sql.Open("postgres", env)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Post("/clientes/:id/transacoes", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		// Funciona né K
		if err != nil || id <= 0 || id > 5 {
			log.Println("Id param", err)
			return c.Status(404).SendString("Cliente não encontrado")
		}

		var body TransactionPayload
		if err := c.BodyParser(&body); err != nil {
			log.Println("Body parser", err)
			return c.Status(400).SendString("Erro ao processar o corpo requisição")
		}

		row := db.QueryRow("select * from create_transaction($1, $2, $3, $4)", id, body.Value, body.Type, body.Description)

		var response json.RawMessage
		if err := row.Scan(&response); err != nil {
			log.Println("Scan err", err)
			return c.Status(400).SendString("Erro ao processar a resposta da requisição ao banco")
		}

		return c.Status(200).Send(response)
	})

	app.Get("/clientes/:id/extrato", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		// Funciona né K
		if err != nil || id <= 0 || id > 5 {
			log.Println("Id param", err)
			return c.Status(404).SendString("Cliente não encontrado")
		}

		row := db.QueryRow("select * from get_extract($1)", id)

		var response json.RawMessage
		if err := row.Scan(&response); err != nil {
			log.Println("Scan err", err)
			return c.Status(400).SendString("Erro ao processar a resposta da requisição ao banco")
		}

		return c.Status(200).Send(response)
	})

	app.Listen(":3000")
}
