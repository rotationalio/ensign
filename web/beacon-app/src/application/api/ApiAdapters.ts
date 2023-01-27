
import type { NewUserAccount, User } from '@/features/auth/types/RegisterService';
import type { UserAuthResponse } from '@/features/auth/types/LoginService';
export interface ApiAdapters {
    createNewAccount(user: NewUserAccount): Promise<User>;
    authenticateUser(user: Pick<User, 'username' | 'password'>): Promise<UserAuthResponse>;
}