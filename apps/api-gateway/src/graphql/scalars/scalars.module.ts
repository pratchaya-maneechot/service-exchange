import { Module } from '@nestjs/common';
import { DateScalar } from './date.scalar';
import { JSONScalar } from './json.scalar';
import { UUIDScalar } from './uuid.scalar';

@Module({
  providers: [DateScalar, JSONScalar, UUIDScalar],
  exports: [DateScalar, JSONScalar, UUIDScalar],
})
export class ScalarsModule {}
