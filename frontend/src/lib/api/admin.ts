import { apiRequest } from './client'
import type { PaginatedResponse, SortDirection } from './pagination'

export interface AdminUserRow {
    id: string
    identifier: string
    email: string
    display_name: string
    role: string
    active: boolean
    last_login: string
    created_at: string
}

export interface ListAdminUsersParams {
    page?: number
    page_size?: number
    search?: string
    sort?: string
    order?: SortDirection
}

function buildUsersQuery(params: ListAdminUsersParams) {
    const query = new URLSearchParams()

    if (params.page) {
        query.set('page', String(params.page))
    }

    if (params.page_size) {
        query.set('page_size', String(params.page_size))
    }

    if (params.search?.trim()) {
        query.set('search', params.search.trim())
    }

    if (params.sort) {
        query.set('sort', params.sort)
    }

    if (params.order) {
        query.set('order', params.order)
    }

    const search = query.toString()
    return search ? `?${search}` : ''
}

export const adminApi = {
    listUsers: async (
        params: ListAdminUsersParams = {}
    ): Promise<PaginatedResponse<AdminUserRow>> => {
        return apiRequest<PaginatedResponse<AdminUserRow>>(
            `/api/admin/users${buildUsersQuery(params)}`,
            {
                method: 'GET',
                requiresAuth: true
            }
        )
    }
}
