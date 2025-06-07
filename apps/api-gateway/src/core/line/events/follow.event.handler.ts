import { Injectable, Logger } from '@nestjs/common';
import { IEventHandler } from './events.type';
import { FollowEvent } from '@line/bot-sdk';
import { UserService } from 'src/core/user/user.service';
import { LineService } from '../line.service';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';

@Injectable()
export class FollowEventHandler implements IEventHandler<FollowEvent> {
  readonly eventType = 'follow';
  constructor(
    @InjectPinoLogger(FollowEventHandler.name)
    private readonly logger: PinoLogger,
    private readonly userService: UserService,
    private readonly lineService: LineService,
  ) {}

  async handle(event: FollowEvent): Promise<void> {
    try {
      if (!event.source?.userId) {
        this.logger.warn('Received follow event without userId', {
          timestamp: event.timestamp,
          mode: event.mode,
        });
        return;
      }

      const userId = event.source.userId;
      this.logger.info(
        {
          userId,
          timestamp: event.timestamp,
        },
        `User ${userId} followed the bot`,
      );

      const lineProfile = await this.lineService.getProfile(userId);
      await this.userService.lineRegister({
        lineId: lineProfile.userId,
        displayName: lineProfile.displayName,
        pictureUrl: lineProfile.pictureUrl,
        role: 'POSTER',
      });
    } catch (error) {
      this.logger.error('Failed to handle follow event', {
        error: error.message,
        userId: event.source?.userId,
        timestamp: new Date().toISOString(),
      });
    }
  }
}
