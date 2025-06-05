import { Injectable } from '@nestjs/common';
import { GrpcClientRepository } from './grpc-client.repository';
import { GrpcClientFactory } from './grpc-client.factory';
import { UserServiceClient } from './types/generated/user/UserService';

@Injectable()
export class UserClientService {
  private readonly package = 'user';
  constructor(private readonly repository: GrpcClientRepository) {}

  get user() {
    return this.repository.getClient<UserServiceClient>(
      GrpcClientFactory.key(this.package, 'UserService'),
    );
  }
}
