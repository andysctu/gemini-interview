## Overview
* This service provides an endpoint `/getDepositAddressForMixing` that users can call to mix their jobcoins.
* Users should provide a list of jobcoin addresses that they own.
* The endpoint will return a deposit address.
* Users can make a one-time deposit into the deposit address of any non-0 amount.
* The mixer will distribute that amount into the list of provided addresses over a period of time.

### Files

* `main.go` - start here, defines service endpoints
* `mixer/mixer.go` - main logic for the `/getDepositAddressForMixing` endpoint
* `jobcoin_gateway/jobcoin_gateway.go` - contains functions to interact with the Jobcoin API to read/write transactions

## Commands

### How to run
    go run main.go

### How to mix Jobcoins
    curl localhost:9001/getDepositAddressForMixing --data '{"addresses":[<ADDRESSES_YOU_OWN>]}'

Example:
    
    curl localhost:9001/getDepositAddressForMixing --data '{"addresses":["address1","address2","address3"]}'
    
    // Response
    {"depositAddress":"85020336-03ad-11ec-8dbb-acde48001122"}
    
Your first deposit to `85020336-03ad-11ec-8dbb-acde48001122` will get mixed equally into `address1`,`address2`, and `address3` over some time.

## Future improvements
* Support custom values for mixed amounts into each address
* Support multiple deposits into the initial deposit address 
* Error handling
* Unit tests
