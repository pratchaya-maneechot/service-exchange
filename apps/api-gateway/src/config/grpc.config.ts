import { getOsEnv, getOsEnvOptional } from 'src/common/utils/env.util';
import { registerAs } from '@nestjs/config';
import { ConfigType } from 'src/common/types/config.type';
import { join } from 'path';

export default registerAs<ConfigType['grpc']>('grpc', () => ({
  packages: {
    user: {
      protoPath: join(__dirname, '..', '..', './proto/user.proto'),
      endpoint: getOsEnv('GRPC_USER_ENDPOINT', 'localhost:50051'),
      version: getOsEnvOptional('GRPC_USER_VERSION'),
    },
  },
}));
