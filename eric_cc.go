package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"

	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var minimalTxStr = "_minimaltx"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Property struct {
	id string `json:"ID"`
	value string `json:"VALUE"`
}

type Iot struct {
	id string `json:"IOT_ID"`
	model string `json:"MODEL"`
	property string `json:"PROPERTY"`
	id_event string `json:"EVENT_ID"`
}

type A_Event struct {
	id string `json:"EVENT_ID"`
	id_car string `json:"CAR_ID"`
	owner string `json:"OWNER"`
	day_code string `json:"DAY_CODE"`
	location string `json:"LOCATION"`
	image string `json:"IMAGE"`
	describe string `json:"DESCRIBE"`
	iot string `json:"IOT"`
}

type AllEvent struct {
	events []A_Event `json:"EVENTS"`
}

type Transaction struct{
	Id string `json:"txID"`					//user who created the open trade order
	Timestamp string `json:"EX_TIME"`			//utc timestamp of creation
	TraderA string  `json:"USER_A_ID"`				//description of desired marble
	TraderB string  `json:"USER_B_ID"`
	SellerA string  `json:"SELLER_A_ID"`				//description of desired marble
	SellerB string  `json:"SELLER_B_ID"`
	PointA string  `json:"POINT_A"`
	PointB string  `json:"POINT_B"`
	Related []Point `json:"related"`		//array of marbles willing to trade away
}

type AllTx struct{
	TXs []Transaction `json:"tx"`
}

//global variable of indexs and values
var iot_key = "_iot_index"
var event_key = "_event_index"

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}



	//just prepare for search
	var all_event AllEvent

	//init the values
	jsonAsBytes, _ := json.Marshal(all_event)
	err = stub.PutState(event_key, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	err = stub.PutState(minimalTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	//name = args[0]															//rename for funsies
	//value = args[1]

	pp := Property{}

	pp.id = "propid"
	pp.value = "propvalue"
	jsonAsBytes, _ := json.Marshal(pp)
	

	err = stub.PutState(iot_key, jsonAsBytes) //rewrite open orders
	if err != nil {
		return nil, err
	}
	/*err = stub.PutState(name, []byte(value))								//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}*/
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "PutEvent" {
		return t.PutEvent(stub, args)
	} else if function == "write" {											//writes a value to the chaincode state
		return t.Write(stub, args)
	}else if function == "init_transaction" {									//create a new trade order
		return t.init_transaction(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) init_transaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error	
	//	0        1      2     3      4      5       6
	//["bob", "blue", "16", "red", "16"] *"blue", "35*

	/*
	Id string `json:"txID"`					//user who created the open trade order
	Timestamp string `json:"EX_TIME"`			//utc timestamp of creation
	TraderA string  `json:"USER_A_ID"`				//description of desired marble
	TraderB string  `json:"USER_B_ID"`
	SellerA string  `json:"SELLER_A_ID"`				//description of desired marble
	SellerB string  `json:"SELLER_B_ID"`
	PointA string  `json:"POINT_A"`
	PointB string  `json:"POINT_B"`
	Related []Point `json:"related"`
}
	*/


	open := Transaction{}
	open.Id = args[0]
	open.TraderA = args[1]
	open.TraderB = args[2]
	open.SellerA = args[3]
	open.SellerB = args[4]
	open.PointA = args[5]
	open.PointB = args[6]
	open.Timestamp = args[7]
	
	fmt.Println("- start open trade")
	jsonAsBytes, _ := json.Marshal(open)
	err = stub.PutState("_debug1", jsonAsBytes)

	//get the open trade struct
	tradesAsBytes, err := stub.GetState(minimalTxStr)
	if err != nil {
		return nil, errors.New("Failed to get TXs")
	}
	var trades AllTx
	json.Unmarshal(tradesAsBytes, &trades)										//un stringify it aka JSON.parse()
	
	trades.TXs = append(trades.TXs, open);						//append to open trades
	fmt.Println("! appended open to trades")
	jsonAsBytes, _ = json.Marshal(trades)
	err = stub.PutState(minimalTxStr, jsonAsBytes)								//rewrite open orders
	if err != nil {
		return nil, err
	}
	fmt.Println("- end open trade")
	return nil, nil
}


// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "GetTimeline" {
		return t.GetTimeline(stub, args)
	} else if function == "GetInsuranceEvent" {
		return t.GetInsuranceEvent(stub, args)
	} else if function == "read" {
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query: " + function)
}


// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}

func (t *SimpleChaincode) PutEvent(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("running PutEvent()")

	//put all parameters to event
	event_input := A_Event{}

	event_input.id = args[0]
	event_input.id_car = args[1]
	event_input.owner = args[2]
	event_input.day_code = args[3]
	event_input.location = args[4]
	event_input.image = args[5]
	event_input.describe = args[6]
	event_input.iot = args[7]

	
	//split Iot informations, get the number of IOTs
	iot_infos := strings.Split(event_input.iot, "|")
	fmt.Printf("There are %d IOTs.", len(iot_infos))

	//save event to BlockChain
	tmpBytes, err := stub.GetState(event_key)
	if err != nil {
		return nil, errors.New("Failed to get events")
	}
	var all_events AllEvent

	json.Unmarshal(tmpBytes, &all_events)

	all_events.events = append(all_events.events, event_input)
	jsonAsBytes, _ := json.Marshal(all_events)
	

	err = stub.PutState(event_key, jsonAsBytes) //rewrite open orders
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) GetTimeline(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var car_id, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting id of the event to query")
	}

	car_id = args[0]

	//get all of the car_ids here
	tmpBytes, err := stub.GetState(event_key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + car_id + "\"}"
		return nil, errors.New(jsonResp)
	}
	var all_events AllEvent
	json.Unmarshal(tmpBytes, &all_events)
	var processed AllEvent
	for i := range all_events.events {
		event_car_id := all_events.events[i].id_car
		if event_car_id == car_id {
			processed.events = append(processed.events, all_events.events[i])
		}
	}
	jsonAsBytes, _ := json.Marshal(processed)

	return jsonAsBytes, nil
}

func (t *SimpleChaincode) GetInsuranceEvent(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var car_id, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting id of the event to query")
	}

	car_id = args[0]

	//get all of the car_ids here
	tmpBytes, err := stub.GetState(event_key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + car_id + "\"}"
		return nil, errors.New(jsonResp)
	}
	var all_events AllEvent
	json.Unmarshal(tmpBytes, &all_events)
	var processed AllEvent
	for i := range all_events.events {
		event_car_id := all_events.events[i].id_car
		if event_car_id == car_id {
			processed.events = append(processed.events, all_events.events[i])
		}
	}

	jsonAsBytes, _ := json.Marshal(processed)

	return jsonAsBytes, nil
}