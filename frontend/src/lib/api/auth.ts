// frontend/src/lib/api/auth.ts

import { apiRequest, setSessionId, clearSession } from './client';

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

		// Store session ID on successful login
		setSessionId(response.session_id);
		return response;
	},

	// Register new user
	register: async (data: RegisterRequest): Promise<AuthResponse> => {
		const response = await apiRequest<AuthResponse>('/auth/register', {
			method: 'POST',
			body: JSON.stringify(data)
		});

		// Store session ID on successful registration
		setSessionId(response.session_id);
		return response;
	},

	// Logout user
	logout: async (): Promise<void> => {
		try {
			await apiRequest('/api/logout', { method: 'POST' });
		} finally {
			clearSession();
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
	},

	// Get current user (new endpoint)
	getCurrentUser: async () => {
		return apiRequest<AuthResponse['user']>('/api/me', { method: 'GET' });
	}
};
