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
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("mylogger")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var consumProdIndexStr = "_consumprodindex" //name for the key/value that will store a list of all known consumer product records
var consumerIndexStr = "_consuindex"        //name for the key/value that will store a list of all known consumer records
var alertIndexStr = "_alertindex"           //name for the key/value that will store a list of all known consumer product alert records

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
	mySQLUsername = "root"
	mySQLPassword = "Test12321"
	mySQLDatabase = "arsystems"
)

// Consumer record
type Consumer struct {
	Consumerid      int64  `json:"00001"`
	ID              int64  `json:"idnumber"`
	Dobdttm         string `json:"01041970"`
	Gendercd        string `json:"female"`
	Cntryid         int    `json:"00"`
	Stateprovid     int    `json:"11"`
	Cityid          int    `json:"01"`
	Postalzipcdabbr string `json:"m1r"`
	Txnuserid       string `json:"barbarella"`
	Txndttm         string `json: "txmdttm"`
}

// Manufacturer record
type Manufacturer struct {
	Manufactureid  string `json:"100000"`
	Manufacturemn  string `json:"100001"`
	Manufactureurl string `json:"100002"`
	Phonenmbrtxt   string `json:"4165559999"`
	Txnuserid      string `json:"barbarella2"`
	Txndttm        string `json: "txndttm2"`
}

// Product record
type Product struct {
	Prodid        string `json:"0000"`
	Upcnum        string `json:"11113"`
	Prodnm        string `json:"111151"`
	Proddesc      string `json:"111152"`
	Stylebrandnm  string `json:"11116"`
	Prodpic       byte   `json:""`
	Manufactureid string `json:"100000"`
	Txnuserid     string `json:"barbarella3"`
	Txndttm       string `json: "txndttm3"`
}

// ConsumerProductAlert record
type ConsumerProductAlert struct {
	Productalertid string `json:"11114"`
	Consumerid     string `json:"10105"`
	Prodid         string `json:"11110"`
	Manufactureid  string `json:"11112"`
	Txnuserid      string `json:"barbarella4"`
	Txndttm        string `json: "txndttm4"`
}

// ProductAlert RECORD
type ProductAlert struct {
	Productalertid       string `json:"11111"`
	Prodid               string `json:"18111"`
	Arsprodalertmsg      string `json:"5678"`
	Prodalertserveritycd string `json:"1234"`
	Healthycdnurl        string `json:"11711"`
	Txnuserid            string `json:"barbara"`
	Txndttm              string `json: "txndttm5"`
}

// ConsumerProduct record -"purchase" table
type ConsumerProduct struct {
	Consumerid   string `json:"1010135435"`
	Prodid       string `json:"11111099898"`
	Totalunitnum string `json:"total"`
	Unitprice    string `json:"1000.00"`
	Purchasedate string `json:"01252015"`
	Sellercd     string `json:"6940385034"`
	Txnuserid    string `json:"barbarella"`
	TxnDttm      string `json: "txndttm6"`
}

/*
type country struct {
	cntryid   string `json:"001"`
	cntrynm   string `json:"001"`
	cntryabbr string `json:"001"`
	cntrycd   string `json:"001"`
	txnuserid string `json:"barbarella"`
	txndttm   string `json: "01081980"`
}

type stateprov struct {
	stateprovid string `json:"00000"`
	cntryid     string `json:"001"`
	stateprovNm string `json: "001"`
	stateprovcd string `json: "111"`
	txnuserid   string `json:"barbarella"`
	txndttm     int    `json: 01081980`
}

type city struct {
	cityid      string `json:"02"`
	stateprovid string `json:"ontario"`
	citynm      string `json:"01"`
	txnuserid   string `json:"barbarella"`
	txndttm     int    `json: 01081980`
}*/

type customEvent struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func main() {
	lld, _ := shim.LogLevel("DEBUG")
	fmt.Println(lld)

	logger.SetLevel(lld)
	fmt.Println(logger.IsEnabledFor(lld))

	err := shim.Start(new(SimpleChaincode)) // Start chaincode system
	if err != nil {
		//fmt.Printf("Error starting Simple chaincode: %s", err)
		logger.Error("Could not start SampleChaincode")
	} else {
		logger.Info("SampleChaincode successfully started")
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
			fmt.Println("Recovered in init \n\r", r)
		}
	}()

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty) //marshal an emtpy array of strings to clear the index
	err = stub.PutState(consumerIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	/*	// Load data from SQL
		fmt.Println("-> mysqlSSHtunnel")

		tunnel := sshTunnel() // Initialize sshTunnel
		go tunnel.Start()     // Start the sshTunnel

		// Declare the dsn (aka database connection string)
		// dsn := "opensim:h0tgrits@tcp(localhost:9000)/opensimdb"
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			mySQLUsername, mySQLPassword, sshLocalHost, sshLocalPort, mySQLDatabase)

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
				unitPrice    float64
				purchaseDate string
				sellerCd     string
				txnUserid    string
				txnDttm      string
			)
			if err := rows.Scan(&consumerID, &prodID, &totalUnitNum, &unitPrice, &purchaseDate, &sellerCd, &txnUserid, &txnDttm); err != nil {
				dbErrorHandler(err)
			}
			fmt.Printf("%d, %d, %d, %.2f, %s, %s, %s, %s\n", consumerID, prodID, totalUnitNum, unitPrice, purchaseDate, sellerCd, txnUserid, txnDttm)
			consumProd := ConsumerProduct{consumerID, prodID, totalUnitNum, unitPrice, purchaseDate, sellerCd, txnUserid, txnDttm}
			consumProdJSONasBytes, err := json.Marshal(consumProd)
			if err != nil {
				return nil, errors.New(err.Error())
			}
			//Place in block
			err = stub.PutState(string(consumerID), consumProdJSONasBytes)
			if err != nil {
				return nil, errors.New(err.Error())
			}
		}
		if err := rows.Err(); err != nil {
			dbErrorHandler(err)
		}*/

	/*	//consumer query
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
			consumer := Consumer{consumerID, imei, dobDttm, genderCD, cntryID, stateProvID, cityID, postalZipCdAbbr, txnUserid, txnDttm}
			consumerJSONasBytes, err := json.Marshal(consumer)
			if err != nil {
				return nil, errors.New(err.Error())
			}
			//Place in block
			err = stub.PutState(string(consumerID), consumerJSONasBytes)
			if err != nil {
				return nil, errors.New(err.Error())
			}
		}
		if err := rowsC.Err(); err != nil {
			dbErrorHandler(err)
		}*/
	fmt.Println("<- mysqlSSHtunnel")

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function, multiple function call API
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.Write(stub, args)
	} else if function == "deleteConsum" {
		return t.DeleteConsum(stub, args)
	} else if function == "search" {
		//return t.Search(stub, args)
	} else if function == "deleteConsumProd" {
		return t.DeleteConsumProd(stub, args)
	} else if function == "ConsumerProduct" {
		return t.ConsumerProduct(stub, args)
	} else if function == "consumer" {
		return t.Consumer(stub, args)
	} else if function == "ConsumerProductAlert" {
		return t.ConsumerProductAlert(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	return nil, errors.New("Received unknown function invocation")
}

// Query -call for checking blockchain value
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	// Handle different functions
	if function == "read" {
		return t.read(stub, args)
	} else if function == "output" {
		return t.output(stub, args)
	} else if function == "Match" {
		return t.output(stub, args)
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
	// Get state from ledger
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

//output modified function to read blockchain and output results to SQL server
func (t *SimpleChaincode) output(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key1, key2, jsonResp string

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the keys to query")
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in init \n\r", r)
		}
	}()

	key1 = args[0]
	key2 = args[1]
	// Get state from ledger, last entry
	valAsbytes, err := stub.GetState(key1)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key1 + "\"}"
		return nil, errors.New(jsonResp)
	}
	if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Nil amount for " + key1 + "\"}"
		return nil, errors.New(jsonResp)
	}
	jsonResp = "{\"Name\":\"" + key1 + "\",\"Value\":\"" + string(valAsbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	valAsbytes, err = stub.GetState(key2)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key2 + "\"}"
		return nil, errors.New(jsonResp)
	}
	if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Nil amount for " + key2 + "\"}"
		return nil, errors.New(jsonResp)
	}
	jsonResp = "{\"Name\":\"" + key2 + "\",\"Value\":\"" + string(valAsbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	keysIter, err := stub.RangeQueryState(key1, key2)
	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
	}
	defer keysIter.Close()

	var keys []string
	for keysIter.HasNext() {
		key, valAsbytes, iterErr := keysIter.Next()
		if iterErr != nil {
			return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
		}
		keys = append(keys, key)
		jsonResp = "{\"Name\":\"" + key + "\",\"Value\":\"" + string(valAsbytes) + "\"}"
		fmt.Printf("Query Response:%s\n", jsonResp)
	}
	//Publishing custom event
	var event = customEvent{"Output", "Successfully read records"}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		return nil, err
	}
	err = stub.SetEvent("evtSender", eventBytes)
	if err != nil {
		fmt.Println("Could not set event for Match record creation", err)
	}
	return valAsbytes, nil
}

// Match function for matching ...consumer "purchase"record and match given "recall alert request data" consumer product alert table
func (t *SimpleChaincode) Match(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key1, key2, jsonResp string
	fmt.Println("running Match()")
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Search", r)
		}
	}()

	key1 = args[0] //consumerID
	key2 = args[1] //product ID
	//key3 = args[2] //product alert ID
	//key4 = args[3] // manufacture ID
	// Get state from ledger, last entry
	fmt.Println("Input key check")
	valAsbytes, err := stub.GetState(key1)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key1 + "\"}"
		return nil, errors.New(jsonResp)
	}
	if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Nil amount for " + key1 + "\"}"
		return nil, errors.New(jsonResp)
	}
	jsonResp = "{\"Name\":\"" + key1 + "\",\"Value\":\"" + string(valAsbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	valAsbytes, err = stub.GetState(key2)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key2 + "\"}"
		return nil, errors.New(jsonResp)
	}
	if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Nil amount for " + key2 + "\"}"
		return nil, errors.New(jsonResp)
	}
	jsonResp = "{\"Name\":\"" + key2 + "\",\"Value\":\"" + string(valAsbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	fmt.Println("Search between product ID and consumer ID")
	keysIter, err := stub.RangeQueryState(key1, key2)
	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
	}
	defer keysIter.Close()

	var keys []string
	for keysIter.HasNext() {
		key, valAsbytes, iterErr := keysIter.Next()
		if iterErr != nil {
			return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
		}
		keys = append(keys, key)
		jsonResp = "{\"Name\":\"" + key + "\",\"Value\":\"" + string(valAsbytes) + "\"}"
		fmt.Printf("Query Response:%s\n", jsonResp)
	}
	//Publishing custom event
	var event = customEvent{"Match", "Successfully found Match between records"}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		return nil, err
	}
	err = stub.SetEvent("evtSender", eventBytes)
	if err != nil {
		fmt.Println("Could not set event for Match record creation", err)
	}
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
func (t *SimpleChaincode) DeleteConsum(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var consumJSON Consumer
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	recordName := args[0]
	err := stub.DelState(recordName) //remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the record index
	consumerAsBytes, err := stub.GetState(consumerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get record index")
	}
	var CIndex []string
	json.Unmarshal(consumerAsBytes, &consumJSON) //un stringify it aka JSON.parse()

	//remove record from index
	for i, val := range CIndex {
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + recordName)
		if val == recordName { //find the correct record
			fmt.Println("found record")
			CIndex = append(CIndex[:i], CIndex[i+1:]...) //remove it
			for x := range CIndex {                      //debug prints...
				fmt.Println(string(x) + " - " + CIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(CIndex) //save new index
	err = stub.PutState(consumerIndexStr, jsonAsBytes)

	//Publishing custom event
	var event = customEvent{"DeleteC", "Successfully deleted consumer record" + recordName}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		return nil, err
	}
	err = stub.SetEvent("evtSender", eventBytes)
	if err != nil {
		fmt.Println("Could not set event for DelectConsumer record ", err)
	}
	return nil, nil
}

// DeleteConsumProd -removes a record key/value pair from state
func (t *SimpleChaincode) DeleteConsumProd(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var conprodJSON ConsumerProduct
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	recordName := args[0]
	err := stub.DelState(recordName) //remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the record index
	consumProdAsBytes, err := stub.GetState(consumProdIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get record index")
	}
	var CPIndex []string
	json.Unmarshal(consumProdAsBytes, &conprodJSON) //un stringify it aka JSON.parse()

	//remove record from index
	for i, val := range CPIndex {
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + recordName)
		if val == recordName { //find the correct record
			fmt.Println("found record")
			CPIndex = append(CPIndex[:i], CPIndex[i+1:]...) //remove it
			for x := range CPIndex {                        //debug prints...
				fmt.Println(string(x) + " - " + CPIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(CPIndex) //save new index
	err = stub.PutState(consumerIndexStr, jsonAsBytes)

	//Publishing custom event
	var event = customEvent{"DeleteCP", "Successfully deleted consumer product record"}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		return nil, err
	}
	err = stub.SetEvent("evtSender", eventBytes)
	if err != nil {
		fmt.Println("Could not set event for DelectConsumProd", err)
	}
	return nil, nil
}

//makeTimestamp time indexing
func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

//ConsumerProduct ..input function
func (t *SimpleChaincode) ConsumerProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("Loading Consumer Product Record...")

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Consumer Product \n\r", r)
		}
	}()

	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}

	consumerid := args[0]
	/*if err != nil {
		return nil, errors.New("1st argument must be a numeric string")
	}*/
	prodID := args[1]
	/*if err != nil {
		return nil, errors.New("2nd argument must be a numeric string")
	}*/
	totalunitnum := args[2]
	/*if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}*/
	unitprice := args[3]
	//unitprice, err := strconv.ParseFloat(args[3], 32)
	/*if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}*/
	purchaseDate := args[4]
	sellerCD := args[5]
	txnUserID := args[6]
	txndttm := args[7]

	//check for duplicate object
	custProdAsBytes, err := stub.GetState(string(prodID))
	if err != nil {
		return nil, errors.New("Failed to get Consumer Product Record " + err.Error())
	} else if custProdAsBytes != nil {
		fmt.Println("Consumer Product Record already exists " + string(prodID))
		return nil, errors.New("Consumer Product Record already exists " + string(prodID))
	}

	// Create object and marshal to JSON
	consumerProduct := &ConsumerProduct{consumerid, prodID, totalunitnum, unitprice, purchaseDate, sellerCD, txnUserID, txndttm}
	consumerProdJSONasBytes, err := json.Marshal(consumerProduct)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// === Save record to state ===
	err = stub.PutState(prodID, consumerProdJSONasBytes)
	//err = stub.PutState(strconv.FormatInt(prodID, 64), consumerProdJSONasBytes)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	//get the record index
	consumProdAsBytes, err := stub.GetState(consumProdIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get record index")
	}
	var CPIndex []string
	json.Unmarshal(consumProdAsBytes, &CPIndex) //un stringify it aka JSON.parse()

	//append
	//prodID := strconv.FormatInt(prodid, 10)
	CPIndex = append(CPIndex, prodID) //add record name to index list
	fmt.Println("! consumer product index: ", CPIndex)
	jsonAsBytes, _ := json.Marshal(CPIndex)
	err = stub.PutState(consumProdIndexStr, jsonAsBytes) //store name of record

	// ==== objects(records) saved and indexed. Return success ====
	fmt.Println("- end init objects")
	return nil, nil
}

//Consumer - input write function
func (t *SimpleChaincode) Consumer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//function, args := stub.GetFunctionAndParameters()
	fmt.Println("Loading Consumer record...")

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Consumer \n\r", r)
		}
	}()

	if len(args) != 10 {
		return nil, errors.New("Incorrect number of arguments. Expecting 10")
	}

	consumerid, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return nil, errors.New("1st argument must be a numeric string")
	}
	id, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return nil, errors.New("2nd argument must be a numeric string")
	}
	dobdttm := args[2]
	gendercd := args[3]
	cntryid, err := strconv.Atoi(args[4])
	if err != nil {
		return nil, errors.New("5th argument must be a numeric string")
	}
	stateprovid, err := strconv.Atoi(args[5])
	if err != nil {
		return nil, errors.New("6th argument must be a numeric string")
	}
	cityid, err := strconv.Atoi(args[6])
	if err != nil {
		return nil, errors.New("7th argument must be a numeric string")
	}
	postalzipcdabbr := args[7]
	txnUserID := args[8]
	txndttm := args[9]

	//check for duplicate object
	consumerID := strconv.FormatInt(consumerid, 10)
	ConsumerAsBytes, err := stub.GetState(consumerID)
	if err != nil {
		return nil, errors.New("Failed to get Consumer Product Record" + err.Error())
	} else if ConsumerAsBytes != nil {
		fmt.Println("Consumer Product Record already exists" + string(consumerid))
		return nil, errors.New("Consumer Product Record already exists" + string(consumerid))
	}

	// Create object and marshal to JSON
	consumerRecord := &Consumer{consumerid, id, dobdttm, gendercd, cntryid, stateprovid, cityid, postalzipcdabbr, txnUserID, txndttm}
	consumerJSONasBytes, err := json.Marshal(consumerRecord)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// === Save record to state ===
	err = stub.PutState(consumerID, consumerJSONasBytes)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	//get the record index
	consumerAsBytes, err := stub.GetState(consumerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get record index")
	}
	var CIndex []string
	json.Unmarshal(consumerAsBytes, &CIndex) //un stringify it aka JSON.parse()

	//append
	CIndex = append(CIndex, consumerID) //add record name to index list
	fmt.Println("! consumer index: ", CIndex)
	jsonAsBytes, _ := json.Marshal(CIndex)
	err = stub.PutState(consumerIndexStr, jsonAsBytes) //store name of record

	// ==== objects(records) saved and indexed. Return success ====
	fmt.Println("- end init objects")
	return nil, nil
}

//ConsumerProductAlert recall data table
func (t *SimpleChaincode) ConsumerProductAlert(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Loading Consumer Product Alert record...")

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Consumer Product Alert\n\r", r)
		}
	}()

	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}
	productalertid := args[0]
	consumerid := args[1]
	prodid := args[2]
	manufactureid := args[3]
	txnUserID := args[4]
	txndttm := args[5]

	//check for duplicate object)
	alertAsBytes, err := stub.GetState(productalertid)
	if err != nil {
		return nil, errors.New("Failed to get Consumer Product Record" + err.Error())
	} else if alertAsBytes != nil {
		fmt.Println("Consumer Product Record already exists" + string(productalertid))
		return nil, errors.New("Consumer Product Record already exists" + string(productalertid))
	}

	// Create object and marshal to JSON
	alertRecord := &ConsumerProductAlert{productalertid, consumerid, prodid, manufactureid, txnUserID, txndttm}
	alertJSONasBytes, err := json.Marshal(alertRecord)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// === Save record to state ===
	err = stub.PutState(productalertid, alertJSONasBytes)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	//get the record index
	alertprodAsBytes, err := stub.GetState(alertIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get record index")
	}
	var CIndex []string
	json.Unmarshal(alertprodAsBytes, &CIndex) //un stringify it aka JSON.parse()

	//append
	CIndex = append(CIndex, productalertid) //add record name to index list
	fmt.Println("! consumer index: ", CIndex)
	jsonAsBytes, _ := json.Marshal(CIndex)
	err = stub.PutState(alertIndexStr, jsonAsBytes) //store name of record

	// ==== objects(records) saved and indexed. Return success ====
	fmt.Println("- end product alert objects")
	return nil, nil
}

//dbErrorHandler - Simple mySql error handling (yet to implement)
func dbErrorHandler(err error) {
	switch err := err.(type) {
	default:
		fmt.Printf("Error %s\n", err)
		os.Exit(-1)
	}
}

// Endpoint - Define an endpoint with ip and port
type Endpoint struct {
	Host string
	Port int
}

// Returns an endpoint as ip:port formatted string
func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// SSHtunnel - Define the endpoints along the tunnel
type SSHtunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint
	Config *ssh.ClientConfig
}

// Start the tunnel
func (tunnel *SSHtunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

// Port forwarding
func (tunnel *SSHtunnel) forward(localConn net.Conn) {
	// Establish connection to the intermediate server
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
		return
	}

	// access the target server
	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	// Transfer the data between  and the remote server
	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

// DecryptPEMkey : Decrypt encrypted PEM key data with a passphrase and embed it to key prefix
// and postfix header data to make it valid for further private key parsing.
func DecryptPEMkey(buffer []byte, passphrase string) []byte {
	block, _ := pem.Decode(buffer)
	der, err := x509.DecryptPEMBlock(block, []byte(passphrase))
	if err != nil {
		fmt.Println("decrypt failed: ", err)
	}
	encoded := base64.StdEncoding.EncodeToString(der)
	encoded = "-----BEGIN RSA PRIVATE KEY-----\n" + encoded +
		"\n-----END RSA PRIVATE KEY-----\n"
	return []byte(encoded)
}

// PublicKeyFile - Get the signers from the OpenSSH key file (.pem) and return them for use in
// the Authentication method. Decrypt encrypted key data with the passphrase.*/
func PublicKeyFile(file string, passphrase string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	if bytes.Contains(buffer, []byte("ENCRYPTED")) {
		// Decrypt the key with the passphrase if it has been encrypted
		buffer = DecryptPEMkey(buffer, passphrase)
	}

	// Get the signers from the key
	signers, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(signers)
}

// Define the ssh tunnel using its endpoint and config data
func sshTunnel() *SSHtunnel {
	localEndpoint := &Endpoint{
		Host: sshLocalHost,
		Port: sshLocalPort,
	}

	serverEndpoint := &Endpoint{
		Host: sshServerHost,
		Port: sshServerPort,
	}

	remoteEndpoint := &Endpoint{
		Host: sshRemoteHost,
		Port: sshRemotePort,
	}

	sshConfig := &ssh.ClientConfig{
		User: sshUserName,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(sshPrivateKeyFile, sshKeyPassphrase)},
	}

	return &SSHtunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}
}
