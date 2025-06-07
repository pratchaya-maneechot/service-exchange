import { Request } from 'express';
import { Identity } from 'src/core/auth/types/auth.type';
export type IAppContext = {
  req: Request & { user: Identity };
};
