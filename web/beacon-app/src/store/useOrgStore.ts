import { create } from 'zustand';
import { createJSONStorage, devtools, persist } from 'zustand/middleware';

import { MemberResponse } from '@/features/members/types/memberServices';
//TODO: clean up this store and remove unused properties
const useOrgStore = create(
  persist(
    devtools((set: any) => ({
      orgID: null,
      userID: null,
      name: null,
      orgName: null,
      email: null,
      isAuthenticated: false,
      permissions: null,
      application: {
        isProjectActive: false,
      },
      onboarding: {
        currentStep: null,
      },
      tempData: null,
      setIsProjectActive: (isProjectActive: boolean) =>
        set({
          application: {
            isProjectActive,
          },
        }),
      setOrg: (org: string) => set({ org }),
      setTenantID: (currentTenantID: string) => set({ currentTenantID }),
      setOrgName: (orgName: string) => set({ orgName }),
      setAuthUser: (token: any, isAuthed: boolean) =>
        set({
          orgID: token?.org,
          userID: token?.sub,
          orgName: token?.name,
          name: token?.name,
          email: token?.email,
          picture: token?.picture,
          permissions: token?.permissions,
          isAuthenticated: isAuthed,
        }),
      setUserProfile: (member: MemberResponse) => set({ ...member }),
      setTempData: (tempData: any) => set({ ...tempData }),
      setOnboardingStep: (currentStep: number) => set({ onboarding: { currentStep } }),
      decrementStep: () =>
        set((state: any) => ({ onboarding: { currentStep: state.onboarding.currentStep - 1 } })),
      increaseStep: () =>
        set((state: any) => ({ onboarding: { currentStep: state.onboarding.currentStep + 1 } })),

      setState: (state: any) =>
        set({
          orgID: state.org,
          userID: state.user,
          name: state.name,
          email: state.email,
          orgName: state.orgName,
        }),
      resetOnboarding: () => set({ onboarding: { currentStep: null } }),
      resetTempData: () => set({ tempData: null }),
      reset: () =>
        set({
          orgID: null,
          userID: null,
          orgName: null,
          name: null,
          email: null,
          isAuthenticated: false,
          permissions: null,
          project: {
            id: null,
            name: null,
          },
          application: {
            isProjectActive: false,
          },
          userProfile: null,
          onboarding: {
            currentStep: null,
          },
          tempData: null,
        }),
    })),
    { name: 'org', storage: createJSONStorage(() => sessionStorage) }
  )
);

export default useOrgStore;
