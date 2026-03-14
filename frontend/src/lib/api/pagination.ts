export type SortDirection = 'asc' | 'desc'

export interface PaginatedSort {
    field: string
    direction: SortDirection
}

export interface PaginatedResponse<T> {
    items: T[]
    page: number
    page_size: number
    total_items: number
    total_pages: number
    sort: PaginatedSort
    search?: string
}
