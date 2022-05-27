/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

// Token stores a value
type TokenInfo struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	TotalSupply string `json:"totalSupply"`
}

const TOKEN_INFO = "token"

type Account struct {
	Value   float64    `json:"value"`
	Address string     `json:"address"`
}

type Transaction struct {
	TxId      string
	From      string
	To        string
	Value     float64
	Timestamp string
}

type Receipt struct {
	ReceiptId     string `json:"receiptId"`
	TxId          string `json:"txId"`
	NextReceiptId string `json:"nextReceiptId"`
	PrevReceiptId string `json:"prevReceiptId"`
	Status        string `json:"status"`
}

type RootReceipt struct {
	ReceiptId string `json:"receiptId"`
}

type LastReceipt struct {
	ReceiptId string `json:"receiptId"`
}

const ROOT_RECEIPT = "ROOT_"
const NODE_RECEIPT = "NODE_"
const LAST_RECEIPT = "LAST_"
const RECEIPT = "RECEIPT"

