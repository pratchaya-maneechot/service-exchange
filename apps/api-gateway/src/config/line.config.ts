import { getOsEnv } from 'src/common/utils/env.util';
import { registerAs } from '@nestjs/config';
import { ConfigType } from 'src/common/types/config.type';

export default registerAs<ConfigType['line']>('line', () => ({
  channelId: getOsEnv('LINE_CHANNEL_ID'),
  channelSecret: getOsEnv('LINE_CHANNEL_SECRET'),
  channelAccessToken: getOsEnv('LINE_CHANNEL_ACCESS_TOKEN'),
}));
