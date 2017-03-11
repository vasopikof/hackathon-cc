package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"

	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Property struct {
	id    string `json:"ID"`
	value string `json:"VALUE"`
}

type Iot struct {
	id       string    `json:"IOT_ID"`
	model    string `json:"MODEL"`
	property string `json:"PROPERTY"`
	id_event string    `json:"EVENT_ID"`
}

type Event struct {
	id       string `json:"EVENT_ID"`
	id_car   string `json:"CAR_ID"`
	owner    string `json:"OWNER"`
	day_code string `json:"DAY_CODE"`
	location string `json:"LOCATION"`
	image    string `json:"IMAGE"`
	describe string `json:"DESCRIBE"`
	iot      string `json:"IOT"`
}

type AllEvent struct {
	events []Event `json:"EVENTS"`
}

//global variable of indexs and values
var iot_key = "_iot_key"
var event_key = "_event_key"

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
	return nil, nil
}

func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]															//rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value))								//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
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
	}
	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation: " + function)
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
	event := Event{}
	err = stub.PutState("_debug2", []byte("enter PutState"))
	event.id = args[0]
	event.id_car = args[1]
	event.owner = args[2]
	event.day_code = args[3]
	event.location = args[4]
	event.image = args[5]
	event.describe = args[6]
	event.iot = args[7]

	
	inAsBytes, _ := json.Marshal(event)
	
	err = stub.PutState("_debug0", inAsBytes)

	err = stub.PutState("_debug1", []byte("debug_this "+event.id+" "+event.id_car+" "+event.owner+" "+event.day_code+" "+event.location+" "+event.image+" "+event.describe+" "+event.iot))

	//split Iot informations, get the number of IOTs
	iot_infos := strings.Split(event.iot, "|")
	fmt.Printf("There are %d IOTs.", len(iot_infos))

	//save event to BlockChain
	tmpBytes, err := stub.GetState(event_key)
	if err != nil {
		return nil, errors.New("Failed to get events")
	}
	err = stub.PutState("_debug3", []byte("enter loop"))
	var all_events AllEvent

	json.Unmarshal(tmpBytes, &all_events)

	all_events.events = append(all_events.events, event)
	jsonAsBytes, _ := json.Marshal(all_events)
	
	err = stub.PutState("_debug4", []byte("enter Resulting"))

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