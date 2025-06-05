import { Injectable, Logger } from '@nestjs/common';
import { WebhookEvent } from '@line/bot-sdk';
import { IEventHandler } from './events.type';

@Injectable()
export class EventHandlerFactory {
  private readonly logger = new Logger(EventHandlerFactory.name);
  private readonly eventHandlers = new Map<
    IEventHandler['eventType'],
    IEventHandler
  >();

  constructor(...handlers: IEventHandler[]) {
    handlers.forEach((handler) => {
      this.eventHandlers.set(handler.eventType, handler);
    });
  }

  async handleEvent(event: WebhookEvent): Promise<void> {
    const eventHandlers = this.eventHandlers.get(event.type);
    if (!eventHandlers) {
      this.logger.warn(`Unsupported event type: ${event.type}`);
      return;
    }

    try {
      await eventHandlers.handle(event);
      this.logger.log(`${event.type} completed`);
    } catch (error) {
      this.logger.error(`Error handling ${event.type} event:`, {
        error: error.message,
        eventType: event.type,
        timestamp: new Date().toISOString(),
      });
    }
  }
}
