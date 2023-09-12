import Cookies from 'universal-cookie';

const cookie = new Cookies();

export const setCookie = (key: string, value: any, path = '/', secure = false) => {
  cookie.set(key, value, { path, secure });
};

export const getCookie = (key: string) => {
  return cookie.get(key);
};

export const removeCookie = (key: string, path = '/') => {
  cookie.remove(key, { path });
};

// get all cookies

export const getAllCookies = () => {
  return cookie.getAll();
};

// clear all cookies

export const clearCookies = () => {
  Object.keys(cookie.getAll()).forEach((key) => {
    removeCookie(key);
  });
  // remove all session storage
  sessionStorage.clear();
};

export const clearSessionStorage = () => {
  sessionStorage.clear();
};
