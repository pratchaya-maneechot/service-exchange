import { ObjectType, Field, registerEnumType } from '@nestjs/graphql';
import { JSONScalar } from '../../graphql/scalars/json.scalar';

export const EnumUserStatus = {
  USER_STATUS_UNSPECIFIED: 'USER_STATUS_UNSPECIFIED',
  ACTIVE: 'ACTIVE',
  INACTIVE: 'INACTIVE',
  SUSPENDED: 'SUSPENDED',
  PENDING_VERIFICATION: 'PENDING_VERIFICATION',
} as const;

export type EnumUserStatus = keyof typeof EnumUserStatus;

registerEnumType(EnumUserStatus, {
  name: 'EnumUserStatus',
  description: 'User status enumeration',
});

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
  lineUserId: string;

  @Field(() => String, { nullable: true })
  email?: string | null;

  @Field(() => String)
  displayName: string;

  @Field(() => String, { nullable: true })
  firstName?: string | null;

  @Field(() => String, { nullable: true })
  lastName?: string | null;

  @Field(() => String, { nullable: true })
  bio?: string | null;

  @Field(() => String, { nullable: true })
  avatarUrl?: string | null;

  @Field(() => String, { nullable: true })
  phoneNumber?: string | null;

  @Field(() => String, { nullable: true })
  address?: string | null;

  @Field(() => JSONScalar, { nullable: true })
  preferences?: any;

  @Field(() => EnumUserStatus)
  status: EnumUserStatus;

  @Field(() => Boolean)
  isVerified: boolean;

  @Field(() => Date, { nullable: true })
  lastLoginAt?: Date | null;

  @Field(() => Date)
  createdAt: Date;

  @Field(() => [EnumUserRole])
  roles: EnumUserRole[];
}
