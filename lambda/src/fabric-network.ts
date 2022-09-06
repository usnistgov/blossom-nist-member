import { Gateway, Wallets, Wallet, Network, Contract } from 'fabric-network';
import path from 'path';
import * as fs from 'fs';
import YAML from 'yaml';
import { getSecret } from './aws';

async function buildIdentity(username: string) {
    const identity = {
        credentials: {
            certificate: await getSecret(`blossom/${username}/cert`),
            privateKey: await getSecret(`blossom/${username}/pk`),
        },
        mspId: await getSecret(`blossom/${username}/mspId`),
        type: 'X.509'
    };
    const wallet = await Wallets.newInMemoryWallet();
    wallet.put(username, identity);
    return { wallet, identity };
}

const CONN_PROFILE_PATH = path.join(__dirname, "./connection-profile.yaml")

export async function setupNetwork(username: string, channel: string) {
    const { identity, wallet } = await buildIdentity(username);
    const profile = YAML.parse(fs.readFileSync(CONN_PROFILE_PATH).toString());

    const gateway = new Gateway();
    await gateway.connect(profile, {
        wallet,
        identity,
        discovery: {
            asLocalhost: false,
            enabled: false,
        }
    });

    return await gateway.getNetwork(channel);
}
