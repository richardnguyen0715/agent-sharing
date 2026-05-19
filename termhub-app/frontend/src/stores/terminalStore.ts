import { create } from 'zustand';
import type { Terminal, TerminalTab } from '@/types';
import * as api from '@/services/api';

interface TerminalState {
  terminals: Terminal[];
  tabs: TerminalTab[];
  activeTabId: string | null;
  isLoading: boolean;

  fetchTerminals: () => Promise<void>;
  createTerminal: (data: api.CreateTerminalRequest) => Promise<Terminal>;
  deleteTerminal: (id: string) => Promise<void>;
  setTerminals: (terminals: Terminal[]) => void;
  addTab: (tab: TerminalTab) => void;
  removeTab: (terminalId: string) => void;
  setActiveTab: (terminalId: string) => void;
}

export const useTerminalStore = create<TerminalState>((set, get) => ({
  terminals: [],
  tabs: [],
  activeTabId: null,
  isLoading: false,

  fetchTerminals: async () => {
    set({ isLoading: true });
    try {
      const terminals = await api.listTerminals();
      set({ terminals, isLoading: false });
    } catch {
      set({ isLoading: false });
    }
  },

  createTerminal: async (data) => {
    const terminal = await api.createTerminal(data);
    set((state) => ({
      terminals: [...state.terminals, terminal],
    }));
    // Auto-add a tab and activate it
    get().addTab({
      terminalId: terminal.id,
      title: terminal.title,
      status: 'active',
      colorTag: terminal.colorTag,
    });
    return terminal;
  },

  deleteTerminal: async (id) => {
    await api.deleteTerminal(id);
    set((state) => {
      const nextTerminals = state.terminals.filter((t) => t.id !== id);
      const nextTabs = state.tabs.filter((t) => t.terminalId !== id);
      const nextActiveId =
        state.activeTabId === id
          ? (nextTabs[nextTabs.length - 1]?.terminalId ?? null)
          : state.activeTabId;
      return { terminals: nextTerminals, tabs: nextTabs, activeTabId: nextActiveId };
    });
  },

  setTerminals: (terminals) =>
    set({ terminals }),

  addTab: (tab) =>
    set((state) => {
      // Don't add duplicate tabs
      const exists = state.tabs.some((t) => t.terminalId === tab.terminalId);
      if (exists) {
        return { activeTabId: tab.terminalId };
      }
      return {
        tabs: [...state.tabs, tab],
        activeTabId: tab.terminalId,
      };
    }),

  removeTab: (terminalId) =>
    set((state) => {
      const idx = state.tabs.findIndex((t) => t.terminalId === terminalId);
      const nextTabs = state.tabs.filter((t) => t.terminalId !== terminalId);
      let nextActiveId = state.activeTabId;
      if (state.activeTabId === terminalId) {
        // Try to select the tab to the left, then right, then null
        const prev = nextTabs[idx - 1];
        const next = nextTabs[idx];
        nextActiveId = prev?.terminalId ?? next?.terminalId ?? null;
      }
      return { tabs: nextTabs, activeTabId: nextActiveId };
    }),

  setActiveTab: (terminalId) =>
    set({ activeTabId: terminalId }),
}));
