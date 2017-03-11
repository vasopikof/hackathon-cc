package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Property struct {
	id    int
	value string
}

type Iot struct {
	id       int    `json:"IOT_ID"`
	model    string `json:"MODEL"`
	property string `json:"PROPERTY"`
	id_event int    `json:"EVENT_ID"`
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
	//just prepare for search
	var all_event AllEvent

	//init the values
	jsonAsBytes, _ := json.Marshal(all_event)
	err := stub.PutState(event_key, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) Engrave(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//just prepare for search

	//init the values
	jsonAsBytes, _ := json.Marshal(args[0])
	err := stub.PutState(args[0], jsonAsBytes)
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
	}else if(function =="Engrave"){
		return t.Engrave(stub,args);
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
		return t.Read(stub, args)
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) PutEvent(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var event Event
	var err error
	fmt.Println("running PutEvent()")

	if len(args) != 8 {
		return nil, errors.New("[PutEvent] Incorrect number of arguments. Expecting 11.")
	}

	//put all parameters to event
	event.id = args[0]
	event.id_car = args[1]
	event.owner = args[2]
	event.day_code = args[3]
	event.location = args[4]
	event.image = args[5]
	event.describe = args[6]
	event.iot = args[7]

	//split Iot informations, get the number of IOTs
	iot_infos := strings.Split(event.iot, "|")
	fmt.Printf("There are %d IOTs.", len(iot_infos))

	//save event to BlockChain
	tmpBytes, err := stub.GetState(event_key)
	if err != nil {
		return nil, errors.New("Failed to get events")
	}
	var all_events AllEvent
	json.Unmarshal(tmpBytes, &all_events)
	all_events.events = append(all_events.events, event)
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

func (t *SimpleChaincode) Read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fcn = args[0]
	valAsbytes, err := stub.GetState(fcn)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}
	return valAsbytes, nil 
}
