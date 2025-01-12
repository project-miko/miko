package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/pgvector/pgvector-go"
	"github.com/project-miko/miko/models"
)

func TestVectorSearch(t *testing.T) {
	initTester()
	var greeting string
	conn := models.GetPGInst("pg_main")
	err := conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		t.Error(err)
	}
	t.Log(greeting)

	// the vector to search
	queryVector := pgvector.NewVector([]float32{0, 1, 0})

	// the vector distance operator has the following types (in pgvector):
	//   <-> : Euclidean distance (default)
	//   <#> : negative inner product
	//   <=> : cosine distance
	//   see pgvector documentation for more details
	//
	// this example uses the default <-> for similarity search
	rows, err := conn.Query(
		context.Background(),
		`
        SELECT id, name, embedding
        FROM items
        ORDER BY embedding <-> $1
        LIMIT 3
        `,
		queryVector,
	)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   int
			name string
			// pgx can map the vector column to a []float32
			// or other type (depending on the pgvector version and pgx driver version)
			embedding pgvector.Vector
		)
		if err := rows.Scan(&id, &name, &embedding); err != nil {
			log.Fatalf("Scan failed: %v\n", err)
		}
		fmt.Printf("id=%d, name=%s, embedding=%v\n", id, name, embedding.String())
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Row iteration error: %v\n", err)
	}
}

func TestDemo(t *testing.T) {
	t.Log("test")
}
