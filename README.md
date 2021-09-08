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

### Initialize NGAC
To initialize the NGAC component the **super user** must call the chaincode function `InitNGAC`.  This function
will initialize the NGAC graph on the blockchain.

## Usage

#### 1. Initialize the NGAC graph

   - Function: InitNGAC

   - User: super

   - Args: none


#### 2. Onboard a sample asset

   - Function: OnboardAsset

   - Username: super

   - Args:

      - `
        {
           "id": "test-asset-id",
           "name": "test-asset",
           "total_amount": 10,
           "available": 10,
           "cost": 100.00,
           "onboarding_date": "",
           "expiration": "2025-01-01",
           "licenses": [
            "test-asset-1",
            "test-asset-2",
            "test-asset-3",
            "test-asset-4",
            "test-asset-5",
            "test-asset-6",
            "test-asset-7",
            "test-asset-8",
            "test-asset-9",
            "test-asset-10"
           ],
           "available_licenses": [],
           "checked_out": {}
        }
        `
   
#### 3. Request a blossom account

   - Function: RequestAccount

   - Username: a1_system_owner

   - Args:

      - `
        {
          "name": "Agency1",
          "ato": "this is a test ato",
          "mspid": "A1MSP",
          "users": {
            "system_owner": "a1_system_owner",
            "acquisition_specialist": "a1_acq_spec",
            "system_administrator": "a1_system_admin"
          },
          "status": "",
          "assets": {}
        }
        `


#### 4. Update Agency1 account status to 'Approved'

   - Function: UpdateAccountStatus

   - Username: super
   
   - Args:
   
      - `"Agency1"`
      - `"Approved"`


#### 5. View available assets

   - Function: Assets
     
   - Username: a1_system_admin
     
   - Args: none


#### 6. Agency1 checks out 2 licenses of the sample asset

   - Function: Checkout

   - Username: a1_system_admin
   
   - Args:

      - `"test-asset-id"`
      - `"Agency1"`
      - `2`


#### 7. Agency1 reports a SwID tag for a license

   - Function: ReportSwID

   - Username: a1_system_admin

   - Args:

      - `
        {
            "primary_tag": "swid-1",
            "xml": "<swid>test</swid>",
            "asset": "test-asset-id",
            "license": "test-asset-1",
            "lease_expiration": "%s"
        }
        `
        
      - `"Agency1"`


#### 8. Get SwIDs that are associated with the sample asset

   - Function: GetSwIDsAssociatedWithAsset

   - Username: a1_system_admin
   
   - Args:

      - `"test-asset-id"`
