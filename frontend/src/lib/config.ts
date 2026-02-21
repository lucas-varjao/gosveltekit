// frontend/src/lib/config.ts

/**
 * Application configuration
 */

// API configuration
export const API_CONFIG = {
    // Base URL for API requests - use environment variable or default
    baseUrl: import.meta.env.VITE_API_URL || '',

    // Default request timeout in milliseconds
    timeout: 30000
}

// Feature flags
export const FEATURES = {
    enableRegistration: true,
    enablePasswordReset: true
}
