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

// // CMindexstr Index of list names for id of manufacture/customer list
// var CMindexstr = "_CMindexstr"

// // MatchState status index
// var MatchState = "_MatchState"

// ManufRetail structure Manufacture recall record with JSON conversion.
type ManufRetail struct {
	Co           string `json:"co"`
	Model        string `json:"model"` //model/style
	Description  string `json:"description"`
	ID           string `json:"id"`     //manufacture/retail ID code
	Serial       string `json:"serial"` //ups, etc...
	PurchaseDate string `json:"purchasedate"`
	Recall       string `json:"recall"`    //date of recall
	Timestamp    int64  `json:"timestamp"` //utc timestamp of creation
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

	/*var empty []string
	jsonAsBytes, _ := json.Marshal(empty) //marshal an emtpy array of strings to clear the index
	err = stub.PutState(CMindexstr, jsonAsBytes)
	if err != nil {
		return nil, err
	}*/

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function, multiple function call API
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.Write(stub, args)
	} else if function == "Manufacture" {
		return t.Manufacture(stub, args)
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

// Manufacture - function for manufacture record input
func (t *SimpleChaincode) Manufacture(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("running Manufacture()")

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	Co := args[0]           //Company
	Model := args[1]        //model #
	Description := args[2]  //Description
	key := args[3]          //Product ID
	Serial := args[4]       //Serial #
	PurchaseDate := args[5] //Purchase Data
	Recall := args[6]       //Recall Date
	Timestamp := makeTimestamp()
	//Check if record already exist
	ManufAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to get Manufacture/Retail Record")
	} else if ManufAsBytes != nil {
		fmt.Println("This manufacture record already exist: " + key)
		return shim.Error("This record already exisit: " + key)
	}

	manuf := &ManufRetail{Co, Model, Description, key, Serial, PurchaseDate, Recall, Timestamp}
	//str := `{"company": "` + co + `", "model": "` + mod + `", "description": "` + des + `", "product ID": "` + id + `", "serial ":"` + serial + `", "purchace data": "` + purc + `", "Recall": "` + recal + `"}`
	manufJSONasbytes, err := json.Marshal(manuf)
	err = stub.PutState(key, manufJSONasbytes) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	//save record index to state
	err = stub.PutState(key, manufJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Index range queries
	indexName := "co~id"
	coidIndexKey, err := stub.CreateCompositeKey(indexName, []string{manuf.Co, manuf.ID})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(coidIndexKey, value)

	// ==== Record (Manuf.) saved and indexed. Return success ====
	fmt.Println("- end init record manufacture")
	return nil, nil
}

// Delete - remove a key/value pair from state
func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string

	var mnaufJSON ManufRetail
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
	coidIndexKey, err := stub.CreateCompositeKey(indexName, []string{recordJSON.Co, recordJSON.ID})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = stub.DelState(coidIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}
	return shim.Success(nil)
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
	//co.Timestamp = makeTimestamp()
	//mo.Timestamp = makeTimestamp()

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
	
	mo := &ManufRetail{}

	//Find match between records
	json.Unmarshal(matchAsBytes, &co)
	//json.Unmarshal(matchAsBytes, &mo)

	//msg := string(co.ProdID) + " " + string(mo.ID)
	//fmt.Println(msg)
	//fmt.Println(key)
	if mo.ID == key {
		fmt.Println("Product ID match found in customer & manufacture records, given search key:" + key)
		jsonmsg := string(matchAsBytes)

		fmt.Println("Manufacture record", jsonmsg)
		return shim.Error("Match found between product ID's")
	}
	fmt.Println("Match not found between product ID's")
	fmt.Println(mo.ID)

	return shim.Error("Received unknown search function invocation")
}
