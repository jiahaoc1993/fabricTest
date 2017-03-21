package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/op/go-logging"
)

var myLogger = logging.MustGetLogger("GPCoin")

type GPCoinChaincode struct {
}

func (t *GPCoinChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	myLogger.Debug("Init Chaincode...")
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	// Create ownership table
	err := stub.CreateTable("GPCoinOwnership", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "GPCoin", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "USD", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "Owner", Type: shim.ColumnDefinition_BYTES, Key: true},
	})
	if err != nil {
		return nil, errors.New("Failed creating GpcoinOnwership table.")
	}

	// Set the admin
	// The metadata will contain the certificate of the administrator
	adminCert, err := stub.GetCallerMetadata()
	if err != nil {
		myLogger.Debug("Failed getting metadata")
		return nil, errors.New("Failed getting metadata.")
	}
	if len(adminCert) == 0 {
		myLogger.Debug("Invalid admin certificate. Empty.")
		return nil, errors.New("Invalid admin certificate. Empty.")
	}

	myLogger.Debug("The administrator is [%x]", adminCert)

	stub.PutState("admin", adminCert)

	myLogger.Debug("Init Chaincode...done")

	return nil, nil
}


func (t *GPCoinChaincode) charge(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	myLogger.Debug("Charging...")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	tmp, _ := strconv.ParseInt(args[0], 10, 32)
	amount := int32(tmp)
	owner, err := base64.StdEncoding.DecodeString(args[1])
	if err != nil {
		return nil, errors.New("Failed decodinf owner")
	}

	// Verify the identity of the caller
	// Only an administrator can invoker assign
	adminCertificate, err := stub.GetState("admin")
	if err != nil {
		return nil, errors.New("Failed fetching admin identity")
	}

	ok, err := t.isCaller(stub, adminCertificate)
	if err != nil {
		return nil, errors.New("Failed checking admin identity")
	}
	if !ok {
		return nil, errors.New("The caller is not an administrator")
	}

	// Register assignment
	myLogger.Debugf("[% x] is charged %d", owner, amount)

	//query first!
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}
	columns = append(columns, col1)

	row, err := stub.GetRow("GPCoinOwnerShip", columns)

	if err != nil {
		myLogger.Debugf("Failed retriving Owner")
		return nil, fmt.Errorf("Failed retriving Owner [%d]: [%s]", amount, err)
	}

	if len(row.Columns) == 0 {
		return nil, errors.New("Can't find user")
	}
		coin   := row.Columns[0].GetInt32()
		amount += row.Columns[1].GetInt32()
		err = stub.DeleteRow(
		"GpcoinOwnership",
		[]shim.Column{shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}},
		)
		if err != nil {
			return nil, errors.New("Failed deliting row.")
		}


	ok, err = stub.InsertRow("GPCoinOwnership", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Int32{Int32: coin}},
			&shim.Column{Value: &shim.Column_Int32{Int32: amount}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Charge was already done.")
	}

	myLogger.Debug("Charge...done!")

	return nil, err
}



func (t *GPCoinChaincode) buy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	myLogger.Debug("Charging...")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	adminCertificate, err := stub.GetState("admin")
	if err != nil {
		return nil, errors.New("Failed fetching admin identity")
	}

	ok, err := t.isCaller(stub, adminCertificate)
	if err != nil {
		return nil, errors.New("Failed checking admin identity")
	}
	if !ok {
		return nil, errors.New("The caller is not an administrator")
	}

	tmp, _ := strconv.ParseInt(args[0], 10, 32)// always returns int64
	amount := int32(tmp)
	owner, err := base64.StdEncoding.DecodeString(args[1])
	if err != nil {
		return nil, errors.New("Failed decodinf owner")
	}

	// Register assignment
	myLogger.Debugf("[% x] is buying %d", owner, amount)

	//query first!
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}
	columns = append(columns, col1)

	row, err := stub.GetRow("GPCoinOwnerShip", columns)

	if err != nil {
		myLogger.Debugf("Failed retriving Owner")
		return nil, fmt.Errorf("Failed retriving Owner [%d]: [%s]", amount, err)
	}

	if len(row.Columns) != 0 {
		if amount > row.Columns[1].GetInt32(){
			return nil, errors.New("You don't have enough money!")
		}
	}else {
		return nil, errors.New("We don't have this users.")
	}

		coin  := row.Columns[0].GetInt32()
		coin  += amount / 10

		amount = row.Columns[1].GetInt32() - amount

		err = stub.DeleteRow(
		"GpcoinOwnership",
		[]shim.Column{shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}},
		)
		if err != nil {
			return nil, errors.New("Failed deliting row.")
		}

	ok, err = stub.InsertRow("GPCoinOwnership", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Int32{Int32: coin}},
			&shim.Column{Value: &shim.Column_Int32{Int32: amount}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Charge was already done.")
	}

	myLogger.Debug("Charge...done!")

	return nil, err
}


func (t *GPCoinChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	myLogger.Debug("Charging...")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	tmp, _ := strconv.ParseInt(args[0], 10, 32)
	amount := int32(tmp)

	from, err := base64.StdEncoding.DecodeString(args[1])
	if err != nil {
		return nil, errors.New("Failed decodinf owner")
	}

	to, err := base64.StdEncoding.DecodeString(args[2])
	if err != nil {
		return nil, errors.New("Failed decodinf owner")
	}

	// Verify the identity of the caller
	// Only an administrator can invoker assign
	adminCertificate, err := stub.GetState("admin")
	if err != nil {
		return nil, errors.New("Failed fetching admin identity")
	}

	ok, err := t.isCaller(stub, adminCertificate)
	if err != nil {
		return nil, errors.New("Failed checking admin identity")
	}
	if !ok {
		return nil, errors.New("The caller is not an administrator")
	}

	// Register assignment
	myLogger.Debugf("[% x] is transform to [%x] : %d", from, to, amount)

	//query first!
	var columns1 []shim.Column
	var columns2 []shim.Column
	col1 := shim.Column{Value: &shim.Column_Bytes{Bytes: from}}
	col2 := shim.Column{Value: &shim.Column_Bytes{Bytes: to}}

	columns1 = append(columns1, col1)
	columns2 = append(columns2, col2)

	row1, err := stub.GetRow("GPCoinOwnerShip", columns1)
	if err != nil {
		myLogger.Debugf("Failed retriving Owner")
		return nil, fmt.Errorf("Failed retriving Owner [%d]: [%s]", amount, err)
	}

	row2, err := stub.GetRow("GPCoinOwnerShip", columns2)
	if err != nil {
		myLogger.Debugf("Failed retriving Owner")
		return nil, fmt.Errorf("Failed retriving Owner [%d]: [%s]", amount, err)
	}

	fromResult := row1
	toResult   := row2

	if len(fromResult.Columns) != 0 && len(toResult.Columns) !=0 {
		if amount > fromResult.Columns[0].GetInt32(){
			return nil, errors.New("You don't have enough coin!")
		}
	}else{
		return nil, errors.New("We can't find one of users!")
	}

		fromCoin := fromResult.Columns[0].GetInt32() - amount
		toCoin   :=  toResult.Columns[0].GetInt32() + amount

		fromUSD  := fromResult.Columns[1].GetInt32()
		toUSD    := toResult.Columns[1].GetInt32()

		err = stub.DeleteRow(
		"GpcoinOwnership",
		[]shim.Column{shim.Column{Value: &shim.Column_Bytes{Bytes: from}}},
		)
		if err != nil {
			return nil, errors.New("Failed deliting row.")
		}

		err = stub.DeleteRow(
		"GpcoinOwnership",
		[]shim.Column{shim.Column{Value: &shim.Column_Bytes{Bytes: to}}},
		)
		if err != nil {
			return nil, errors.New("Failed deliting row.")
		}


		ok, err = stub.InsertRow("GPCoinOwnership", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Int32{Int32: fromCoin}},
			&shim.Column{Value: &shim.Column_Int32{Int32: fromUSD}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: from}}},
		})

		if !ok && err == nil {
			return nil, errors.New("Charge was already done.")
		}

		ok, err = stub.InsertRow("GPCoinOwnership", shim.Row{
			Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Int32{Int32: toCoin}},
			&shim.Column{Value: &shim.Column_Int32{Int32: toUSD}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: to}}},
		})

		if !ok && err == nil {
			return nil, errors.New("Charge was already done.")
		}

	myLogger.Debug("Transfer...done!")

	return nil, err
}




func (t *GPCoinChaincode) isCaller(stub shim.ChaincodeStubInterface, certificate []byte) (bool, error) {
	myLogger.Debug("Check caller...")

	// In order to enforce access control, we require that the
	// metadata contains the signature under the signing key corresponding
	// to the verification key inside certificate of
	// the payload of the transaction (namely, function name and args) and
	// the transaction binding (to avoid copying attacks)

	// Verify \sigma=Sign(certificate.sk, tx.Payload||tx.Binding) against certificate.vk
	// \sigma is in the metadata

	sigma, err := stub.GetCallerMetadata()
	if err != nil {
		return false, errors.New("Failed getting metadata")
	}
	payload, err := stub.GetPayload()
	if err != nil {
		return false, errors.New("Failed getting payload")
	}
	binding, err := stub.GetBinding()
	if err != nil {
		return false, errors.New("Failed getting binding")
	}

	myLogger.Debugf("passed certificate [% x]", certificate)
	myLogger.Debugf("passed sigma [% x]", sigma)
	myLogger.Debugf("passed payload [% x]", payload)
	myLogger.Debugf("passed binding [% x]", binding)

	ok, err := stub.VerifySignature(
		certificate,
		sigma,
		append(payload, binding...),
	)
	if err != nil {
		myLogger.Errorf("Failed checking signature [%s]", err)
		return ok, err
	}
	if !ok {
		myLogger.Error("Invalid signature")
	}

	myLogger.Debug("Check caller...Verified!")

	return ok, err
}


func (t *GPCoinChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "buy" {
		// Assign ownership
		return t.buy(stub, args)
	} else if function == "transfer" {
		// Transfer ownership
		return t.transfer(stub, args)
	}else if function == "charge"{
		return t.charge(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

func (t *GPCoinChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	myLogger.Debugf("Query [%s]", function)

	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting 'query' but found '" + function + "'")
	}

	var err error

	if len(args) != 1 {
		myLogger.Debug("Incorrect number of arguments. Expecting name of an user to query")
		return nil, errors.New("Incorrect number of arguments. Expecting name of an asset to query")
	}

	// Who is the owner of the asset?
	owner, err := base64.StdEncoding.DecodeString(args[0])

	//myLogger.Debugf("Arg [%s]", string(asset))

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}
	columns = append(columns, col1)

	row, err := stub.GetRow("GPCoinOwnership", columns)
	if err != nil {
		myLogger.Debugf("Failed retriving asset [%s]: [%s]", string(owner), err)
		return nil, fmt.Errorf("Failed retriving asset [%s]: [%s]", string(owner), err)
	}

//	myLogger.Debugf("Query done [% x]", row.Columns[1].GetBytes())


	return row.Columns[3].GetBytes(), nil
}

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(GPCoinChaincode))
	if err != nil {
		fmt.Printf("Error starting AssetManagementChaincode: %s", err)
	}
}

