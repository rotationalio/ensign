// useAuth hook base on react-query useLogin hook and add some logic to handle auth token
import { useOrgStore } from '@/store';
import { clearCookies, getCookie } from '@/utils/cookies';
import { decodeToken } from '@/utils/decodeToken';
export const useAuth = () => {
  const org = useOrgStore.getState() as any;

  const isAuthenticated = !!org?.isAuthenticated;

  const logout = () => {
    useOrgStore.setState((state: any) => state.reset());
    useOrgStore.persist.clearStorage();
    clearCookies();
  };

  const token = getCookie('bc_atk');
  const decodedToken = token && decodeToken(token);
  if (decodedToken) {
    const { exp } = decodedToken;
    const now = new Date().getTime() / 1000;
    if (exp < now) {
      // call reset slice
      logout();
    }
  }

  return {
    isAuthenticated,
    logout,
  };
};
