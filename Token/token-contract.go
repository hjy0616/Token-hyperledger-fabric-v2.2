/*
 * SPDX-License-Identifier: Apache-2.0
 */

 package main

 import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
 
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
 )
 
 // TokenContract contract for managing CRUD for Token
 type TokenContract struct {
	contractapi.Contract
 }
 
 // 체인코드 생성 시 토큰정보 생성
 func (t *TokenContract) CreateToken(ctx contractapi.TransactionContextInterface, name, symbol, totalSupply string) error {
 
	tokenInfo := TokenInfo{
	   Name:        name,
	   Symbol:      symbol,
	   TotalSupply: totalSupply,
	}
 
	tokenAsBytes, err := json.Marshal(tokenInfo)
	if err != nil {
	   return fmt.Errorf("err")
	}
 
	return ctx.GetStub().PutState(TOKEN_INFO, tokenAsBytes)
 }
 
 // 체인코드 생성시 만들어진 토큰정보 조회
 func (t *TokenContract) Get_Token_Info(ctx contractapi.TransactionContextInterface) (*TokenInfo, error) {
	value, err := ctx.GetStub().GetState(TOKEN_INFO)
 
	if err != nil {
	   return nil, errors.New("failed to get asset: " + TOKEN_INFO)
	}
 
	if value == nil {
	   return nil, errors.New("asset not found: ")
	}
 
	var token TokenInfo
	err = json.Unmarshal(value, &token)
	if err != nil {
	   return nil, err
	}
 
	return &token, nil
 }
 
 func (t *TokenContract) CreateAccount(ctx contractapi.TransactionContextInterface, key string, alloc float64) (string, error) {
	hash := sha1.New()
 
	hash.Write([]byte(key))
 
	md := hash.Sum(nil)
	createdAccount := hex.EncodeToString(md)
 
	accountInfo, _ := ctx.GetStub().GetState(createdAccount)
	ai := Account{}
	json.Unmarshal([]byte(accountInfo), &ai)
 
	if (ai == (Account{})) == false {
	   return "", errors.New("already account address")
	}
 
	var account = Account{Value: alloc, Address: createdAccount}
 
	accountAsBytes, _ := json.Marshal(account)
	ctx.GetStub().PutState(createdAccount, accountAsBytes)
 
	return createdAccount, nil
 }
 
 // account 정보 가져오기
 func (t *TokenContract) Get_Account(ctx contractapi.TransactionContextInterface, key string) (string, error) {
 
	hash := sha1.New()
 
	hash.Write([]byte(key))
 
	md := hash.Sum(nil)
	account_address := hex.EncodeToString(md)
	addressInfo, err := ctx.GetStub().GetState(account_address)
	address := &Account{}
	json.Unmarshal([]byte(addressInfo), &address)
 
	if err != nil {
	   return "", errors.New("Failed to get asset: " + account_address)
	}
 
	if address == nil {
	   return "", errors.New("Asset not found: " + account_address)
	}
 
	return string(addressInfo[:]), nil
 }
 
 // 토큰 전송
 // from: 보내는사람, to: 받는사람, value: 토큰량
 func (t *TokenContract) Transfer(ctx contractapi.TransactionContextInterface, from string, to string, value float64) (string, error) {
	timestamp, _ := ctx.GetStub().GetTxTimestamp()

	// kst
	loc, _ := time.LoadLocation("Asia/Seoul")
	time1 := time.Unix(timestamp.Seconds, int64(timestamp.Nanos)).UTC().In(loc).String()

	hash := sha256.New()
	hash.Write([]byte(from + to + string(rune(value)) + string(time1)))
	md := hash.Sum(nil)
	txId := hex.EncodeToString(md)
 
	// 1.1 각 account 정보 조회
	fromInfo, err := ctx.GetStub().GetState(from)
	if err != nil {
		return "", errors.New("Failed to get asset: " + from)
	}
	if fromInfo == nil {
		return "", errors.New("Asset not found: " + from)
	}
	fromAddress := &Account{}
	json.Unmarshal([]byte(fromInfo), &fromAddress)
 
	toInfo, err := ctx.GetStub().GetState(to)
	if err != nil {
	   return "", errors.New("Failed to get asset: " + to)
	}
	if toInfo == nil {
	   return "", errors.New("Asset not found: " + to)
	}
	toAddress := &Account{}
	json.Unmarshal([]byte(toInfo), &toAddress)
 
	// 1.2 from account value 확인
	if fromAddress.Value < value {
	   return "", errors.New("잔액부족")
	}
 
	if fromAddress.Value-value < 0 {
	   return "", errors.New("잔액부족")
	}
 
	fromAddress.Value -= value // 보내는 사람 차감
	toAddress.Value += value   // 받는사람 증가
 
	// 1.3 account 정보 업데이트
	fromAddressAsBytes, _ := json.Marshal(fromAddress)
	ctx.GetStub().PutState(from, fromAddressAsBytes)
	toAddressAsBytes, _ := json.Marshal(toAddress)
	ctx.GetStub().PutState(to, toAddressAsBytes)
 
	// 2. transaction 기록
	var transaction = Transaction{
	   TxId:      txId,
	   From:      from,
	   To:        to,
	   Value:     value,
	   Timestamp: time1,
	}
	transactionAsBytes, _ := json.Marshal(transaction)
	ctx.GetStub().PutState(txId, []byte(transactionAsBytes))
 
	// 3. receipt 기록
	inReceipt := addReceipt(ctx, txId, toAddress.Address, "IN", time1)
	outReceipt := addReceipt(ctx, txId, fromAddress.Address, "OUT", time1)
 
	// 완료
	fmt.Println("inReceipt : ", inReceipt)
	fmt.Println("outReceipt: ", outReceipt)
	fmt.Println("*************************")
 
	return string(transactionAsBytes[:]), nil
 }
 
 func addReceipt(ctx contractapi.TransactionContextInterface, txId string, toAccount string, status string, timestamp string) string {
	hash := sha256.New()
	hash.Write([]byte(txId + timestamp + status))
	md := hash.Sum(nil)
	receiptId := string(hex.EncodeToString(md))
 
	rr, _ := ctx.GetStub().GetState(ROOT_RECEIPT + toAccount + RECEIPT) // 첫 번째 노드
	rootReceipt := RootReceipt{}
	json.Unmarshal([]byte(rr), &rootReceipt)
 
	lr, _ := ctx.GetStub().GetState(LAST_RECEIPT + toAccount + RECEIPT) // 마지막 노드
	lastReceipt := LastReceipt{}
	json.Unmarshal([]byte(lr), &lastReceipt)
 
	e, _ := json.Marshal(lastReceipt)
	fmt.Println("lstReceipt: ", string(e))
 
	receipt := Receipt{
	   ReceiptId:     receiptId,
	   TxId:          txId,
	   NextReceiptId: "",
	   PrevReceiptId: lastReceipt.ReceiptId,
	   Status:        status,
	}
 
	receiptAsBytes, _ := json.Marshal(receipt)
 
	if (rootReceipt == (RootReceipt{})) == true {
	   rootReceipt.ReceiptId = receiptId
	   rootReceiptAsBytes, _ := json.Marshal(receipt)
 
	   fmt.Println("root receipt id", string(rootReceiptAsBytes))
	   ctx.GetStub().PutState(ROOT_RECEIPT+toAccount+RECEIPT, rootReceiptAsBytes)
	   ctx.GetStub().PutState(NODE_RECEIPT+receiptId+RECEIPT, receiptAsBytes)
	} else {
	   prevReceiptId := lastReceipt.ReceiptId
	   prevReceipt, _ := ctx.GetStub().GetState(NODE_RECEIPT + prevReceiptId + RECEIPT)
	   pr := Receipt{}
	   json.Unmarshal([]byte(prevReceipt), &pr)
 
	   pr.NextReceiptId = receiptId
	   prAsBytes, _ := json.Marshal(pr)
	   ctx.GetStub().PutState(NODE_RECEIPT+prevReceiptId+RECEIPT, prAsBytes)
 
	   ctx.GetStub().PutState(NODE_RECEIPT+receiptId+RECEIPT, receiptAsBytes)
	}
 
	lastReceipt.ReceiptId = receiptId
	lastReceiptAsBytes, _ := json.Marshal(lastReceipt)
	fmt.Println("last receipt id", string(lastReceiptAsBytes))
	ctx.GetStub().PutState(LAST_RECEIPT+toAccount+RECEIPT, lastReceiptAsBytes)
 
	return receiptId
 }
 
 func (t *TokenContract) Get_tx(ctx contractapi.TransactionContextInterface, value string) (string, error) {
	txInfo, _ := ctx.GetStub().GetState(value)
 
	return string(txInfo), nil
 }
 
 // 가장 처음 발생한 내역 가져오기
 // args[0] : address
 func (t *TokenContract) Get_Root_Receipt(ctx contractapi.TransactionContextInterface, value string) (string, error) {
	address := value
 
	rr, _ := ctx.GetStub().GetState(ROOT_RECEIPT + address + RECEIPT)
	rootReceipt := RootReceipt{}
	json.Unmarshal([]byte(rr), &rootReceipt)
 
	receipt, _ := ctx.GetStub().GetState(NODE_RECEIPT + rootReceipt.ReceiptId + RECEIPT)
 
	return string(receipt), nil
 }
 
 // 가장 마지막에 발생한 내역 가져오기
 // args[0] : address
 func (t *TokenContract) Get_Last_Receipt(ctx contractapi.TransactionContextInterface, value string) (string, error) {
	address := value
 
	lr, _ := ctx.GetStub().GetState(LAST_RECEIPT + address + RECEIPT)
	lastReceipt := LastReceipt{}
	json.Unmarshal([]byte(lr), &lastReceipt)
 
	receipt, _ := ctx.GetStub().GetState(NODE_RECEIPT + lastReceipt.ReceiptId + RECEIPT)
 
	return string(receipt), nil
 }
 
 // 모든내역 가져오기
 // args[0] : address
 func (t *TokenContract) Get_Receipts(ctx contractapi.TransactionContextInterface, value string) (string, error) {
 
	var receipts bytes.Buffer
 
	RootReceiptAsBytes, _ := ctx.GetStub().GetState(ROOT_RECEIPT + value + RECEIPT)
	rootReceipt := RootReceipt{}
	json.Unmarshal(RootReceiptAsBytes, &rootReceipt)
 
	receiptAsBytes, _ := ctx.GetStub().GetState(NODE_RECEIPT + rootReceipt.ReceiptId + RECEIPT)
	receipt := Receipt{}
	json.Unmarshal(receiptAsBytes, &receipt)
 
	receipts.WriteString("[")
 
	bArrayMemberAlreadyWritten := false
 
	for receipt.ReceiptId != "" {
 
	   if bArrayMemberAlreadyWritten == true {
		  receipts.WriteString(",")
	   }
 
	   receipts.WriteString("{\"ReceiptId\":")
	   receipts.WriteString("\"")
	   receipts.WriteString(receipt.ReceiptId)
	   receipts.WriteString("\"")
 
	   receipts.WriteString(", \"TxId\":")
	   receipts.WriteString("\"")
	   receipts.WriteString(receipt.TxId)
	   receipts.WriteString("\"")
 
	   receipts.WriteString(", \"NextReceiptId\":")
	   receipts.WriteString("\"")
	   receipts.WriteString(receipt.NextReceiptId)
	   receipts.WriteString("\"")
 
	   receipts.WriteString(", \"PrevReceiptId\":")
	   receipts.WriteString("\"")
	   receipts.WriteString(receipt.PrevReceiptId)
	   receipts.WriteString("\"")
 
	   receipts.WriteString(", \"Status\":")
	   receipts.WriteString("\"")
	   receipts.WriteString(receipt.Status)
	   receipts.WriteString("\"")
	   receipts.WriteString("}")
 
	   if receipt.NextReceiptId == "" {
		  break
	   }
 
	   d, _ := ctx.GetStub().GetState(NODE_RECEIPT + receipt.NextReceiptId + RECEIPT)
	   receipt = Receipt{}
	   json.Unmarshal(d, &receipt)
 
	   bArrayMemberAlreadyWritten = true
	}
	receipts.WriteString("]")
	return receipts.String(), nil
 }
 
 // 특정 receipt 정보 가져오기
 func (t *TokenContract) Get_Receipt(ctx contractapi.TransactionContextInterface, value string) (string, error) {
	receiptId := value
 
	receipt, _ := ctx.GetStub().GetState(NODE_RECEIPT + receiptId + RECEIPT)
 
	if receipt == nil {
	   return "", errors.New("not found receipt")
	}
 
	return string(receipt), nil
 
 }
