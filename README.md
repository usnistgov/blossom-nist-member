# Blossom Smart Contracts
This package contains the code for the Blossom Smart Contracts.

## Table of Contents
- [Local Testing](#local-testing)
- [Deployment Steps](#chaincode-deployment-steps)
   - [Fabric 2.2](#fabric-22)
     - [Updating Chaincode](#updating-chaincode)
   - [Fabric 1.4](#fabric-14)
- [Adding an Organization](#adding-an-organization) 
- [NGAC](#ngac)
- [Smart Contract Usage](#usage)
- [Private Data Collection Design Doc](docs/pdc-design.pdf)

## Local Testing

To deploy the chaincode locally for testing, you can use the `IBM Blockchain Extension` for VS Code.

In order to test it locally, ensure that:
1. Go and Docker is installed on your machine and your user is part of the `docker` group
2. You have the `IBM Blockchain Extension` for VS Code: https://marketplace.visualstudio.com/items?itemName=IBMBlockchain.ibm-blockchain-platform
3. You have cloned the chaincode to a folder under your `$GOPATH`, or you have created a symlink from the chaincode to a folder under your `$GOPATH` as shown below:
```bash
# ensure $GOPATH is set to /home/<your username>/go
$ echo $GOPATH
# create appropriate folder under $GOPATH
$ mkdir -p $GOPATH/src/github.com/usnistgov/blossom
# working from the chaincode folder, create a symbolic link
$ ln -s $(pwd) $GOPATH/src/github.com/usnistgov/blossom/chaincode
```
4. Your VS Code instance is working from the `$GOPATH` symlink (this is important!)
5. You have downloaded the required go dependencies:
```bash
# working from the chaincode folder !!
go mod vendor
go mod tidy
```

From there you can deploy the test environment using the following steps:
1. Open the `IBM Blockchain Platform` side panel
2. Create a new `Fabric Environment` from the 1 org template
3. Hit `+ Deploy smart contract` from the `my channel` dropdown
4. Select `chaincode (open project)`, give it a name and a version number, and hit follow the prompts.

## Chaincode Deployment Steps

### Setting Admin Membership Service Provider ID (MSPID)
The first step before doing anything else is to set the Administrative MSPID in the code.  This will ensure
all peers that install the chaincode will have the same Admin MSPID set.  If two peers have different values for the Admin MSPID, 
their packages will have different hashes and will fail the commit stage for approving two different packages.

1. In [adminmsp/adminmsp.go](/adminmsp/adminmsp.go) set the value of `AdminMSP` to the Admin MSP of the deployment. 

   **Example:** `const AdminMSP = "SAMS-MSPID"`

### Lifecycle Endorsement Policy

In the `configtx.yaml` used to create the channel.  Modify the `Application > Policies > LifecycleEndorsement` policy to:

```yaml
LifecycleEndorsement:
  Type: Signature
  Rule: “AND('SAMS-MSPID.member', OutOf(2, 'NIST-MSPID.member', 'DHS-MSPID.member'))"
```

The `OutOf` function will need to be updated everytime an organization is added to the network.  The new organization should 
be added to the list (i.e. `OutOf(2, 'NIST-MSPID.member', 'DHS-MSPID.member', NewOrg-MSPID.member)`) and the `2` should be updated
to ensure it is a majority of the members in the list.

### Fabric 2.2

1. Package chaincode on each peer.
   
   ```shell
   peer lifecycle chaincode package blossomcc.tar.gz --path <path to chaincode directory> --lang golang --label blossomcc_1.0
   ```   
   
   This will package the chaincode into a file called `blossomcc.tar.gz`.
   

2. Install chaincode on each peer.
   
   ```shell
   peer lifecycle chaincode install blossomcc.tar.gz
   ```

3. Get chaincode package ID.

   ```shell
   peer lifecycle chaincode queryinstalled
   ```
   
   - Look for the label that matches the label set in step 1.
   

4. Approve chaincode definition.

   ```shell
   lifecycle chaincode approveformyorg \
    -o $ORDERER \
    --tls --cafile $ORDERER_CA \
    --channelID $CHANNEL --name blossomcc --package-id $PACKAGE_ID \
    --collections-config <path to collections_config.json> \
    --version 1.0 --sequence 1
   ```
   
   - This command will need to be executed by enough organizations to satisfy the [policy](#lifecycle-endorsement-policy) defined in the channel's `configtx.yaml`.
   

5. Check commit readiness

   ```shell
   peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL --name blossomcc --version 1.0 --sequence 1 --tls --cafile $ORDERER_CA --output json
   ```
   
   - This command will show which organizations on the channel have approved the chaincode and which ones haven't.


6. Commit chaincode.

   ```shell
   peer lifecycle chaincode commit \
    -o $ORDERER \
    --tls --cafile $ORDERER_CA \
    --channelID $CHANNEL --name blossomcc \
    --peerAddresses <PEER_ADDRESS> --tlsRootCertFiles <path to peer's tls ca cert> \
    --version 1.0 --sequence 1 --collections-config <path to collections_config.json>
   ```

   - `--peerAddresses`
      - 1 or more peers that have approved the chaincode to target for commit.

   - This is when the [lifecycle endorsement policy](#lifecycle-endorsement-policy) will be checked. An endorsement policy error
   will be returned if not enough organizations have approved the chaincode to satisfy the policy.


5. Invoke.

   ```shell
   peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile $ORDERER_CA \
    -C $CHANNEL -n blossomcc \
    --peerAddresses <PEER_ADDRESS> --tlsRootCertFiles <path to peer's tls ca cert> \
    -c '{"function":"test","Args":["hello world"]}'
   ```
   
   - `--peerAddresses`
     - 1 or more peers that have approved the chaincode to target for invoke.
   - If an org did not approve the chaincode in step 3,  they will need to target a org that did or else an error will occur.
   - If an org did approve the chaincode, they do not need to target another peer.


#### Updating Chaincode

Repeat steps 3 and 4 above incrementing the sequence and version flags.

### Fabric 1.4

1. Make sure the Blossom project is cloned on the peer machine.  The path provided in the following `install` command
   assumes the chaincode is located in `$GOPATH`.

2. Install chaincode on the peer
   
   ```
   docker exec cli peer chaincode install -n blossomcc -v {VERSION} -p github.com/usnistgov/blossom/chaincode
   ```


3. Instantiate chaincode the chaincode on the channel `blossom-1`
   
   ```
   docker exec cli peer chaincode instantiate -o $ORDERER -C blossom-1 -n blossomcc -v {VERSION} -c '{"Args":["init", "<ADMIN_MSP>"]}' --cafile /opt/home/managedblockchain-tls-chain.pem --tls --collections-config /opt/gopath/github.com/usnistgov/blossom/chaincode/collections_config.json
   ```
   
    - **IMPORTANT:** Replace <ADMIN_MSP> with the MSPID of the administrative member of the network


4. Check chaincode instantiation

   ```
   docker exec cli peer chaincode list --instantiated -o $ORDERER -C blossom-1 --cafile /opt/home/managedblockchain-tls-chain.pem --tls
   ```

5. Invoke chaincode

   ```
   docker exec cli peer chaincode invoke -C blossom-1 -n blossomcc -c  '{"Args":["test", "awesome blossom"]}' -o $ORDERER --cafile /opt/home/managedblockchain-tls-chain.pem --tls
   ```

## Adding an Organization

When a new organization is added to the network, two things must happen:

1. They must have their own PDC.
2. They must be added to the `catalog_coll` PDC, ONLY when the account status is `AUTHORIZED`. 

Having their own PDC will allow them to upload an ATO, and in the future, checkout licenses.  Having access to the `catalog_coll`, will allow them to view the software assets available for lease.

### Organization Collection
Once a new member is added to the network, we must update the chaincode definition to create a Private Data Collection 
for the new member. Use the below JSON as a template for creating a new PDC for the account in `collections_config.json`.
   
**IMPORTANT: This should be done during the enrollment process, before `RequestAccount` is called.**

  ```json
  {
    "name":"<Account MSPID>_coll",
    "policy":"OR('BlossomMSP.member', '<Account MSPID>.member')",
    "requiredPeerCount":0,
    "maxPeerCount":2,
    "blockToLive":1000000,
    "memberOnlyRead":true,
    "memberOnlyWrite":true
  }
  ```

Once this collection is created, and the chaincode is upgraded, the account will be able to upload an ATO.

### Catalog Collection
Add the org's MSPID to the `catalog_coll` and increase the max peer count. Use the template below to update the `catalog_coll`
collection to include the new member.
   
**IMPORTANT: This should only be done when the account status is set to `AUTHORIZED` via the chaincode function `UpdateAccountStatus`. If an account is set to a status other than `AUTHORIZED`, the account MSPID should be removed from this collection definition, and the chaincode upgraded.**
        
  ```json
  {
    "name":"catalog_coll",
    "policy":"OR('BlossomMSP.member', 'Org1MSP.member', '<Account MSPID>.member')",
    "requiredPeerCount":0,
    "maxPeerCount":3,
    "blockToLive":1000000,
    "memberOnlyRead":true,
    "memberOnlyWrite":true
  }
  ```

### Upgrading Chaincode

There are three situations to upgrade chaincode:

1. New member enrollment
   
   - Use Organization Collection template [above](#organization-collection), and add to `collections_config.json`.
   
2. Account status set to `AUTHORIZED`

   - Add Organization MSPID to `catalog_coll` as shown [above](#catalog-collection).

3. Account status set to NOT `AUTHORIZED`

   - Remove ORganization MSPID from `catalog_coll`.
   
#### Upgrade Chaincode
Tp upgrade chaincode, install on all necessary peers.  Then, call upgrade with new `collections_config.json` file.

```bash
docker exec cli peer chaincode install -n blossomcc -v {VERSION} -p github.com/usnistgov/blossom/chaincode  
docker exec cli peer chaincode upgrade -o $ORDERER -C blossom-1 -n blossomcc -v {VERSION} -c '{"Args":["init", <ADMIN_MSPID>]}' --cafile /opt/home/managedblockchain-tls-chain.pem --tls --collections-config /opt/gopath/github.com/usnistgov/blossom/chaincode/collections_config.json
```

**Note: The updated chaincode must be installed and upgraded on all peers.**

## Building
From the chaincode root directory run `go build`.

## APIs

  - Account: Request a Blossom account and modify account information.
  - Asset: Onboard software assets and transact with them.
  - SwID: Report SwID tags.

## NGAC
### Administrative Users and Graph Initialization
- The user that calls `InitNGAC` must be in the Administrative MSP defined when the chaincode was instantiated.

- The user that calls `ApproveAccount` will be the admin user for that account. They also must be in the MSP defined during instantiation.

### Policy Definition
There are two NGAC policies to be used in the smart contracts found in [ngac/pap/policy.go](ngac/pap/policy.go).  
The first is the `Catalog policy`, which is initialized in the `InitNGAC` smart contract function. This policy allows the super
user to Onboard and Offboard assets in the catalog and gives them administrative control over the graph, allowing them to
delegate to other users.  This policy is stored in the **Catalog PDC**.

The second policy is the `Account policy`.  This policy is created each time a new account is created and saved in the
Account's PDC, meaning each account has its own graph that decisions are executed on.  Users do not have access to graphs
that belong to other accounts.  This policy is loaded when a user calls the RequestAccount function.  This policy grants the super user
full administrative permissions on the account.  It also creates a series of Obligations which define responses to certain
events that can happen. For example, before an Account's status is set to "AUTHORIZED" the System Admin does not have the 
necessary permissions to check out asset licenses.  Setting the Account's status to "AUTHORIZED" is defined as an event, and 
in response to that event, the System Admin is granted these permissions.

### Policy Decisions
NGAC policy decisions are made in the PDP located in [ngac/pdp/pdp.go](ngac/pdp/pdp.go). The functions available in this
package serve as helper functions to call the NGAC decision algorithm on nodes in the NGAC graphs described above.

### Events
NGAC event processing is done in the EPP located in [ngac/epp/epp.go](ngac/epp/epp.go). These functions also serve as helpers
to process events in the underlying NGAC implementation. 

## Usage

### Initialization
- InitNGAC
   - user: any user in blossom admin MSP set during instantiation

### Onboarding an asset
- OnboardAsset
   - user: super (BlossomMSP)
   - args: `["101","asset1","01/01/2025"]`
   - transient data: 
      ```json
      {
        "asset":"{\"licenses\":[\"asset1-license-1\", \"asset1-license-2\", \"asset1-license-3\", \"asset1-license-4\", \"asset1-license-5\"]}"
      }
      ```

### Creating an account and checking out an asset
1. **RequestAccount**
    - user: a1_system_owner (A1MSP)
    - args: `[]`
    - transient data:
      ```json
      {
        "account":"{\"system_owner\":\"a1_system_owner\",\"system_admin\":\"a1_system_admin\",\"acquisition_specialist\": \"a1_acq_spec\",\"ato\": \"a1 test ato\"}"
      }
      ```
   

2. **ApproveAccount**
   - super (BlossomMSP)
   - args: `["A1MSP"]`


2. **UpdateAccountStatus**
    - user: super (BlossomMSP)
    - args: `["A1MSP","AUTHORIZED"]`
    

3. **Add new account to collections config**

    1. Add account to the catalog collection
        ```json
        {
            "name":"catalog_coll",
            "policy":"OR('BlossomMSP.member', 'A1MSP.member', 'A2MSP.member')",
            "requiredPeerCount":0,
            "maxPeerCount":3,
            "blockToLive":1000000,
            "memberOnlyRead":true,
            "memberOnlyWrite":true
        }
        ```
    2. Add account's own PDC
        ```json
        {
            "name":"A1MSP_account_coll",
            "policy":"OR('BlossomMSP.member', 'A1MSP.member')",
            "requiredPeerCount":0,
            "maxPeerCount":2,
            "blockToLive":1000000,
            "memberOnlyRead":true,
            "memberOnlyWrite":true
        } 
        ```

4. **Install and upgrade chaincode on channel**
    ```
    docker exec cli peer chaincode install -n blossomcc -v {VERSION} -p github.com/usnistgov/blossom/chaincode  
    docker exec cli peer chaincode upgrade -o $ORDERER -C blossom-1 -n blossomcc -v {VERSION} -c '{"Args":["init"]}' --cafile /opt/home/managedblockchain-tls-chain.pem --tls
    ```

4. **RequestCheckout**
    - user: a1_system_admin (A1MSP)
    - args: `[]`
    - transient data:
      ```json
      {
        "checkout": "{\"asset_id\":\"101\",\"amount\":2}"
      }
      ```
      
5. **ApproveCheckout**
    - user: super (BlossomMSP)
   - args: `[]`
   - transient data:
     ```json
     {
        "checkout": "{\"account\":\"A1MSP\",\"asset_id\":\"101\"}"
     }
     ```
     
### More examples

- See the [vscode](vscode) directory for how to use the smart contracts using the IBM Blockchain Platform for VSCode.
- See [blossom-transactions.txdata](vscode/transaction_data/blossom-transactions.txdata) for example smart contract function calls.
