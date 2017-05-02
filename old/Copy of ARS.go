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

// CMindexstr Index of list names for id of manufacture/customer list
var CMindexstr = "_CMindexstr"
var MatchState = "_MatchState" // MatchState status index

// Customer structure for purchase record JSON conversion
type Customer struct {
	Firstname string `json:"firstname"` //user first/last name, address not kept in blockchain! MVP
	Lastname  string `json:"lastname"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Provstate string `json:"provstate"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Dob       string `json:"dob"`
	ProdID    string `json:"prodID"`
	IMEI      int    `json:"0000000000000"`
	Timestamp int64  `json:"timestamp"`
}

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

// CrossRef - holds referenced records for matching
type CrossRef struct {
	Search []byte `json:"search_match"`
}

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
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
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
	} else if function == "Manufacture" {
		return t.Manufacture(stub, args)
	} else if function == "delete" {
		return t.Delete(stub, args)
	} else if function == "search" {
		return t.Search(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

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
	if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Nil amount for " + key + "\"}"
		return nil, errors.New(jsonResp)
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
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
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
	//var value, val2, val3, val4, val5, val6, val7 string
	var err error
	fmt.Println("running customer()")

	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 7. Name of the key and six values to set")
	}

	city := args[0]   //city
	provst := args[1] //province/state
	sex := args[2]    //gender
	email := args[3]  //email
	phone := args[4]  //phone number
	dob := args[5]    //DoB
	prodid := args[6] //prod ID
	imei := args[7]   // IMEI

	//Check if record already exist
	customAsBytes, err := stub.GetState(imei)
	if err != nil {
		return nil, errors.New("Failed to get Customer Record")
	}
	cust := Customer{}
	json.Unmarshal(customAsBytes, &cust)
	if cust.ProdID == prodid {
		fmt.Println("Customer Record already exists" + prodid)
		fmt.Println(cust)
		return nil, errors.New("Customer Record already exists")
	}

	//store record
	str := `{"company": "` + city + `", "model": "` + provst + `", "description": "` + sex + `", "E-mail": "` + email + `", "Phone ":"` + phone + `", "Date of Birth ":"` + dob + `", "Product ID": "` + prodid + `", "IMEI": "` + imei + `"}`
	err = stub.PutState(prodid, []byte(str)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	//get record index
	customerAsBytes, err := stub.GetState(CMindexstr)
	if err != nil {
		return nil, errors.New("Failed to get C record index")
	}
	var CrecordIndex []string
	json.Unmarshal(customerAsBytes, &CrecordIndex) //un stringify it aka JSON.parse()

	//append
	CrecordIndex = append(CrecordIndex, prodid) //add customer name to index list
	fmt.Println("! C record index: ", CrecordIndex)
	jsonAsBytes, _ := json.Marshal(CrecordIndex)
	err = stub.PutState(CMindexstr, jsonAsBytes) //store name of customer

	fmt.Println("- end record C input")
	return nil, nil
}

// Manufacture - function for manufacture record input
func (t *SimpleChaincode) Manufacture(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("running Manufacture()")

	if len(args) != 7 {
		return nil, errors.New("Incorrect number of arguments. Expecting 7")
	}

	co := args[0]     //Company
	mod := args[1]    //model #
	des := args[2]    //Description
	id := args[3]     //Product ID
	serial := args[4] //Serial #
	purc := args[5]   //Purchase Data
	recal := args[6]  //Recall Date

	//Check if record already exist
	ManufAsBytes, err := stub.GetState(id)
	if err != nil {
		return nil, errors.New("Failed to get Manufacture/Retail Record")
	}
	manuf := ManufRetail{}
	json.Unmarshal(ManufAsBytes, &manuf)
	if manuf.ID == id {
		fmt.Println("Manufacture Record already exists" + id)
		fmt.Println(manuf)
		return nil, errors.New("Manufacture Record already exists")
	}

	//store record
	str := `{"company": "` + co + `", "model": "` + mod + `", "description": "` + des + `", "product ID": "` + id + `", "serial ":"` + serial + `", "purchace data": "` + purc + `", "Recall": "` + recal + `"}`
	err = stub.PutState(id, []byte(str)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	//get record index
	ManufRAsBytes, err := stub.GetState(CMindexstr)
	if err != nil {
		return nil, errors.New("Failed to get M record index")
	}
	var MrecordIndex []string
	json.Unmarshal(ManufRAsBytes, &MrecordIndex) //un stringify it aka JSON.parse()

	//append
	MrecordIndex = append(MrecordIndex, id) //add manufacture record to index list
	fmt.Println("! M record index: ", MrecordIndex)
	jsonAsBytes, _ := json.Marshal(MrecordIndex)
	err = stub.PutState(CMindexstr, jsonAsBytes) //store record of marble

	fmt.Println("- end record M input")
	return nil, nil
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

	recordsAsBytes, err := stub.GetState(CMindexstr)
	if err != nil {
		return nil, errors.New("Failed to get record index")
	}
	var recordIndex []string
	json.Unmarshal(recordsAsBytes, &recordIndex) //un stringify it aka JSON.parse()

	//remove marble from index
	for i, val := range recordIndex {
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + record)
		if val == record { //find the correct record
			fmt.Println("found record")
			recordIndex = append(recordIndex[:i], recordIndex[i+1:]...) //remove it
			for x := range recordIndex {                                //debug prints...
				fmt.Println(string(x) + " - " + recordIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(recordIndex) //save new index
	err = stub.PutState(CMindexstr, jsonAsBytes)

	return nil, nil
}

//makeTimestamp time indexing
func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// Search function for search comparison & matching ..."search" "prodID"
func (t *SimpleChaincode) Search(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key string

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	fmt.Println("running Search()")

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
		return nil, errors.New("Failed to get record")
	}

	// get string list of chaincode arguments.
	/*data := stub.GetStringArgs()
	for z := range data {
		fmt.Println("in the chain: ", data[z])
	}*/
	co := Customer{}
	mo := ManufRetail{}

	//Find match between records
	json.Unmarshal(matchAsBytes, &co)
	json.Unmarshal(matchAsBytes, &mo)
	if co.ProdID == key && mo.ID == key {
		fmt.Println("Product ID match found in customer & manufacture records:" + key)
		fmt.Println(co)
		fmt.Println(mo)
		return nil, errors.New("Product ID match not found in customer & manufacture records")
	}

	return nil, errors.New("Received unknown search function invocation")
}
