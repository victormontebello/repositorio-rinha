package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

type TransactionPayload struct {
	Value       int32  `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

func main() {
	port := os.Getenv("PORT")
	dsn := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open => ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping => ", err)
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

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

	app.Listen(fmt.Sprintf(":%s", port))
}
