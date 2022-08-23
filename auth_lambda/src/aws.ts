import SecretsManager from "aws-sdk/clients/secretsmanager";

const secrets = new SecretsManager();

export function getSecret(key: string): Promise<string> {
    // aws's secret API is a bit annoying so we're wrapping it in a promise
    return new Promise((resolve, reject) => {
        secrets.getSecretValue({ SecretId: key }, (err, data) => {
            if (err) {
                return reject(err);
            }

            if (data.SecretString) {
                return resolve(data.SecretString);
            } else if (data.SecretBinary) {
                return resolve(data.SecretBinary.toString('ascii'));
            } else {
                return reject('Secret string or secret binary not provided');
            }
        })
    });
}
