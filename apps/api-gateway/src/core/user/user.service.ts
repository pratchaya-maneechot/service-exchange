import { Injectable } from '@nestjs/common';
import { UserClientService } from 'src/grpc-client/user-client.service';
import { LineRegisterRequest } from 'src/grpc-client/types/generated/user/LineRegisterRequest';
import { User } from '../entities/user.entity';
import { UpdateProfileInput } from './inputs/update-profile.input';
import { User__Output } from 'src/grpc-client/types/generated/user/User';
import { toDate } from 'src/common/utils/grpc.util';
import { protoToEnumUserRole } from '../utils';

@Injectable()
export class UserService {
  constructor(private readonly user: UserClientService) {}

  private transformUser(user: User__Output): User {
    return {
      id: user.id,
      name: user.name,
      phone: user.phone,
      password: user.password,
      email: user.email,
      roles: user.roles.map((r) =>
        typeof r === 'number' ? protoToEnumUserRole(r) : r,
      ),
      createdAt: toDate(user.createdAt) ?? new Date(),
      updatedAt: toDate(user.updatedAt),
    };
  }

  async getProfile(userId: string) {
    const profile = await this.user.userService.getProfile({ userId });
    return this.transformUser(profile!.user!);
  }

  async lineRegister(request: LineRegisterRequest) {
    const result = await this.user.userService.lineRegister(request);
    return result;
  }

  async updateProfile(id: string, input: UpdateProfileInput) {
    const result = await this.user.userService.updateProfile({
      userId: id,
      displayName: input.displayName,
      pictureUrl: input.pictureUrl,
    });
    return result;
  }
}
