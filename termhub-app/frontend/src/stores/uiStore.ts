import { create } from 'zustand';

type ModalType = 'newTerminal' | 'settings' | null;

interface UIState {
  sidebarOpen: boolean;
  rightPanelOpen: boolean;
  modal: ModalType;
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  toggleRightPanel: () => void;
  setRightPanelOpen: (open: boolean) => void;
  openModal: (modal: NonNullable<ModalType>) => void;
  closeModal: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  sidebarOpen: true,
  rightPanelOpen: false,
  modal: null,

  toggleSidebar: () =>
    set((state) => ({ sidebarOpen: !state.sidebarOpen })),

  setSidebarOpen: (open) =>
    set({ sidebarOpen: open }),

  toggleRightPanel: () =>
    set((state) => ({ rightPanelOpen: !state.rightPanelOpen })),

  setRightPanelOpen: (open) =>
    set({ rightPanelOpen: open }),

  openModal: (modal) =>
    set({ modal }),

  closeModal: () =>
    set({ modal: null }),
}));
