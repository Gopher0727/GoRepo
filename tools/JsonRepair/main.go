package main

import (
	"fmt"
	"log"

	"github.com/kaptinlin/jsonrepair"
)

func main() {
	json := "{name: 'John'}"

	ans, err := jsonrepair.JSONRepair(json)
	if err != nil {
		log.Fatalf("Failed to repair JSON: %v", err)
	}
	fmt.Println(ans)
}
