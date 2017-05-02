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
	//"bytes"
	//"crypto/x509"
	"database/sql"
	//"encoding/base64"
	//"encoding/pem"
	//"io"
	//"io/ioutil"
	//"net"
	//"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hyperledger/fabric/core/chaincode/shim"
m	//pb "github.com/hyperledger/fabric/protos/peer"
)

// Declare your connection data and user credentials here
const (
	// ssh connection related data
	sshServerHost     = "208.75.75.70"
	sshServerPort     = 22
	sshUserName       = "ubuntu"
	sshPrivateKeyFile = "ars_key_pair.pem" // exported as OpenSSH key from .ppk
	sshKeyPassphrase  = "Test12321"        // key file encrytion password

	// ssh tunneling related data
	sshLocalHost  = "localhost" // local localhost ip (client side)
	sshLocalPort  = 9000        // local port used to forward the connection
	sshRemoteHost = "127.0.0.1" // remote local ip (server side)
	sshRemotePort = 3306        // remote MySQL port

	// MySQL access data
	mySqlUsername = "root"
	mySqlPassword = "Test12321"
	mySqlDatabase = "arsystems"
)

type consumer struct {
	consumerid      string `json:"00001"`
	imei            string `json:"000000000000001"`
	dobdtm          string `json:"01/041970"`
	gendercd        string `json:`
	cntryid         string `json:"00"`
	stateprovid     string `json:"11"`
	cityid          string `json:"01"`
	postalzipcdabbr string `json:"m1r1a2`
	txnuserid       string `json:"barbarella"`
	txndttm         string `json: "01/08/1980"`
}

/*type manufacture struct {
	manufactureid string `json:"100000"`
	manufacturemn string `json:"100000"`
	manufactureurl string `json:"100000"`
	phonenmbrtxt string `json:"4165559999"`
	txnuserid string `json:"barbarella"`
	txndttm   string `json: "01081980"`
}

type product struct {
	prodid string `json:"11111"`
	upcnum string `json:"11111"`
	prodnm string `json:"11111"`
	proddesc string `json:"11111"`
	stylebrandnm string `json:"11111"`
	prodpic byte `json:""`
	manufactureid string `json:"100000"`
	txnuserid string `json:"barbarella"`
	txndttm   string `json: "01081980"`
}

type consumer_product_alert struct {
	productalertid string `json:"11111"`
	consumerid string `json:"10101"`
	prodid string `json:"11111"`
	manufactureid string `json:"11111"`
	txnuserid string `json:"barbarella"`
	txndttm   string `json: "01081980"`
}

type product_alert struct {
	productalertid string `json:"11111"`
	prodid string `json:"11111"`
	arsprodalertmsg string `json:"11111"`
	prodalertserveritycd string `json:"11111"`
	healthycdnurl string `json:"11111"`
	txnuserid string `json:"barbarella"`
	txndttm   string `json: "01081980"`
}*/

type consumer_product struct {
	consumerid   string `json:"10101"`
	prodid       string `json:"11111"`
	totalunitnum string `json:"12345"`
	unitprice    string `json:"$1000"`
	purchasedate string `json:"01252015"`
	sellercd     string `json:"6940385034"`
	txnuserid    string `json:"barbarella"`
	txndttm      string `json: "01081980"`
}

type country struct {
	cntryid   string `json:"001"`
	cntrynm   string `json:"001"`
	cntryabbr string `json:"001"`
	cntrycd   string `json:"001"`
	txnuserid string `json:"barbarella"`
	txndttm   string `json: "01081980"`
}

type stateprov struct {
	stateprovid  string `json:"00000"`
	cntryid      string `json:"001"`
	stateprov_nm string `json: "001"`
	stateprovcd  string `json: "111"`
	txnuserid    string `json:"barbarella"`
	txndttm      int    `json: 01081980"`
}

type city struct {
	cityid      string `json:"02"`
	stateprovid string `json:"ontario"`
	citynm      string `json:"01"`
	txnuserid   string `json:"barbarella"`
	txndttm     int    `json: 01081980"`
}

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
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) ([]byte, error) {
	_, args := stub.GetFunctionAndParameters()
	fmt.Println("Initializing system")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("Startup Sequence", []byte(args[0])) //start up test		//std hello_world
	if err != nil {
		return nil, shim.Error
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in init %s\n\r", r)
		}
	}()

	// Load data from SQL
	dbsshcom()
}

// Invoke is our entry point to invoke a chaincode function, multiple function call API
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) ([]byte, error) {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init")
	} else if function == "write" {
		return t.Write(stub, args)
	} else if function == "deleteConsum" {
		return t.DeleteConsum(stub, args)
	} else if function == "search" {
		return t.Search(stub, args)
	} else if function == "Query" { //read a variable
		return t.Query(stub, args)
	} else if function == "delecteConsumProd" {
		return t.DeleteConsProd(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	return shim.Error("Received unknown function invocation")
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

// read - function to read key/value pair
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

// DeleteConsum -removes a record key/value pair from state
func (t *SimpleChaincode) DeleteConsum(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var consumJSON consumer
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	recordName := args[0]

	// to maintain the color~name index, we need to read the marble first and get its color
	valAsbytes, err := stub.GetState(recordName) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + recordName + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + recordName + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &consumJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + recordName + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(recordName) //remove the marble from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	// maintain the index
	indexName := "color~name"
	recordNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{consumer.consumerid, consumer.cntryid, consumer.stateprovid, consumer.cityid})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = stub.DelState(recordNameIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}
	return shim.Success(nil)
}

// DeleteConsumProd -removes a record key/value pair from state
func (t *SimpleChaincode) DeleteConsumProd(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var conprodJSON consumer_product
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	recordName := args[0]

	// to maintain the color~name index, we need to read the marble first and get its color
	valAsbytes, err := stub.GetState(recordName) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + recordName + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + recordName + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &conprodJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + recordName + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(recordName) //remove the marble from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	// maintain the index
	indexName := "color~name"
	recordNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{conprodJSON.Color, conprodJSON.Name})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = stub.DelState(recordNameIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}
	return shim.Success(nil)
}

//makeTimestamp time indexing
func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

//CustomerProduct
func (t *SimpleChaincode) CustomerProduct(stub shim.ChaincodeInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Loading Customer Product Record...")

	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}

	consumerid := args[0]
	prodid := args[1]
	totalunit := args[2]
	unitprice := args[3]
	purchaseDate := args[4]
	sellerCD := args[5]
	txnUserID := args[6]
	txndttm := args[7]

	//check for duplicate object
	custProdAsBytes, err := stub.GetState(prodid)
	if err != nil {
		return shim.Error("Failed to get Customer Product Record" + err.Error())
	} else if custProdAsBytes != nil {
		fmt.Println("Customer Product Record already exists" + prodid)
		return shim.Error("Customer Product Record already exists" + prodid)
	}

	// Create object and marshal to JSON
	consumerProduct := &consumer_product{consumerid, prodid, totalunitnum, unitprice, purchasedate, sellercd, txnuserid, txndttm}
	consumerProdJSONasBytes, err := json.Marshal(consumerProduct)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save marble to state ===
	err = stub.PutState(prodid, consumerProdJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Index composite key -range based on product & consumer ID
	indexName := "prodconsum_id"
	prodconsumIndexKey, err := stub.CreateCompositeKey(indexName, []string{consumerProduct.prodid, consumerProduct.consumerid})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(prodconsumIndexKey, value)

	// ==== objects(records) saved and indexed. Return success ====
	fmt.Println("- end init objects")
	return shim.Success(nil)
}

//Customer_Product
func (t *SimpleChaincode) consumer(stub shim.ChaincodeInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	fmt.Println("Loading Customer record...")

	if len(args) != 10 {
		return nil, errors.New("Incorrect number of arguments. Expecting 10")
	}

	consumerid := args[0]
	imei := args[1]
	dobdtm := args[2]
	gendercd := args[3]
	cntryid := args[4]
	stateprovid := args[5]
	cityid := args[6]
	postalzipcdabbr := args[7]
	txnUserID := args[8]
	txndttm := args[9]

	//check for duplicate object
	customerAsBytes, err := stub.GetState(consumerid)
	if err != nil {
		return shim.Error("Failed to get Customer Product Record" + err.Error())
	} else if customerAsBytes != nil {
		fmt.Println("Customer Product Record already exists" + consumerid)
		return shim.Error("Customer Product Record already exists" + consumerid)
	}

	// Create object and marshal to JSON
	consumerRecord := &consumer{consumerid, imei, dobdttm, gendercd, stateprovid, cityid, postalzipcdabbr, txnuserid, txndttm}
	consumerJSONasBytes, err := json.Marshal(consumerRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save marble to state ===
	err = stub.PutState(consumerid, consumerJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Index composite key -range based on product & consumer ID
	indexNmae := "consumerloc_id"
	consumerIndexKey, err := stub.CreateCompositeKey(indexName, []string{consumer.consumerid, consumer.cntryid, consumer.stateprovid, consumer.cityid})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(consumerIndexKey, value)

	// ==== objects(records) saved and indexed. Return success ====
	fmt.Println("- end init objects")
	return shim.Success(nil)
}

// ssh connection and query
func dbsshcom() {
	fmt.Println("-> mysqlSSHtunnel")

	tunnel := sshTunnel() // Initialize sshTunnel
	go tunnel.Start()     // Start the sshTunnel

	// Declare the dsn (aka database connection string)
	// dsn := "opensim:h0tgrits@tcp(localhost:9000)/opensimdb"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		mySqlUsername, mySqlPassword, sshLocalHost, sshLocalPort, mySqlDatabase)

	// Open the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		dbErrorHandler(err)
	}
	defer db.Close() // keep it open until we are finished

	stmtOut, err := db.Prepare("SELECT * FROM consumer_product")
	if err != nil {
		fmt.Println("Failed to prepare to read data", err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	stmtOut2, err := db.Prepare("SELECT * FROM consumer")
	if err != nil {
		fmt.Println("Failed to prepare to read data", err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut2.Close()

	// Query to check DB consumer_product
	rows, err := db.Query("SELECT consumer_id as consumer, prod_id as prod, Total_unit_num as unit_num, unit_price, Purchase_date as purch_date, seller_cd as seller, txn_userid as userid, txn_dttm as dttm FROM consumer_product")
	if err != nil {
		dbErrorHandler(err)
	}
	defer rows.Close()

	// Iterate though the rows returned and print them
	for rows.Next() {
		var (
			consumerID   int64
			prodID       int64
			totalUnitNum int
			unitPrice    float32
			purchaseDate string
			sellerCd     string
			txnUserid    string
			txnDttm      string
		)
		if err := rows.Scan(&consumerID, &prodID, &totalUnitNum, &unitPrice, &purchaseDate, &sellerCd, &txnUserid, &txnDttm); err != nil {
			dbErrorHandler(err)
		}
		fmt.Printf("%d, %d, %d, %.2f, %s, %s, %s, %s\n", consumerID, prodID, totalUnitNum, unitPrice, purchaseDate, sellerCd, txnUserid, txnDttm)
		consumProd := consumer_product{consumerID, prodID, totalUnitNum, unitPrice, purchaseDate, sellerCd, txnUserid, txnDttm}
		consumProdJSONasBytes, err := json.Marshal(consumProd)
		if err != nil {
			return shim.Error(err.Error())
		}
		//Place in block
		err = stub.PutState(consumerid, consumProdJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	if err := rows.Err(); err != nil {
		dbErrorHandler(err)
	}

	//consumer query
	rowsC, err := db.Query("SELECT consumer_id as consumer, imei, dob_dttm, gender_cd, cntry_id as cntryid, state_prov_id, city_id, postal_zip_cd_abbr, txn_userid as userid, txn_dttm as dttm FROM consumer WHERE cntry_id='41' ")
	if err != nil {
		dbErrorHandler(err)
	}
	defer rows.Close()

	// Iterate though the rows returned and print them
	for rowsC.Next() {
		var (
			consumerID      int64
			imei            int64
			dobDttm         string
			genderCD        string
			cntryID         int
			stateProvID     int
			cityID          int
			postalZipCdAbbr string
			txnUserid       string
			txnDttm         string
		)
		if err := rowsC.Scan(&consumerID, &imei, &dobDttm, &genderCD, &cntryID, &stateProvID, &cityID, &postalZipCdAbbr, &txnUserid, &txnDttm); err != nil {
			dbErrorHandler(err)
		}
		fmt.Printf("%d, %d, %s, %s, %d, %d, %d, %s, %s, %s\n", consumerID, imei, dobDttm, genderCD, cntryID, stateProvID, cityID, postalZipCdAbbr, txnUserid, txnDttm)
		consumer := consumer{consumerID, imei, dobDttm, genderCD, cntryID, stateProvID, cityID, postalZipCdAbbr, txnUserid, txnDttm}
		consumerJSONasBytes, err := json.Marshal(consumer)
		if err != nil {
			return shim.Error(err.Error())
		}
		//Place in block
		err = stub.PutState(consumerid, consumerJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	if err := rowsC.Err(); err != nil {
		dbErrorHandler(err)
	}

	fmt.Println("<- mysqlSSHtunnel")
	return nil
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

	/*key = args[0] // ProdID
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
		}
	} else {
		fmt.Printf("the query fails given key %s", key)
	}

	//Select the manufacturing_id, given the prod_id key from the product table
	tablepd, err := GetTable("product")
	if err != nil {
		return nil, errors.New(err.Error())
	}
	fmt.Printf(tablepd)
	/*col1 = shim.Column{Value: &shim.Column_String_{String_: key}}
	columns = append(columns, col1)

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
	*/

	return nil, errors.New("Received unknown search function invocation")
}
