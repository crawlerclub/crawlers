package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"extractor/extractor"
)

// https://racing.hkjc.com/racing/information/Chinese/Horse/Horse.aspx?HorseId=HK_2021_G372

func main() {
	schemaFile := flag.String("schema", "", "Path to the schema JSON file")
	url := flag.String("url", "", "URL to extract data from")
	flag.Parse()

	if *schemaFile == "" || *url == "" {
		log.Fatal("Both schema file and URL are required")
	}

	schemaData, err := os.ReadFile(*schemaFile)
	if err != nil {
		log.Fatalf("Error reading schema file: %v", err)
	}

	var schema extractor.Schema
	if err := json.Unmarshal(schemaData, &schema); err != nil {
		log.Fatalf("Error parsing schema JSON: %v", err)
	}

	extractor := extractor.NewExtractor(schema)
	items, err := extractor.Extract(*url)
	if err != nil {
		log.Fatalf("Error extracting data: %v", err)
	}

	jsonData, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		log.Fatalf("Error converting results to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
