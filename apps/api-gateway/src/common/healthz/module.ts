import { Module } from '@nestjs/common';
import { TerminusModule } from '@nestjs/terminus';
import { HealthService } from './service';
import { HealthController, ReadyController } from './controller';

@Module({
  imports: [TerminusModule],
  controllers: [HealthController, ReadyController],
  providers: [HealthService],
})
export class HealthModule {}
