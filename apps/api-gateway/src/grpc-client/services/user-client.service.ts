import { Injectable } from '@nestjs/common';
import { GrpcClientRepository } from '../grpc-client.repository';
import { GrpcClientFactory } from '../grpc-client.factory';
import { UserServiceClient } from '../types/generated/user/v1/UserService';
import { GrpcClientService } from '../grpc-client.service';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';

@Injectable()
export class UserClientService extends GrpcClientService {
  constructor(
    @InjectPinoLogger()
    logger: PinoLogger,
    private readonly repository: GrpcClientRepository,
  ) {
    super(logger);
  }

  get raw() {
    return {
      userService: this.repository.getClient<UserServiceClient>(
        GrpcClientFactory.key('user', 'UserService'),
      ),
    };
  }

  get userService() {
    return this.caller(this.raw.userService);
  }
}
