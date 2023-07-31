package main

import (
	"context"
	"log"

	"github.com/elliot14A/practice/search"
)

func main() {
	ctx := context.Background()
	result, err := search.SearchGoogle(ctx, "bbc.com")
	if err != nil {
		log.Println(err)
	}
	log.Println(result)
}
