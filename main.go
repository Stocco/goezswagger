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
	FieldOne	             *string `json:"field_one,omitempty" description:"field is the number one" validate:"required"`
	FieldTwoNotMapped        *string `validate:"required"`
	FieldThreeForgotten      string
	FieldFour			      string `json:"field_four"`
	FieldFive			      NestedVal `json:"aeho"`
	FieldFiveArray			  []*NestedVal `json:"aeho_arrat"`
	IncReq					  IncomingTransferRequest `json:"inc_request"`
}

type NestedVal struct {
	FieldX		string `json:"field_x"`
	FieldY		int `json:"field_y"`
	FieldZeta	uint64  `json:"field_zeta"`
	FieldBool	bool `json:"field_bool"`
}

type MyDummyModelResponse struct {
	Amount int		`json:"amount"`
	AmountDouble float64 `json:"amount_double"`
}

type IncomingTransferRequest struct {
	Holden		      *AccountHolder     `json:"source_holder"`
	Paymentz	      *Paymentz          `json:"payment_details"`
	Longitude         *float64           `json:"longitude"`
	Metadata          map[string]string  `json:"metadata"`
}

type AccountHolder struct {
	Name             string       `json:"name,omitempty"`
	BankAccount      *BankAc `json:"bankaccount"`
}

type BankAc struct {
	Destination   *bool  `json:"destination,omitempty"`
}

type Paymentz struct {
	Type             string  `json:"type,omitempty"`
}


type Health struct {
	Name       string           `json:"name"`
	Components HealthComponents `json:"components"`
}

type HealthComponents struct {
	Http     bool `json:"Http"`
	Cache    bool `json:"Cache"`
}