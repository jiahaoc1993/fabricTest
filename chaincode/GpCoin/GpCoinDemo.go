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
		&shim.ColumnDefinition{Name: "GPCoin", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "USD", Type: shim.ColumnDefinition_STRING, Key: false},
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

	ok, err := stub.InsertRow("GPCoinOwnership", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "1000000.00000"}},
			&shim.Column{Value: &shim.Column_String_{String_: "1000000.00000"}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: adminCert}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Charge was already done.")
	}

	myLogger.Debug("Init Chaincode...done")

	return nil, nil
}


func (t *GPCoinChaincode) topup(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	myLogger.Debug("Topup...")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	amount, _ := strconv.ParseFloat(args[0], 64)

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

	if len(row.Columns) != 0 {

		coinString   := row.Columns[0].GetString_()
		USDString    := row.Columns[1].GetString_()

		USD, _	:= strconv.ParseFloat(USDString, 64)
		USD += amount

		USDString   = strconv.FormatFloat(USD, 'f', 5, 64)

		err = stub.DeleteRow(
		"GpcoinOwnership",
		[]shim.Column{shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}},
		)
		if err != nil {
			return nil, errors.New("Failed deliting row.")
		}
	}

	ok, err = stub.InsertRow("GPCoinOwnership", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: coinString}},
			&shim.Column{Value: &shim.Column_String_{String_: USDString}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Charge was already done.")
	}

	myLogger.Debug("Charge...done!")

	return nil, err
}



func (t *GPCoinChaincode) invest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	myLogger.Debug("invest...")

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

	amount, _ := strconv.ParseFloat(args[0], 64)// always returns Float64
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


	//update the account
		USDString := row.Columns[1].GetString_()
		coinString := row.Columns[0].GetString_()

		USD, _ := strconv.ParseFloat(USDString, 64)
		coin, _ := strconv.ParseFloat(coinString, 64)

		if amount > USD {
			return nil, errors.New("You don't have enough money!")
		}
		sentCoin := amount/ 123
		coin  += sentCoin
		fee := amount * 0.05
		USD -= amount + fee

		USDString   = strconv.FormatFloat(USD, 'f', 5, 64)

		coinString   = strconv.FormatFloat(coin, 'f', 5, 64)

		err = stub.DeleteRow(
		"GpcoinOwnership",
		[]shim.Column{shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}},
		)
		if err != nil {
			return nil, errors.New("Failed deliting row.")
		}
	}

	ok, err = stub.InsertRow("GPCoinOwnership", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: coinString}},
			&shim.Column{Value: &shim.Column_String_{String_: USDString}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Charge was already done.")
	}

	//don't forget the Gp account !  adminCertificate from here 

	//adminCertificate, _ := stub.GetState("admin")

	var GPcolumns []shim.Column
	GPcol := shim.Column{Value: &shim.Column_Bytes{Bytes: adminCertificate}}
	GPcolumns = append(GPcolumns, GPcol)

	GProw, err := stub.GetRow("GPCoinOwnerShip", GPcolumns)

	if err != nil {
		myLogger.Debugf("Failed retriving Owner")
		return nil, fmt.Errorf("Failed retriving Owner [%d]: [%s]", amount, err)
	}


	GPUSDString := GProw.Columns[1].GetString_()
	GPCoinString := GProw.Columns[0].GetString_()

	GPUSD, _ := strconv.ParseFloat(GPUSDString, 64)
	GPUSD += fee


	GPCoin, _ := strconv.ParseFloat(GPCoinString, 64)
	GPCoin -= sentCoin

	GPUSDString   = strconv.FormatFloat(GPUSD, 'f', 5, 64)
	GPCoinString  = strconv.FormatFloat(GPCoin, 'f', 5, 64)

	err = stub.DeleteRow(
	"GpcoinOwnership",
	[]shim.Column{shim.Column{Value: &shim.Column_Bytes{Bytes: adminCertificate}}},
	)

	if err != nil {
		return nil, errors.New("Failed deliting row.")
	}

	ok, err = stub.InsertRow("GPCoinOwnership", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: GPCoinString}},
			&shim.Column{Value: &shim.Column_String_{String_: GPUSDString}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: adminCertificate}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Charge was already done.")
	}


	myLogger.Debug("Invest...done!")

	return nil, err
}


func (t *GPCoinChaincode) cashout(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	myLogger.Debug("invest...")

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

	amount, _ := strconv.ParseFloat(args[0], 64)// always returns Float64
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

	//update the account
		USDString := row.Columns[1].GetString_()
		coinString := row.Columns[0].GetString_()

		USD, _ := strconv.ParseFloat(USDString, 64)
		coin, _ := strconv.ParseFloat(coinString, 64)

		if amount > coin {
			return nil, errors.New("You don't have enough coin!")
		}

		coin  -= amount
		cash := amount * 123
		fee := cash * 0.05

		USD += cash - fee

		USDString   = strconv.FormatFloat(USD, 'f', 5, 64)

		coinString   = strconv.FormatFloat(coin, 'f', 5, 64)

		err = stub.DeleteRow(
		"GpcoinOwnership",
		[]shim.Column{shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}},
		)
		if err != nil {
			return nil, errors.New("Failed deliting row.")
		}
	}
	ok, err = stub.InsertRow("GPCoinOwnership", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: coinString}},
			&shim.Column{Value: &shim.Column_String_{String_: USDString}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Charge was already done.")
	}

	//don't forget the Gp account !  adminCertificate from here 

	//adminCertificate, _ := stub.GetState("admin")

	var GPcolumns []shim.Column
	GPcol := shim.Column{Value: &shim.Column_Bytes{Bytes: adminCertificate}}
	GPcolumns = append(GPcolumns, GPcol)

	GProw, err := stub.GetRow("GPCoinOwnerShip", GPcolumns)

	if err != nil {
		myLogger.Debugf("Failed retriving Owner")
		return nil, fmt.Errorf("Failed retriving Owner [%d]: [%s]", amount, err)
	}


	GPUSDString := GProw.Columns[1].GetString_()
	GPCoinString := GProw.Columns[0].GetString_()

	GPUSD, _ := strconv.ParseFloat(GPUSDString, 64)
	GPUSD += fee

	GPCoin, _ := strconv.ParseFloat(GPCoinString, 64)
	GPCoin += amount

	GPUSDString   = strconv.FormatFloat(GPUSD, 'f', 5, 64)
	GPCoinString   = strconv.FormatFloat(GPCoin, 'f', 5, 64)
	err = stub.DeleteRow(
	"GpcoinOwnership",
	[]shim.Column{shim.Column{Value: &shim.Column_Bytes{Bytes: adminCertificate}}},
	)

	if err != nil {
		return nil, errors.New("Failed deliting row.")
	}

	ok, err = stub.InsertRow("GPCoinOwnership", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: GPCoinString}},
			&shim.Column{Value: &shim.Column_String_{String_: GPUSDString}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: adminCertificate}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Charge was already done.")
	}


	myLogger.Debug("Invest...done!")

	return nil, err
}

func (t *GPCoinChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	myLogger.Debug("Charging...")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	amount, _ := strconv.ParseFloat(args[0], 64)

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

	if len(fromResult.Columns) == 0 || len(toResult.Columns) ==0 {
		return nil, errors.New("We can't find one of users!")
	}
		fromCoinString := fromResult.Columns[0].GetString_()
		toCoinString   :=  toResult.Columns[0].GetString_()

		fromCoin, _ := strconv.ParseFloat(fromCoinString, 64)
		toCoin, _   := strconv.ParseFloat(toCoinString, 64)

		if amount > fromCoin{
			return nil, errors.New("You don't have enough coin!")
		}

		fromCoin -= amount
		toCoin   += amount

		fromUSDString  := fromResult.Columns[1].GetString_()
		toUSDString    := toResult.Columns[1].GetString_()



		fromCoinString = strconv.FormatFloat(fromCoin, 'f', 5, 64)
		toCoinString   = strconv.FormatFloat(toCoin, 'f', 5, 64)

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
			&shim.Column{Value: &shim.Column_String_{String_: fromCoinString}},
			&shim.Column{Value: &shim.Column_String_{String_: fromUSDString}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: from}}},
		})

		if !ok && err == nil {
			return nil, errors.New("Charge was already done.")
		}

		ok, err = stub.InsertRow("GPCoinOwnership", shim.Row{
			Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: toCoinString}},
			&shim.Column{Value: &shim.Column_String_{String_: toUSDString}},
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
	if function == "invest" {
		// Assign ownership
		return t.invest(stub, args)
	} else if function == "transfer" {
		// Transfer ownership
		return t.transfer(stub, args)
	}else if function == "topup"{
		return t.topup(stub, args)
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

	var result string
	if len(row.Columns) != 0 {
	fmt.Sprintf(result,"[%x] has %s GpCoins, %s USD", row.Columns[2].GetBytes(), row.Columns[0].GetString_(), row.Columns[1].GetString_())
	}else{
		result = "No ansawer!"
	}

	return []byte(result), nil
}

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(GPCoinChaincode))
	if err != nil {
		fmt.Printf("Error starting AssetManagementChaincode: %s", err)
	}
}

