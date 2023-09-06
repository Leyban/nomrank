package app

import "nomrank/pkg/infra/storage/postgres"

type server struct {
	db     postgres.DB
	router any
}
