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
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}


// Customer structure for purchase record JSON conversion
type Customer struct {
	//Firstname string `json:"firstname"` //user first/last name, address not kept in blockchain! MVP
	//Lastname  string `json:"lastname"`
	//Address   string `json:"address"`
	City      string `json:"city"`
	Provstate string `json:"provstate"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Dob       string `json:"dob"`
	ProdID    string `json:"prodID"`
	IMEI      int    `json:"000000000000000"`
	Timestamp int64  `json:"timestamp"`
}


// Main
func main() {
	err := shim.Start(new(SimpleChaincode)) // Start chaincode system
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init sets all internal values and preparing the blockchain
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//var key string
	//key = args[0]

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("Startup Sequence", []byte(args[0])) //start up test		//std hello_world
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in init %s\n\r", r)
		}
	}()

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty) //marshal an emtpy array of strings to clear the index
	err = stub.PutState(CMindexstr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	var SearchMatch CrossRef
	jsonAsBytes, _ = json.Marshal(SearchMatch) //clear the open search state struct
	err = stub.PutState(MatchState, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function, multiple function call API
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.Write(stub, args)
	} else if function == "Customer" {
		return t.Customer(stub, args)
	} else if function == "delete" {
		return t.Delete(stub, args)
	} else if function == "search" {
		return t.Search(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return shim.Error("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return shim.Error("Received unknown function query")
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}
	if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Nil amount for " + key + "\"}"
		return shim.Error(jsonResp)
	}
	jsonResp = "{\"Name\":\"" + key + "\",\"Value\":\"" + string(valAsbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return valAsbytes, nil
}

// Write - write variable into chaincode state
func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]
	value = args[1]
	err = stub.PutState(name, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Customer - invoke function to write for customers
func (t *SimpleChaincode) Customer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var err error
	fmt.Println("running customer()")

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8. Name of the key and six values to set")
	}

	city := args[0]   //city
	provst := args[1] //province/state
	sex := args[2]    //gender
	email := args[3]  //email
	phone := args[4]  //phone number
	dob := args[5]    //DoB
	prodid := args[6] //product ID
	imei := strconv.Atoi(args[7])   // IMEI
	if err != nil {
		return shim.Error("Expecting integer value for IMEI value")
	}
	Timestamp := makeTimestamp()
	
	//Check if record already exist
	customAsBytes, err := stub.GetState(prodid)
	if err != nil {
		return shim.Error("Failed to get Customer Record")
	} else if customAsBytes != nil {
		fmt.Println("This customer record already exist:" + prodid)
		return shim.Error("This customer record already exist:" + prodid)
	}
	
	
	cust := &Customer{city, provst, sex, email, phone, dob, prodid, imei, Timestamp}
	//str := `{"company": "` + city + `", "model": "` + provst + `", "description": "` + sex + `", "E-mail": "` + email + `", "Phone ":"` + phone + `", "Date of Birth ":"` + dob + `", "Product ID": "` + prodid + `", "IMEI": "` + imei + `"}`
	customJSONasBytes, err := json.Marshal(cust)
	err = stub.PutState(prodid, customJSONasBytes) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	//save record index to state
	customerAsBytes, err := stub.PutState(prodid, customJSONasBytes)
	if err != nil {
		return shim.Error("Failed to get C record index")
	}
	//index key/value range queries
	indexName := "imei~id"
	imeidIndexKey, err := stub.CreateCompositeKey(indexName, []string{cust.prodid, cust.imei})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(imeidIndexKey, value)

	// ==== Marble saved and indexed. Return success ====
	fmt.Println("- end init marble")

	return nil, nil
}

//makeTimestamp time indexing
func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// Search function for search comparison & matching ..."search" "prodID"
func (t *SimpleChaincode) Search(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key string

	fmt.Println("running Search()")
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Search", r)
		}
	}()

	key = args[0] // ProdID
	

	//Check if ProdID exist
	matchAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to get record")
	}

	// get string list of chaincode arguments.
	/*data := stub.GetStringArgs()
	for z := range data {
		fmt.Println("in the chain: ", data[z])
	}*/
	co := &Customer{}
	

	//Find match between records
	json.Unmarshal(matchAsBytes, &co)
	//json.Unmarshal(matchAsBytes, &mo)

	//msg := string(co.ProdID) + " " + string(mo.ID)
	//fmt.Println(msg)
	//fmt.Println(key)
	if co.ProdID == key {
		fmt.Println("Product ID match found in customer & manufacture records, given search key:" + key)
		jsonmsgA := string(matchAsBytes)
		fmt.Println("Customer record", jsonmsg)

		//return shim.Error("Match not found between product ID's")
	} else {
		fmt.Println("Match not found between product ID's")
		fmt.Println(co.ProdID)
	}

	return shim.Error("Received unknown search function invocation")
}

// Delete - remove a key/value pair from state
func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var customJSON customer
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	recordName := args[0]

	// to maintain the record~name index, we need to read the record first and get its ID
	valAsbytes, err := stub.GetState(recordName) //get the record from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + recordName + "\"}"
		return shim.E(jsonRresp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + recordName + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &customJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + recordName + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(recordName) //remove the marble from chaincode state
	if err != nil {
		return nil, shim.Error("Failed to delete state:" + err.Error())
	}

	// maintain the index
	indexName := "co~id"
	imeidIndexKey, err := stub.CreateCompositeKey(indexName, []string{customJSON.IMEI, customJSON.ProdID})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = stub.DelState(imeidIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}
	return shim.Success(nil)
}