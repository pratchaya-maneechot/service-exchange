import { Module } from '@nestjs/common';
import { UserClientService } from './services/user-client.service';
import { GrpcClientFactory } from './grpc-client.factory';
import { GrpcClientRepository } from './grpc-client.repository';
@Module({
  providers: [GrpcClientRepository, GrpcClientFactory, UserClientService],
  exports: [UserClientService],
})
export class GrpcClientModule {}
