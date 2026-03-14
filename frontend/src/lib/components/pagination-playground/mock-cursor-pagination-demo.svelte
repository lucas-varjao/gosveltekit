<script lang="ts">
    import { onMount } from 'svelte'
    import { mockPaginationApi, type MockPaginationItem } from '$lib/api/mock-pagination'
    import { isCursorPaginatedResponse, type SortDirection } from '$lib/api/pagination'
    import MockPaginationTable from '$lib/components/data-table/mock-pagination-table.svelte'
    import { Alert, AlertDescription } from '$lib/components/ui/alert'

    const pageSize = 8

    let initialized = $state(false)
    let requestCounter = 0

    let rows = $state<MockPaginationItem[]>([])
    let searchValue = $state('')
    let debouncedSearchValue = $state('')
    let sortField = $state('created_at')
    let sortDirection = $state<SortDirection>('desc')
    let nextCursor = $state<string | undefined>(undefined)
    let prevCursor = $state<string | undefined>(undefined)
    let hasNext = $state(false)
    let hasPrev = $state(false)
    let requestedDirection = $state<'next' | 'prev' | null>(null)
    let requestedValue = $state('')
    let isLoading = $state(false)
    let errorMessage = $state('')

    async function loadItems(params: {
        search: string
        sortField: string
        sortDirection: SortDirection
        cursorDirection: 'next' | 'prev' | null
        cursor: string
    }) {
        const currentRequest = ++requestCounter
        isLoading = true
        errorMessage = ''

        try {
            const response = await mockPaginationApi.listItems({
                pagination_mode: 'cursor',
                page_size: pageSize,
                search: params.search,
                sort: params.sortField,
                order: params.sortDirection,
                after: params.cursorDirection === 'next' ? params.cursor : undefined,
                before: params.cursorDirection === 'prev' ? params.cursor : undefined
            })

            if (!isCursorPaginatedResponse(response)) {
                throw new Error('Resposta inesperada para o modo cursor')
            }

            if (currentRequest !== requestCounter) {
                return
            }

            rows = response.items
            nextCursor = response.pagination.next_cursor
            prevCursor = response.pagination.prev_cursor
            hasNext = response.pagination.has_next
            hasPrev = response.pagination.has_prev
        } catch (error) {
            if (currentRequest !== requestCounter) {
                return
            }

            errorMessage =
                error instanceof Error ? error.message : 'Falha ao carregar o exemplo cursor'
            rows = []
            nextCursor = undefined
            prevCursor = undefined
            hasNext = false
            hasPrev = false
        } finally {
            if (currentRequest === requestCounter) {
                isLoading = false
            }
        }
    }

    function resetCursorWindow() {
        requestedDirection = null
        requestedValue = ''
        nextCursor = undefined
        prevCursor = undefined
        hasNext = false
        hasPrev = false
    }

    function handleSearchChange(value: string) {
        searchValue = value
        resetCursorWindow()
    }

    function handleSortChange(field: string, direction: SortDirection) {
        sortField = field
        sortDirection = direction
        resetCursorWindow()
    }

    function handleCursorChange(direction: 'next' | 'prev', cursor: string) {
        requestedDirection = direction
        requestedValue = cursor
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
            search: debouncedSearchValue,
            sortField,
            sortDirection,
            cursorDirection: requestedDirection,
            cursor: requestedValue
        })
    })
</script>

<section class="rounded-2xl border border-slate-800/80 bg-slate-950/40 p-4 sm:p-6">
    <div class="space-y-2">
        <h2 class="text-xl font-semibold text-white">Cursor-Based</h2>
        <p class="text-sm text-slate-400">
            Mock backend com <code>pagination_mode=cursor</code>, sem totais e com sorts estaveis
            limitados a <code>title</code>, <code>category</code> e <code>created_at</code>.
        </p>
    </div>

    {#if errorMessage}
        <Alert variant="destructive" class="mt-6 border-red-500/60 bg-red-950/50 text-red-200">
            <AlertDescription>{errorMessage}</AlertDescription>
        </Alert>
    {/if}

    <MockPaginationTable
        mode="cursor"
        {rows}
        {pageSize}
        {searchValue}
        {sortField}
        {sortDirection}
        {hasNext}
        {hasPrev}
        {nextCursor}
        {prevCursor}
        {isLoading}
        sortableFields={['title', 'category', 'created_at']}
        statusText="Cursores opacos vindos do backend mockado, sem contagem total."
        onSearchChange={handleSearchChange}
        onSortChange={handleSortChange}
        onCursorChange={handleCursorChange}
    />
</section>
