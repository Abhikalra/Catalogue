package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
	"time"
	)


///////////////////////////////////////////////////////////////////////////////////////
// Createing a record of the items available in the inventory
// the items have an associated id, name,description, quantity and price per quantity
///////////////////////////////////////////////////////////////////////////////////////

type ItemObject struct {
	ItemID         string
	ItemName       string
	ItemDesc       string
	ItemPrice      string 
	ItemQuantity   string
	}

////////////////////////////////////////////////////////////////////////////////
// Log of items purchased by the customers
////////////////////////////////////////////////////////////////////////////////
type ItemPurchase struct {
	PurchaseID   string // Purchase Code
    ItemID       string // Item ID
	BuyerID      string // Customer Name
	ItemQuantity string // Amount of items
	ItemCost     string // Total Price
	Date         string // Date when status changed
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// A Map that holds TableNames and the number of Keys
// This information is used to dynamically Create, Update
// Replace , and Query the Ledger
/////////////////////////////////////////////////////////////////////////////////////////////////////

func GetNumberOfKeys(tname string) int {
	TableMap := map[string]int{
		"ItemObject":        1,
		"ItemPurchase":      1,
		}
	return TableMap[tname]
}

//////////////////////////////////////////////////////////////
// Invoke Functions based on Function name
// The function name gets resolved to one of the following calls
// during an invoke
//
//////////////////////////////////////////////////////////////
func InvokeFunction(fname string) func(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	InvokeFunc := map[string]func(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error){
		"AddItem":            AddItem,
		"BuyItem":            BuyItem,
		}
	return InvokeFunc[fname]
}

//////////////////////////////////////////////////////////////
// Query Functions based on Function name
//
//////////////////////////////////////////////////////////////
func QueryFunction(fname string) func(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	QueryFunc := map[string]func(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error){
		"GetItem":               GetItem,
		"GetPurchase":           GetPurchase,
		"GetPurchasebyBuyer":    GetPurchasebyBuyer,
		"GetPurchaselog":        GetPurchaselog,
		"GetItemlog":            GetItemlog,
		}
	return QueryFunc[fname]
}


type SimpleChaincode struct {
}
var gopath string
var ccPath string
////////////////////////////////////////////////////////////////////////////////
// Chain Code Kick-off Main function
////////////////////////////////////////////////////////////////////////////////
func main() {

	// maximize CPU usage for maximum performance
	//runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("This is a demo product catalogue chaincode example......")

	//gopath = os.Getenv("GOPATH")
	//ccPath = fmt.Sprintf("%s/src/github.com/ITPeople-Blockchain/auction/art/artchaincode/", gopath)
	//ccPath = fmt.Sprintf("%s/src/github.com/hyperledger/fabric/examples/chaincode/go/artfun/", gopath)
	// Start the shim -- running the fabric
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("Oops !! Something went wrong.Error starting chaincode application: %s", err)
	
}

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// SimpleChaincode - Init Chaincode implementation - The following sequence of transactions can be used to test the Chaincode
// Initialize the initial data values
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	// TODO - Include all initialization to be complete before Invoke and Query
	// Uses aucTables to delete tables if they exist and re-create them

	//myLogger.Info("[Trade and Auction Application] Init")
	fmt.Println("Simple product Catalogue application INIT")
	var err error

	//for _, val := range aucTables {
	//	err = stub.DeleteTable(val)
	//	if err != nil {
	//		return nil, fmt.Errorf("Init(): DeleteTable of %s  Failed ", val)
	//	}
		err = InitLedger(stub, "ItemObject")
		if err != nil {
			return nil, fmt.Errorf("Init(): InitLedger of %s  Failed ", val)
		}
        err = InitLedger(stub, "ItemPurchase")
		if err != nil {
			return nil, fmt.Errorf("Init(): InitLedger of %s  Failed ", val)
		}
	//}

	fmt.Println("Init() Initialization Complete  : ", args)
	return []byte("Init(): Initialization Complete"), nil
}

////////////////////////////////////////////////////////////////
// SimpleChaincode - INVOKE Chaincode implementation
// User Can Invoke
////////////////////////////////////////////////////////////////

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var err error
	var buff []byte

	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)}

		InvokeRequest := InvokeFunction(function)
		if InvokeRequest != nil {
			buff, err = InvokeRequest(stub, function, args)
		}

		fmt.Println("Invoke() Invalid recType : " + args)
		//return nil, errors.New("Invoke() : Invalid recType : " + args[0])

	return buff, err
}

//////////////////////////////////////////////////////////////////////////////////////////
// SimpleChaincode - QUERY Chaincode implementation
//////////////////////////////////////////////////////////////////////////////////////////

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var err error
	var buff []byte
	if len(args) < 1 {
		fmt.Println("Query() : Invalid number of arguments provided")
		return nil, errors.New("Query() : Error occured in Query Function.")
	}

	QueryRequest := QueryFunction(function)
	if QueryRequest != nil {
		buff, err = QueryRequest(stub, function, args)
	} else {
		fmt.Println("Query() Invalid function call : ", function)
		return nil, errors.New("Query() : Invalid function call : " + function)
	}

	if err != nil {
		fmt.Println("Query() Object not found : ", args[0])
		return nil, errors.New("Query() : Object not found : " + args[0])
	}
	return buff, err
}

//////////////////////////////////////////////////////////////////////////////////////////
// Retrieve User Information
// example:
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetUser", "Args": ["100"]}'
//
//////////////////////////////////////////////////////////////////////////////////////////
func GetItem(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	var err error

	// Get the Object and Display it
	Avalbytes, err := QueryLedger(stub, "ItemObject", args[0])
	if err != nil {
		fmt.Println("GetItem() : Failed to Query Object ")
		jsonResp := "{\"Error\":\"Failed to get  Object Data for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		fmt.Println("GetItem() : Incomplete Query Object ")
		jsonResp := "{\"Error\":\"Incomplete information about the key for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Println("GetItem() : Response : Successfull -")
	return Avalbytes, nil
}

//////////////////////////////////////////////////////////////////////////////////////////
// Retrieve User Information
// example:
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetUser", "Args": ["100"]}'
//
//////////////////////////////////////////////////////////////////////////////////////////
func GetPurchase(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	var err error

	// Get the Object and Display it
	Avalbytes, err := QueryLedger(stub, "ItemPurchase", args[0])
	if err != nil {
		fmt.Println("GetPurchase() : Failed to Query Object ")
		jsonResp := "{\"Error\":\"Failed to get  Object Data for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		fmt.Println("GetPurchase() : Incomplete Query Object ")
		jsonResp := "{\"Error\":\"Incomplete information about the key for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Println("GetPurchase() : Response : Successfull -")
	return Avalbytes, nil
}

/////////////////////////////////////////////////////////////////////////////////////////
// Validates The Ownership of an Asset using ItemID, OwnerID, and HashKey
//
// ./peer chaincode query -l golang -n mycc -c '{"Function": "ValidateItemOwnership", "Args": ["1000", "100", "tGEBaZuKUBmwTjzNEyd+nr/fPUASuVJAZ1u7gha5fJg="]}'
//
/////////////////////////////////////////////////////////////////////////////////////////
func GetPurchasebyBuyer(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	var err error

	// Get the Object and Display it
	Avalbytes, err := QueryLedger(stub, "ItemPurchase", args[0])
	if err != nil {
		fmt.Println("GetPurchasebyBuyer() : Failed to Query Object ")
		jsonResp := "{\"Error\":\"Failed to get  Object Data for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		fmt.Println("GetPurchasebyBuyer() : Incomplete Query Object ")
		jsonResp := "{\"Error\":\"Incomplete information about the key for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Println("GetPurchasebyBuyer() : Response : Successfull -")
	return Avalbytes, nil
}

///////////////////////////////////////////////////////
//
//
//////////////////////////////////////////////////////
func GetPurchaselog(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	var err error

	// Get the Object and Display it
	Avalbytes, err := QueryLedger(stub, "ItemPurchase", args[0])
	if err != nil {
		fmt.Println("GetPurchaselog() : Failed to Query Object ")
		jsonResp := "{\"Error\":\"Failed to get  Object Data for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		fmt.Println("GetPurchaselog() : Incomplete Query Object ")
		jsonResp := "{\"Error\":\"Incomplete information about the key for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Println("GetPurchaselog() : Response : Successfull -")
	return Avalbytes, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// Retrieve a Bid based on two keys - AucID, BidNo
// A Bid has two Keys - The Auction Request Number and Bid Number
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetLastBid", "Args": ["1111"], "1"}'
//
///////////////////////////////////////////////////////////////////////////////////////////////////
func GetItemlog(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	var err error

	// Get the Object and Display it
	Avalbytes, err := QueryLedger(stub, "ItemPurchase", args[0])
	if err != nil {
		fmt.Println("GetItemlog() : Failed to Query Object ")
		jsonResp := "{\"Error\":\"Failed to get  Object Data for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		fmt.Println("GetItemlog() : Incomplete Query Object ")
		jsonResp := "{\"Error\":\"Incomplete information about the key for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Println("GetItemlog() : Response : Successfull -")
	return Avalbytes, nil
}



///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Create a User Object. The first step is to have users
// registered
// There are different types of users - Traders (TRD), Auction Houses (AH)
// Shippers (SHP), Insurance Companies (INS), Banks (BNK)
// While this version of the chain code does not enforce strict validation
// the business process recomends validating each persona for the service
// they provide or their participation on the auction blockchain, future enhancements will do that
// ./peer chaincode invoke -l golang -n mycc -c '{"Function": "PostUser", "Args":["100", "USER", "Ashley Hart", "TRD",  "Morrisville Parkway, #216, Morrisville, NC 27560", "9198063535", "ashley@itpeople.com", "SUNTRUST", "00017102345", "0234678"]}'
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func AddItem(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	record, err := []string{arg[0],arg[1],arg[2],arg[3],arg[4],arg[5]} //
	if err != nil {
		return nil, err
	}
//	buff, err := UsertoJSON(record) //

	if err != nil {
		fmt.Println("AddItem() : Failed Cannot create object buffer for write : ", args[1])
		return nil, errors.New("AddItem(): Failed Cannot create object buffer for write : " + args[1])
	} else {
		// Update the ledger with the Buffer Data
		// err = stub.PutState(args[0], buff)
		keys := []string{args[0]}
		err = UpdateLedger(stub, "ItemObject", keys, record)
		if err != nil {
			fmt.Println("AddItem() : write error while inserting record")
			return nil, err
		}
	}

	return record, err
}


/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Create a master Object of the Item
// Since the Owner Changes hands, a record has to be written for each
// Transaction with the updated Encryption Key of the new owner
// Example
//./peer chaincode invoke -l golang -n mycc -c '{"Function": "PostItem", "Args":["1000", "ARTINV", "Shadows by Asppen", "Asppen Messer", "20140202", "Original", "Landscape" , "Canvas", "15 x 15 in", "sample_7.png","$600", "100"]}'
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

func BuyItem(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	record, err := []string{arg[0],arg[1],arg[2],arg[3],arg[4],time.Now().String()} //
	if err != nil {
		return nil, err
	}
//	buff, err := UsertoJSON(record) //

	if err != nil {
		fmt.Println("BuyItem() : Failed Cannot create object buffer for write : ", args[1])
		return nil, errors.New("BuyItem(): Failed Cannot create object buffer for write : " + args[1])
	} else {
		// Update the ledger with the Buffer Data
		// err = stub.PutState(args[0], buff)
		keys := []string{args[0]}
		err = UpdateLedger(stub, "ItemPurchase", keys, record)
		if err != nil {
			fmt.Println("BuyItem() : write error while inserting record")
			return nil, err
		}

		
	}

	return record, err
}


////////////////////////////////////////////////////////////////////////////
// Query a User Object by Table Name and Key
////////////////////////////////////////////////////////////////////////////
func QueryLedger(stub *shim.ChaincodeStub, tableName string, args []string) ([]byte, error) {

	var columns []shim.Column
	nCol := GetNumberOfKeys(tableName)
	for i := 0; i < nCol; i++ {
		colNext := shim.Column{Value: &shim.Column_String_{String_: args[i]}}
		columns = append(columns, colNext)
	}

	row, err := stub.GetRow(tableName, columns)
	fmt.Println("Length or number of rows retrieved ", len(row.Columns))

	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed retrieving data " + args[0] + ". \"}"
		fmt.Println("Error retrieving data record for Key = ", args[0], "Error : ", jsonResp)
		return nil, errors.New(jsonResp)
	}

	//fmt.Println("User Query Response:", row)
	jsonResp := "{\"Output\":\"" + string(row.Columns[nCol].GetBytes()) + "\"}"
	fmt.Println("User Query Response:%s\n", jsonResp)

	Avalbytes := row.Columns[nCol].GetBytes()

	// Perform Any additional processing of data
	fmt.Println("QueryLedger() : Successful - Proceeding to ProcessRequestType ")
	err = ProcessQueryResult(stub, Avalbytes, args)
	if err != nil {
		fmt.Println("QueryLedger() : Cannot create object  : ", args[1])
		jsonResp := "{\"QueryLedger() Error\":\" Cannot create Object for key " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	return Avalbytes, nil
}


////////////////////////////////////////////////////////////////////////////
// Open a Ledgers if one does not exist
// These ledgers will be used to write /  read data
// Use names are listed in aucTables {}
// THIS FUNCTION REPLACES ALL THE INIT Functions below
//  - InitUserReg()
//  - InitAucReg()
//  - InitBidReg()
//  - InitItemReg()
//  - InitItemMaster()
//  - InitTransReg()
//  - InitAuctionTriggerReg()
//  - etc. etc.
////////////////////////////////////////////////////////////////////////////
func InitLedger(stub *shim.ChaincodeStub, tableName string) error {

	// Generic Table Creation Function - requires Table Name and Table Key Entry
	// Create Table - Get number of Keys the tables supports
	// This version assumes all Keys are String and the Data is Bytes
	// This Function can replace all other InitLedger function in this app such as InitItemLedger()

	nKeys := GetNumberOfKeys(tableName)
	if nKeys < 1 {
		fmt.Println("Atleast 1 Key must be provided")
		fmt.Println("Auction_Application: Failed creating Table ", tableName)
		return errors.New("Auction_Application: Failed creating Table " + tableName)
	}

	var columnDefsForTbl []*shim.ColumnDefinition

	for i := 0; i < nKeys; i++ {
		columnDef := shim.ColumnDefinition{Name: "keyName" + strconv.Itoa(i), Type: shim.ColumnDefinition_STRING, Key: true}
		columnDefsForTbl = append(columnDefsForTbl, &columnDef)
	}

	columnLastTblDef := shim.ColumnDefinition{Name: "Details", Type: shim.ColumnDefinition_BYTES, Key: false}
	columnDefsForTbl = append(columnDefsForTbl, &columnLastTblDef)

	// Create the Table (Nil is returned if the Table exists or if the table is created successfully
	err := stub.CreateTable(tableName, columnDefsForTbl)

	if err != nil {
		fmt.Println("Auction_Application: Failed creating Table ", tableName)
		return errors.New("Auction_Application: Failed creating Table " + tableName)
	}

	return err
}

////////////////////////////////////////////////////////////////////////////
// Open a User Registration Table if one does not exist
// Register users into this table
////////////////////////////////////////////////////////////////////////////
func UpdateLedger(stub *shim.ChaincodeStub, tableName string, keys []string, args []byte) error {

	nKeys := GetNumberOfKeys(tableName)
	if nKeys < 1 {
		fmt.Println("Atleast 1 Key must be provided \n")
	}

	var columns []*shim.Column

	for i := 0; i < nKeys; i++ {
		col := shim.Column{Value: &shim.Column_String_{String_: keys[i]}}
		columns = append(columns, &col)
	}

	lastCol := shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(args)}}
	columns = append(columns, &lastCol)

	row := shim.Row{columns}
	ok, err := stub.InsertRow(tableName, row)
	if err != nil {
		return fmt.Errorf("UpdateLedger: InsertRow into "+tableName+" Table operation failed. %s", err)
	}
	if !ok {
		return errors.New("UpdateLedger: InsertRow into " + tableName + " Table failed. Row with given key " + keys[0] + " already exists")
	}

	fmt.Println("UpdateLedger: InsertRow into ", tableName, " Table operation Successful. ")
	return nil
}

