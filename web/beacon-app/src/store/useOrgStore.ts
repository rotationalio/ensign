import { create } from 'zustand';
import { createJSONStorage, devtools, persist } from 'zustand/middleware';
const useOrgStore = create(
  persist(
    devtools((set) => ({
      org: null,
      user: null,
      name: null,
      orgName: null,
      email: null,
      isAuthenticated: false,
      picture: null,
      projectID: null, // should remove this in favor of project.id
      permissions: null,
      project: {
        id: null,
        name: null,
      },
      setOrg: (org: string) => set({ org }),
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
          permissions: null,
          project: {
            id: null,
            name: null,
          },
        }),
    })),
    { name: 'org', storage: createJSONStorage(() => sessionStorage) }
  )
);

export default useOrgStore;
