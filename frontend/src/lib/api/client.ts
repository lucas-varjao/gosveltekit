// frontend/src/lib/api/client.ts

/**
 * API Client for communicating with the backend
 * Handles authentication, request formatting, and error handling
 */
import { browser } from '$app/environment';

// Base URL for API requests
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// Token storage keys
const ACCESS_TOKEN_KEY = 'accessToken';
const REFRESH_TOKEN_KEY = 'refreshToken';

// Error types
export class ApiError extends Error {
	status: number;

	constructor(message: string, status: number) {
		super(message);
		this.status = status;
		this.name = 'ApiError';
	}
}

// Token management
export const getAccessToken = (): string | null => {
	if (!browser) return null;
	return localStorage.getItem(ACCESS_TOKEN_KEY);
};

export const getRefreshToken = (): string | null => {
	if (!browser) return null;
	return localStorage.getItem(REFRESH_TOKEN_KEY);
};

export const setTokens = (accessToken: string, refreshToken: string): void => {
	if (!browser) return;
	localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
	localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
};

export const clearTokens = (): void => {
	if (!browser) return;
	localStorage.removeItem(ACCESS_TOKEN_KEY);
	localStorage.removeItem(REFRESH_TOKEN_KEY);
};

// API request function with automatic token handling
export async function apiRequest<T = unknown>(
	endpoint: string,
	options: RequestInit = {}
): Promise<T> {
	const url = `${API_BASE_URL}${endpoint}`;
	const accessToken = getAccessToken();

	// Set default headers
	const headers = new Headers(options.headers);
	headers.set('Content-Type', 'application/json');

	// Add authorization header if token exists
	if (accessToken) {
		headers.set('Authorization', `Bearer ${accessToken}`);
	}

	// Create request with headers
	const request: RequestInit = {
		...options,
		headers
	};

	try {
		const response = await fetch(url, request);

		// Handle 401 Unauthorized - attempt token refresh
		if (response.status === 401 && getRefreshToken()) {
			const newToken = await refreshAccessToken();
			if (newToken) {
				// Retry the request with new token
				headers.set('Authorization', `Bearer ${newToken}`);
				const retryResponse = await fetch(url, { ...request, headers });
				return handleResponse(retryResponse);
			}
		}

		return handleResponse(response);
	} catch (error) {
		console.error('API request failed:', error);
		throw new ApiError('Network error', 0);
	}
}

// Handle API response
async function handleResponse<T>(response: Response): Promise<T> {
	const data = await response.json().catch(() => ({}));

	if (!response.ok) {
		const message = data.error || data.message || 'Something went wrong';
		throw new ApiError(message, response.status);
	}

	return data as T;
}

// Refresh the access token
async function refreshAccessToken(): Promise<string | null> {
	const refreshToken = getRefreshToken();
	if (!refreshToken) return null;

	try {
		const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ refresh_token: refreshToken })
		});

		if (!response.ok) {
			clearTokens();
			return null;
		}

		const data = await response.json();
		setTokens(data.access_token, data.refresh_token);
		return data.access_token;
	} catch {
		clearTokens();
		return null;
	}
}
