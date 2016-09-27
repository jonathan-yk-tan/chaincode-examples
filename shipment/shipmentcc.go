package main

import (
  "errors"
  "fmt"
  "strconv"

  "github.com/hyperledger/fabric/core/chaincode/shim"
)
type ShipmentChaincode struct{}

func (t *ShipmentChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

  // Create shiptment table
	err := stub.CreateTable("Shipment", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "id", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "status", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating Shipment table.")
	}

  return nil, nil
}

func (t *ShipmentChaincode) assign(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error){
  id := args[0]
  status := args[1]
  ok, err = stub.InsertRow("Shipment", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: id}},
			&shim.Column{Value: &shim.Column_Bytes{String: status}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Asset was already assigned.")
	}
}

func (t *ShipmentChaincode) update(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error){
  id := args[0]
  status := args[1]
  err = stub.DeleteRow(
		"Shipment",
		[]shim.Column{shim.Column{Value: &shim.Column_String_{String_: id}}},
	)
	if err != nil {
		return nil, errors.New("Failed deliting row.")
	}

	_, err = stub.InsertRow(
		"Shipment",
		shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: id}},
				&shim.Column{Value: &shim.Column_String{String: status}},
			},
		})
	if err != nil {
		return nil, errors.New("Failed inserting row.")
	}
}

func (t *ShipmentChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
  if function == "update" {
    return t.update(stub,args)
  }
  return nil, errors.New("Received unknown function invocation")

}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
  var err error

	if len(args) != 1 {
		myLogger.Debug("Incorrect number of arguments. Expecting name of an asset to query")
		return nil, errors.New("Incorrect number of arguments. Expecting name of an asset to query")
	}

	// Who is the owner of the asset?
	id := args[0]

  var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: id}}
	columns = append(columns, col1)

	row, err := stub.GetRow("Shipment", columns)
	if err != nil {
		myLogger.Debugf("Failed retriving asset [%s]: [%s]", string(id), err)
		return nil, fmt.Errorf("Failed retriving asset [%s]: [%s]", string(id), err)
	}


	return row.Columns[1], nil
}


func main() {
	err := shim.Start(new(ShipmentChaincode))
	if err != nil {
		fmt.Printf("Error starting Shipment chaincode: %s", err)
	}
}
