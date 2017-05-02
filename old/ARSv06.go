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
	"errors"
	"fmt"
	"time"
	//"strings"
	//"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode)) // Start chaincode system
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init sets all internal values and preparing the blockchain
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Initializing system")

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

	err = stub.CreateTable("customer", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "consumer_id", Type: shim.ColumnDefinition_UINT64, Key: true},
		&shim.ColumnDefinition{Name: "IMEI", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "dob_dttm", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "gender_cd", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "cntry_id", Type: shim.ColumnDefinition_INT32, Key: true},
		&shim.ColumnDefinition{Name: "state_prov_id", Type: shim.ColumnDefinition_INT32, Key: true},
		&shim.ColumnDefinition{Name: "city_id", Type: shim.ColumnDefinition_INT32, Key: true},
		&shim.ColumnDefinition{Name: "postal_zip_cd_abbr", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_userid", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_dttm", Type: shim.ColumnDefinition_UINT64, Key: false},
	})
	if err != nil {
		return nil, err
	}

	err = stub.CreateTable("manufacture", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "manufacturer_id", Type: shim.ColumnDefinition_UINT64, Key: true},
		&shim.ColumnDefinition{Name: "manufacture_nm", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "manufacture_url", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "phone_nmbr_txt", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_userid", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_dttm", Type: shim.ColumnDefinition_UINT64, Key: false},
	})
	if err != nil {
		return nil, err
	}

	err = stub.CreateTable("product", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "prod_id_id", Type: shim.ColumnDefinition_UINT64, Key: true},
		&shim.ColumnDefinition{Name: "upc_num", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "prod_nm", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "prod_desc", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "style_brand_nm", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "prod_pic", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "manufacture_id", Type: shim.ColumnDefinition_UINT64, Key: true},
		&shim.ColumnDefinition{Name: "txn_userid", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_dttm", Type: shim.ColumnDefinition_UINT64, Key: false},
	})
	if err != nil {
		return nil, err
	}

	err = stub.CreateTable("consumer_product_alert", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "product_alert_id", Type: shim.ColumnDefinition_UINT64, Key: true},
		&shim.ColumnDefinition{Name: "consumer_id", Type: shim.ColumnDefinition_UINT64, Key: true},
		&shim.ColumnDefinition{Name: "prod_id", Type: shim.ColumnDefinition_UINT64, Key: true},
		&shim.ColumnDefinition{Name: "manufacture_id", Type: shim.ColumnDefinition_UINT64, Key: true},
		&shim.ColumnDefinition{Name: "txn_userid", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_dttm", Type: shim.ColumnDefinition_UINT64, Key: false},
	})
	if err != nil {
		return nil, err
	}

	err = stub.CreateTable("product_alert", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "product_alert_id", Type: shim.ColumnDefinition_UINT64, Key: true},
		&shim.ColumnDefinition{Name: "prod_id", Type: shim.ColumnDefinition_UINT64, Key: true},
		&shim.ColumnDefinition{Name: "ars_prod_alert_msg", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "prod_alert_serverity_cd", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "healthy_cdn_url", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_userid", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_dttm", Type: shim.ColumnDefinition_UINT64, Key: false},
	})
	if err != nil {
		return nil, err
	}

	err = stub.CreateTable("consumer_product", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "consumer_id", Type: shim.ColumnDefinition_UINT64, Key: true},
		&shim.ColumnDefinition{Name: "prod_id", Type: shim.ColumnDefinition_UINT64, Key: true},
		&shim.ColumnDefinition{Name: "unit_price", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "purchase_date", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "seller_cd", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_userid", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_dttm", Type: shim.ColumnDefinition_UINT64, Key: false},
	})
	if err != nil {
		return nil, err
	}

	err = stub.CreateTable("country", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "cntry_id", Type: shim.ColumnDefinition_INT32, Key: true},
		&shim.ColumnDefinition{Name: "cntry_nm", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "cntry_abbr", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "cntry_cd", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_userid", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_dttm", Type: shim.ColumnDefinition_UINT64, Key: false},
	})
	if err != nil {
		return nil, err
	}

	err = stub.CreateTable("state_prov", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "state_prov_id", Type: shim.ColumnDefinition_INT32, Key: true},
		&shim.ColumnDefinition{Name: "cntry_id", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "state_prov_nm", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "state_prov_cd", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_userid", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_dttm", Type: shim.ColumnDefinition_UINT64, Key: false},
	})
	if err != nil {
		return nil, err
	}

	err = stub.CreateTable("city", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "city_id", Type: shim.ColumnDefinition_INT32, Key: true},
		&shim.ColumnDefinition{Name: "state_prov_id", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "city_nm", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_userid", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "txn_dttm", Type: shim.ColumnDefinition_UINT64, Key: false},
	})
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
	} else if function == "delete" {
		//return t.Delete(stub, args)
	} else if function == "search" {
		return t.Search(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	return nil, errors.New("Received unknown function invocation")
}

// Query -call for checking blockchain value
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query")
}

// read- function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	// GEt state from ledger
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
		return nil, errors.New(err.Error())
	}
	return nil, nil
}

//makeTimestamp time indexing
func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// Search function for search comparison & matching ..."search" "prodID"
func (t *SimpleChaincode) Search(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key string
	//key = prod_id

	fmt.Println("running Search()")
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Search", r)
		}
	}()

	key = args[0] // ProdID
	//co.Timestamp = makeTimestamp()

	//Select all columns, from named tables JOIN between tables
	tablepa, err := GetTable("product_alert")
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tablecp, err := GetTable("consumer_product")
	if err != nil {
		return nil, errors.New(err.Error())
	}

	//On key found equal to prod_id in both tables; define column, get row(s)
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: key}}
	columns = append(columns, col1)

	rowA, err := GetRows("product_alert", columns)
	if err != nil {
		return nil, err
	}

	rowB, err := GetRows("consumer_product", columns)
	if err != nil {
		return nil, err
	}

	//confirm key is equal
	if rowA == rowB {
		key = rowA
		fmt.Printf("the query given key %s is found", key) //print out rows with matching key from both rows that match
		fmt.Printf("%s", rowA)
		fmt.Printf("%s", rowB)
		/*for i, v := rowA {
			fmt.Printf(v)
			//fmt.Printf(rowB)
		}*/
	} else {
		fmt.Printf("the query fails given key", key)
	}

	//Select the manufacturing_id, given the prod_id key from the product table
	tablepd, err := GetTable("product")
	if err != nil {
		return nil, errors.New(err.Error())
	}
	fmt.Printf(tablepd)
	/*col1 = shim.Column{Value: &shim.Column_String_{String_: key}}
	columns = append(columns, col1)*/

	row3, err := GetRows("product", columns)
	if err != nil {
		return nil, err
	}
	fmt.Printf(row3)

	//Insert keys (x4) into consumer product alert table
	state, err := InsertRow("consumer_product_alert", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: manufacturer_id}},
			&shim.Column{Value: &shim.Column_String_{String_: prod_id}},
			&shim.Column{Value: &shim.Column_String_{String_: product_alert_id}},
			&shim.Column{Value: &shim.Column_String_{String_: consumer_id}},
		}})
	if !state && err != nil {
		return nil, err

	}

	tablecpa, err := GetTable("consumer_product_alert")
	if err != nil {
		return nil, errors.New(err.Error())
	}
	fmt.Printf(table_cpa)

	return nil, errors.New("Received unknown search function invocation")
}

// Customer - invoke function to write for customers
/*func (t *SimpleChaincode) Customer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("running customer()")

	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8. Name of the key and six values to set")
	}

	city := args[0]                    //city
	provst := args[1]                  //province/state
	sex := args[2]                     //gender
	email := args[3]                   //email
	phone := args[4]                   //phone number
	dob := args[5]                     //DoB
	prodid := args[6]                  //product ID
	imei, err := strconv.Atoi(args[7]) // IMEI
	if err != nil {
		return nil, errors.New("Expecting integer value for IMEI value")
	}
	Timestamp := makeTimestamp()

	//Check if record already exist
	customAsBytes, err := stub.GetState(prodid)
	if err != nil {
		return nil, errors.New("Failed to get Customer Record")
	} else if customAsBytes != nil {
		fmt.Println("This customer record already exist:" + prodid)
		return nil, errors.New("This customer record already exist:" + prodid)
	}

	cust := &Customer{city, provst, sex, email, phone, dob, prodid, imei, Timestamp}
	//str := `{"company": "` + city + `", "model": "` + provst + `", "description": "` + sex + `", "E-mail": "` + email + `", "Phone ":"` + phone + `", "Date of Birth ":"` + dob + `", "Product ID": "` + prodid + `", "IMEI": "` + imei + `"}`
	customJSONasBytes, err := json.Marshal(cust)
	err = stub.PutState(prodid, customJSONasBytes) //write the variable into the chaincode state
	if err != nil {
		return nil, errors.New(err.Error())
	}

	//save record index to state
	err = stub.PutState(prodid, customJSONasBytes)
	if err != nil {
		return nil, errors.New("Failed to get C record index")
	}
	//index key/value range queries
	indexName := "imei~id"
	imeidIndexKey, err := stub.CreateCompositeKey(indexName, []string{cust.ProdID, string(cust.IMEI)})
	if err != nil {
		return nil, errors.New(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(imeidIndexKey, value)

	// ==== Marble saved and indexed. Return success ====
	fmt.Println("- end init marble")

	return nil, nil
}

// Manufacture - function for manufacture record input
func (t *SimpleChaincode) Manufacture(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("running Manufacture()")

	if len(args) != 7 {
		return nil, errors.New("Incorrect number of arguments. Expecting 7")
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
		return nil, errors.New("Failed to get Manufacture/Retail Record")
	} else if ManufAsBytes != nil {
		fmt.Println("This manufacture record already exist: " + key)
		return nil, errors.New("This record already exisit: " + key)
	}

	manuf := &ManufRetail{Co, Model, Description, key, Serial, PurchaseDate, Recall, Timestamp}
	//str := `{"company": "` + co + `", "model": "` + mod + `", "description": "` + des + `", "product ID": "` + id + `", "serial ":"` + serial + `", "purchace data": "` + purc + `", "Recall": "` + recal + `"}`
	manufJSONasBytes, err := json.Marshal(manuf)
	err = stub.PutState(key, manufJSONasBytes) //write the variable into the chaincode state
	if err != nil {
		return nil, errors.New(err.Error())
	}

	//save record index to state
	err = stub.PutState(key, manufJSONasBytes)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	//Index range queries
	indexName := "co~id"
	coidIndexKey, err := stub.CreateCompositeKey(indexName, []string{manuf.Co, manuf.ID})
	if err != nil {
		return nil, errors.New(err.Error())
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
func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var customJSON Customer
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	recordName := args[0]

	// to maintain the record~name index, we need to read the record first and get its ID
	valAsbytes, err := stub.GetState(recordName) //get the record from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + recordName + "\"}"
		return nil, errors.New(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + recordName + "\"}"
		return nil, errors.New(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &customJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + recordName + "\"}"
		return nil, errors.New(jsonResp)
	}

	err = stub.DelState(recordName) //remove the marble from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state:" + err.Error())
	}

	// maintain the index
	indexName := "co~id"
	imeidIndexKey, err := stub.CreateCompositeKey(indexName, []string{string(customJSON.IMEI), customJSON.ProdID})
	if err != nil {
		return nil, errors.New(err.Error())
	}

	//  Delete index entry to state.
	err = stub.DelState(imeidIndexKey)
	if err != nil {
		return nil, errors.New("Failed to delete state:" + err.Error())
	}
	return nil, nil
}*/
