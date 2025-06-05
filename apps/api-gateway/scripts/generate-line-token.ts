import { VerifyIDToken } from '@line/bot-sdk';
import { JwtService } from '@nestjs/jwt';
import * as dotenv from 'dotenv';
import { v4 } from 'uuid';

// eslint-disable-next-line @typescript-eslint/no-unsafe-call
dotenv.config();

const secret = process.env.LINE_CHANNEL_SECRET;

const jwt = new JwtService({
  secret,
  signOptions: {
    algorithm: 'HS256',
  },
});

const payload = {
  iss: 'https://access.line.me',
  sub: 'U1234567890abcdef1234567890abcdef',
  aud: process.env.LINE_CHANNEL_ID,
  exp: Math.floor(Date.now() / 1000) + 3600,
  iat: Math.floor(Date.now() / 1000),
  nonce: v4(),
  name: 'Test User', //optional
  picture: 'https://example.com/profile.jpg', //optional
  email: 'test@example.com', //optional
};

const token = jwt.sign(payload);
console.log('ID Token:', token);
