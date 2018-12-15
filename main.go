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


//@path /v2/dummy/route
//@method post
//@summary create a dummy route WITH ARRAY IN BODY
//@tags create
//@request []MyDummyModel
//@response 200:MyDummyModelResponse 400:MyDummyModelResponse
func dummyRouteV2() {
	log.Println("HEY THERE SQUARED")
}

//@model MyDummyModel
type MyDummyModel struct {
	FieldOne            *string `json:"field_one,omitempty" description:"field is the number one"`
	FieldTwoNotMapped   *string `validate:"required"`
	FieldThreeForgotten string
	FieldFour           string                  `json:"field_four" description:"field is the number four"`
	FieldFive           NestedVal               `json:"aeho" description:"this is a composite value"`
	FieldFiveArray      []*NestedVal            `json:"aeho_arrat"`
	IncReq              IncomingTransferRequest `json:"inc_request"`
}

//@model NestedVal
type NestedVal struct {
	FieldX    string `json:"field_x" description:"field x is the best" example:"best"`
	FieldY    int    `json:"field_y" description:"field y is integer" example:"-42"`
	FieldZeta uint64 `json:"field_zeta" description:"field is always positive" example:"1"`
	FieldBool bool   `json:"field_bool" description:"field bool is always boolean" example:"true"`
}

//@model MyDummyModelResponse
type MyDummyModelResponse struct {
	Amount       int     `json:"amount"`
	AmountDouble float64 `json:"amount_double" description:"field double is of course doubl" example:"53.21"`
}

//@model IncomingTransferRequest
type IncomingTransferRequest struct {
	Holden    *AccountHolder    `json:"source_holder"`
	Paymentz  *Paymentz         `json:"payment_details"`
	Longitude *float64          `json:"longitude"`
	Metadata  map[string]string `json:"metadata"`
}

//@model AccountHolder
type AccountHolder struct {
	Name        string  `json:"name,omitempty" description:"the name of the person" example:"renato"`
	Surname     string  `json:"surname,omitempty" description:"the surname of the person" example:"stocco"`
	BankAccount *BankAc `json:"bankaccount" description:"if person has bankaccount" example:"true"`
	Weight      *float64 `json:"weight,omitempty" description:"the weight of this person" example:"73.2"`
	Age         int     `json:"age,omitempty" description:"the age of the person" example:"25"`
}

//@model BankAc
type BankAc struct {
	Destination *bool `json:"destination,omitempty"`
}

//@model Paymentz
type Paymentz struct {
	Type string `json:"type,omitempty"`
}

type Health struct {
	Name       string           `json:"name"`
	Components HealthComponents `json:"components"`
}

type HealthComponents struct {
	Http  bool `json:"Http"`
	Cache bool `json:"Cache"`
}
