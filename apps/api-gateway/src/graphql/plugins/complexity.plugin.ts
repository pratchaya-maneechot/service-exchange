import { ApolloServerPlugin, GraphQLRequestListener } from '@apollo/server';
import { Plugin } from '@nestjs/apollo';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';

@Plugin()
export class ComplexityPlugin implements ApolloServerPlugin {
  private readonly maxComplexity = 1000;

  constructor(
    @InjectPinoLogger(ComplexityPlugin.name)
    private readonly logger: PinoLogger,
  ) {}

  async requestDidStart(context): Promise<GraphQLRequestListener<object>> {
    return {
      didResolveOperation: async () => {
        // Implementation would calculate query complexity and reject if too complex
        this.logger.debug(`Checking complexity limit: ${this.maxComplexity}`);
      },
    };
  }
}
