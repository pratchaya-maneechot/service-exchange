import { ObjectType, Field, registerEnumType } from '@nestjs/graphql';
export const EnumUserRole = {
  POSTER: 'POSTER',
  TASKER: 'TASKER',
} as const;

export type EnumUserRole = keyof typeof EnumUserRole;

registerEnumType(EnumUserRole, {
  name: 'EnumUserRole',
  description: 'User roles (POSTER or TASKER)',
});

@ObjectType()
export class User {
  @Field(() => String)
  id: string;

  @Field(() => String)
  name: string;

  @Field(() => String, { nullable: true })
  phone?: string | null;

  @Field(() => String, { nullable: true })
  password?: string | null;

  @Field(() => String, { nullable: true })
  email?: string | null;

  @Field(() => [EnumUserRole])
  roles: EnumUserRole[];

  @Field(() => Date)
  createdAt: Date;

  @Field(() => Date, { nullable: true })
  updatedAt?: Date | null;
}
