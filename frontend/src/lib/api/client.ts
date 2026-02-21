// frontend/src/lib/api/client.ts

/**
 * API Client for communicating with the backend
 * Handles session-based authentication, request formatting, and error handling
 */
// Base URL for API requests
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// Error types
export class ApiError extends Error {
	status: number;

	constructor(message: string, status: number) {
		super(message);
		this.status = status;
		this.name = 'ApiError';
	}
}

// API request function with automatic session handling
export async function apiRequest<T = unknown>(
	endpoint: string,
	options: RequestInit = {}
): Promise<T> {
	const url = `${API_BASE_URL}${endpoint}`;

	// Set default headers
	const headers = new Headers(options.headers);
	headers.set('Content-Type', 'application/json');

	// Create request with headers
	const request: RequestInit = {
		...options,
		headers,
		credentials: 'include'
	};

	try {
		const response = await fetch(url, request);
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
