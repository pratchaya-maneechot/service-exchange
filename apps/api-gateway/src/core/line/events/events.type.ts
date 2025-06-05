import { WebhookEvent } from '@line/bot-sdk';

export interface IEventHandler<T extends WebhookEvent = WebhookEvent> {
  handle(event: T): Promise<void>;
  readonly eventType: T['type'];
}
