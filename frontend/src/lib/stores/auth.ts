// frontend/src/lib/stores/auth.ts

/**
 * Authentication store for managing user state
 */
import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { authApi } from '$lib/api/auth';

export interface User {
	ID: number;
	Username: string;
	Email: string;
	DisplayName: string;
	Role: string;
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

		// Initialize auth state from localStorage
		init: async () => {
			if (!browser) return;

			update((state) => ({ ...state, isLoading: true }));

			try {
				// Check if we have a stored user
				const userJson = localStorage.getItem('user');
				if (userJson) {
					const user = JSON.parse(userJson);
					update((state) => ({
						...state,
						user,
						isAuthenticated: true,
						isLoading: false
					}));
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

				// Store user data
				if (browser) {
					localStorage.setItem('user', JSON.stringify(response.user));
				}

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

				// Store user data
				if (browser) {
					localStorage.setItem('user', JSON.stringify(response.user));
				}

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
				// Clear user data
				if (browser) {
					localStorage.removeItem('user');
				}

				update((state) => ({
					...state,
					user: null,
					isAuthenticated: false,
					isLoading: false
				}));
			}
		}
	};
}

// Export the store
export const authStore = createAuthStore();

// Initialize the store when the app loads
if (browser) {
	authStore.init();
}
