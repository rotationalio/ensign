import { create } from 'zustand';
import { createJSONStorage, devtools, persist } from 'zustand/middleware';
const useOrgStore = create(
  persist(
    devtools((set) => ({
      org: null,
      user: null,
      name: null,
      email: null,
      isAuthenticated: false,
      setOrg: (org: string) => set({ org }),
      setUser: (user: string) => set({ user }),
      setName: (name: string) => set({ name }),
      setEmail: (email: string) => set({ email }),
      setIsAuthenticated: (isAuthenticated: boolean) => set({ isAuthenticated }),
      setState: (state: any) =>
        set({
          org: state.org,
          user: state.user,
          name: state.name,
          email: state.email,
        }),
      reset: () => set({ org: null, user: null, name: null, email: null, isAuthenticated: false }),
    })),
    { name: 'org', storage: createJSONStorage(() => sessionStorage) }
  )
);

export default useOrgStore;
