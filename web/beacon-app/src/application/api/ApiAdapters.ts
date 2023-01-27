import type { UserAuthResponse } from '@/features/auth/types/LoginService';
import type { NewUserAccount, User } from '@/features/auth/types/RegisterService';
export interface ApiAdapters {
  createNewAccount(user: NewUserAccount): Promise<User>;
  authenticateUser(user: Pick<User, 'username' | 'password'>): Promise<UserAuthResponse>;
}
