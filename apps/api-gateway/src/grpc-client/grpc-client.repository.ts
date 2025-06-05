import { Injectable, OnModuleInit } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { GrpcClientFactory } from './grpc-client.factory';
import { EnvConfig } from 'src/config/config.type';
import * as grpc from '@grpc/grpc-js';

@Injectable()
export class GrpcClientRepository implements OnModuleInit {
  private clients: Map<string, grpc.Client> = new Map();

  constructor(
    private readonly factory: GrpcClientFactory,
    private readonly config: ConfigService<EnvConfig>,
  ) {}

  onModuleInit() {
    const packages = this.config.getOrThrow('grpc.packages', {
      infer: true,
    });

    this.clients = this.factory.createClients([
      {
        packageName: 'user',
        serviceName: 'UserService',
        protoPath: packages.user.protoPath,
        serviceAddress: packages.user.endpoint,
      },
    ]);
  }

  getClient<T extends grpc.Client>(key: string): T {
    const client = this.clients.get(key) as T;
    if (client) {
      return client;
    }

    const packages = this.config.getOrThrow('grpc.packages', {
      infer: true,
    });

    const [packageName, serviceName] = GrpcClientFactory.splitKey(key);
    if (!packageName || !serviceName)
      throw new Error(`gRPC client for ${key} not found`);

    const config = packages[packageName];
    if (!config) throw new Error(`config gRPC client for ${key} not found`);

    const _client = this.factory.createClient<T>({
      packageName,
      serviceName,
      protoPath: config.protoPath,
      serviceAddress: config.endpoint,
    });
    this.clients.set(key, _client);

    return _client;
  }
}
