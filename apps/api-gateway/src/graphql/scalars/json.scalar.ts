import { Scalar, CustomScalar } from '@nestjs/graphql';
import { ValueNode, Kind } from 'graphql';

@Scalar('JSON')
export class JSONScalar implements CustomScalar<any, any> {
  description = 'JSON custom scalar type';

  parseValue(value: any): any {
    return typeof value === 'object' ? value : JSON.parse(value);
  }

  serialize(value: any): any {
    return typeof value === 'object' ? value : JSON.parse(value);
  }

  parseLiteral(ast: ValueNode): any {
    if (ast.kind === Kind.STRING) {
      return JSON.parse(ast.value);
    }
    return null;
  }
}
