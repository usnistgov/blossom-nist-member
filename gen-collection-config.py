#!/usr/bin/env python3
# Generate a blossom collection config

BLOCKS_TO_LIVE = 1000000

def gen_or_signature_policy(members: 'list[str]')->str:
    if len(members) == 0:
        raise Exception('Invalid policy produced with no member access')
    return f"OR('{', '.join([f'{member}.member' for member in members])}')"

def gen_single_collection_config(name: str, participants: 'list[str]'):
    return {
        'name': name,
        'policy': gen_or_signature_policy(participants),
        'requiredPeerCount': 1 if len(participants) > 1 else 0,
        'maxPeerCount': len(participants) - 1,
        'blocksToLive': BLOCKS_TO_LIVE,
        'memberOnlyRead': True,
        'memberOnlyWrite': True
    }

def gen_collection_config(admin: str, approved: 'list[str]', unapproved: 'list[str]'):
    return [
        gen_single_collection_config('catalog_coll', [admin, *approved]),
        gen_single_collection_config('licenses_coll', [admin]),
        *[
            gen_single_collection_config(f'{member}_account_coll', [admin, member])
            for member in [*approved, *unapproved]
        ]
    ]

if __name__ == '__main__':
    from argparse import ArgumentParser
    import json

    parser = ArgumentParser('gen-collection-config.py',
        description='Generate a collection config')
    parser.add_argument('--admin', type=str, required=True,
        help="The NGAC admin member's ID")
    parser.add_argument('--approved', default=[], nargs="*",
        help='IDs of members who have an account')
    parser.add_argument('--unapproved', default=[], nargs="*",
        help='IDs of members who do not have an account yet')
    args = parser.parse_args()

    collection_config = gen_collection_config(**args.__dict__)
    print(json.dumps(collection_config, indent=4))
