# Blossom Smart Contracts
This package contains the code for the Blossom Smart Contracts.

## Deployment Steps
In the below commands to deploy the chaincode, `blossom-1` is the name of the channel and `blossomcc` is the name of the chaincode.

0. Make sure the Blossom project is cloned on the peer machine.  The path provided in the following `install` command
   assumes the chaincode is located in `$GOPATH`.
   
1. Install chaincode on the peer
   
   ```
   docker exec cli peer chaincode install -n blossomcc -v {VERSION} -p github.com/usnistgov/blossom/chaincode
   ```

2. Instantiate chaincode the chaincode on the channel `blossom-1`
   
   ```
   docker exec cli peer chaincode instantiate -o $ORDERER -C blossom-1 -n blossomcc -v {VERSION} -c '{"Args":["init"]}' --cafile /opt/home/managedblockchain-tls-chain.pem --tls
   ```

3. Check chaincode instantiation

   ```
   docker exec cli peer chaincode list --instantiated -o $ORDERER -C blossom-1 --cafile /opt/home/managedblockchain-tls-chain.pem --tls
   ```

4. Invoke chaincode

   ```
   docker exec cli peer chaincode invoke -C blossom-1 -n blossomcc -c  '{"Args":["test", "awesome blossom"]}' -o $ORDERER --cafile /opt/home/managedblockchain-tls-chain.pem --tls
   ```

## Upgrading Chaincode
To upgrade the chaincode, run the following commands

```
docker exec cli peer chaincode install -n blossomcc -v {VERSION} -p github.com/usnistgov/blossom/chaincode  
docker exec cli peer chaincode upgrade -o $ORDERER -C blossom-1 -n blossomcc -v {VERSION} -c '{"Args":["init"]}' --cafile /opt/home/managedblockchain-tls-chain.pem --tls
```

## Building
From the chaincode root directory run `go build`.

## APIs

  - Account: Request a Blossom account and modify account information.
  - Asset: Onboard software assets and transact with them.
  - SwID: Report SwID tags.

## NGAC
There is a NGAC Policy Enforcement Point (PEP) controlling access to each API function.  The user sending the request must
have the necessary permissions to carry out the request or an error will occur. The NGAC access control policies are 
administered manually using the [policy-machine-go](https://github.com/PM-Master/policy-machine-go) library.
The **pap** package contains the code to build the initial NGAC graph configuration and to update the graph in response 
to API functions being called.

### Super User
NGAC requires a super user to create the initial configuration. This user will be responsible for accepting Blossom account 
requests and managing assets. Users in NGAC are defined using their username and Membership Service Provider (MSP) ID in the format:
`<username>:<mspid>`.  

- The Blossom super user is defined as: `super:BlossomMSP`.

On initial start up the super user must call the InitBlossom chaincode function.  This function initializes the NGAC graph
which is needed for any subsequent chaincode calls.
