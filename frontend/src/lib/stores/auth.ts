// frontend/src/lib/stores/auth.ts

/**
 * Authentication store for managing user state
 * Uses session-based authentication
 */
import { writable } from 'svelte/store'
import { browser } from '$app/environment'
import { authApi } from '$lib/api/auth'
import { setUnauthorizedHandler } from '$lib/api/client'

// User interface matching backend response from /api/me
export interface User {
    id: string
    identifier: string // username
    email: string
    display_name: string
    role: string
    active: boolean
}

interface AuthState {
    user: User | null
    isAuthenticated: boolean
    isLoading: boolean
    error: string | null
}

interface ValidateSessionOptions {
    silent?: boolean
}

const initialState: AuthState = {
    user: null,
    isAuthenticated: false,
    isLoading: true,
    error: null
}

function createAuthStore() {
    const { subscribe, update } = writable<AuthState>(initialState)
    let validationRequestID = 0

    function cancelPendingValidation() {
        validationRequestID += 1
    }

    function setAuthenticated(user: User) {
        update((state) => ({
            ...state,
            user,
            isAuthenticated: true,
            isLoading: false,
            error: null
        }))
    }

    function setUnauthenticated(error: string | null = null) {
        update((state) => ({
            ...state,
            user: null,
            isAuthenticated: false,
            isLoading: false,
            error
        }))
    }

    async function validateSession(options: ValidateSessionOptions = {}) {
        if (!browser) return false

        const { silent = false } = options
        const requestID = ++validationRequestID

        if (!silent) {
            update((state) => ({ ...state, isLoading: true, error: null }))
        }

        try {
            const user = await authApi.getCurrentUser()

            if (requestID !== validationRequestID) {
                return true
            }

            setAuthenticated(user)
            return true
        } catch {
            if (requestID !== validationRequestID) {
                return false
            }

            setUnauthenticated()
            return false
        }
    }

    return {
        subscribe,

        // Initialize auth state by validating server-side session cookie.
        init: async () => {
            return validateSession()
        },

        refreshSession: async () => {
            return validateSession({ silent: true })
        },

        invalidateSession: () => {
            cancelPendingValidation()
            setUnauthenticated()
        },

        login: async (username: string, password: string) => {
            cancelPendingValidation()
            update((state) => ({ ...state, isLoading: true, error: null }))

            try {
                const response = await authApi.login({ username, password })

                setAuthenticated(response.user)

                return response.user
            } catch (error) {
                const message = error instanceof Error ? error.message : 'Login failed'
                update((state) => ({
                    ...state,
                    error: message,
                    isLoading: false
                }))
                throw error
            }
        },

        // Registration creates user but does not authenticate.
        register: async (data: {
            username: string
            email: string
            password: string
            display_name: string
        }) => {
            cancelPendingValidation()
            update((state) => ({ ...state, isLoading: true, error: null }))

            try {
                const registeredUser = await authApi.register({
                    username: data.username,
                    email: data.email,
                    password: data.password,
                    display_name: data.display_name
                })

                update((state) => ({
                    ...state,
                    isLoading: false,
                    isAuthenticated: false,
                    user: null
                }))

                return registeredUser
            } catch (error) {
                const message = error instanceof Error ? error.message : 'Registration failed'
                update((state) => ({
                    ...state,
                    error: message,
                    isLoading: false
                }))
                throw error
            }
        },

        logout: async () => {
            cancelPendingValidation()
            update((state) => ({ ...state, isLoading: true }))

            try {
                await authApi.logout()
            } finally {
                setUnauthenticated()
            }
        },

        clearError: () => {
            update((state) => ({ ...state, error: null }))
        }
    }
}

export const authStore = createAuthStore()
setUnauthorizedHandler(() => authStore.invalidateSession())
