import { Module } from '@nestjs/common';
import { UserResolver } from './user.resolver';
import { UserService } from './user.service';
import { GrpcClientModule } from 'src/grpc-client/grpc-client.module';
import { ScalarsModule } from 'src/graphql/scalars/scalars.module';

@Module({
  imports: [GrpcClientModule, ScalarsModule],
  providers: [UserService, UserResolver],
  exports: [UserService, UserResolver],
})
export class UserModule {}
