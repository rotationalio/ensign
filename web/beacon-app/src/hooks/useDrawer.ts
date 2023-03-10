import { create } from 'zustand';

type State = {
  isOpen: boolean;
  openDrawer: () => void;
  closeDrawer: () => void;
};

const useDrawer = create<State>((set) => ({
  isOpen: false,
  openDrawer: () => set((state) => ({ ...state, isOpen: true })),
  closeDrawer: () => set((state) => ({ ...state, isOpen: false })),
}));

export default useDrawer;
