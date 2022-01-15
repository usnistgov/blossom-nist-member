#!/usr/bin/env python3
# Generate a connection profile

import boto3
from itertools import chain

TLS_CA_CERT_PATH = '/tmp/managedblockchain-tls-chain.pem'

def gen_channels(channels: 'list', orderer_name: str, peers_dict: 'dict'):
    # flatten dict of lists into list
    peers = list(chain(*(peers_dict.values())))

    return {
        channel: {
            'orderers': [orderer_name],
            'peers': {
                node['Id']: {
                    'chaincodeQuery': True,
                    'ledgerQuery': True,
                    'endorsingPeer': True,
                    'eventSource': True
                }
                for node in peers
            }
        }
        for channel in channels
    }

def gen_orderers(network):
    # looks like orderer.(...).amazon.com:300xx
    endpoint = network['FrameworkAttributes']['Fabric']['OrderingServiceEndpoint']

    return {
        f'orderer-{network["Name"]}': {
            'url': f'grpcs://{endpoint}',
            'grpcsOptions': {
                # strip port
                'ssl-target-name-override': endpoint.split(':')[0],
            },
            'tlsCACerts': {
                'path': TLS_CA_CERT_PATH
            }
        }
    }

def gen_organizations(members, peers_dict):
    return {
        member['Name']: {
            'mspid': member['Id'],
            'peers': [
                peer['Id']
                for peer in peers_dict[member['Name']]
            ],
            "certificateAuthorities": [f'ca-{member["Name"]}']
        }
        for member in members
    }

def gen_peers(peers_dict):
    return {
        peer['Id']: {
            'url': f'grpcs://{peer["FrameworkAttributes"]["Fabric"]["PeerEndpoint"]}',
            'eventUrl': f'grpcs://{peer["FrameworkAttributes"]["Fabric"]["PeerEventEndpoint"]}',
            'grpcsOptions': {
                'ssl-target-name-override': peer["FrameworkAttributes"]["Fabric"]["PeerEndpoint"].split(':')[0]
            },
            'tlsCACerts': {
                'path': TLS_CA_CERT_PATH
            }
        }
        for peer in list(chain(*(peers_dict.values())))
    }

def gen_certificate_authorities(members):
    return {
        f'ca-{member["Id"]}': {
            'url': member['FrameworkAttributes']['Fabric']['CaEndpoint'],
            'httpOptions': {
                'verify': False
            },
            'tlsCACerts': {
                'path': TLS_CA_CERT_PATH
            },
            'caName': member['Id']
        }
        for member in members
    }

def gen_connection_profile(network_id: str, channels: 'list[str]'):
    client = boto3.client('managedblockchain')
    
    network = client.get_network(NetworkId=network_id)['Network']
    # get a list of member summaries, then get the actual member objects
    members = [
        client.get_member(NetworkId=network_id, MemberId=summary['Id'])['Member']
        for summary in client.list_members(NetworkId=network_id)['Members']
    ]
    # for each member, get a list of node summaries, then get the actual node objects
    nodes = {
        member['Name']: [
            client.get_node(NetworkId=network_id, MemberId=member['Id'], NodeId=summary['Id'])['Node']
            for summary in client.list_nodes(NetworkId=network_id, MemberId=member['Id'])['Nodes']
        ]
        for member in members
    }

    network_name = network['Name']
    orderer_name = f'orderer-{network["Name"]}'

    return {
        'name': network_name,
        'x-type': 'hlfv1',
        'description': f'Generated connecction profile',
        'version': '1.0',
        'channels': gen_channels(channels, orderer_name, nodes),
        'orderers': gen_orderers(network),
        'organizations': gen_organizations(members, nodes),
        'peers': gen_peers(nodes),
        'certificateAuthorities': gen_certificate_authorities(members),
    }

if __name__ == '__main__':
    from argparse import ArgumentParser
    import json
    
    parser = ArgumentParser('gen-connection-profile.py',
        description='Generate a connection profile')
    parser.add_argument('--network_id', type=str, required=True,
        help="The network id (starts with n_...)")
    parser.add_argument('--channels', default=[], nargs="*",
        help='Channels to include in the profile')
    args = parser.parse_args()

    connection_profile = gen_connection_profile(**args.__dict__)
    print(json.dumps(connection_profile, indent=4))
