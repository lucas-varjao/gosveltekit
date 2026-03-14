import { apiRequest } from './client'
import type { PaginatedResponse, PaginationMode, SortDirection } from './pagination'

export interface MockPaginationItem {
    id: string
    title: string
    category: string
    priority: string
    active: boolean
    created_at: string
}

interface BaseMockPaginationParams {
    pagination_mode: PaginationMode
    page_size?: number
    search?: string
    sort?: string
    order?: SortDirection
}

export interface MockPaginationOffsetParams extends BaseMockPaginationParams {
    pagination_mode: 'offset'
    page?: number
}

export interface MockPaginationCursorParams extends BaseMockPaginationParams {
    pagination_mode: 'cursor'
    after?: string
    before?: string
}

export type MockPaginationParams = MockPaginationOffsetParams | MockPaginationCursorParams

function buildMockPaginationQuery(params: MockPaginationParams) {
    const query = new URLSearchParams()

    query.set('pagination_mode', params.pagination_mode)

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

    if (params.pagination_mode === 'offset' && params.page) {
        query.set('page', String(params.page))
    }

    if (params.pagination_mode === 'cursor' && params.after) {
        query.set('after', params.after)
    }

    if (params.pagination_mode === 'cursor' && params.before) {
        query.set('before', params.before)
    }

    const search = query.toString()
    return search ? `?${search}` : ''
}

export const mockPaginationApi = {
    listItems: async (
        params: MockPaginationParams
    ): Promise<PaginatedResponse<MockPaginationItem>> => {
        return apiRequest<PaginatedResponse<MockPaginationItem>>(
            `/api/examples/pagination/items${buildMockPaginationQuery(params)}`,
            {
                method: 'GET',
                requiresAuth: true
            }
        )
    }
}
