import {
    Context,
    APIGatewayProxyCallback,
    APIGatewayEvent,
} from 'aws-lambda';
import { HandlerFunc, invokeHandler, queryHandler, pinError } from './handlers';

/**
 * 
 * @param event AWS gateway event
 * @param context BloSS🌻M context to properly dispatch event with more details
 * @param callback Gateway callback to communicate returned information
 */
export const handler = async (
    event: APIGatewayEvent,
    context: Context,
    callback: APIGatewayProxyCallback
) => {
    console.log(`index.ts-L19-Event: ${JSON.stringify(event, null, 2)}`);
    console.log(`index.ts-L20-Context: ${JSON.stringify(context, null, 2)}`);


    const bodyJson = JSON.parse(event.body ?? '');
    console.log(`index.ts-L24: ${bodyJson['functionType']}`);
    console.log(`index.ts-L25: ${JSON.stringify(context, null, 2)}`);
    let handlerFunc: HandlerFunc;
    switch (bodyJson['functionType']) {
        case 'query':
            handlerFunc = queryHandler;
            break;
        case 'invoke':
            handlerFunc = invokeHandler;
            break;
        default:
            throw new Error(`${pinError(new Error())} Request body "functionType" must be one of "query" or "invoke"`);
    }

    try {
        const result = await handlerFunc(event, bodyJson);
        callback(null, result);
    } catch (error) {
        callback(`${pinError(new Error())} ${error}`);
    }
};
