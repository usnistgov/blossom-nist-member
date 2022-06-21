# Blossom Smart Contracts
This package contains the code for the Blossom Smart Contracts.

## Table of Contents
- [User Identities, Roles, and Access Control](#access-control-and-user-identities)
- [Local Testing](#local-testing)
- [Deployment Steps](#chaincode-deployment-steps)
   - [Fabric 2.2](#fabric-22)
     - [Updating Chaincode](#updating-chaincode)
   - [Fabric 1.4](#fabric-14)
- [Adding an Organization](#adding-an-organization) 
- [NGAC](#ngac)
- [Smart Contract Usage](#usage)
- [Private Data Collection Design Doc](docs/pdc-design.pdf)

## Access Control and User Identities

### NGAC
Next Generation Access Control (NGAC) controls access to chaincode functions. Users are assigned to attributes reflecting
their organization and their role within the organization. Organization accounts are assigned to attributes reflecting 
their current status in the blossom system. The organization and role a user is assigned to is determined by the attributes
in their Fabric identity. 

There are two Fabric attributes supported by Blossom:

   - `blossom.admin`: Specifies if the user is a blossom admin user (will also need to be in the Admin MSP)
     
      - true | false
      
   - `blossom.role`: Specifies the role for this user in an organization
     
      - SystemOwner | SystemAdministrator | AcquisitionSpecialist

### Roles
- **SystemOwner**: Can upload account ATOs
- **SystemAdministrator**: Can check out/check in licenses and report/delete SWID tags
- **AcquisitionSpecialist**: Can audit their account's licenses

### User Registration
Below are examples of registering a user with Blossom attributes.

*Note: The users MSPID is determined by the Fabric CA the identity is registered with*

- Using the node sdk
   ```javascript
   // create a blossom admin
   const secret = await caClient.register({
      affiliation: '',
      enrollmentID: 'admin1', // the username
      role: 'client',
      [
          {name: 'blossom.admin', value: 'true', ecert: true} // additional attribute for the blossom role
      ]
   }, adminUser);
  
   // create an organization system owner
   const secret = await caClient.register({
      affiliation: '',
      enrollmentID: 'org1_sys_owner', // the username
      role: 'client',
      [
          {name: 'blossom.role', value: 'SystemOwner', ecert: true} // additional attribute for the blossom role
      ]
   }, adminUser);
   ```

- Using the CLI
   ```
   # Create a blossom admin
   ./fabric-ca-client register ... --id.attrs 'blossom.admin=true' ...
  
   # Create a system owner
   ./fabric-ca-client register ... --id.attrs 'blossom.role=SystemOwner' ...
   ```

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
   peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL --name blossomcc --version 1.0 --sequence 1 --tls --cafile $ORDERER_CA --output json --collections-config <path to collections_config.json>
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
     - This is only needed if more than one peer is needed for endorsement.
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

Once this collection is created, and the **chaincode is upgraded**, the account will be able to upload an ATO.
   
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
    

3. **Install and upgrade chaincode on channel**

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
