import { Plugin } from '@nestjs/apollo';
import {
  ApolloServerPlugin,
  GraphQLRequestContext,
  GraphQLRequestListener,
} from '@apollo/server';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';

@Plugin()
export class MetricsPlugin implements ApolloServerPlugin {
  constructor(
    @InjectPinoLogger(MetricsPlugin.name)
    private readonly logger: PinoLogger,
  ) {}
  async requestDidStart(
    ctx: GraphQLRequestContext<object>,
  ): Promise<GraphQLRequestListener<object>> {
    const startTime = Date.now();
    const _logger = this.logger;
    return {
      didResolveOperation: async () => {
        if (ctx.operationName) {
          _logger.info(`Increment Query Count: ${ctx.operationName}`);
        }
      },
      willSendResponse: async () => {
        const duration = Date.now() - startTime;
        if (ctx.operationName) {
          _logger.info(`Record Query Duration: ${ctx.operationName}`);
        }
      },
      didEncounterErrors: async (err) => {
        err.errors?.forEach((error) => {
          if (ctx.operationName) {
            _logger.info(`Record Error Count: ${ctx.operationName}`);
          }
        });
      },
    };
  }
}
