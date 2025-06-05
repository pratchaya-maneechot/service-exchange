import {
  OAuth,
  messagingApi,
  validateSignature,
  WebhookRequestBody,
  VerifyIDToken,
} from '@line/bot-sdk';
import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { EnvConfig } from 'src/config/config.type';

@Injectable()
export class LineService {
  private readonly logger = new Logger(LineService.name);

  constructor(
    private readonly authClient: OAuth,
    private readonly config: ConfigService<EnvConfig>,
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
