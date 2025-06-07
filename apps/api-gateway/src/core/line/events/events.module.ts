import { Module } from '@nestjs/common';
import { EventHandlerFactory } from './events.factory';
import { FollowEventHandler } from './follow.event.handler';
import { UnfollowEventHandler } from './unfollow.event.handler';
import { IEventHandler } from './events.type';
import { UserModule } from 'src/core/user/user.module';
import { LineModule } from '../line.module';

@Module({
  imports: [UserModule, LineModule],
  providers: [
    {
      provide: EventHandlerFactory,
      useFactory: (...handler: IEventHandler[]) => {
        return new EventHandlerFactory(...handler);
      },
      inject: [FollowEventHandler, UnfollowEventHandler],
    },
    FollowEventHandler,
    UnfollowEventHandler,
  ],
  exports: [EventHandlerFactory],
})
export class LineEventModule {}
