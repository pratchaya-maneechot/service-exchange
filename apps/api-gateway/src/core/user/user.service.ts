import { Injectable } from '@nestjs/common';
import { UserClientService } from 'src/grpc-client/services/user-client.service';
import { UpdateProfileInput } from './dtos/update-profile.input';
import { LineRegisterInput } from './dtos/line-register.input';
import { UserProfile__Output } from '../../grpc-client/@types/generated/user/v1/UserProfile';
import { EnumUserRole, EnumUserStatus, User } from '../entities/user.entity';
import { toDate, toStrVal } from '../../common/utils/grpc.util';
import { NotFoundError } from '../../common/errors';

@Injectable()
export class UserService {
  constructor(private readonly user: UserClientService) {}

  private toUser(resp: UserProfile__Output | undefined): User | undefined {
    if (!resp) {
      return resp;
    }
    const res: User = {
      id: resp.userId,
      lineUserId: resp.lineUserId,
      displayName: resp.displayName,
      status: resp.status as EnumUserStatus,
      isVerified: resp.isVerified,
      createdAt: toDate(resp.createdAt)!,
      roles: resp.roles as EnumUserRole[],
      email: toStrVal(resp.email),
      firstName: toStrVal(resp.firstName),
      lastName: toStrVal(resp.lastName),
      bio: toStrVal(resp.bio),
      avatarUrl: toStrVal(resp.avatarUrl),
      phoneNumber: toStrVal(resp.phoneNumber),
      address: toStrVal(resp.address),
      preferences: resp.preferences,
      lastLoginAt: toDate(resp.lastLoginAt),
    };
    return res;
  }

  async getProfile(userId: string) {
    const result = await this.user.userService.getUserProfile({
      userId,
    });
    const profile = this.toUser(result);
    if (!profile) {
      throw new NotFoundError('User not found');
    }
    return profile;
  }

  async lineRegister(input: LineRegisterInput) {
    const result = await this.user.userService.LineRegister({
      lineUserId: input.lineUserId,
      displayName: input.displayName,
      avatarUrl: { value: input.avatarUrl },
      email: { value: input.email },
      password: { value: input.password },
    });
    return result;
  }

  async updateProfile(id: string, input: UpdateProfileInput) {
    //
  }
}
