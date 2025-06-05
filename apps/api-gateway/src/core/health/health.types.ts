import { ObjectType, Field } from '@nestjs/graphql';
import { HealthCheckStatus, HealthIndicatorResult } from '@nestjs/terminus';
import { JSONScalar } from 'src/graphql/scalars/json.scalar';

@ObjectType()
export class HealthStatus {
  @Field()
  status: HealthCheckStatus;
  @Field(() => JSONScalar, { nullable: true })
  info?: HealthIndicatorResult;
  @Field(() => JSONScalar, { nullable: true })
  error?: HealthIndicatorResult;
  @Field(() => JSONScalar)
  details: HealthIndicatorResult;
}
