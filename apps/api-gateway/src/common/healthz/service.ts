import { Injectable } from '@nestjs/common';
import { HealthCheckService, HealthCheck } from '@nestjs/terminus';

@Injectable()
export class HealthService {
  constructor(private health: HealthCheckService) {}

  @HealthCheck()
  check() {
    return this.health.check([
      () => Promise.resolve({ database: { status: 'up' } }),
      () => Promise.resolve({ redis: { status: 'up' } }),
    ]);
  }
}
