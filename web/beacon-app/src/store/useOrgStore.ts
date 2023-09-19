import { create } from 'zustand';
import { createJSONStorage, devtools, persist } from 'zustand/middleware';

import { MemberResponse } from '@/features/members/types/memberServices';
//TODO: clean up this store and remove unused properties
const useOrgStore = create(
  persist(
    devtools((set) => ({
      orgID: null,
      userID: null,
      name: null,
      orgName: null,
      email: null,
      isAuthenticated: false,
      picture: null,
      currentTenantID: null,
      projectID: null, // should remove this in favor of project.id
      permissions: null,
      userProfile: null,
      project: {
        id: null,
        name: null,
      },
      application: {
        isProjectActive: false,
      },
      onboarding: {
        currentStep: null,
      },
      setIsProjectActive: (isProjectActive: boolean) =>
        set({
          application: {
            isProjectActive,
          },
        }),
      setOrg: (org: string) => set({ org }),
      setTenantID: (currentTenantID: string) => set({ currentTenantID }),
      setAuthUser: (token: any, isAuthed: boolean) =>
        set({
          org: token.org,
          user: token.sub,
          name: token.name,
          isAuthenticated: isAuthed,
          email: token.email,
          picture: token?.picture,
          permissions: token?.permissions,
        }),
      setUserProfile: (member: MemberResponse) => set({ ...member }),
      setOnboardingStep: (currentStep: number) => set({ onboarding: { currentStep } }),
      decrementStep: () =>
        set((state: any) => ({ onboarding: { currentStep: state.onboarding.currentStep - 1 } })),
      increaseStep: () =>
        set((state: any) => ({ onboarding: { currentStep: state.onboarding.currentStep + 1 } })),
      setUser: (user: string) => set({ user }),
      setName: (name: string) => set({ name }),
      setEmail: (email: string) => set({ email }),
      setPicture: (picture: string) => set({ picture }),
      setProjectID: (projectID: string) => set({ projectID }), // should remove this in favor of setProject
      setOrgName: (orgName: string) => set({ orgName }),
      setIsAuthenticated: (isAuthenticated: boolean) => set({ isAuthenticated }),
      setPermissions: (permissions: string[]) => set({ permissions }),
      setProject: (project: { id?: string; name?: string }) => set({ project }),
      setState: (state: any) =>
        set({
          org: state.org,
          user: state.user,
          name: state.name,
          email: state.email,
        }),
      resetOnboarding: () => set({ onboarding: { currentStep: null } }),
      reset: () =>
        set({
          org: null,
          user: null,
          name: null,
          email: null,
          isAuthenticated: false,
          picture: null,
          orgName: null,
          projectID: null,
          currentTenantID: null,
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
        }),
    })),
    { name: 'org', storage: createJSONStorage(() => sessionStorage) }
  )
);

export default useOrgStore;
