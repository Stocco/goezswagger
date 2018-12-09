package main

import (
	"ezswag/services"
	"log"
	"os"
	"testing"
)

func TestSwagger(t *testing.T) {

	string, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	services.GenerateFile(string)
}
