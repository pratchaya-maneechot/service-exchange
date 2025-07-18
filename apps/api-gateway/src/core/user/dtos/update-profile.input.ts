import { Field, InputType } from '@nestjs/graphql';

@InputType()
export class UpdateProfileInput {
  @Field()
  displayName: string;

  @Field()
  pictureUrl: string;
}
