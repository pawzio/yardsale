package main

import (
	"context"
	"log"

	"github.com/pawzio/yardsale/catalog/pkg/executor"
	"github.com/pawzio/yardsale/catalog/pkg/httpsvc"
)

func main() {
	log.Println("Initializing Catalog svc")

	ctx := context.Background()

	srv, err := httpsvc.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initialized Catalog svc")

	executor.Run(ctx, srv.Start)
}
