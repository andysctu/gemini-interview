// jobcoingateway defines functions to interact with the jobcoin API
// https://jobcoin.gemini.com/unify-yelling/api
package jobcoingateway

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

const (
    BaseURL             = "https://jobcoin.gemini.com/unify-yelling/api"
    AddressesEndpoint   = BaseURL + "/addresses"
    TransactionEndpoint = BaseURL + "/transactions"
)

// AddressInfoResponse defines the response object of the GET addresses endpoint
type AddressInfoResponse struct {
    Balance string `json:"balance"`
    Transactions []Transaction `json:"transactions"`
}

// Transaction defines the transaction object that the jobcoin API uses
type Transaction struct {
    Timestamp string `json:"timestamp,omitempty"`
    FromAddress string `json:"fromAddress,omitempty"`
    ToAddress string `json:"toAddress"`
    Amount string `json:"amount"`
}

// GetTransactionsForAddress calls the jobcoin API to fetch transactions for a given address
func GetTransactionsForAddress(address string) ([]Transaction, error) {
    // Make HTTP call
    resp, err := http.Get(fmt.Sprintf("%v/%v", AddressesEndpoint, address))
    defer resp.Body.Close()
    if err != nil {
        return nil, err
    }

    // Parse response
    body, err := ioutil.ReadAll(resp.Body)
    var res AddressInfoResponse
    err = json.Unmarshal(body, &res)
    if err != nil {
        return nil, err
    }

    return res.Transactions, nil
}

// PostTransaction calls the jobcoin API to create a new transaction
func PostTransaction(from, to, amount string) error {
    // Create transaction object
    transaction := Transaction{
        FromAddress: from,
        ToAddress: to,
        Amount: amount,
    }
    b, err := json.Marshal(transaction)
    if err != nil {
        return err
    }

    // Make HTTP call
    _, err = http.Post(TransactionEndpoint, "application/json", bytes.NewBuffer(b))
    
    // TODO: handle when resp.status code != 200
    return err
}
