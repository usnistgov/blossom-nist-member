# Blossom Smart Contracts
This package contains the code for the Blossom Smart Contracts.

## APIs

  - Agency: Request a Blossom account and modify account information.
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
