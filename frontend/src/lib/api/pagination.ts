export type SortDirection = 'asc' | 'desc'
export type PaginationMode = 'offset' | 'cursor'

export interface PaginatedSort {
    field: string
    direction: SortDirection
}

export interface OffsetPaginationMetadata {
    page: number
    page_size: number
    total_items: number
    total_pages: number
}

export interface CursorPaginationMetadata {
    page_size: number
    next_cursor?: string
    prev_cursor?: string
    has_next: boolean
    has_prev: boolean
}

interface BasePaginatedResponse<T, TMode extends PaginationMode, TPagination> {
    items: T[]
    sort: PaginatedSort
    pagination_mode: TMode
    pagination: TPagination
    search?: string
}

export type OffsetPaginatedResponse<T> = BasePaginatedResponse<
    T,
    'offset',
    OffsetPaginationMetadata
>

export type CursorPaginatedResponse<T> = BasePaginatedResponse<
    T,
    'cursor',
    CursorPaginationMetadata
>

export type PaginatedResponse<T> = OffsetPaginatedResponse<T> | CursorPaginatedResponse<T>

export function isOffsetPaginatedResponse<T>(
    response: PaginatedResponse<T>
): response is OffsetPaginatedResponse<T> {
    return response.pagination_mode === 'offset'
}

export function isCursorPaginatedResponse<T>(
    response: PaginatedResponse<T>
): response is CursorPaginatedResponse<T> {
    return response.pagination_mode === 'cursor'
}
