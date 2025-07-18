/* eslint-disable @typescript-eslint/no-redundant-type-constituents */
import { HttpStatus } from '@nestjs/common';

/**
 * Custom base error class for the application.
 * All application-specific errors should extend this class.
 */
export class BaseError extends Error {
  public name: string;
  public statusCode: HttpStatus;
  public errorCode: string;
  public details?: any;
  public isOperational: boolean;

  // Store the original underlying error for debugging (e.g., grpc.ServiceError)
  public originalError?: Error | unknown;

  constructor(
    message: string,
    statusCode: HttpStatus = HttpStatus.INTERNAL_SERVER_ERROR,
    errorCode: string = 'GENERIC_ERROR',
    details?: any,
    isOperational: boolean = true,
    originalError?: Error | unknown,
  ) {
    super(message);
    this.name = this.constructor.name;
    this.statusCode = statusCode;
    this.errorCode = errorCode;
    this.details = details;
    this.isOperational = isOperational;
    this.originalError = originalError;

    Error.captureStackTrace(this, this.constructor);
  }

  static isBaseError(error: unknown): error is BaseError {
    return error instanceof BaseError;
  }
}
