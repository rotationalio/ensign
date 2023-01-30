export interface User {
  id: string;
  full_name?: string;
  email: string;
  username: string;
  password: string;
}

export type NewUserAccount = Omit<User, 'id'>;
