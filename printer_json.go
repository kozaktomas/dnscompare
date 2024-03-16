package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type JsonResultPrinter struct{}

func (printer JsonResultPrinter) Print(results []DnsResult) {
	data, err := json.Marshal(results)
	if err != nil {
		log.Fatal("could not marshal results to JSON")
	}

	fmt.Println(string(data))
}
