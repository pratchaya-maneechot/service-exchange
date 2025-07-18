/* eslint-disable @typescript-eslint/no-redundant-type-constituents */
// src/errors/index.ts
import { HttpStatus } from '@nestjs/common';
import { BaseError } from './base.error';

export class BadRequestError extends BaseError {
  constructor(
    message: string = 'Bad Request',
    details?: any,
    originalError?: Error | unknown,
  ) {
    super(
      message,
      HttpStatus.BAD_REQUEST,
      'BAD_REQUEST',
      details,
      true,
      originalError,
    );
    this.name = 'BadRequestError';
  }
}

export class ValidationError extends BadRequestError {
  constructor(
    message: string = 'Validation Failed',
    details?: any,
    originalError?: Error | unknown,
  ) {
    super(message, details, originalError);
    this.name = 'ValidationError';
    this.errorCode = 'VALIDATION_FAILED';
  }
}

export class NotFoundError extends BaseError {
  constructor(
    message: string = 'Resource Not Found',
    details?: any,
    originalError?: Error | unknown,
  ) {
    super(
      message,
      HttpStatus.NOT_FOUND,
      'NOT_FOUND',
      details,
      true,
      originalError,
    );
    this.name = 'NotFoundError';
  }
}

export class UnauthorizedError extends BaseError {
  constructor(
    message: string = 'Unauthorized',
    details?: any,
    originalError?: Error | unknown,
  ) {
    super(
      message,
      HttpStatus.UNAUTHORIZED,
      'UNAUTHORIZED',
      details,
      true,
      originalError,
    );
    this.name = 'UnauthorizedError';
  }
}

export class ForbiddenError extends BaseError {
  constructor(
    message: string = 'Forbidden',
    details?: any,
    originalError?: Error | unknown,
  ) {
    super(
      message,
      HttpStatus.FORBIDDEN,
      'FORBIDDEN',
      details,
      true,
      originalError,
    );
    this.name = 'ForbiddenError';
  }
}

export class ServiceUnavailableError extends BaseError {
  constructor(
    message: string = 'Service Unavailable',
    details?: any,
    originalError?: Error | unknown,
  ) {
    super(
      message,
      HttpStatus.SERVICE_UNAVAILABLE,
      'SERVICE_UNAVAILABLE',
      details,
      true,
      originalError,
    );
    this.name = 'ServiceUnavailableError';
  }
}

export class GatewayTimeoutError extends BaseError {
  constructor(
    message: string = 'Gateway Timeout',
    details?: any,
    originalError?: Error | unknown,
  ) {
    super(
      message,
      HttpStatus.GATEWAY_TIMEOUT,
      'GATEWAY_TIMEOUT',
      details,
      true,
      originalError,
    );
    this.name = 'GatewayTimeoutError';
  }
}

export class InternalServerError extends BaseError {
  constructor(
    message: string = 'Internal Server Error',
    details?: any,
    originalError?: Error | unknown,
  ) {
    super(
      message,
      HttpStatus.INTERNAL_SERVER_ERROR,
      'INTERNAL_SERVER_ERROR',
      details,
      false,
      originalError,
    );
    this.name = 'InternalServerError';
  }
}
