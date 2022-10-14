#!/usr/bin/env python3
# Generate a blossom collection config

BLOCKS_TO_LIVE = 1000000


def gen_or_signature_policy(members: 'list[str]') -> str:
    if len(members) == 0:
        raise Exception('Invalid policy produced with no member access')
    return f"OR('{', '.join([f'{member}.member' for member in members])}')"


def gen_single_collection_config(name: str, participants: 'list[str]', blocksToLive: int, memberOnlyRead: bool):
    return {
        'name': name,
        'policy': gen_or_signature_policy(participants),
        'requiredPeerCount': 1 if len(participants) > 1 else 0,
        'maxPeerCount': len(participants) - 1,
        'blocksToLive': blocksToLive,
        'memberOnlyRead': memberOnlyRead,
        'memberOnlyWrite': True
    }


def gen_collection_config(admin: str, approved: 'list[str]', unapproved: 'list[str]', blocksToLive: int, migrate: bool):
    migration = [
        # generate stub catalog_coll that only allows the admin member to read.
        gen_single_collection_config(
            'catalog_coll', [admin], blocksToLive, True)
    ] if migrate else []

    return migration + [
        gen_single_collection_config(
            'catalog_coll_v2', [admin], blocksToLive, False),
        gen_single_collection_config(
            'licenses_coll', [admin], blocksToLive, True),
        *[
            gen_single_collection_config(f'{member}_account_coll', [
                                         admin, member], blocksToLive, True)
            for member in [*approved, *unapproved]
        ]
    ]


if __name__ == '__main__':
    from argparse import ArgumentParser
    import json

    parser = ArgumentParser('gen-collection-config.py',
                            description='Generate a collection config')
    parser.add_argument('--admin', type=str, required=True,
                        help="The NGAC admin member's ID (usually starts with m-...)")
    parser.add_argument('--approved', default=[], nargs="*",
                        help='IDs of members who have an account (usually starts with m-...)')
    parser.add_argument('--unapproved', default=[], nargs="*",
                        help='IDs of members who do not have an account yet (usually starts with m-...)')
    parser.add_argument('--blocksToLive', default=0,
                        help='How many blocks data should live in a collection (0 for forever)')
    parser.add_argument('--migrate', action='store_true',
                        help='Migrate from a chaincode version that had older catalog configs')
    args = parser.parse_args()

    collection_config = gen_collection_config(**args.__dict__)
    print(json.dumps(collection_config, indent=4))
