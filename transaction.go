package chain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Input struct {
	TransactionHash string `json:"transaction_hash"`
	OutputHash      string `json:"output_hash"`
	OutputIndex     int64  `json:"output_index"`
	Value           int64
	Addresses       []string
	ScriptSignature string `json:"script_signature"`
}

type Output struct {
	TransactionHash    string `json:"transaction_hash"`
	OutputIndex        uint32 `json:"output_index"`
	Value              int64
	Addresses          []string
	Script             string
	ScriptHex          string `json:"script_hex"`
	RequiredSignatures int64  `json:"required_signatures"`
	Spent              bool

	// Only populated with GetAddressUnspentOutputs.
	Confirmations int64
}

type Transaction struct {
	Hash          string
	BlockHash     string `json:"block_hash"`
	BlockHeight   int64  `json:"block_height"`
	BlockTime     string `json:"block_time"`
	Inputs        []Input
	Outputs       []Output
	Amount        int64
	Fees          int64
	Confirmations int64
}

func (c *Chain) transactionURL(hash string) string {
	return fmt.Sprintf("%s/transactions/%s", c.network, hash)
}

func (c *Chain) GetTransaction(hash string) (Transaction, error) {
	url, tx := c.transactionURL(hash), &Transaction{}
	return *tx, c.httpGetJSON(url, tx)
}

func (c *Chain) sendTransactionURL() string {
	return fmt.Sprintf("%s/transactions", c.network)
}

func (c *Chain) SendTransaction(hex string) (string, error) {
	url := c.sendTransactionURL()

	hexMessage := struct {
		Hex string `json:"hex"`
	}{hex}

	data, err := json.Marshal(hexMessage)
	if err != nil {
		return "", err
	}
	response, err := c.httpPut(url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	hash, err := ioutil.ReadAll(response)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}