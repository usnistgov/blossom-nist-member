import { Context, APIGatewayProxyCallback, APIGatewayEvent } from 'aws-lambda';
import axios from 'axios';
import { getSecret } from './aws';

export const handler = async (
    event: APIGatewayEvent,
    context: Context,
    callback: APIGatewayProxyCallback
) => {
    console.log(`Event: ${JSON.stringify(event, null, 2)}`);
    console.log(`Context: ${JSON.stringify(context, null, 2)}`);

    const forwardUrl = process.env['FORWARD_URL'];

    // grab secret from aws secrets
    const cert = await getSecret(`/todo`);
    const pk = await getSecret(`/todo`);
    const mspid = await getSecret(`/todo`);

    const results = await axios.request({
        url: forwardUrl,
        method: event.httpMethod,
        headers: Object.keys(event.headers).reduce<Record<string, string>>((acc, curr, _) => {
            // remove undefined header entries
            const headerEntry = event.headers[curr];
            if (headerEntry) {
                if (acc[curr]) { // check to see the user isn't putting in a bad request
                    const error = new Error('Request attempting to override injected hyperledger fabric identity data');
                    callback(error);
                    throw error;
                }
                acc[curr] = headerEntry;
            }
            return acc;
        }, {
            HLFI_CERT: cert,
            HLFI_PK: pk,
            HLFI_MSPID: mspid
        }),
        data: event.body
    })

    // send result back to apigw
    callback(null, {
        statusCode: results.status,
        headers: results.headers,
        body: results.data,
    });
};
