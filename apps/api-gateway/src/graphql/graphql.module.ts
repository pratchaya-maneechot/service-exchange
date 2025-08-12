import { ApolloDriver, ApolloDriverConfig } from '@nestjs/apollo';
import { Module } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { GraphQLModule as NestGraphQLModule } from '@nestjs/graphql';
import { join } from 'path';
import { IAppContext } from 'src/common/types/context.type';
import { ConfigType } from 'src/common/types/config.type';
import { GraphQLLoggerPlugin } from './plugins/logger.plugin';

@Module({
  providers: [GraphQLLoggerPlugin],
  imports: [
    NestGraphQLModule.forRootAsync<ApolloDriverConfig>({
      driver: ApolloDriver,
      useFactory: (config: ConfigService<ConfigType>) => {
        const isProd =
          config.get('app.nodeEnv', { infer: true }) === 'production';
        return {
          autoSchemaFile: join(process.cwd(), 'src/graphql/generated.gql'),
          sortSchema: true,
          playground: !isProd,
          introspection: !isProd,
          debug: !isProd,
          context: ({ req }): IAppContext => ({ req }),
          formatError: (error, err) => {
            return {
              message: error.message,
              code: error.extensions?.code,
              timestamp: new Date().toISOString(),
              path: error.path,
              ...(isProd
                ? {}
                : {
                    locations: error.locations,
                    originalError: error.extensions?.originalError,
                    stack: (error as any).stack,
                  }),
            };
          },
        };
      },
      inject: [ConfigService],
    }),
  ],
})
export class GraphQLModule {}
