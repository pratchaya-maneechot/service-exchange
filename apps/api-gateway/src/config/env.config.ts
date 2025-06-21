import { toNumber, getOsEnv } from 'src/common/utils/env.util';
import { EnvConfig } from './config.type';
import { join } from 'path';

export const env = (): EnvConfig => ({
  app: {
    port: toNumber(getOsEnv('PORT', '8080')),
    host: getOsEnv('HOST', '0.0.0.0'),
    nodeEnv: getOsEnv('NODE_ENV', 'development') as
      | 'production'
      | 'development'
      | 'test',
    corsOrigins: getOsEnv('CORS_ORIGINS', '*')?.split(',').filter(Boolean),
    appName: getOsEnv('APP_NAME', 'api-gateway'),
    appVersion: getOsEnv('APP_VERSION', 'vLocal-dev'),
  },
  redis: {
    host: getOsEnv('REDIS_HOST', 'localhost'),
    port: toNumber(getOsEnv('REDIS_PORT', '6379')),
    password: getOsEnv('REDIS_PASSWORD', ''),
    db: toNumber(getOsEnv('REDIS_DB', '0')),
  },
  grpc: {
    packages: {
      user: {
        protoPath: join(__dirname, '..', '..', './proto/user.proto'),
        endpoint: getOsEnv('GRPC_USER_ENDPOINT', 'localhost:50051'),
      },
    },
  },
  line: {
    channelId: getOsEnv('LINE_CHANNEL_ID'),
    channelSecret: getOsEnv('LINE_CHANNEL_SECRET'),
    channelAccessToken: getOsEnv('LINE_CHANNEL_SECRET'),
  },
});
