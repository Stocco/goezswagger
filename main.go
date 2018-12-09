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
	workingProject, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	services.GenerateFile(workingProject)
}

//@path /v1/dummy/route
//@method post
//@summary create a dummy route
//@tags create
//@request MyDummyModel
//@response 200:MyDummyModelResponse 400:MyDummyModelResponse
func dummyRoute() {
	log.Println("HEY THERE")
}


type MyDummyModel struct {
	FieldOne	             *string `json:"field_one" description:"field is the number one" validate:"required"`
	FieldTwoNotMapped        *string `validate:"required"`
	FieldThreeForgotten      string
	FieldFour			      string `json:"field_four"`
	FieldFive			      NestedVal `json:"aeho"`
}

type NestedVal struct {
	FieldX		string `json:"field_x"`
	FieldY		int `json:"field_y"`
	Field		string  `json:"field"`
}

type MyDummyModelResponse struct {
	Amount int		`json:"amount"`
	AmountDouble float64 `json:"amount_double"`
}