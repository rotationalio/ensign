// useAuth hook base on react-query useLogin hook and add some logic to handle auth token
import { queryClient } from '@/application/config/react-query';
import { useOrgStore } from '@/store';
import { clearCookies, getCookie } from '@/utils/cookies';
import { decodeToken } from '@/utils/decodeToken';
export const useAuth = () => {
  const org = useOrgStore.getState() as any;

  const isAuthenticated = !!org.isAuthenticated;

  function logout() {
    org.reset();
    useOrgStore.persist.clearStorage();
    clearCookies();
    queryClient.clear();
  }

  const token = getCookie('bc_atk');
  console.log('[useAuth] token', getCookie('access_token'));
  const decodedToken = token && decodeToken(token);
  if (decodedToken) {
    const { exp } = decodedToken;
    const now = new Date().getTime() / 1000;
    if (exp < now) {
      // token expired so logout user and clear cookies
      // we could refresh token later on
      logout();
    }
  }

  return {
    isAuthenticated,
    logout,
  };
};
