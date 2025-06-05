export interface AppConfig {
  port: number;
  host: string;
  nodeEnv: 'production' | 'development' | 'test';
  corsOrigins?: string[];
  appName: string;
  appVersion: string;
}

export interface RedisConfig {
  host: string;
  port: number;
  password: string;
  db: number;
}

export interface JwtConfig {
  secret: string;
  expiresIn: string;
  refreshSecret: string;
  refreshExpiresIn: string;
}
export interface LineConfig {
  channelAccessToken: string;
  channelId: string;
  channelSecret: string;
}

export interface GrpcConfig {
  packages: Record<
    string,
    {
      endpoint: string;
      protoPath: string;
    }
  >;
}

export interface EnvConfig {
  app: AppConfig;
  redis: RedisConfig;
  // jwt: JwtConfig;
  grpc: GrpcConfig;
  line: LineConfig;
}
