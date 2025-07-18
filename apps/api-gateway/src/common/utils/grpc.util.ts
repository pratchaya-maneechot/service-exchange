import { Timestamp__Output } from '@grpc/grpc-js/build/src/generated/google/protobuf/Timestamp';
import { StringValue__Output } from '../../grpc-client/types/generated/google/protobuf/StringValue';

export function toDate(timestamp: Timestamp__Output | null): Date | null {
  if (!timestamp) return null;
  const milliseconds =
    Number(timestamp.seconds) * 1000 + Math.floor(timestamp.nanos / 1_000_000);
  return new Date(milliseconds);
}

export function toStrVal(val: StringValue__Output | null): string | null {
  if (!val) return null;
  return val.value;
}
