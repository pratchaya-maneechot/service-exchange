import { Injectable, Logger } from '@nestjs/common';
import { UserClientService } from 'src/grpc-client/user-client.service';
import { LineRegisterRequest } from 'src/grpc-client/types/generated/user/LineRegisterRequest';
import { User } from '../entities/user.entity';
import { UpdateProfileInput } from './inputs/update-profile.input';
import { GetProfileResponse__Output } from 'src/grpc-client/types/generated/user/GetProfileResponse';
import { User__Output } from 'src/grpc-client/types/generated/user/User';
import { toDate } from 'src/common/utils/grpc.util';

@Injectable()
export class UserService {
  private readonly logger = new Logger(UserService.name);

  constructor(private readonly client: UserClientService) {}

  private transformUser(user: User__Output): User {
    return {
      id: user.id,
      name: user.name,
      phone: user.phone,
      password: user.password,
      email: user.email,
      roles: user.roles,
      createdAt: toDate(user.createdAt)!,
      updatedAt: toDate(user.updatedAt),
    };
  }

  getProfile(userId: string) {
    return new Promise<User>((resolve, reject) => {
      this.client.user.getProfile(
        { userId },
        (error, response: GetProfileResponse__Output) => {
          if (error) {
            this.logger.error(
              `Failed to get profile: ${error?.message}`,
              error?.stack,
            );
            return reject(error);
          }
          if (!response.user) {
            this.logger.error('No user found in response');
            return reject(new Error('User not found'));
          }
          resolve(this.transformUser(response.user));
        },
      );
    });
  }

  lineRegister(request: LineRegisterRequest) {
    return new Promise<{ userId: string }>((resolve, reject) => {
      this.client.user.lineRegister(request, (error, response) => {
        if (error) {
          this.logger.error(
            `Failed to register via Line: ${error.message}`,
            error.stack,
          );
          return reject(error);
        }
        resolve({ userId: response!.userId });
      });
    });
  }

  updateProfile(id: string, input: UpdateProfileInput) {
    return new Promise<User>((resolve, reject) => {
      this.client.user.updateProfile(
        { ...input, userId: id },
        (error, response: User__Output) => {
          if (error) {
            this.logger.error(
              `Failed to update profile: ${error.message}`,
              error.stack,
            );
            return reject(error);
          }
          resolve(this.transformUser(response));
        },
      );
    });
  }
}
