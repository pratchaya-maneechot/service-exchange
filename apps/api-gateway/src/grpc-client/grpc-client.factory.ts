import { Injectable } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';

interface GrpcServiceConfig {
  protoPath: string;
  packageName: string;
  serviceName: string;
  serviceAddress: string;
}

@Injectable()
export class GrpcClientFactory {
  static key(packageName: string, serviceName: string) {
    return `${packageName}_${serviceName}`;
  }

  static splitKey(key: string) {
    return key.split('_');
  }

  private loadProtoFile(protoPath: string): grpc.GrpcObject {
    try {
      const packageDefinition = protoLoader.loadSync(protoPath, {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true,
      });
      return grpc.loadPackageDefinition(packageDefinition);
    } catch (error) {
      throw new Error(
        `Failed to load proto file at ${protoPath}: ${error.message}`,
      );
    }
  }

  createClient<T>(config: GrpcServiceConfig): T {
    const proto = this.loadProtoFile(config.protoPath);
    const packageObj = proto[config.packageName];

    if (!packageObj || !packageObj[config.serviceName]) {
      throw new Error(
        `Service ${config.serviceName} not found in package ${config.packageName}`,
      );
    }

    // eslint-disable-next-line @typescript-eslint/no-unsafe-call
    return new packageObj[config.serviceName](
      config.serviceAddress,
      grpc.credentials.createInsecure(),
    ) as T;
  }

  createClients<T extends grpc.Client>(
    configs: GrpcServiceConfig[],
  ): Map<string, T> {
    return configs.reduce(
      (clients, config) => {
        clients[GrpcClientFactory.key(config.packageName, config.serviceName)] =
          this.createClient<T>(config);
        return clients;
      },
      {} as Map<string, T>,
    );
  }
}
