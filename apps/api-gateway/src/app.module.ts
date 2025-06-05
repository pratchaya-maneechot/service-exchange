import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { ThrottlerModule } from '@nestjs/throttler';
import { TerminusModule } from '@nestjs/terminus';
import { CoreModule } from './core/core.module';
import { env } from './config/env.config';
import { GraphQLModule } from './graphql/graphql.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      load: [env],
      envFilePath: ['.env.local', '.env'],
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
