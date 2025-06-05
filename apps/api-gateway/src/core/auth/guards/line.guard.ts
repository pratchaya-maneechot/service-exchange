import {
  Injectable,
  ExecutionContext,
  UnauthorizedException,
} from '@nestjs/common';
import { AuthGuard } from '@nestjs/passport';
import { GqlExecutionContext } from '@nestjs/graphql';
import { Request } from 'express';
import { IAppContext } from 'src/common/types/context.type';
import { Identity } from '../types/auth.type';

@Injectable()
export class LineAuthGuard extends AuthGuard('line') {
  constructor() {
    super();
  }

  getRequest(context: ExecutionContext): Request {
    const ctx = GqlExecutionContext.create(context);
    const gqlReq = ctx.getContext<IAppContext>().req;
    return gqlReq;
  }

  canActivate(context: ExecutionContext) {
    return super.canActivate(context);
  }

  handleRequest<T = Identity>(err: any, user: any): T {
    if (err || !user) {
      throw (
        err ||
        new UnauthorizedException(
          'Unauthorized: Invalid or missing LINE ID Token',
        )
      );
    }
    return user as T;
  }
}
