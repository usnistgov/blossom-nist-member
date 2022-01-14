#!/bin/bash

setOrg1Admin() {
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="Org1MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
  export CORE_PEER_ADDRESS=localhost:7051
}

setOrg2Admin() {
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="Org2MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
  export CORE_PEER_ADDRESS=localhost:9051
}

setOrg3Admin() {
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="Org3MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp
  export CORE_PEER_ADDRESS=localhost:11051
}

setOrg1User1() {
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="Org1MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp
  export CORE_PEER_ADDRESS=localhost:7051
}

#setOrg1User2() {
#  export CORE_PEER_LOCALMSPID=Org1MSP
#  export PEER0_ORG1_CA=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
#  export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
#  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/User2@org1.example.com/msp
#  export CORE_PEER_ADDRESS=peer0.org1.example.com:9051
#}

setOrg2User1() {
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="Org2MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp
  export CORE_PEER_ADDRESS=localhost:9051
}

#setOrg2User2() {
#  export CORE_PEER_LOCALMSPID=Org2MSP
#  export PEER0_ORG2_CA=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
#  export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
#  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/User2@org2.example.com/msp
#  export CORE_PEER_ADDRESS=peer0.org2.example.com:9051
#}

setOrg3User1() {
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="Org3MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org3.example.com/users/User1@org3.example.com/msp
  export CORE_PEER_ADDRESS=localhost:11051
}

#setOrg3User2() {
#  export CORE_PEER_LOCALMSPID=Org3MSP
#  export PEER0_ORG3_CA=${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt
#  export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG3_CA
#  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org3.example.com/users/User2@org3.example.com/msp
#  export CORE_PEER_ADDRESS=peer0.org3.example.com:11051
#}

setUser() {
  user=$1
  if [ "$1" == "Org1Admin" ]; then
    setOrg1Admin
  elif [ "$1" == "Org2Admin" ]; then
    setOrg2Admin
  elif [ "$1" == "Org3Admin" ]; then
    setOrg3Admin
  elif [ "$1" == "Org2User1" ]; then
    setOrg2User1
  elif [ "$1" == "Org2User2" ]; then
    setOrg2User2
  elif [ "$1" == "Org3User1" ]; then
    setOrg3User1
  elif [ "$1" == "Org3User2" ]; then
    setOrg3User2
  fi
}

InitNGAC() {
  setUser $1
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c '{"Args":["InitNGAC"]}'
}

OnboardAsset() {
  setUser $1
  asset=$2
  export LICENSES=$(echo -n "{\"licenses\":[{\"license_id\": \"asset$asset-license-1\", \"expiration\": \"exp\"}, {\"license_id\": \"asset$asset-license-2\", \"expiration\": \"exp\"}, {\"license_id\": \"asset$asset-license-3\", \"expiration\": \"exp\"}, {\"license_id\": \"asset$asset-license-4\", \"expiration\": \"exp\"}]}" | base64 | tr -d \\n)
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["OnboardAsset", "10'"$asset"'", "asset'"$asset"'", "01/01/2022", "01/01/2025"]}' --transient "{\"asset\":\"$LICENSES\"}"
}

Assets() {
  setUser $1
  peer chaincode query -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc -c  '{"Args":["Assets"]}'
}

AssetInfo() {
  setUser $1
  asset=$2
  peer chaincode query -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n blossomcc -c  '{"Args":["AssetInfo", "10'"$asset"'"]}'
}

RequestAccount() {
  setUser $1
  sysOwner=$2
  sysAdmin=$3
  acqSpec=$4
  export ACCOUNT=$(echo -n "{\"system_owner\":\"$sysOwner\",\"system_admin\":\"$sysAdmin\",\"acquisition_specialist\": \"$acqSpec\"}" | base64 | tr -d \\n)
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc -c  \
    '{"Args":["RequestAccount"]}' --transient "{\"account\":\"$ACCOUNT\"}"
}

ApproveOrg2Account() {
  setUser $1
  account=$2
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["ApproveAccount", "Org2MSP"]}'
}

ApproveOrg3Account() {
  setUser $1
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["ApproveAccount", "Org3MSP"]}'
}

UploadATOOrg2() {
  setUser $1
  export ATO=$(echo -n "{\"ato\":\"org2 test ato\"}" | base64 | tr -d \\n)
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["UploadATO"]}' --transient "{\"ato\":\"$ATO\"}"
}

UploadATOOrg3() {
  setUser $1
  export ATO=$(echo -n "{\"ato\":\"org3 test ato\"}" | base64 | tr -d \\n)
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["UploadATO"]}' --transient "{\"ato\":\"$ATO\"}"
}

Accounts() {
  setUser $1
  peer chaincode query -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc -c  '{"Args":["Accounts"]}'
}

Account() {
  setUser $1
  account=$2
  peer chaincode query -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n blossomcc -c  '{"Args":["Account", "'"$account"'"]}'
}

UpdateOrg2Status() {
  setUser $1
  status=$2
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["UpdateAccountStatus", "Org2MSP", "'"$status"'"]}'
}

UpdateOrg3Status() {
  setUser $1
  status=$2
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["UpdateAccountStatus", "Org3MSP", "'"$status"'"]}'
}

Org2RequestCheckout() {
  setUser $1
  asset=$2
  amount=$3
  export CHECKOUT=$(echo -n "{\"asset_id\":\"10$asset\",\"amount\":$amount}" | base64 | tr -d \\n)
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["RequestCheckout"]}' --transient "{\"checkout\":\"$CHECKOUT\"}"
}

Org3RequestCheckout() {
  setUser $1
  asset=$2
  amount=$3
  export CHECKOUT=$(echo -n "{\"asset_id\":\"10$asset\",\"amount\":$amount}" | base64 | tr -d \\n)
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["RequestCheckout"]}' --transient "{\"checkout\":\"$CHECKOUT\"}"
}

CheckoutRequests() {
  setUser $1
  account=$2
  peer chaincode query -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c '{"Args":["CheckoutRequests", "'"$account"'"]}'
}

ApproveCheckout() {
  setUser $1
  account=$2
  asset=$3
  export CHECKOUT=$(echo -n "{\"account\":\"$account\",\"asset_id\":\"10$asset\"}" | base64 | tr -d \\n)
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["ApproveCheckout"]}' --transient "{\"checkout\":\"$CHECKOUT\"}"
}

InitiateCheckin() {
  setUser $1
  asset=$2
  export CHECKIN=$(echo -n "{\"asset_id\":\"10$asset\",\"licenses\":[\"asset1-license-1\", \"asset1-license-2\"]}" | base64 | tr -d \\n)
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["InitiateCheckin"]}' --transient "{\"checkin\":\"$CHECKIN\"}"
}

InitiatedCheckins() {
  setUser $1
  account=$2
  peer chaincode query -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["InitiatedCheckins", "'"$account"'"]}'
}

ProcessCheckin() {
  setUser $1
  asset_id=$2
  account=$3
  export CHECKIN=$(echo -n "{\"asset_id\":\"10$asset_id\",\"account\":\"$account\"}" | base64 | tr -d \\n)
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["ProcessCheckin"]}' --transient "{\"checkin\":\"$CHECKIN\"}"
}

Licenses() {
  setUser $1
  account=$2
  asset=$3
  peer chaincode query -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n blossomcc -c  '{"Args":["Licenses", "'"$account"'", "10'"$asset"'"]}'
}

ReportSwID() {
  setUser $1
  account=$2
  export SWID=$(echo -n "{\"account\":\"$account\",\"primary_tag\":\"123\",\"asset\":\"101\",\"license\":\"asset1-license-1\",\"xml\":\"<swid></swid>\"}" | base64 | tr -d \\n)
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["ReportSwID"]}' --transient "{\"swid\":\"$SWID\"}"
}

DeleteSwID() {
  setUser $1
  account=$2
  tag=$3
  export SWID=$(echo -n "{\"account\":\"$account\",\"primary_tag\":\"$tag\"}" | base64 | tr -d \\n)
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["DeleteSwID"]}' --transient "{\"swid\":\"$SWID\"}"
}

GetSwID() {
  setUser $1
  account=$2
  tag=$3
  export SWID=$(echo -n "{\"account\":\"$account\",\"primary_tag\":\"$tag\"}" | base64 | tr -d \\n)
  peer chaincode query -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n blossomcc -c  '{"Args":["GetSwID"]}' --transient "{\"swid\":\"$SWID\"}"
}

GetSwIDsAssociatedWithAsset() {
  setUser $1
  account=$2
  asset_id=$3
  peer chaincode query -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
    -C mychannel -n blossomcc \
    -c  '{"Args":["GetSwIDsAssociatedWithAsset", "'"$account"'", "10'"$asset_id"'"]}'
}

func=$1

if [ "$func" == "InitNGAC" ]; then
  InitNGAC $2
elif [ "$func" == "OnboardAsset" ]; then
  OnboardAsset $2 $3
elif [ "$func" == "Assets" ]; then
  Assets $2 | python -m json.tool
elif [ "$func" == "AssetInfo" ]; then
  AssetInfo $2 $3 | python -m json.tool
elif [ "$func" == "RequestAccount" ]; then
  # user, system owner, system admin, acq spec
  RequestAccount $2 $3 $4 $5
elif [ "$func" == "ApproveOrg2Account" ]; then
  ApproveOrg2Account $2
elif [ "$func" == "ApproveOrg3Account" ]; then
  ApproveOrg3Account $2
elif [ "$func" == "UploadATOOrg2" ]; then
  UploadATOOrg2 $2
elif [ "$func" == "UploadATOOrg3" ]; then
  UploadATOOrg3 $2
elif [ "$func" == "Accounts" ]; then
  Accounts $2 | python -m json.tool
elif [ "$func" == "Account" ]; then
  Account $2 $3 | python -m json.tool
elif [ "$func" == "UpdateOrg2Status" ]; then
  UpdateOrg2Status $2 $3
elif  [ "$func" == "UpdateOrg3Status" ]; then
  UpdateOrg3Status $2 $3
elif [ "$func" == "Org2RequestCheckout" ]; then
  Org2RequestCheckout $2 $3 $4
elif [ "$func" == "Org3RequestCheckout" ]; then
  Org3RequestCheckout $2 $3 $4
elif [ "$func" == "CheckoutRequests" ]; then
  CheckoutRequests $2 $3
elif [ "$func" == "ApproveCheckout" ]; then
  ApproveCheckout $2 $3 $4
elif [ "$func" == "InitiateCheckin" ]; then
  InitiateCheckin $2 $3 | python -m json.tool
elif [ "$func" == "InitiatedCheckins" ]; then
  InitiatedCheckins $2 $3 $4
elif [ "$func" == "ProcessCheckin" ]; then
  ProcessCheckin $2 $3 $4
elif [ "$func" == "Licenses" ]; then
  Licenses $2 $3 $4 | python -m json.tool
elif [ "$func" == "ReportSwID" ]; then
  ReportSwID $2 $3
elif [ "$func" == "DeleteSwID" ]; then
  DeleteSwID $2 $3 $4
elif [ "$func" == "GetSwID" ]; then
  GetSwID $2 $3 $4 | python -m json.tool
elif [ "$func" == "GetSwIDsAssociatedWithAsset" ]; then
  GetSwIDsAssociatedWithAsset $2 $3 $4 | python -m json.tool
fi