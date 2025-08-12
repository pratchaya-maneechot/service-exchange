import {
  CallOptions,
  ClientUnaryCall,
  Metadata,
  ServiceError,
  status,
} from '@grpc/grpc-js';
import { HttpStatus, Injectable } from '@nestjs/common';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';
import {
  BadRequestError,
  UnauthorizedError,
  ServiceUnavailableError,
  GatewayTimeoutError,
  InternalServerError,
  NotFoundError,
  ForbiddenError,
} from '../common/errors';
import { BaseError } from '../common/errors/base.error';

type GrpcCallback<T> = (error: ServiceError | null, response: T | null) => void;
type GrpcMethod<TRequest, TResponse> = (
  request: TRequest,
  metadataOrCallback?: Metadata | GrpcCallback<TResponse>,
  optionsOrCallback?: CallOptions | GrpcCallback<TResponse>,
  callback?: GrpcCallback<TResponse>,
) => ClientUnaryCall;

type CallerReturnType<TService> = {
  [K in keyof TService]: TService[K] extends GrpcMethod<
    infer TRequest,
    infer TResponse
  >
    ? (
        request: TRequest,
        metadata?: Metadata,
        options?: CallOptions,
      ) => Promise<TResponse>
    : TService[K];
};

@Injectable()
export abstract class GrpcClientService {
  constructor(
    @InjectPinoLogger()
    private readonly logger: PinoLogger,
  ) {
    logger.setContext(this.constructor.name);
  }

  private createClient<TRequest, TResponse>(
    method: GrpcMethod<TRequest, TResponse>,
    request: TRequest,
    metadata: Metadata,
    options?: CallOptions,
    name?: string,
  ): Promise<TResponse> {
    return new Promise<TResponse>((resolve, reject) => {
      const callback: GrpcCallback<TResponse> = (error, response) => {
        if (error) {
          this.logger.error(
            `Raw gRPC error for ${name || 'unknown method'}: ${error.message}`,
            {
              code: error.code,
              details: error.details,
              metadata: error.metadata?.toJSON(),
              stack: error.stack,
            },
          );
          const mappedError = this.mapError(error);
          reject(mappedError);
        } else if (response) {
          this.logger.debug(
            `gRPC ${name || 'unknown method'} completed successfully`,
          );
          resolve(response);
        } else {
          const errorMessage = `gRPC ${name || 'unknown method'} failed: empty response received from gRPC service`;
          this.logger.error(errorMessage);
          const emptyResponseError = new InternalServerError(errorMessage, {
            context: 'empty_grpc_response',
            serviceMethod: name,
          });
          reject(emptyResponseError);
        }
      };

      if (options) {
        method(request, metadata, options, callback);
      } else {
        method(request, metadata, callback);
      }
    });
  }

  protected caller<TService extends Record<string, any>>(
    service: TService,
  ): CallerReturnType<TService> {
    const promisified: Record<string, any> = {};

    for (const key of Object.getOwnPropertyNames(
      Object.getPrototypeOf(service),
    )) {
      const prop = service[key as keyof TService];
      if (typeof prop === 'function' && key !== 'constructor') {
        promisified[key] = (
          request: any,
          metadata?: Metadata,
          options?: CallOptions,
        ) => {
          return this.createClient(
            // eslint-disable-next-line @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-unsafe-call
            prop.bind(service),
            request,
            metadata ?? new Metadata(),
            options,
            key,
          );
        };
      } else {
        promisified[key] = prop;
      }
    }

    return promisified as CallerReturnType<TService>;
  }

  private mapError(grpcError: ServiceError): BaseError {
    const message = grpcError.details || grpcError.message;
    const grpcMetadata = grpcError.metadata?.toJSON();

    switch (grpcError.code) {
      case status.NOT_FOUND:
        return new NotFoundError(message, grpcMetadata, grpcError);
      case status.INVALID_ARGUMENT:
        return new BadRequestError(message, grpcMetadata, grpcError);
      case status.UNAUTHENTICATED:
        return new UnauthorizedError(message, grpcMetadata, grpcError);
      case status.PERMISSION_DENIED:
        return new ForbiddenError(message, grpcMetadata, grpcError);
      case status.UNAVAILABLE:
        return new ServiceUnavailableError(message, grpcMetadata, grpcError);
      case status.DEADLINE_EXCEEDED:
        return new GatewayTimeoutError(message, grpcMetadata, grpcError);
      case status.INTERNAL:
      case status.UNKNOWN:
        return new InternalServerError(message, grpcMetadata, grpcError);
      case status.CANCELLED:
        return new BaseError(
          'Request Cancelled',
          HttpStatus.CONFLICT,
          'REQUEST_CANCELLED',
          grpcMetadata,
          true,
          grpcError,
        );
      case status.RESOURCE_EXHAUSTED:
        return new BaseError(
          'Too Many Requests / Resource Exhausted',
          HttpStatus.TOO_MANY_REQUESTS,
          'RESOURCE_EXHAUSTED',
          grpcMetadata,
          true,
          grpcError,
        );
      default:
        return new InternalServerError(
          `Unhandled gRPC error: ${message} (Code: ${grpcError.code})`,
          {
            grpcCode: grpcError.code,
            originalDetails: grpcError.details,
            grpcMetadata,
          },
          grpcError,
        );
    }
  }
}
