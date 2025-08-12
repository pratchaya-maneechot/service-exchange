import { getOsEnv, toNumber } from 'src/common/utils/env.util';
import { registerAs } from '@nestjs/config';
import { ConfigType } from 'src/common/types/config.type';

export default registerAs<ConfigType['redis']>('redis', () => ({
  host: getOsEnv('REDIS_HOST', 'localhost'),
  port: toNumber(getOsEnv('REDIS_PORT', '6379')),
  password: getOsEnv('REDIS_PASSWORD', ''),
  db: toNumber(getOsEnv('REDIS_DB', '0')),
}));
