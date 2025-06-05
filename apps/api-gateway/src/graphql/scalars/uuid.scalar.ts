import { Scalar, CustomScalar } from '@nestjs/graphql';
import { ValueNode, Kind } from 'graphql';

@Scalar('UUID')
export class UUIDScalar implements CustomScalar<string, string> {
  description = 'UUID custom scalar type';

  parseValue(value: string): string {
    if (!this.isValidUUID(value)) {
      throw new Error('Invalid UUID format');
    }
    return value;
  }

  serialize(value: string): string {
    if (!this.isValidUUID(value)) {
      throw new Error('Invalid UUID format');
    }
    return value;
  }

  parseLiteral(ast: ValueNode): string {
    if (ast.kind === Kind.STRING) {
      if (!this.isValidUUID(ast.value)) {
        throw new Error('Invalid UUID format');
      }
      return ast.value;
    }
    throw new Error('Invalid data type');
  }

  private isValidUUID(value: string): boolean {
    const uuidRegex =
      /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
    return uuidRegex.test(value);
  }
}
