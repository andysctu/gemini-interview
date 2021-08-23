# Golang Jobcoin
This service provides an endpoint `/getDepositAddressForMixing` that users can call to mix their jobcoins.
Users should provide a list of jobcoin addresses that they own.
The endpoint will return a deposit address.
Users can make a one-time deposit into the deposit address of any non-0 amount.
The mixer will distribute that amount into the list of provided addresses over a period of time.

### How to run
    `go run main.go`

### Mix jobcoins
    `curl localhost:9001/getDepositAddressForMixing --data '{"addresses":[<ADDRESSES_YOU_OWN>]}'

    example:
    `curl localhost:9001/getDepositAddressForMixing --data '{"addresses":["address1","address2","address3"]}'

### Future improvements
