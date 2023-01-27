
import type { NewUserAccount, User } from '@/features/registration/AccountService';
export interface ApiAdapters {
    createNewAccount(user: NewUserAccount): Promise<User>;
}