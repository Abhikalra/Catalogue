package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	)

type SimpleChaincode struct {
}

////////////////////////////////////////////////////////////////////////////////////////
// main function
// start program execution here
////////////////////////////////////////////////////////////////////////////////////////

func main() {

	fmt.Println("This is a demo student record table chaincode")
    err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Oops !! Something went wrong...Error starting chaincode application: %s", err)
	}

}


// ////////////////////////////////////////////////////////////////////////////////////
// Init Function
// Initializes the chaincode and default parameters
// ///////////////////////////////////////////////////////////////////////////////////
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments in Init(). Expecting 1")
	}

	// Create student record table
	err := stub.CreateTable("Student_Record", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "BannerID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Name", Type: shim.ColumnDefinition_STRING, Key: false},
        	&shim.ColumnDefinition{Name: "Subject1", Type: shim.ColumnDefinition_STRING, Key: false},
        	&shim.ColumnDefinition{Name: "Subject2", Type: shim.ColumnDefinition_STRING, Key: false},
        	&shim.ColumnDefinition{Name: "Subject3", Type: shim.ColumnDefinition_STRING, Key: false},

	})
	if err != nil {
		return nil, errors.New("Failed creating Student record table.")
	}

	return []byte("Initialization complete for " + args[0]), nil
}

// ////////////////////////////////////////////////////////////////////////////////////
// Invoke Function
// Perform transaction with the chaincode
// ///////////////////////////////////////////////////////////////////////////////////


func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "addDetail" {
		// trigger add detail function to add student info
		return t.addDetail(stub,"addDetail", args)
	}else if function == "init" {
		return t.Init(stub,"Init",args)
	} 

	return nil, errors.New("Received unknown function invocation")
}

// ////////////////////////////////////////////////////////////////////////////////////
// Query Function
// Query data on the chaincode
// ///////////////////////////////////////////////////////////////////////////////////
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	
	if function != "getDetail" {
		return nil, errors.New("Invalid query function name. Expecting getDetail")
	}
	if function == "getDetail" {
		return t.getDetail(stub,"getDetail", args)
	}
	return nil, errors.New("Received unknown function in query")
}

// ////////////////////////////////////////////////////////////////////////////////////
// GetDetail Function
// Fetches matching row corresponding to provided key value of student name
// User defined function triggered in Query 
// ///////////////////////////////////////////////////////////////////////////////////

func (t *SimpleChaincode) getDetail(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var err error
	var output string
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of student to fetch data")
	}

	
	name := args[0]

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: name}}
	columns = append(columns, col1)

	row, err := stub.GetRow("Student_Record", columns)
	if err != nil {
				return nil, fmt.Errorf("Failed retriving value")
	}

	
	output = string(row.Columns[0].GetBytes()) + string(row.Columns[1].GetBytes()) + string(row.Columns[2].GetBytes()) + string(row.Columns[3].GetBytes()) + string(row.Columns[4].GetBytes())
	return []byte(output), nil
}


func (t *SimpleChaincode) addDetail(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var err error
	var success string
	var ok bool
	if len(args) !=5 {
		return nil, errors.New("Incorrect number of arguments. Expecting five arguments in addDetail()")
	}

	success = " Record successfully added\nBannerID :"+ args[0] +" Name: "+ args[1] +" Marks1 : " + args[2] +" Marks2 : "+ args[3] +" Marks3 : "+ args[4]
	bannerID := args[0]
	name := args[1]
	marks1 := args[2]
	marks2 := args[3]
	marks3 := args[4]
	
	ok , err = stub.InsertRow("Student_Record", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: bannerID}},
			&shim.Column{Value: &shim.Column_String_{String_: name}},
			&shim.Column{Value: &shim.Column_String_{String_: marks1}},
			&shim.Column{Value: &shim.Column_String_{String_: marks2}},
			&shim.Column{Value: &shim.Column_String_{String_: marks3}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Record already added...")
	}

	return []byte(success), err
}
