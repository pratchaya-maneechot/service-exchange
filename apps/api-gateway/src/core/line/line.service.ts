import {
  OAuth,
  messagingApi,
  validateSignature,
  WebhookRequestBody,
  VerifyIDToken,
} from '@line/bot-sdk';
import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';
import { ConfigType } from 'src/common/types/config.type';

@Injectable()
export class LineService {
  constructor(
    @InjectPinoLogger(LineService.name)
    private readonly logger: PinoLogger,
    private readonly authClient: OAuth,
    private readonly config: ConfigService<ConfigType>,
    private readonly client: messagingApi.MessagingApiClient,
  ) {}

  public async getProfile(
    lineUserId: string,
  ): Promise<messagingApi.UserProfileResponse> {
    const profile = await this.client.getProfile(lineUserId);
    return profile;
  }

  public verifyWebhook(
    requestBody: WebhookRequestBody,
    signature: string,
  ): boolean {
    const body = JSON.stringify(requestBody);
    const isValid = validateSignature(
      body,
      this.config.getOrThrow('line.channelSecret', { infer: true }),
      signature,
    );
    return isValid;
  }

  public async verifyIdToken(id_token: string): Promise<VerifyIDToken> {
    const verify = await this.authClient.verifyIdToken(
      id_token,
      this.config.getOrThrow('line.channelId', { infer: true }),
    );
    return verify;
  }
}
