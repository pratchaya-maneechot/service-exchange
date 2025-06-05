import { NestFactory } from '@nestjs/core';
import { ConfigService } from '@nestjs/config';
import { Logger, ValidationPipe } from '@nestjs/common';
import { AppModule } from './app.module';
import { EnvConfig } from './config/config.type';

async function bootstrap() {
  const logger = new Logger('Bootstrap');

  try {
    const app = await NestFactory.create(AppModule, {
      logger: ['error', 'warn', 'log', 'debug', 'verbose'],
    });

    const config = app.get(ConfigService<EnvConfig>);

    app.useGlobalPipes(
      new ValidationPipe({
        transform: true,
        whitelist: true,
        forbidNonWhitelisted: true,
        disableErrorMessages:
          config.get('app.nodeEnv', { infer: true }) === 'production',
      }),
    );

    app.enableCors({
      origin: config.get('app.corsOrigins', { infer: true }),
      credentials: true,
      methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
      allowedHeaders: [
        'Origin',
        'X-Requested-With',
        'Content-Type',
        'Accept',
        'Authorization',
        'Apollo-Require-Preflight',
      ],
    });

    const port = config.getOrThrow('app.port', { infer: true });
    const host = config.getOrThrow('app.host', { infer: true });

    await app.listen(port, host);

    logger.log(`🚀 Application is running on: http://${host}:${port}`);
    logger.log(`🎯 GraphQL Playground: http://${host}:${port}/graphql`);
    logger.log(`🏥 Health Check: http://${host}:${port}/healthz`);
    logger.log(`📊 Environment: ${config.get('app.nodeEnv', { infer: true })}`);

    process.on('SIGTERM', () => {
      logger.log('SIGTERM received, shutting down gracefully');
      void app.close().then(() => process.exit(0));
    });

    process.on('SIGINT', () => {
      logger.log('SIGINT received, shutting down gracefully');
      void app.close().then(() => process.exit(0));
    });
  } catch (error) {
    logger.error('Error starting application:', error);
    process.exit(1);
  }
}

process.on('unhandledRejection', (reason, promise) => {
  const logger = new Logger('UnhandledRejection');
  logger.error('Unhandled Rejection at:', promise, 'reason:', reason);
  process.exit(1);
});

process.on('uncaughtException', (error) => {
  const logger = new Logger('UncaughtException');
  logger.error('Uncaught Exception:', error);
  process.exit(1);
});

void bootstrap();
