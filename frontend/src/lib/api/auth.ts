// frontend/src/lib/api/auth.ts

import { apiRequest } from './client';

export interface LoginRequest {
	username: string;
	password: string;
}

export interface RegisterRequest {
	username: string;
	email: string;
	password: string;
	display_name: string;
}

// Updated to match new session-based response from backend
export interface AuthResponse {
	session_id: string;
	expires_at: string;
	user: {
		id: string;
		identifier: string;
		email: string;
		display_name: string;
		role: string;
		active: boolean;
	};
}

export interface RegisterResponse {
	id: number;
	username: string;
	email: string;
	display_name: string;
	role: string;
	active: boolean;
}

export interface PasswordResetRequest {
	email: string;
}

export interface PasswordResetConfirmRequest {
	token: string;
	new_password: string;
	confirm_password: string;
}

export const authApi = {
	// Login user
	login: async (data: LoginRequest): Promise<AuthResponse> => {
		return apiRequest<AuthResponse>('/auth/login', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	// Register new user
	register: async (data: RegisterRequest): Promise<RegisterResponse> => {
		return apiRequest<RegisterResponse>('/auth/register', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	// Logout user
	logout: async (): Promise<void> => {
		await apiRequest('/api/logout', { method: 'POST' });
	},

	// Request password reset
	requestPasswordReset: async (data: PasswordResetRequest): Promise<void> => {
		await apiRequest('/auth/password-reset-request', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	// Reset password with token
	resetPassword: async (data: PasswordResetConfirmRequest): Promise<void> => {
		await apiRequest('/auth/password-reset', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	// Get current user (new endpoint)
	getCurrentUser: async () => {
		return apiRequest<AuthResponse['user']>('/api/me', { method: 'GET' });
	}
};
