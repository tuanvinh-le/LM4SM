
// Chaincode

type SimpleChaincode struct {
}

func pow(b string) int {
	var hashInt big.Int
	var hash [32]byte

	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	nonce := 0

	// t1 := time.Now() // get current time
	for nonce < maxNonce {
		data := prepareData(nonce, b)

		hash = sha256.Sum256(data)
		// fmt.Printf("#%d = %x\r", nonce, hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			// elapsed := time.Since(t1)
			// fmt.Print("\nApp elapsed: ", elapsed)
			break
		} else {
			nonce++
		}
	}

	return nonce
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("cgublock Init")
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("cgublock Invoke")

	function, args := stub.GetFunctionAndParameters()
	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "add" {
		// the old "Query" is now implemtned in invoke
		return t.add(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"delete\" \"query\" \"add\"")
}

func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("cgublock delete")
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ID := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(ID)

	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) add(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("cgublock add")

	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Write the state back to the ledger
	err = stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("cgublock query")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ID := args[0]

	// Delete the key from the state in ledger
	fileblock, err := stub.GetState(ID)
	if err != nil {
		return shim.Error("Could not fetch application with id")
	}

	if fileblock == nil {
		jsonResp := "{\"Error\":\"Nil data for " + ID + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Printf("Query Response:\n")

	return shim.Success(fileblock)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}


// Public blocks (batch signatures) generation

sigall := pairing.NewG1().Set1()

 for i:=0 ; i<=3 ; i++{ 
 signature[i] := pairing.NewG1().SetBytes(ct.signature[i])

  h[i] := pairing.NewG1().SetFromStringHash(message[i], sha256.New())

sigall .Mul(sigall, signature[i])
}


// Smart contract

client, err := ethclient.Dial("https://ropsten.infura.io/v3/a565d0dc884c476fa2e25636ea19fa82")
  if err != nil {
    log.Fatal(err)
  }

  privateKey, err := crypto.HexToECDSA("5F03F06E2B524F4D8FF6135967899992B6B609F8A37B4D0015A1C0154E1A4FDB")
  if err != nil {
    log.Fatal(err)
  }

  publicKey := privateKey.Public()
  publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
  if !ok {
    log.Fatal("error casting public key to ECDSA")
  }

  fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
  nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
  if err != nil {
    log.Fatal(err)
  }
 value := new(big.Int)
  value.SetString("30000000000000000", 10) // in wei (0.3 eth)
  gasLimit := uint64(40000)                 // in units

  gasPrice, err := client.SuggestGasPrice(context.Background())
  if err != nil {
    log.Fatal(err)
  }

  toAddress := common.HexToAddress("0xBa7adA49BffDc8c641D1cB8f3f9aF29F9BD9C66e")
  data := []byte(PB)
  tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
  signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
  if err != nil {
    log.Fatal(err)
  }
  err = client.SendTransaction(context.Background(), signedTx)
  if err != nil {
    log.Fatal(err)
  }
