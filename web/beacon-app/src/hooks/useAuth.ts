// useAuth hook base on react-query useLogin hook and add some logic to handle auth token
import { queryClient } from '@/application/config/react-query';
import { useOrgStore } from '@/store';
import { clearCookies } from '@/utils/cookies';
export const useAuth = () => {
  const org = useOrgStore.getState() as any;

  const isAuthenticated = !!org.isAuthenticated;

  function logout() {
    org.reset();
    clearCookies();
    queryClient.clear();
  }

  return {
    isAuthenticated,
    logout,
  };
};
