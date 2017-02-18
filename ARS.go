/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var customtindex = "_customindex" //customer id
var manuindex = "_manuindex"      //name for id of manufacture list
var MatchState = "_matchstate"    //match state status

//Customer purchase record
type Customer struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Provstate string `json:"provstate"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Dob       string `json:"dob"`
	ProdID    string `json:"prodID"`
}

//"ManufRetail... " Manufacture recall record
type ManufRetail struct {
	Co           string `json:"co"`
	Model        string `json:"model"` //model/style
	Description  string `json:"description"`
	ID           string `json:"id"`     //manufacture/retail ID code
	Serial       string `json:"serial"` //ups, etc...
	PurchaseDate string `json:"purchasedate"`
	Recall       string `json:"recall"` //date of recall
}

type storage struct {
	store []byte
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var key string

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("Startup Sequence", []byte(args[0])) //start up test		//std hello_world
	if err != nil {
		return nil, err
	}

	key = args[0]

	//system test subsection -MVP only, customer 1 = manu#2, model of ext database
	c0 := Customer{
		Firstname: `json:"John"`,
		Lastname:  `json:"Smith"`,
		Address:   `json:"123 main st."`,
		City:      `json:"Toronto"`,
		Provstate: `json:"Ontario"`,
		Gender:    `json:"male"`,
		Email:     `json:"js@yahoo.ca"`,
		Phone:     `json:"416-555-9988"`,
		Dob:       `json:"15/04/85"`,
		ProdID:    `json:"0000001A"`,
	}

	c1 := Customer{
		Firstname: "Lydia",
		Lastname:  "Elas",
		Address:   "55 Danforth Ave",
		City:      "Toronto",
		Provstate: "Ontario",
		Gender:    "female",
		Email:     "le@hotmail.com",
		Phone:     "416-222-1112",
		Dob:       "",
		ProdID:    `json:"0000001B"`,
	}

	c2 := Customer{
		Firstname: "Bob",
		Lastname:  "Borg",
		Address:   "700 Mathesson ave",
		City:      "Mississauga",
		Provstate: "Ontario",
		Gender:    "male",
		Email:     "bb3@gmail.com",
		Phone:     "905-123-8974",
		Dob:       "",
		ProdID:    `json:"0000001C"`,
	}

	c3 := Customer{
		Firstname: "Cynthia",
		Lastname:  "Nyquist",
		Address:   "100 Country Rd. NE.",
		City:      "Calgary",
		Provstate: "Alberta",
		Gender:    "female",
		Email:     "cnq@ieee.org",
		Phone:     "285-552-4578",
		Dob:       "",
		ProdID:    `json:"0000001D"`,
	}

	c4 := Customer{
		Firstname: "Mathew",
		Lastname:  "Johns",
		Address:   "10 Steel ave.",
		City:      "Vaughan",
		Provstate: "Ontario",
		Gender:    "male",
		Email:     "mj@gmail.com",
		Phone:     "905-593-3345",
		Dob:       "",
		ProdID:    `json:"0000001E"`,
	}

	m0 := ManufRetail{
		Co:           `json:"Roots"`,
		Model:        `json:"1305-0131"`, //model/style
		Description:  `json:"Roots Cabin Long Johns"`,
		ID:           `json:"0000001E"`, //manufacture/retail product ID code, to be auto generated
		Serial:       `json:""`,         //ups, etc...
		PurchaseDate: `json:"08/20/2016"`,
		Recall:       `json:"10/7/2016"`,
	}

	m1 := ManufRetail{
		Co:           `json:"Canada Varsity"`,
		Model:        `json:"1105-0226"`, //model/style
		Description:  `json:"Canada Varsity Jacket Black Pepper"`,
		ID:           `json:"0000001A"`, //manufacture/retail product ID code
		Serial:       `json:""`,         //ups, etc...
		PurchaseDate: `json:"08/29/2016"`,
		Recall:       `json:"11/30/2016"`,
	}

	m2 := ManufRetail{
		Co:           `json:"Royloo Educational"`,
		Model:        `json:"R59601"`,
		Description:  `json:"Royloo Educational Light Cube"`,
		ID:           `json:"0000001B"`,
		Serial:       `json:"66960596014"`, //ups
		PurchaseDate: `json:"04/02/2016"`,
		Recall:       `json:"11/30/2016"`,
	}

	m3 := ManufRetail{
		Co:           `json:"L'Atelier Cheval de Bois"`,
		Model:        `json:""`,
		Description:  `json:"L'Atelier Cheval de Bois -Wood Rattle"`,
		ID:           `json:"0000001C"`,
		Serial:       `json:""`,
		PurchaseDate: `json:"06/19/2016"`,
		Recall:       `json:"11/16/2016"`,
	}

	m4 := ManufRetail{
		Co:           `json:"Specialized Bicycle Components Cdn. Inc"`,
		Model:        `json:"RA-61136"`,
		Description:  `json:"Specialized Bicycle Components"`,
		ID:           `json:"0000001D"`,
		Serial:       `json:""`,
		PurchaseDate: `json:"07/09/2016"`,
		Recall:       `json:"11/18/2016"`,
	}

	//use for loop for non MVP, need system to read in data
	//ci := []Customer{}
	//mi := []ManufRetail{}
	//for j := 0; j < 5; j++ {}

	store1, err := json.Marshal(c0)  //parse json members
	err = stub.PutState(key, store1) //put in blockchain, given key
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("test field in blockchain -purchase record")
	//fmt.Println(string(store1))

	store2, err := json.Marshal(m0)
	err = stub.PutState(key, store2)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("test field in blockchain manufacture record")
	//fmt.Println(string(store2))

	store3, err := json.Marshal(c1)
	err = stub.PutState(key, store3)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("test field in blockchain -purchase record")
	//fmt.Println(string(store1))

	store4, err := json.Marshal(m1)
	err = stub.PutState(key, store4)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("test field in blockchain manufacture record")
	//fmt.Println(string(store2))

	store5, err := json.Marshal(c2)
	err = stub.PutState(key, store5)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("test field in blockchain -purchase record")
	//fmt.Println(string(store1))

	store6, err := json.Marshal(m2)
	err = stub.PutState(key, store6)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("test field in blockchain manufacture record")
	//fmt.Println(string(store2))
	store7, err := json.Marshal(c3)
	err = stub.PutState(key, store7)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("test field in blockchain -purchase record")
	//fmt.Println(string(store1))

	store8, err := json.Marshal(m3)
	err = stub.PutState(key, store8)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("test field in blockchain manufacture record")
	//fmt.Println(string(store2))

	store9, err := json.Marshal(c4)
	err = stub.PutState(key, store9)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("test field in blockchain -purchase record")
	//fmt.Println(string(store1))

	store10, err := json.Marshal(m4)
	err = stub.PutState(key, store10)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("test field in blockchain manufacture record")
	//fmt.Println(string(store2))

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var key string

	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "delete" {
		return t.Delete(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	key = args[0]
	//functions for search comparison & matching
	output, err := stub.GetState(key)
	co := []Customer{}
	mo := []ManufRetail{}

	//test two for MVP-sim production
	for i := 0; i < 5; i++ {
		err = json.Unmarshal(output, &co[i])
		if err != nil {
			fmt.Println("error:", err)
		}

		err = json.Unmarshal(output, &mo[i])
		if err != nil {
			fmt.Println("error:", err)
		}
		if co[i].ProdID == mo[i].ID {
			fmt.Println("BlockChain matched recall customer purchase record by product ID ", co[i].Firstname, co[i].Lastname)
			//console.log("BlockChain matched recall customer purchase record by product ID %+v", co[i].Firstname, co[i].Lastname)
		}
	}

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query")
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0]
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// Delete - remove a key/value pair from state
func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	record := args[0]
	err := stub.DelState(record) //remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

//time stamp
func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
