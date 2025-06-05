import { EnumUserRole } from './entities/user.entity';

export function protoToEnumUserRole(protoRole: number): EnumUserRole {
  switch (protoRole) {
    case 1:
      return EnumUserRole.POSTER;
    case 2:
      return EnumUserRole.TASKER;
    default:
      throw new Error(`Unknown proto role: ${protoRole}`);
  }
}
