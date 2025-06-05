import { Plugin } from '@nestjs/apollo';
import {
  ApolloServerPlugin,
  GraphQLRequestContext,
  GraphQLRequestListener,
} from '@apollo/server';
import { Logger } from '@nestjs/common';

@Plugin()
export class MetricsPlugin implements ApolloServerPlugin {
  private readonly logger = new Logger(MetricsPlugin.name);

  async requestDidStart(
    ctx: GraphQLRequestContext<object>,
  ): Promise<GraphQLRequestListener<object>> {
    const startTime = Date.now();
    const _logger = this.logger;
    return {
      didResolveOperation: async () => {
        if (ctx.operationName) {
          _logger.log(`Increment Query Count: ${ctx.operationName}`);
        }
      },
      willSendResponse: async () => {
        const duration = Date.now() - startTime;
        if (ctx.operationName) {
          _logger.log(`Record Query Duration: ${ctx.operationName}`);
        }
      },
      didEncounterErrors: async (err) => {
        err.errors?.forEach((error) => {
          if (ctx.operationName) {
            _logger.log(`Record Error Count: ${ctx.operationName}`);
          }
        });
      },
    };
  }
}
