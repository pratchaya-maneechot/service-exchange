import { Module } from '@nestjs/common';
import { LineService } from './line.service';
import { messagingApi, OAuth } from '@line/bot-sdk';
import { ConfigService } from '@nestjs/config';
import { ConfigType } from 'src/common/types/config.type';
import { EventHandlerFactory } from './events/events.factory';
import { IEventHandler } from './events/events.type';
import { FollowEventHandler } from './events/follow.event.handler';
import { UnfollowEventHandler } from './events/unfollow.event.handler';
import { LineController } from './line.controller';
import { UserModule } from 'src/core/user/user.module';

@Module({
  imports: [UserModule],
  controllers: [LineController],
  providers: [
    {
      provide: OAuth,
      useClass: OAuth,
    },
    {
      provide: messagingApi.MessagingApiClient,
      useFactory(config: ConfigService<ConfigType>) {
        return new messagingApi.MessagingApiClient({
          channelAccessToken: config.getOrThrow('line.channelAccessToken', {
            infer: true,
          }),
        });
      },
      inject: [ConfigService],
    },
    {
      provide: EventHandlerFactory,
      useFactory: (...handler: IEventHandler[]) => {
        return new EventHandlerFactory(...handler);
      },
      inject: [FollowEventHandler, UnfollowEventHandler],
    },
    FollowEventHandler,
    UnfollowEventHandler,
    LineService,
  ],
  exports: [LineService],
})
export class LineModule {}
