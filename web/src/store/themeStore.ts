/**
 * Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
 * This source code is proprietary. Confidential and private.
 * Contact: iletisim@alibuyuk.net | Website: alibuyuk.net
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface ThemeState {
    isDark: boolean;
    toggle: () => void;
    setDark: (value: boolean) => void;
}

export const useThemeStore = create<ThemeState>()(
    persist(
        (set) => ({
            isDark: false,
            toggle: () => set((state) => ({ isDark: !state.isDark })),
            setDark: (value) => set({ isDark: value }),
        }),
        {
            name: 'subflow-theme',
        }
    )
);

// Auth store for user session
interface User {
    id: string;
    email: string;
    firstName: string;
    lastName: string;
    role: 'ADMIN' | 'MANAGER' | 'ACCOUNTANT' | 'VIEWER';
}

interface AuthState {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
    login: (user: User, token: string) => void;
    logout: () => void;
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set) => ({
            user: null,
            token: null,
            isAuthenticated: false,
            login: (user, token) => set({ user, token, isAuthenticated: true }),
            logout: () => set({ user: null, token: null, isAuthenticated: false }),
        }),
        {
            name: 'subflow-auth',
        }
    )
);

// Application store for global state
interface AppState {
    sidebarOpen: boolean;
    currentTenantId: string | null;
    toggleSidebar: () => void;
    setTenantId: (id: string) => void;
}

export const useAppStore = create<AppState>((set) => ({
    sidebarOpen: true,
    currentTenantId: null,
    toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
    setTenantId: (id) => set({ currentTenantId: id }),
}));
