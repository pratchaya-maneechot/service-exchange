import { InputType } from '@nestjs/graphql';

@InputType()
export class LineRegisterInput {
  lineUserId: string;
  email?: string;
  password?: string;
  displayName: string;
  avatarUrl?: string;
}
