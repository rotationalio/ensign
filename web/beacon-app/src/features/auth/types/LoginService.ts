
import { User } from './RegisterService';

export interface UserAuthResponse {
    access_token: string;
    refresh_token: string;
}

export interface LoginMutation {
    authenticate: (user: Pick<User, 'username' | 'password'>) => void;
    reset: () => void;
    auth: UserAuthResponse;
    error: any;
    isAuthenticating: boolean;
    authenticated: boolean;
    hasAuthFailed: boolean;
}

export const isAuthenticated = (mutation: LoginMutation): mutation is Required<LoginMutation> =>
    mutation.authenticated && mutation.auth !== undefined;
