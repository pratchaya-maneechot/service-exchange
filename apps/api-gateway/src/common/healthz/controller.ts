import { HealthService } from './service';
import { HealthStatus } from './types';
import { Controller, Get } from '@nestjs/common';

@Controller('healthz')
export class HealthController {
  constructor(private healthService: HealthService) {}

  @Get()
  async index(): Promise<HealthStatus> {
    const result = await this.healthService.check();
    return {
      status: result.status,
      info: result.info,
      error: result.error,
      details: result.details,
    };
  }
}

@Controller('ready')
export class ReadyController {
  @Get()
  index(): HealthStatus {
    return {
      status: 'ok',
      info: {},
      error: undefined,
      details: {},
    };
  }
}
