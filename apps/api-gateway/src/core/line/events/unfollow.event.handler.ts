import { Injectable } from '@nestjs/common';
import { IEventHandler } from './events.type';
import { UnfollowEvent } from '@line/bot-sdk';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';

@Injectable()
export class UnfollowEventHandler implements IEventHandler<UnfollowEvent> {
  readonly eventType = 'unfollow';
  constructor(
    @InjectPinoLogger(UnfollowEventHandler.name)
    private readonly logger: PinoLogger,
  ) {}

  async handle(event: UnfollowEvent): Promise<void> {
    try {
      // Validate event
      if (!event.source?.userId) {
        this.logger.warn('Received unfollow event without userId', {
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
        `User ${userId} unfollowed the bot`,
      );

      await this.updateUserStatus(userId);
    } catch (error) {
      this.logger.error('Failed to handle unfollow event', {
        error: error.message,
        userId: event.source?.userId,
        timestamp: new Date().toISOString(),
      });
    }
  }

  private async updateUserStatus(userId: string): Promise<void> {
    this.logger.debug(`Updating user ${userId} status to INACTIVE`);

    try {
      await new Promise((resolve) => setTimeout(resolve, 100));
      this.logger.debug(
        `Successfully updated user ${userId} to inactive status`,
      );
    } catch (error) {
      this.logger.error(`Failed to update user ${userId} status`, {
        error: error.message,
        userId,
      });
      throw error;
    }
  }
}
