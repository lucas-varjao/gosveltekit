<script lang="ts">
    import { onMount } from 'svelte'
    import { mockPaginationApi, type MockPaginationItem } from '$lib/api/mock-pagination'
    import { isOffsetPaginatedResponse, type SortDirection } from '$lib/api/pagination'
    import MockPaginationTable from '$lib/components/data-table/mock-pagination-table.svelte'
    import { Alert, AlertDescription } from '$lib/components/ui/alert'

    const pageSize = 8

    let initialized = $state(false)
    let requestCounter = 0

    let rows = $state<MockPaginationItem[]>([])
    let page = $state(1)
    let totalItems = $state(0)
    let totalPages = $state(1)
    let searchValue = $state('')
    let debouncedSearchValue = $state('')
    let sortField = $state('created_at')
    let sortDirection = $state<SortDirection>('desc')
    let isLoading = $state(false)
    let errorMessage = $state('')

    async function loadItems(params: {
        page: number
        search: string
        sortField: string
        sortDirection: SortDirection
    }) {
        const currentRequest = ++requestCounter
        isLoading = true
        errorMessage = ''

        try {
            const response = await mockPaginationApi.listItems({
                pagination_mode: 'offset',
                page: params.page,
                page_size: pageSize,
                search: params.search,
                sort: params.sortField,
                order: params.sortDirection
            })

            if (!isOffsetPaginatedResponse(response)) {
                throw new Error('Resposta inesperada para o modo offset')
            }

            if (currentRequest !== requestCounter) {
                return
            }

            rows = response.items
            totalItems = response.pagination.total_items
            totalPages = response.pagination.total_pages
        } catch (error) {
            if (currentRequest !== requestCounter) {
                return
            }

            errorMessage =
                error instanceof Error ? error.message : 'Falha ao carregar o exemplo offset'
            rows = []
            totalItems = 0
            totalPages = 1
        } finally {
            if (currentRequest === requestCounter) {
                isLoading = false
            }
        }
    }

    function handleSearchChange(value: string) {
        searchValue = value
        page = 1
    }

    function handlePageChange(nextPage: number) {
        page = nextPage
    }

    function handleSortChange(field: string, direction: SortDirection) {
        sortField = field
        sortDirection = direction
        page = 1
    }

    onMount(() => {
        initialized = true
    })

    $effect(() => {
        if (!initialized) {
            return
        }

        const currentValue = searchValue
        const timeout = window.setTimeout(() => {
            debouncedSearchValue = currentValue.trim()
        }, 300)

        return () => window.clearTimeout(timeout)
    })

    $effect(() => {
        if (!initialized || searchValue.trim() !== debouncedSearchValue) {
            return
        }

        void loadItems({
            page,
            search: debouncedSearchValue,
            sortField,
            sortDirection
        })
    })
</script>

<section class="rounded-2xl border border-slate-800/80 bg-slate-950/40 p-4 sm:p-6">
    <div class="space-y-2">
        <h2 class="text-xl font-semibold text-white">Offset / Limit</h2>
        <p class="text-sm text-slate-400">
            Mock backend com <code>pagination_mode=offset</code>, ideal para experiências com total
            de registros e páginas numeradas.
        </p>
    </div>

    {#if errorMessage}
        <Alert variant="destructive" class="mt-6 border-red-500/60 bg-red-950/50 text-red-200">
            <AlertDescription>{errorMessage}</AlertDescription>
        </Alert>
    {/if}

    <MockPaginationTable
        mode="offset"
        {rows}
        {page}
        {pageSize}
        {totalItems}
        {totalPages}
        {searchValue}
        {sortField}
        {sortDirection}
        {isLoading}
        sortableFields={['title', 'category', 'priority', 'created_at']}
        onSearchChange={handleSearchChange}
        onPageChange={handlePageChange}
        onSortChange={handleSortChange}
    />
</section>
