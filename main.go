package main

import (
	"ezswag/services"
	"log"
	"os"
)

//@title goezswag api
//@version 1.0.0-beta
//@description goezswag is a easy swagger
func main() {
	string, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	services.GenerateFile(string)
}

//@path /v1/dummy/route
//@method post
//@summary create a dummy route
//@tags create
//@request MyDummyModel
//@response 200:MyDummyModelResponse
func dummyRoute() {
	log.Println("HEY THERE")
}


type MyDummyModel struct {
	FieldOne	             string `json:"field_one" validate:"required"`
	FieldTwoNotMapped        string `validate:"required"`
	FieldThreeForgotten      string
	FieldFour			      string `json:"field_four"`
}

type MyDummyModelResponse struct {
	Amount int		`json:"amount"`
	AmountDouble float64 `json:"amount_double"`
}