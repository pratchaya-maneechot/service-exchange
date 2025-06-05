import { VerifyIDToken } from '@line/bot-sdk';
import { Injectable } from '@nestjs/common';
import { PassportStrategy } from '@nestjs/passport';
import { ExtractJwt, Strategy } from 'passport-jwt';
import { Identity } from '../types/auth.type';
import { getOsEnv } from 'src/common/utils/env.util';

@Injectable()
export class LineStrategy extends PassportStrategy(Strategy, 'line') {
  constructor() {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      ignoreExpiration: true,
      secretOrKey: getOsEnv('LINE_CHANNEL_SECRET'),
    });
  }

  validate(payload: VerifyIDToken): Identity {
    return {
      id: payload.sub,
      roles: [],
    };
  }
}
