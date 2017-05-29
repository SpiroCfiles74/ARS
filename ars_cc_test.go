//TDD file

package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	fmt.Println("Entering Init function")

}

func TestInvoke(t *testing.T) {
	fmt.Println("Entering Invoke function")

	_, err := MockStub.MockInvoke("t123", "InvalidFunctionName", []string{})
	if err == nil {
		t.Fatalf("Expected invalid function name error")
	}

	bytes, err := MockStub.MockInvoke("t123", "CConsumerProduct", []string{loanApplicationID, loanApplication})
	if err != nil {
		t.Fatalf("Expected CreateLoanApplication function to be invoked")
	}
	//A spy could have been used here to ensure CreateLoanApplication method actually got invoked.
	var la LoanApplication
	err = json.Unmarshal(bytes, &la)
	if err != nil {
		fmt.Fatalf("Expected valid loan application JSON string to be returned from CreateLoanApplication method")
	}

}

func TestQuery(t *testing.T) {
	fmt.Println("Entering Query function")
	//function, value nill
	//then function output, match with value
}

func TestConsumer(t *testing.T) {
	fmt.Println("Entering TestConsumer function")
	//attributes := make(map[string][]byte) //Create a custom MockStub that internally uses shim.MockStub
	//stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	//if stub == nil {
	//	t.Fatalf("MockStub creation failed")
	//}
}

func TestConsumerProduct(t *testing.T) {
	fmt.Println("Entering TestConsumerProdut function")
}

func TestConsumerProductAlert(t *testing.T) {
	fmt.Println("Entering TestConsumerProductAlert function")
}

func TestDeleteConsum(t *testing.T) {
	fmt.Println("Entering TestDeleteConsum function")
}

func TestDeleteConsumProd(t *testing.T) {
	fmt.Println("Entering TestDeleteConsumProd function")
}
