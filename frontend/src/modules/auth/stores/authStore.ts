import { create } from "zustand";

export type Permission = {
  uid: string;
  label: string;
};

type AuthState = {
  isAuthenticated: boolean;
  userId: string | null;
  permissions: Permission[] | null;

  setAuth: (userId: string) => void;
  setPermissions: (p: Permission[]) => void;
  logout: () => void;
};

export const useAuthStore = create<AuthState>((set) => ({
  isAuthenticated: false,
  userId: null,
  permissions: null,

  setAuth: (userId) =>
    set({
      isAuthenticated: true,
      userId,
    }),

  setPermissions: (permissions) => set({ permissions }),

  logout: () =>
    set({
      isAuthenticated: false,
      userId: null,
      permissions: null,
    }),
}));
