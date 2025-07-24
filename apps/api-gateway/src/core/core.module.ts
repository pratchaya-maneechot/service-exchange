import { Module } from '@nestjs/common';
import { UserModule } from './user/user.module';
import { AuthModule } from './auth/auth.module';
import { HealthModule } from './health/health.module';
import { LineModule } from './line/line.module';
@Module({
  imports: [UserModule, HealthModule, AuthModule, LineModule],
})
export class CoreModule {}
