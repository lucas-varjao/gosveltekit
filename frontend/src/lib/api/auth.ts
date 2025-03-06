// frontend/src/lib/api/auth.ts

import { apiRequest, setTokens, clearTokens } from './client';

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

export interface AuthResponse {
	access_token: string;
	refresh_token: string;
	expires_at: string;
	user: {
		ID: number;
		Username: string;
		Email: string;
		DisplayName: string;
		Role: string;
	};
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
		const response = await apiRequest<AuthResponse>('/auth/login', {
			method: 'POST',
			body: JSON.stringify(data)
		});

		// Store tokens on successful login
		setTokens(response.access_token, response.refresh_token);
		return response;
	},

	// Register new user
	register: async (data: RegisterRequest): Promise<AuthResponse> => {
		const response = await apiRequest<AuthResponse>('/auth/register', {
			method: 'POST',
			body: JSON.stringify(data)
		});

		// Store tokens on successful registration
		setTokens(response.access_token, response.refresh_token);
		return response;
	},

	// Logout user
	logout: async (): Promise<void> => {
		try {
			await apiRequest('/api/logout', { method: 'POST' });
		} finally {
			clearTokens();
		}
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
	}
};
