import {
  Body,
  Controller,
  Headers,
  HttpCode,
  HttpStatus,
  Logger,
  Post,
  UnauthorizedException,
} from '@nestjs/common';
import { LineService } from './line.service';
import { EventHandlerFactory } from './events/events.factory';
import { WebhookRequestBody } from '@line/bot-sdk';

@Controller('line')
export class LineController {
  private readonly logger = new Logger(LineController.name);

  constructor(
    private readonly lineService: LineService,
    private readonly eventHandler: EventHandlerFactory,
  ) {}

  @Post('/webhook')
  @HttpCode(HttpStatus.OK)
  async webhook(
    @Body() payload: WebhookRequestBody,
    @Headers('x-line-signature') signature: string,
  ) {
    if (!this.lineService.verifyWebhook(payload, signature)) {
      this.logger.warn('Invalid webhook signature', { payload, signature });
      throw new UnauthorizedException('Invalid webhook signature');
    }

    if (!payload.events.length) {
      return { status: true };
    }

    for (const event of payload.events) {
      await this.eventHandler.handleEvent(event);
    }

    return { status: true };
  }
}
