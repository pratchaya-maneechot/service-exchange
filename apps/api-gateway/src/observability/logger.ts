import { Module } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { LoggerModule } from 'nestjs-pino';
import { EnvConfig } from 'src/config/config.type';
import { v4 } from 'uuid';
@Module({
  imports: [
    LoggerModule.forRootAsync({
      useFactory(config: ConfigService<EnvConfig>) {
        return {
          pinoHttp: {
            name: 'Backend',
            level:
              config.getOrThrow('app.nodeEnv', { infer: true }) !== 'production'
                ? 'debug'
                : 'info',
            genReqId: (req) =>
              req.id ?? (req.headers['x-request-id'] as string) ?? v4(),
            redact: [
              // `req.headers.x-line-signature`,
              'req.headers.authorization',
              'req.headers.cookie',
            ],
            autoLogging: {
              ignore: (req) => req.url === '/health' || req.url === '/metrics',
            },
            transport:
              config.getOrThrow('app.nodeEnv', { infer: true }) === 'production'
                ? {
                    targets: [
                      // {
                      //   target: 'pino-loki',
                      //   level: 'info',
                      //   options: {
                      //     host: getOsEnv('LOKI_HOST'),
                      //     labels: {
                      //       application: `backend:${getOsEnv('APP_VERSION')}`,
                      //     },
                      //     batching: true,
                      //     interval: 5,
                      //   },
                      // },
                      {
                        target: 'pino-pretty',
                        level: 'debug',
                        options: {
                          colorize: true,
                          translateTime: 'SYS:yyyy-mm-dd HH:MM:ss',
                        },
                      },
                    ],
                  }
                : {
                    target: 'pino-pretty',
                    options: {
                      colorize: true,
                      translateTime: 'SYS:yyyy-mm-dd HH:MM:ss',
                    },
                  },
          },
        };
      },
      inject: [ConfigService],
    }),
  ],
})
export class AppLoggerModule {}
