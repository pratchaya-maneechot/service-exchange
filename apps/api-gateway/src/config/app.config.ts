import { toNumber, getOsEnv } from 'src/common/utils/env.util';
import { registerAs } from '@nestjs/config';
import { ConfigType } from 'src/common/types/config.type';

export default registerAs<ConfigType['app']>('app', () => ({
  port: toNumber(getOsEnv('PORT', '8080')),
  host: getOsEnv('HOST', '0.0.0.0'),
  nodeEnv: getOsEnv('NODE_ENV', 'development') as
    | 'production'
    | 'development'
    | 'test',
  corsOrigins: getOsEnv('CORS_ORIGINS', '*')?.split(',').filter(Boolean),
  appName: getOsEnv('APP_NAME', 'api-gateway'),
  appVersion: getOsEnv('APP_VERSION', 'vLocal-dev'),
}));
