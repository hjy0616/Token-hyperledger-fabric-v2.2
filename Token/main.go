/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	tokenContract := new(TokenContract)
	tokenContract.Info.Version = "0.0.1"
	tokenContract.Info.Description = "My Smart Contract"
	tokenContract.Info.License = new(metadata.LicenseMetadata)
	tokenContract.Info.License.Name = "Apache-2.0"
	tokenContract.Info.Contact = new(metadata.ContactMetadata)
	tokenContract.Info.Contact.Name = "John Doe"

	chaincode, err := contractapi.NewChaincode(tokenContract)
	chaincode.Info.Title = "Token chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from TokenContract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
