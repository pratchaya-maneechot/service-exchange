import { JwtService } from '@nestjs/jwt';
import { v4 } from 'uuid';

const jwt = new JwtService({
  secret: 'LINE_CHANNEL_SECRETLINE_CHANNEL_SECRET',
  signOptions: {
    algorithm: 'HS256',
  },
});

const payload = {
  iss: 'https://access.line.me',
  sub: '1a2bac2c-57be-457d-98a2-d91788461a56',
  aud: 'LINE_CHANNEL_IDLINE_CHANNEL_ID',
  exp: Math.floor(Date.now() / 1000) + 3600,
  iat: Math.floor(Date.now() / 1000),
  nonce: v4(),
  name: 'Test User', //optional
  picture: 'https://example.com/profile.jpg', //optional
  email: 'test@example.com', //optional
};

const token = jwt.sign(payload);
console.log('ID Token:', token);
