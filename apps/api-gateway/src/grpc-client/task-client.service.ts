import { Injectable } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { join } from 'path';
import { Observable } from 'rxjs';
import { ProtoGrpcType } from './types/generated/user';
import { UserServiceClient } from './types/generated/user/UserService';
import { GetProfileRequest } from './types/generated/user/GetProfileRequest';
import { GetProfileResponse__Output } from './types/generated/user/GetProfileResponse';
import { LineRegisterResponse__Output } from './types/generated/user/LineRegisterResponse';
import { UpdateProfileResponse__Output } from './types/generated/user/UpdateProfileResponse';
import { LineRegisterRequest } from './types/generated/user/LineRegisterRequest';
import { UpdateProfileRequest } from './types/generated/user/UpdateProfileRequest';
@Injectable()
export class UserGrpcClientService {
  private client: UserServiceClient;

  constructor() {
    const PROTO_PATH = join(__dirname, 'user.proto');
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });
    const proto = grpc.loadPackageDefinition(
      packageDefinition,
    ) as unknown as ProtoGrpcType;
    this.client = new proto.user.UserService(
      'localhost:50051',
      grpc.credentials.createInsecure(),
    );
  }

  getProfile(
    request: GetProfileRequest,
  ): Observable<GetProfileResponse__Output> {
    return new Observable((observer) => {
      this.client.getProfile(
        request,
        (error, response: GetProfileResponse__Output) => {
          if (error) {
            observer.error(error);
            return;
          }
          observer.next(response);
          observer.complete();
        },
      );
    });
  }

  lineRegister(
    request: LineRegisterRequest,
  ): Observable<LineRegisterResponse__Output> {
    return new Observable((observer) => {
      this.client.lineRegister(
        request,
        (error, response: LineRegisterResponse__Output) => {
          if (error) {
            observer.error(error);
            return;
          }
          observer.next(response);
          observer.complete();
        },
      );
    });
  }

  updateProfile(
    request: UpdateProfileRequest,
  ): Observable<UpdateProfileResponse__Output> {
    return new Observable((observer) => {
      this.client.updateProfile(
        request,
        (error, response: UpdateProfileResponse__Output) => {
          if (error) {
            observer.error(error);
            return;
          }
          observer.next(response);
          observer.complete();
        },
      );
    });
  }
}
