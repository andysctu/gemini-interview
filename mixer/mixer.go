package mixer

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "time"

    "github.com/andysctu/jobcoin/jobcoin_gateway"
    "github.com/google/uuid"
)


const (
    // Jobcoin address of the House
    HouseAddress = "TestHouseAddress"

    // Delay between polling the deposit address for transactions
    DepositTransactionSleepTimeMillis = 100

    // Delay between paying out to the provided mixed addresses
    PayoutSleepTimeMillis = 100
)

// Request for the GetDepositAddressForMixing endpoint
type GetDepositAddressForMixingRequest struct {
    Addresses []string `json:"addresses"`
}

// Response for the GetDepositAddressForMixing endpoint
type GetDepositAddressForMixingResponse struct {
    DepositAddress string `json:"depositAddress,omitempty"`
    Error string `json:"error,omitempty"`
}

// Mapping of (deposit addresses -> [ mixed addresses ... ])
// TODO: Move to database
var AddressMap map[string][]string = map[string][]string{}

// GetDepositAddressForMixing will
// 1. Generate a random deposit address
// 2. Start a goroutine to listen for transactions to the deposit address
// 3. Return the deposit address to the caller
// We will distribute the amount from only the first deposit into the address among the given mixing addresses
// Users should call this endpoint again if they want to mix more deposits
func GetDepositAddressForMixing(w http.ResponseWriter, req *http.Request) {
    // Set content-type of response
    w.Header().Set("Content-Type", "application/json")

    // Parse request
    var r GetDepositAddressForMixingRequest
    err := json.NewDecoder(req.Body).Decode(&r)
    if err != nil {
        json.NewEncoder(w).Encode(GetDepositAddressForMixingResponse{
            Error: err.Error(),
        })
        return
    }

    if len(r.Addresses) < 2 {
        json.NewEncoder(w).Encode(GetDepositAddressForMixingResponse{
            Error: "Please provide at least 2 addresses for mixing.",
        })
        return
    }

    // Get deposit address
    depositAddress := getDepositAddress()
    
    // Add mapping of (depositAddress -> [ mixed addresses ... ])
    AddressMap[depositAddress] = r.Addresses

    // Start listening for transactions to deposit address
    go watchForDepositTransaction(depositAddress)

    // Return deposit address
    json.NewEncoder(w).Encode(GetDepositAddressForMixingResponse{
        DepositAddress: depositAddress,
    })
}

// Returns a random UUID to be used as a deposit address
func getDepositAddress() string {
    depositAddress, _ := uuid.NewUUID()
    return depositAddress.String()
}

// Poll the jobcoin API to listen for transactions to the given address
// We only care about the first transaction
func watchForDepositTransaction(address string) {
    fmt.Printf("Waiting for transactions to deposit address: %v\n", address)
    for {
        // Wait some time before checking the jobcoin API again
        fmt.Printf(".")
        time.Sleep(DepositTransactionSleepTimeMillis * time.Millisecond)

        // Fetch all transactions for the deposit address
        transactions, err := jobcoingateway.GetTransactionsForAddress(address)
        if err != nil {
            // TODO: handle errors
            continue
        }

        // If there are no transactions yet, continue polling
        if len(transactions) == 0 {
            continue
        }

        // Otherwise, just look at the first transaction
        fmt.Println()
        firstTransaction := transactions[0]
        fmt.Printf("Found transaction: %+v\n", firstTransaction)

        // Extract the amount from the transaction into the House Account
        fmt.Printf("Sending %v from deposit address to House Address\n", firstTransaction.Amount)
        err = jobcoingateway.PostTransaction(address, HouseAddress, firstTransaction.Amount)
        if err != nil {
            // TODO: handle errors
            break
        }

        // Payout to the correct provided addresses for mixing
        go payoutFromHouse(AddressMap[address], firstTransaction.Amount)
        
        break
    }   
}

// payoutFromHouse will post transactions from the House address to the provided addresses
// in equal increments based on (amount / # addresses)
func payoutFromHouse(addresses []string, amount string) {
    // Parse amount
    amt, err := strconv.ParseFloat(amount, 64)
    if err != nil {
        // TODO: handle errors
        fmt.Printf("invalid amount string: %v\n", amount)
        return
    }

    increment := amt / float64(len(addresses))

    // TODO: Handle rounding errors
    fmt.Printf("Paying out (%v) over (%v) payments of (%v) each from House Address to each of the mixing addresses.\n", amt, len(addresses), increment)
    for i, address := range addresses {
        fmt.Println("Waiting some time before next payment...")
        fmt.Printf("Payment (%v of %v) to address (%v) of amount (%v)\n", i+1, len(addresses), address, increment)
        err = jobcoingateway.PostTransaction(HouseAddress, address, fmt.Sprintf("%f",increment))
        if err != nil {
            // TODO: handle errors
            fmt.Printf("error paying out to %v\n", address)
        }

        // Wait some time before paying out each address
        time.Sleep(PayoutSleepTimeMillis * time.Millisecond)
    }
}