import { Module } from '@nestjs/common';
import { LineService } from './line.service';
import { LineController } from './line.controller';
import { messagingApi, OAuth } from '@line/bot-sdk';
import { ConfigService } from '@nestjs/config';
import { EnvConfig } from 'src/config/config.type';
import { LineEventModule } from './events/events.module';

@Module({
  imports: [LineEventModule],
  controllers: [LineController],
  providers: [
    {
      provide: OAuth,
      useClass: OAuth,
    },
    {
      provide: messagingApi.MessagingApiClient,
      useFactory(config: ConfigService<EnvConfig>) {
        return new messagingApi.MessagingApiClient({
          channelAccessToken: config.getOrThrow('line.channelAccessToken', {
            infer: true,
          }),
        });
      },
      inject: [ConfigService],
    },
    LineService,
  ],
  exports: [LineService],
})
export class LineModule {}
