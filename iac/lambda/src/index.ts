import {
    Context,
    APIGatewayProxyCallback,
    APIGatewayEvent,
} from 'aws-lambda';
import { HandlerFunc, invokeHandler, queryHandler } from './handlers';


export const handler = async (
    event: APIGatewayEvent,
    context: Context,
    callback: APIGatewayProxyCallback
) => {
    console.log(`Event: ${JSON.stringify(event, null, 2)}`);
    console.log(`Context: ${JSON.stringify(context, null, 2)}`);

    const bodyJson = JSON.parse(event.body ?? '');

    let handlerFunc: HandlerFunc;
    switch (bodyJson['functionType']) {
        case 'query':
            handlerFunc = queryHandler;
            break;
        case 'invoke':
            handlerFunc = invokeHandler;
            break;
        default:
            throw new Error('Request body "functionType" must be one of "query" or "invoke"');
    }

    try {
        const result = await handlerFunc(event, bodyJson);
        callback(null, result);
    } catch (error) {
        callback(`${error}`);
    }
};
