import { ApolloServerPlugin, GraphQLRequestListener } from '@apollo/server';
import { Plugin } from '@nestjs/apollo';
import type { GraphQLRequestContext } from '@apollo/server';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';
import { stdSerializers } from 'pino-http';
import { IAppContext } from 'src/common/types/context.type';
import { ConfigService } from '@nestjs/config';
import { EnvConfig } from 'src/config/config.type';

interface RequestTimings {
  requestStart: number;
  parsingStart?: number;
  validationStart?: number;
  executionStart?: number;
  resolverTimings: Record<string, number>;
}

const LOG_SLOW_RESOLVERS_THRESHOLD = 1500; // ms

@Plugin()
export class GraphQLLoggerPlugin implements ApolloServerPlugin {
  constructor(
    @InjectPinoLogger(GraphQLLoggerPlugin.name)
    private readonly logger: PinoLogger,
    private readonly config: ConfigService<EnvConfig>,
  ) {}

  async requestDidStart(
    requestContext: GraphQLRequestContext<IAppContext>,
  ): Promise<GraphQLRequestListener<IAppContext>> {
    const { request, contextValue } = requestContext;
    const operationName = request.operationName || 'anonymous';
    const userId = contextValue.req.user?.id || 'anonymous';
    const timings: RequestTimings = {
      requestStart: performance.now(),
      resolverTimings: {},
    };

    this.logger.info(
      {
        operationName,
        userId,
        variables: this.sanitizeVariables(request.variables),
      },
      `GraphQL ${operationName} Started`,
    );

    const _logger = this.logger;

    return {
      async parsingDidStart() {
        timings.parsingStart = performance.now();
        return async (err) => {
          const duration = err
            ? undefined
            : performance.now() - (timings.parsingStart || 0);
          if (err) {
            _logger.warn({ err, operationName }, 'Query Parsing Failed');
          } else if (duration) {
            if (duration > 50) {
              _logger.debug(
                { operationName, duration: `${duration.toFixed(2)}ms` },
                'Query Parsing Completed',
              );
            }
          }
        };
      },

      async validationDidStart() {
        timings.validationStart = performance.now();
        return async (err) => {
          const duration = err
            ? undefined
            : performance.now() - (timings.validationStart || 0);
          if (err) {
            _logger.warn({ err, operationName }, 'Query Validation Failed');
          } else if (duration) {
            if (duration > 50) {
              _logger.debug(
                { operationName, duration: `${duration.toFixed(2)}ms` },
                'Query Validation Completed',
              );
            }
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

              // Log only slow resolvers or errors to reduce log volume
              if (duration > LOG_SLOW_RESOLVERS_THRESHOLD || error) {
                const logLevel =
                  error || duration > LOG_SLOW_RESOLVERS_THRESHOLD
                    ? 'warn'
                    : 'debug';
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
