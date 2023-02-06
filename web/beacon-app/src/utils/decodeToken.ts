/* eslint-disable prettier/prettier */

import jwt from 'jwt-decode';

export const decodeToken = (token: string) => {
    return jwt(token);
};
