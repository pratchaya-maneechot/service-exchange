import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { ThrottlerModule } from '@nestjs/throttler';
import { TerminusModule } from '@nestjs/terminus';
import { CoreModule } from './core/core.module';
import { configs } from './config';
import { GraphQLModule } from './graphql/graphql.module';
import { AppLoggerModule } from './observability/logger';
import { PrometheusModule } from '@willsoto/nestjs-prometheus';

@Module({
  imports: [
    AppLoggerModule,
    ConfigModule.forRoot({
      isGlobal: true,
      load: configs,
      envFilePath: ['.env.local', '.env'],
    }),
    PrometheusModule.register({
      path: '/metrics',
    }),
    GraphQLModule,
    ThrottlerModule.forRoot({
      // 1000 requests per minute
      throttlers: [
        {
          ttl: 60 * 1000,
          limit: 1000,
        },
      ],
    }),
    TerminusModule,
    CoreModule,
  ],
})
export class AppModule {}
