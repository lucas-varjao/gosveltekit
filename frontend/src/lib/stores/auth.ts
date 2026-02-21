// frontend/src/lib/stores/auth.ts

/**
 * Authentication store for managing user state
 * Uses session-based authentication
 */
import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { authApi } from '$lib/api/auth';

// User interface matching backend response from /api/me
export interface User {
    id: string;
    identifier: string; // username
    email: string;
    display_name: string;
    role: string;
    active: boolean;
}

interface AuthState {
    user: User | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    error: string | null;
}

const initialState: AuthState = {
    user: null,
    isAuthenticated: false,
    isLoading: true,
    error: null
};

function createAuthStore() {
    const { subscribe, update } = writable<AuthState>(initialState);

    return {
        subscribe,

        // Initialize auth state by validating server-side session cookie.
        init: async () => {
            if (!browser) return;

            update((state) => ({ ...state, isLoading: true }));

            try {
                const user = await authApi.getCurrentUser();
                update((state) => ({
                    ...state,
                    user,
                    isAuthenticated: true,
                    isLoading: false
                }));
            } catch {
                update((state) => ({
                    ...state,
                    user: null,
                    isAuthenticated: false,
                    isLoading: false
                }));
            }
        },

        login: async (username: string, password: string) => {
            update((state) => ({ ...state, isLoading: true, error: null }));

            try {
                const response = await authApi.login({ username, password });

                update((state) => ({
                    ...state,
                    user: response.user,
                    isAuthenticated: true,
                    isLoading: false
                }));

                return response.user;
            } catch (error) {
                const message = error instanceof Error ? error.message : 'Login failed';
                update((state) => ({
                    ...state,
                    error: message,
                    isLoading: false
                }));
                throw error;
            }
        },

        // Registration creates user but does not authenticate.
        register: async (data: {
            username: string;
            email: string;
            password: string;
            display_name: string;
        }) => {
            update((state) => ({ ...state, isLoading: true, error: null }));

            try {
                const registeredUser = await authApi.register({
                    username: data.username,
                    email: data.email,
                    password: data.password,
                    display_name: data.display_name
                });

                update((state) => ({
                    ...state,
                    isLoading: false,
                    isAuthenticated: false,
                    user: null
                }));

                return registeredUser;
            } catch (error) {
                const message = error instanceof Error ? error.message : 'Registration failed';
                update((state) => ({
                    ...state,
                    error: message,
                    isLoading: false
                }));
                throw error;
            }
        },

        logout: async () => {
            update((state) => ({ ...state, isLoading: true }));

            try {
                await authApi.logout();
            } finally {
                update((state) => ({
                    ...state,
                    user: null,
                    isAuthenticated: false,
                    isLoading: false
                }));
            }
        },

        clearError: () => {
            update((state) => ({ ...state, error: null }));
        }
    };
}

export const authStore = createAuthStore();

if (browser) {
    authStore.init();
}
