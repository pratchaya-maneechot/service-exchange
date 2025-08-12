export interface RawEnv {
  PORT?: string;
  HOST?: string;
  NODE_ENV?: string;
  CORS_ORIGINS?: string;
  APP_NAME?: string;
  APP_VERSION?: string;
  REDIS_HOST?: string;
  REDIS_PORT?: string;
  REDIS_PASSWORD?: string;
  REDIS_DB?: string;
  GRPC_USER_ENDPOINT?: string;
  GRPC_USER_VERSION?: string;
  LINE_CHANNEL_ID?: string;
  LINE_CHANNEL_SECRET?: string;
  LINE_CHANNEL_ACCESS_TOKEN?: string;
}
