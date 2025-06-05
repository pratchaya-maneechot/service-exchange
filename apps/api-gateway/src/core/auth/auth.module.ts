import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { LineStrategy } from './strategies/line.strategy';
import { LineAuthGuard } from './guards/line.guard';
import { LineModule } from 'src/core/line/line.module';
import { PassportModule } from '@nestjs/passport';

@Module({
  imports: [ConfigModule, LineModule, PassportModule],
  providers: [LineStrategy, LineAuthGuard],
  exports: [LineAuthGuard],
})
export class AuthModule {}
