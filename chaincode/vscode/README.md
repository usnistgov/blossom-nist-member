# Blossom Demo Using IBM Blockchain Platform

## Microfab Network
Start the microfab network with three organizations using `demo.sh`. The three organizations are:

- Blossom
- A1
- A2

### Connect the VSCode Extension to Network
Under `Fabric Environments` press `+`, then `Add a Microfab network`.  The URL should be `http://console.127-0-0-1.nip.io:8080`.

### Register Organization Users
Unfortunately I haven't found a way to automate this process using the vscode extension yet.  In order to add users to 
the organizations, press `command + shift + p`, then select `Create Identity (register and enroll)`.
The users that need to be added are as follows:

- super (BlossomMSP)
- a1_system_owner (A1MSP)
- a1_system_admin (A1MSP)
- a1_acq_spec (A1MSP)
- a2_system_owner (A2MSP)
- a2_system_admin (A2MSP)
- a2_acq_spec (A2MSP)

Note: In my experience restarting the network causes an issue with these identities and you will need to 
recreate them.

## Chaincode Deployment
Open the `blossom/chaincode` directory in VSCode. Press `command + shift + p` and select `Package Open Project`. Choose 
a name and version for the chaincode and press `Enter`.  In the `Fabric Environments` panel, expand `channel1` and press
`Deploy smart contract`.  Select the package and add `collections_config.json` as the collections config file. This file 
defines the private data collections.

## Chaincode Usage
In `Fabric Gateways` right-click on the deployed smart contract and select `Transact with smart contract`.  In the window
that appears, select `Transaction data directory` at the top (next to "Manual input"). From there, press `Add directory`,
and select the `transaction_data` folder. Once the directory is selected, you can browse through predefined transactions.
To update these transactions, just update the `blossom-transactions.txdata` file and reselect the directory.