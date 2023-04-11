// hook that allows you to check if the user has a specific permission
import { useOrgStore } from '@/store';
export const usePermissions = () => {
  const { permissions } = useOrgStore.getState() as any;

  const hasOneOfPermissions = (permissionsToCheck: string[]) => {
    return permissionsToCheck.some((permission) => permissions.includes(permission));
  };

  const hasPermissions = (permissionsToCheck: string[]) => {
    return permissionsToCheck.every((permission) => permissions.includes(permission));
  };

  const hasPermission = (permissionToCheck: string) => {
    return permissions.includes(permissionToCheck);
  };

  return { hasPermission, hasOneOfPermissions, hasPermissions };
};
