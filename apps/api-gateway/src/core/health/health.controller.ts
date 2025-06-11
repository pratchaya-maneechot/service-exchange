import { HealthService } from './health.service';
import { HealthStatus } from './health.types';
import { Controller, Get } from '@nestjs/common';

@Controller('')
export class HealthController {
  constructor(private healthService: HealthService) {}

  @Get('healthz')
  async health(): Promise<HealthStatus> {
    const result = await this.healthService.check();
    return {
      status: result.status,
      info: result.info,
      error: result.error,
      details: result.details,
    };
  }
  @Get('ready')
  async ready(): Promise<HealthStatus> {
    const result = await this.healthService.check();
    return {
      status: result.status,
      info: result.info,
      error: result.error,
      details: result.details,
    };
  }
}
