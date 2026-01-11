// frontend/src/lib/api/client.ts

/**
 * API Client for communicating with the backend
 * Handles session-based authentication, request formatting, and error handling
 */
import { browser } from '$app/environment';

// Base URL for API requests
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// Session storage key
const SESSION_ID_KEY = 'sessionId';

// Error types
export class ApiError extends Error {
	status: number;

	constructor(message: string, status: number) {
		super(message);
		this.status = status;
		this.name = 'ApiError';
	}
}

// Session management
export const getSessionId = (): string | null => {
	if (!browser) return null;
	return localStorage.getItem(SESSION_ID_KEY);
};

export const setSessionId = (sessionId: string): void => {
	if (!browser) return;
	localStorage.setItem(SESSION_ID_KEY, sessionId);
};

export const clearSession = (): void => {
	if (!browser) return;
	localStorage.removeItem(SESSION_ID_KEY);
};

// API request function with automatic session handling
export async function apiRequest<T = unknown>(
	endpoint: string,
	options: RequestInit = {}
): Promise<T> {
	const url = `${API_BASE_URL}${endpoint}`;
	const sessionId = getSessionId();

	// Set default headers
	const headers = new Headers(options.headers);
	headers.set('Content-Type', 'application/json');

	// Add authorization header if session exists
	if (sessionId) {
		headers.set('Authorization', `Bearer ${sessionId}`);
	}

	// Create request with headers
	const request: RequestInit = {
		...options,
		headers
	};

	try {
		const response = await fetch(url, request);

		// Handle 401 Unauthorized - session expired or invalid
		if (response.status === 401) {
			clearSession();
			// Optionally redirect to login
			if (browser) {
				window.location.href = '/login';
			}
			throw new ApiError('Session expired', 401);
		}

		return handleResponse(response);
	} catch (error) {
		if (error instanceof ApiError) throw error;
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
