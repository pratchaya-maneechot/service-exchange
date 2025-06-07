import { ApolloServerPlugin, GraphQLRequestListener } from '@apollo/server';
import { Plugin } from '@nestjs/apollo';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';

@Plugin()
export class QueryDepthPlugin implements ApolloServerPlugin {
  private readonly maxDepth = 10;
  constructor(
    @InjectPinoLogger(QueryDepthPlugin.name)
    private readonly logger: PinoLogger,
  ) {}

  async requestDidStart(context): Promise<GraphQLRequestListener<object>> {
    return {
      didResolveOperation: async () => {
        // Implementation would calculate query depth and reject if too deep
        this.logger.debug(`Checking query depth limit: ${this.maxDepth}`);
      },
    };
  }
}
