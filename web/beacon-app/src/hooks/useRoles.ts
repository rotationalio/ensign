import { useFetchMember } from '@/features/members/hooks/useFetchMember';
import { useOrgStore } from '@/store';
export const useRoles = () => {
  const orgDataState = useOrgStore.getState() as any;
  const { member } = useFetchMember(orgDataState?.user);

  const hasRoles = (rolesToCheck: string[]) => {
    if (member?.role) {
      return rolesToCheck.some((role) => role === member.role);
    }
    return false;
  };

  const hasRole = (roleToCheck: string) => {
    if (member?.role) {
      return member.role === roleToCheck;
    }
    return false;
  };

  const hasOneOfRoles = (rolesToCheck: string[]) => {
    if (member?.role) {
      return rolesToCheck.some((role) => role === member.role);
    }
    return false;
  };

  return {
    hasRoles,
    hasRole,
    hasOneOfRoles,
  };
};
