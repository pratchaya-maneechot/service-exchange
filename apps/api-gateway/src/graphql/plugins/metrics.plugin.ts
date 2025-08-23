import { ApolloServerPlugin, GraphQLRequestListener } from '@apollo/server';
import { Plugin } from '@nestjs/apollo';
import type { GraphQLRequestContext } from '@apollo/server';
import { IAppContext } from 'src/common/types/context.type';
import { Histogram, Counter } from 'prom-client';

const graphqlRequestDurationSeconds = new Histogram({
  name: 'graphql_request_duration_seconds',
  help: 'Duration of GraphQL requests in seconds',
  labelNames: ['operation_name', 'operation_type', 'status'],
  buckets: [0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10],
});

const graphqlRequestsTotal = new Counter({
  name: 'graphql_requests_total',
  help: 'Total number of GraphQL requests',
  labelNames: ['operation_name', 'operation_type', 'status'],
});

const graphqlErrorsTotal = new Counter({
  name: 'graphql_errors_total',
  help: 'Total number of GraphQL errors by error type',
  labelNames: ['operation_name', 'operation_type', 'error_code', 'error_type'],
});

const graphqlResolverDurationSeconds = new Histogram({
  name: 'graphql_resolver_duration_seconds',
  help: 'Duration of GraphQL field resolvers in seconds',
  labelNames: ['parent_type', 'field_name'],
  buckets: [0.0001, 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5],
});

const graphqlResolverErrorsTotal = new Counter({
  name: 'graphql_resolver_errors_total',
  help: 'Total number of GraphQL resolver errors',
  labelNames: ['parent_type', 'field_name', 'error_type'],
});

@Plugin()
export class GraphQLMetricsPlugin implements ApolloServerPlugin {
  async requestDidStart(
    requestContext: GraphQLRequestContext<IAppContext>,
  ): Promise<GraphQLRequestListener<IAppContext>> {
    const { request, operation, operationName: oPn } = requestContext;

    // Handle batching or lack of operation name
    const operationName = request.operationName || oPn || 'anonymous';
    const operationType = operation?.operation ?? 'query';

    const requestStart = performance.now();

    graphqlRequestsTotal.labels(operationName, operationType, 'started').inc();

    return {
      async executionDidStart() {
        return {
          willResolveField({ info }) {
            const start = performance.now();

            return (error) => {
              const duration = performance.now() - start;

              graphqlResolverDurationSeconds
                .labels(info.parentType.name, info.fieldName)
                .observe(duration / 1000);

              if (error) {
                graphqlResolverErrorsTotal
                  .labels(
                    info.parentType.name,
                    info.fieldName,
                    error.name || 'UnknownError',
                  )
                  .inc();
              }
            };
          },
        };
      },
      async didEncounterErrors({ errors, operationName: opName, operation }) {
        const opNameLabel = opName || operation?.name?.value || 'anonymous';
        const opTypeLabel = operation?.operation ?? 'query';

        errors?.forEach((error) => {
          graphqlErrorsTotal
            .labels(
              opNameLabel,
              opTypeLabel,
              (error.extensions?.code as string) || 'UNKNOWN_CODE',
              error.name || 'UnknownError',
            )
            .inc();
        });
      },

      async willSendResponse({ response, operationName: opName, operation }) {
        const duration = performance.now() - requestStart;
        const hasErrors =
          response.body.kind === 'single' && response.body.singleResult.errors;
        const status = hasErrors ? 'error' : 'success';

        const opNameLabel = opName || operation?.name?.value || 'anonymous';
        const opTypeLabel = operation?.operation ?? 'query';

        graphqlRequestDurationSeconds
          .labels(opNameLabel, opTypeLabel, status)
          .observe(duration / 1000);

        graphqlRequestsTotal.labels(opNameLabel, opTypeLabel, status).inc();
      },
    };
  }
}
