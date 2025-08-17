import { Module } from '@nestjs/common';
import { UserModule } from './user/user.module';
import { AuthModule } from './auth/auth.module';
import { HealthModule } from '../common/healthz/module';
import { LineModule } from './line/line.module';
@Module({
  imports: [UserModule, HealthModule, AuthModule, LineModule],
})
export class CoreModule {}
