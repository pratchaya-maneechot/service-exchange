/* eslint-disable @typescript-eslint/no-redundant-type-constituents */
import {
  ExceptionFilter,
  Catch,
  ArgumentsHost,
  HttpException,
  HttpStatus,
} from '@nestjs/common';
import { Request, Response } from 'express';
import { GraphQLError } from 'graphql';
import { PinoLogger } from 'nestjs-pino';
import { BaseError } from '../errors/base.error';
import { getOsEnv } from '../utils/env.util';

@Catch()
export class AllExceptionsFilter implements ExceptionFilter {
  constructor(private readonly logger: PinoLogger) {}

  catch(exception: unknown, host: ArgumentsHost) {
    const isGraphql = host.getType<'graphql'>() === 'graphql';

    let statusCode: HttpStatus;
    let errorCode: string;
    let message: string;
    let details: any;
    let isOperational: boolean = false;
    let originalError: unknown | undefined;

    if (exception instanceof HttpException) {
      statusCode = exception.getStatus();
      const exceptionResponse = exception.getResponse();

      if (typeof exceptionResponse === 'string') {
        message = exceptionResponse;
        errorCode = exception.constructor.name
          .replace(/Error$/, '')
          .toUpperCase();
      } else {
        message = (exceptionResponse as any).message || exception.message;
        errorCode =
          (exceptionResponse as any).error ||
          exception.constructor.name.replace(/Error$/, '').toUpperCase();
        details =
          (exceptionResponse as any).details ||
          (exceptionResponse as any).errors ||
          undefined;
      }
      isOperational = true;
      originalError = exception;
      this.logger.warn(
        `[HttpException] ${message} - Path: ${isGraphql ? 'GraphQL' : host.switchToHttp().getRequest().url}`,
        { statusCode, errorCode, details, stack: exception.stack },
      );
    } else if (BaseError.isBaseError(exception)) {
      statusCode = exception.statusCode;
      errorCode = exception.errorCode;
      message = exception.message;
      details = exception.details;
      isOperational = exception.isOperational;
      originalError = exception.originalError || exception;
    } else {
      statusCode = HttpStatus.INTERNAL_SERVER_ERROR;
      errorCode = 'UNHANDLED_ERROR';
      message = 'An unexpected internal server error occurred.';
      details =
        getOsEnv('NODE_ENV') === 'production'
          ? undefined
          : (exception as Error)?.message;
      isOperational = false;
      originalError = exception;

      this.logger.error(
        `[UNHANDLED EXCEPTION] ${message} - Path: ${isGraphql ? 'GraphQL' : host.switchToHttp().getRequest().url}`,
        {
          stack: (exception as Error)?.stack || 'No stack trace available',
          originalException: exception,
        },
      );
    }

    if (isGraphql) {
      const graphQLError = new GraphQLError(message, {
        extensions: {
          code: errorCode,
          statusCode: statusCode,
          details: details,
          timestamp: new Date().toISOString(),
          originalError:
            originalError instanceof Error
              ? originalError.message
              : String(originalError),
          stacktrace:
            getOsEnv('NODE_ENV') === 'production'
              ? undefined
              : (originalError as Error)?.stack?.split('\n'),
        },
      });

      throw graphQLError;
    } else {
      const ctx = host.switchToHttp();
      const response = ctx.getResponse<Response>();
      const request = ctx.getRequest<Request>();

      response.status(statusCode).json({
        statusCode: statusCode,
        errorCode: errorCode,
        message: message,
        details: details,
        timestamp: new Date().toISOString(),
        path: request.url,
      });
    }
  }
}
