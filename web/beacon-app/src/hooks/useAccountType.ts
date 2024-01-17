import { useOrgStore } from '@/store';

function useAccountType() {
  const state = useOrgStore((state: any) => state) as any;
  const accountType = state.account as string;

  return accountType;
}

export default useAccountType;
