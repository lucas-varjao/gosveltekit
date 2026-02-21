// frontend/src/lib/api/client.ts

import { browser } from '$app/environment'
import { goto } from '$app/navigation'
import { resolve } from '$app/paths'

/**
 * API Client for communicating with the backend
 * Handles session-based authentication, request formatting, and error handling
 */
// Base URL for API requests
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

interface ApiRequestOptions extends RequestInit {
    requiresAuth?: boolean
}

// Error types
export class ApiError extends Error {
    status: number

    constructor(message: string, status: number) {
        super(message)
        this.status = status
        this.name = 'ApiError'
    }
}

// API request function with automatic session handling
export async function apiRequest<T = unknown>(
    endpoint: string,
    options: ApiRequestOptions = {}
): Promise<T> {
    const { requiresAuth = false, ...requestOptions } = options
    const url = `${API_BASE_URL}${endpoint}`

    // Set default headers
    const headers = new Headers(requestOptions.headers)
    headers.set('Content-Type', 'application/json')

    // Create request with headers
    const request: RequestInit = {
        ...requestOptions,
        headers,
        credentials: 'include'
    }

    try {
        const response = await fetch(url, request)
        return handleResponse(response, requiresAuth)
    } catch (error) {
        if (error instanceof ApiError) throw error
        console.error('API request failed:', error)
        throw new ApiError('Network error', 0)
    }
}

// Handle API response
async function handleResponse<T>(response: Response, requiresAuth: boolean): Promise<T> {
    const data = await response.json().catch(() => ({}))

    if (!response.ok) {
        const message = data.error || data.message || 'Something went wrong'

        if (requiresAuth && browser) {
            if (response.status === 401) {
                void goto(resolve('/session-expired'))
            } else if (response.status === 403) {
                void goto(resolve('/forbidden'))
            }
        }

        throw new ApiError(message, response.status)
    }

    return data as T
}
