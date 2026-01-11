// frontend/src/lib/stores/auth.ts

/**
 * Authentication store for managing user state
 * Uses session-based authentication
 */
import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { authApi, type AuthResponse } from '$lib/api/auth';
import { getSessionId, clearSession } from '$lib/api/client';

// User interface matching new backend response
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

// Initial state
const initialState: AuthState = {
	user: null,
	isAuthenticated: false,
	isLoading: true,
	error: null
};

// Create the store
function createAuthStore() {
	const { subscribe, update } = writable<AuthState>(initialState);

	return {
		subscribe,

		// Initialize auth state - check if we have a valid session
		init: async () => {
			if (!browser) return;

			update((state) => ({ ...state, isLoading: true }));

			try {
				const sessionId = getSessionId();
				if (sessionId) {
					// Try to fetch current user to validate session
					try {
						const user = await authApi.getCurrentUser();
						update((state) => ({
							...state,
							user,
							isAuthenticated: true,
							isLoading: false
						}));
					} catch {
						// Session invalid, clear it
						clearSession();
						update((state) => ({ ...state, isLoading: false }));
					}
				} else {
					update((state) => ({ ...state, isLoading: false }));
				}
			} catch (error) {
				console.error('Failed to initialize auth state:', error);
				update((state) => ({ ...state, isLoading: false }));
			}
		},

		// Login user
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

		// Register user
		register: async (data: {
			username: string;
			email: string;
			password: string;
			display_name: string;
		}) => {
			update((state) => ({ ...state, isLoading: true, error: null }));

			try {
				const response = await authApi.register({
					username: data.username,
					email: data.email,
					password: data.password,
					display_name: data.display_name
				});

				update((state) => ({
					...state,
					user: response.user,
					isAuthenticated: true,
					isLoading: false
				}));

				return response.user;
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

		// Logout user
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

		// Clear error
		clearError: () => {
			update((state) => ({ ...state, error: null }));
		}
	};
}

// Export the store
export const authStore = createAuthStore();

// Initialize the store when the app loads
if (browser) {
	authStore.init();
}
