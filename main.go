package main

import (
	"context"
	"log"

	"github.com/propro-productions/go-utils/search"
)

func main() {
	ctx := context.Background()
	result, err := search.SearchGoogle(ctx, "bbc.com")
	if err != nil {
		log.Println(err)
	}
	log.Println(result)
}
