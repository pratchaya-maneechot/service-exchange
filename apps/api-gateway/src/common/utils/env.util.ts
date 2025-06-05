export function getOsEnv(key: string, defaultValue?: string): string {
  return (
    process.env[key] ??
    defaultValue ??
    (() => {
      throw new Error(`Environment variable ${key} is not set`);
    })()
  );
}

export function getOsEnvOptional(key: string): string | undefined {
  return process.env[key];
}

export function toNumber(value: string): number {
  const num = parseInt(value, 10);
  if (isNaN(num)) {
    throw new Error(`Cannot convert "${value}" to number`);
  }
  return num;
}

export function toBoolean(value: string): boolean {
  return value.toLowerCase() === 'true';
}
