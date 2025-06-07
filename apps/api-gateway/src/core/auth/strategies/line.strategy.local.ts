import { VerifyIDToken } from '@line/bot-sdk';
import { Injectable } from '@nestjs/common';
import { PassportStrategy } from '@nestjs/passport';
import { Request } from 'express';
import { Strategy } from 'passport-local';
import { catchError, EMPTY, from, tap } from 'rxjs';
import { Identity } from '../types/auth.type';
import { LineService } from 'src/core/line/line.service';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';

@Injectable()
export class LineStrategy extends PassportStrategy(Strategy, 'line') {
  constructor(
    @InjectPinoLogger(LineStrategy.name)
    private readonly logger: PinoLogger,
    private readonly lineService: LineService,
  ) {
    super();
  }

  authenticate(req: Request) {
    const idToken = req.headers.authorization?.replace('Bearer ', '');
    if (idToken) {
      from(this.lineService.verifyIdToken(idToken))
        .pipe(
          tap((user) => this.success(user)),
          catchError((error) => {
            this.logger.error(`Authentication failed: ${error.message}`);
            this.fail('Invalid LINE ID Token', 401);
            return EMPTY;
          }),
        )
        .subscribe();
    } else {
      this.logger.error(
        `Authentication failed empty idToken: ${req.headers.authorization}`,
      );
      this.fail('No ID token provided', 401);
    }
  }

  validate(payload: VerifyIDToken): Identity {
    return {
      id: payload.sub,
      roles: [],
    };
  }
}
