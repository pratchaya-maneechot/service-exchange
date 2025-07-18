import { Resolver, Query, Mutation, Args } from '@nestjs/graphql';
import { UseGuards } from '@nestjs/common';
import { User } from '../entities/user.entity';
import { UpdateProfileInput } from './dtos/update-profile.input';
import { LineAuthGuard } from 'src/core/auth/guards/line.guard';
import { CurrentIdentity } from 'src/core/auth/decorators/current-identity';
import { Identity } from 'src/core/auth/types/auth.type';
import { UserService } from './user.service';
import { UUIDScalar } from 'src/graphql/scalars/uuid.scalar';

@Resolver(() => User)
@UseGuards(LineAuthGuard)
export class UserResolver {
  constructor(private readonly userService: UserService) {}

  @Query(() => User, { name: 'profile' })
  async getProfile(@CurrentIdentity() user: Identity): Promise<User> {
    return this.userService.getProfile('ab4d3a5a-b7b9-4720-bf6a-45eaf6cb8b2e');
  }

  @Mutation(() => UUIDScalar)
  async updateProfile(
    @CurrentIdentity() user: Identity,
    @Args('input') input: UpdateProfileInput,
  ): Promise<string> {
    await this.userService.updateProfile(user.id, input);
    return user.id;
  }
}
