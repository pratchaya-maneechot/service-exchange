import { ApolloServerPlugin, GraphQLRequestListener } from '@apollo/server';
import { Plugin } from '@nestjs/apollo';
import type { GraphQLRequestContext } from '@apollo/server';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';
import { stdSerializers } from 'pino-http';
import { IAppContext } from 'src/common/types/context.type';
import { ConfigService } from '@nestjs/config';
import { ConfigType } from 'src/common/types/config.type';

interface RequestTimings {
  requestStart: number;
  parsingStart?: number;
  validationStart?: number;
  executionStart?: number;
  resolverTimings: Record<string, number>;
}

@Plugin()
export class GraphQLLoggingPlugin implements ApolloServerPlugin {
  constructor(
    @InjectPinoLogger(GraphQLLoggingPlugin.name)
    private readonly logger: PinoLogger,
    private readonly config: ConfigService<ConfigType>,
  ) {}

  async requestDidStart(
    requestContext: GraphQLRequestContext<IAppContext>,
  ): Promise<GraphQLRequestListener<IAppContext>> {
    const {
      request,
      contextValue,
      operation,
      operationName: oPn,
    } = requestContext;
    const operationName = request.operationName || oPn;
    const operationType = operation?.operation ?? 'query';
    const userId = contextValue.req.user?.id || 'anonymous';
    const timings: RequestTimings = {
      requestStart: performance.now(),
      resolverTimings: {},
    };

    this.logger.info(
      {
        operationName,
        operationType,
        userId,
        variables: this.sanitizeVariables(request.variables),
      },
      `GraphQL ${operationName} Started`,
    );

    const _logger = this.logger;
    const slowThreshold = this.config.get('app.slowThreshold', {
      infer: true,
    });

    return {
      async parsingDidStart() {
        timings.parsingStart = performance.now();
        return async (err) => {
          if (err) {
            _logger.warn({ err, operationName }, 'Query Parsing Failed');
          }
        };
      },

      async validationDidStart() {
        timings.validationStart = performance.now();
        return async (err) => {
          if (err) {
            _logger.warn({ err, operationName }, 'Query Validation Failed');
          }
        };
      },

      async executionDidStart() {
        timings.executionStart = performance.now();
        return {
          willResolveField({ info, args }) {
            const start = performance.now();
            const fieldKey = `${info.parentType.name}.${info.fieldName}`;

            return (error) => {
              const duration = performance.now() - start;
              timings.resolverTimings[fieldKey] =
                (timings.resolverTimings[fieldKey] || 0) + duration;

              const isSlow = slowThreshold && duration > slowThreshold;

              if (isSlow || error) {
                const logLevel = error ? 'warn' : 'debug';
                _logger[logLevel](
                  {
                    fieldName: info.fieldName,
                    parentType: info.parentType.name,
                    operationName,
                    duration: `${duration.toFixed(2)}ms`,
                    args,
                    ...(error && { error: stdSerializers?.err(error) }),
                  },
                  `Field Resolution ${fieldKey}`,
                );
              }
            };
          },
        };
      },

      async didEncounterErrors({ errors }) {
        errors?.forEach((error) => {
          _logger.error(
            {
              operationName,
              userId,
              message: error.message,
              stack: error.stack,
              path: error.path,
              code: error.extensions?.code,
            },
            `GraphQL ${operationName} Error`,
          );
        });
      },

      async willSendResponse() {
        const duration = performance.now() - timings.requestStart;

        _logger.info(
          {
            operationName,
            userId,
            duration: `${duration.toFixed(2)}ms`,
            resolverTimings: Object.fromEntries(
              Object.entries(timings.resolverTimings).map(([key, value]) => [
                key,
                `${value.toFixed(2)}ms`,
              ]),
            ),
          },
          `GraphQL ${operationName} Completed`,
        );
      },
    };
  }

  private sanitizeVariables(
    variables: Record<string, unknown> | undefined,
  ): Record<string, unknown> {
    if (!variables) return {};
    const sanitized = { ...variables };
    for (const key in sanitized) {
      if (/password|token/i.test(key)) {
        sanitized[key] = '[Redacted]';
      }
    }
    return sanitized;
  }
}
