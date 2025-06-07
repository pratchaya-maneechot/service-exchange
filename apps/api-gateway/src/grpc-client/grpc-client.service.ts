import {
  CallOptions,
  ClientUnaryCall,
  Metadata,
  ServiceError,
} from '@grpc/grpc-js';
import { Status } from '@grpc/grpc-js/build/src/constants';
import { Injectable } from '@nestjs/common';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';

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

  protected prepareMetadata(userMetadata?: Metadata): Metadata {
    const meta = new Metadata();
    const requestId =
      (this.logger.logger.bindings().req.id as string) ?? 'unknown';

    meta.set('x-request-id', requestId);
    meta.set('x-service', this.constructor.name);
    meta.set('x-timestamp', new Date().toISOString());

    if (userMetadata) {
      Object.entries(userMetadata.getMap()).forEach(([key, value]) => {
        if (Array.isArray(value)) {
          value.forEach((v) => meta.add(key, v.toString()));
        } else {
          meta.set(key, value);
        }
      });
    }

    return meta;
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
            `gRPC ${name} failed: ${error.message}`,
            error.stack,
          );
          reject(error);
        } else if (response) {
          this.logger.debug(`gRPC ${name} completed successfully`);
          resolve(response);
        } else {
          this.logger.error(
            `gRPC ${name} failed: empty response received from gRPC service`,
          );
          const err = new Error(
            'empty response received from gRPC service',
          ) as ServiceError;
          err.code = Status.INTERNAL;
          err.details = 'The service returned a null or undefined response';
          reject(err);
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
          const meta = this.prepareMetadata(metadata);
          return this.createClient(
            // eslint-disable-next-line @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-unsafe-call
            prop.bind(service),
            request,
            meta,
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
}
