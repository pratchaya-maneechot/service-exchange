export function toDate(
  timestamp: {
    seconds: string;
    nanos: number;
  } | null,
): Date | null {
  if (!timestamp) return null;
  const milliseconds =
    Number(timestamp.seconds) * 1000 + Math.floor(timestamp.nanos / 1_000_000);
  return new Date(milliseconds);
}
