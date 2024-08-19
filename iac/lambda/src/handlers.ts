import { APIGatewayEvent, APIGatewayProxyResult } from "aws-lambda";
import { setupNetwork } from "./fabric-network";

// const CHANNEL_NAME = process.env.CHANNEL_NAME ?? 'acquisition';
// const CONTRACT_NAME = process.env.CONTRACT_NAME ?? 'blossom';
const CHANNEL_NAME = process.env.CHANNEL_NAME ?? 'asset';
const CONTRACT_NAME = process.env.CONTRACT_NAME ?? 'authorization';

export type HandlerFunc = (event: APIGatewayEvent, bodyJson: any) => Promise<APIGatewayProxyResult>;

function getUsername(event: APIGatewayEvent): string {
    const username = event.requestContext.authorizer?.claims.username;
    if (username === undefined || username === null) {
        throw new Error(`Could not get username from requestContext (got ${JSON.stringify(event.requestContext.authorizer)})`);
    }
    return username as string;
}

type TransactionRequestBody = {
    channel: string;
    contract: string;
    functionName: string;
    args: string[];
    transient?: Record<string, string>;
}

/**
 * Convert string-string map to string-buffer
 */
function convertTransientToBuffer(transient: Record<string, string>) {
    return Object.keys(transient).reduce<{
        [key: string]: Buffer;
    }>((acc, key) => {
        acc[key] = Buffer.from(transient[key]);
        return acc;
    }, {})
}

const transactionHandler = async (event: APIGatewayEvent, bodyJson: any, type: 'query' | 'invoke'): ReturnType<HandlerFunc> => {
    console.log('Getting username...');
    const body = bodyJson as TransactionRequestBody;
    const username = getUsername(event);
    console.log('Setting up network...');
    const network = await setupNetwork(username, body.channel);
    console.log('Setting up contract...');
    const transaction = network.getContract(body.contract).createTransaction(body.functionName);
    if (body.transient) {
        transaction.setTransient(convertTransientToBuffer(body.transient));
    }

    transaction.setEndorsingOrganizations(network.getGateway().getIdentity().mspId);

    console.log('Evaluating/submitting transaction...');
    try {
        let result;
        if (type === 'query') {
            result = await transaction.evaluate(...body.args);
        } else {
            result = await transaction.submit(...body.args);
        }
        return {
            body: result.toString(),
            headers: {
                'Content-Type': 'application/json'
            },
            statusCode: 200
        };
    } catch (e) {
        return {
            body: `Error: ${e}`,
            headers: {},
            statusCode: 500,
        }
    } finally {
        network.getGateway().disconnect();
    }
}

export const queryHandler: HandlerFunc = (event, bodyJson) => transactionHandler(event, bodyJson, 'query');
export const invokeHandler: HandlerFunc = (event, bodyJson) => transactionHandler(event, bodyJson, 'invoke');
