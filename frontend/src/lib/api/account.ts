import { apiRequest } from './client'

export interface AccountProfile {
    id: string
    identifier: string
    email: string
    display_name: string
    role: string
    active: boolean
    first_name?: string
    last_name?: string
    email_verified: boolean
    last_login: string
    last_active: string
}

export interface UpdateProfileRequest {
    display_name?: string
    first_name?: string
    last_name?: string
}

export interface ChangePasswordRequest {
    current_password: string
    new_password: string
    confirm_password: string
}

export interface AccountSession {
    id: string
    created_at: string
    expires_at: string
    user_agent?: string
    ip?: string
    is_current: boolean
}

interface MessageResponse {
    message: string
}

export const accountApi = {
    getProfile: async (): Promise<AccountProfile> => {
        return apiRequest<AccountProfile>('/api/account/profile', {
            method: 'GET',
            requiresAuth: true
        })
    },

    updateProfile: async (data: UpdateProfileRequest): Promise<AccountProfile> => {
        return apiRequest<AccountProfile>('/api/account/profile', {
            method: 'PATCH',
            body: JSON.stringify(data),
            requiresAuth: true
        })
    },

    changePassword: async (data: ChangePasswordRequest): Promise<MessageResponse> => {
        return apiRequest<MessageResponse>('/api/account/change-password', {
            method: 'POST',
            body: JSON.stringify(data),
            requiresAuth: true
        })
    },

    listSessions: async (): Promise<AccountSession[]> => {
        return apiRequest<AccountSession[]>('/api/account/sessions', {
            method: 'GET',
            requiresAuth: true
        })
    },

    revokeSession: async (sessionID: string): Promise<MessageResponse> => {
        return apiRequest<MessageResponse>(`/api/account/sessions/${sessionID}`, {
            method: 'DELETE',
            requiresAuth: true
        })
    }
}
