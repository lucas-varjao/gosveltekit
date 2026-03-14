<script lang="ts">
    import { onMount } from 'svelte'
    import { mockPaginationApi, type MockPaginationItem } from '$lib/api/mock-pagination'
    import {
        isCursorPaginatedResponse,
        isOffsetPaginatedResponse,
        type SortDirection
    } from '$lib/api/pagination'
    import MockPaginationTable from '$lib/components/data-table/mock-pagination-table.svelte'
    import PageHeader from '$lib/components/layout/page-header.svelte'
    import { Alert, AlertDescription } from '$lib/components/ui/alert'
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card'

    const offsetPageSize = 8
    const cursorPageSize = 8

    let initialized = $state(false)
    let requestCounter = 0

    let offsetRows = $state<MockPaginationItem[]>([])
    let offsetPage = $state(1)
    let offsetTotalItems = $state(0)
    let offsetTotalPages = $state(1)
    let offsetSearchValue = $state('')
    let offsetDebouncedSearchValue = $state('')
    let offsetSortField = $state('created_at')
    let offsetSortDirection = $state<SortDirection>('desc')
    let offsetIsLoading = $state(false)
    let offsetErrorMessage = $state('')

    let cursorRows = $state<MockPaginationItem[]>([])
    let cursorSearchValue = $state('')
    let cursorDebouncedSearchValue = $state('')
    let cursorSortField = $state('created_at')
    let cursorSortDirection = $state<SortDirection>('desc')
    let cursorNextCursor = $state<string | undefined>(undefined)
    let cursorPrevCursor = $state<string | undefined>(undefined)
    let cursorHasNext = $state(false)
    let cursorHasPrev = $state(false)
    let cursorRequestedDirection = $state<'next' | 'prev' | null>(null)
    let cursorRequestedValue = $state('')
    let cursorIsLoading = $state(false)
    let cursorErrorMessage = $state('')

    async function loadOffsetItems(params: {
        page: number
        search: string
        sortField: string
        sortDirection: SortDirection
    }) {
        const currentRequest = ++requestCounter
        offsetIsLoading = true
        offsetErrorMessage = ''

        try {
            const response = await mockPaginationApi.listItems({
                pagination_mode: 'offset',
                page: params.page,
                page_size: offsetPageSize,
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

            offsetRows = response.items
            offsetTotalItems = response.pagination.total_items
            offsetTotalPages = response.pagination.total_pages
        } catch (error) {
            if (currentRequest !== requestCounter) {
                return
            }

            offsetErrorMessage =
                error instanceof Error ? error.message : 'Falha ao carregar o exemplo offset'
            offsetRows = []
            offsetTotalItems = 0
            offsetTotalPages = 1
        } finally {
            if (currentRequest === requestCounter) {
                offsetIsLoading = false
            }
        }
    }

    async function loadCursorItems(params: {
        search: string
        sortField: string
        sortDirection: SortDirection
        cursorDirection: 'next' | 'prev' | null
        cursor: string
    }) {
        const currentRequest = ++requestCounter
        cursorIsLoading = true
        cursorErrorMessage = ''

        try {
            const response = await mockPaginationApi.listItems({
                pagination_mode: 'cursor',
                page_size: cursorPageSize,
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

            cursorRows = response.items
            cursorNextCursor = response.pagination.next_cursor
            cursorPrevCursor = response.pagination.prev_cursor
            cursorHasNext = response.pagination.has_next
            cursorHasPrev = response.pagination.has_prev
        } catch (error) {
            if (currentRequest !== requestCounter) {
                return
            }

            cursorErrorMessage =
                error instanceof Error ? error.message : 'Falha ao carregar o exemplo cursor'
            cursorRows = []
            cursorNextCursor = undefined
            cursorPrevCursor = undefined
            cursorHasNext = false
            cursorHasPrev = false
        } finally {
            if (currentRequest === requestCounter) {
                cursorIsLoading = false
            }
        }
    }

    function handleOffsetSearchChange(value: string) {
        offsetSearchValue = value
        offsetPage = 1
    }

    function handleOffsetPageChange(nextPage: number) {
        offsetPage = nextPage
    }

    function handleOffsetSortChange(field: string, direction: SortDirection) {
        offsetSortField = field
        offsetSortDirection = direction
        offsetPage = 1
    }

    function resetCursorWindow() {
        cursorRequestedDirection = null
        cursorRequestedValue = ''
        cursorNextCursor = undefined
        cursorPrevCursor = undefined
        cursorHasNext = false
        cursorHasPrev = false
    }

    function handleCursorSearchChange(value: string) {
        cursorSearchValue = value
        resetCursorWindow()
    }

    function handleCursorSortChange(field: string, direction: SortDirection) {
        cursorSortField = field
        cursorSortDirection = direction
        resetCursorWindow()
    }

    function handleCursorChange(direction: 'next' | 'prev', cursor: string) {
        cursorRequestedDirection = direction
        cursorRequestedValue = cursor
    }

    onMount(() => {
        initialized = true
    })

    $effect(() => {
        if (!initialized) {
            return
        }

        const currentValue = offsetSearchValue
        const timeout = window.setTimeout(() => {
            offsetDebouncedSearchValue = currentValue.trim()
        }, 300)

        return () => window.clearTimeout(timeout)
    })

    $effect(() => {
        if (!initialized) {
            return
        }

        const currentValue = cursorSearchValue
        const timeout = window.setTimeout(() => {
            cursorDebouncedSearchValue = currentValue.trim()
        }, 300)

        return () => window.clearTimeout(timeout)
    })

    $effect(() => {
        if (!initialized || offsetSearchValue.trim() !== offsetDebouncedSearchValue) {
            return
        }

        void loadOffsetItems({
            page: offsetPage,
            search: offsetDebouncedSearchValue,
            sortField: offsetSortField,
            sortDirection: offsetSortDirection
        })
    })

    $effect(() => {
        if (!initialized || cursorSearchValue.trim() !== cursorDebouncedSearchValue) {
            return
        }

        void loadCursorItems({
            search: cursorDebouncedSearchValue,
            sortField: cursorSortField,
            sortDirection: cursorSortDirection,
            cursorDirection: cursorRequestedDirection,
            cursor: cursorRequestedValue
        })
    })
</script>

<section class="page-shell">
    <PageHeader
        title="Pagination Playground"
        description="Exemplo independente de admin, com dados mockados no backend e as duas estratégias de paginação na mesma tela."
        eyebrow="Examples"
    />

    <Card class="surface-card mt-8">
        <CardHeader>
            <CardTitle>Backend mockado</CardTitle>
        </CardHeader>
        <CardContent class="space-y-3 text-sm text-slate-300">
            <p>
                Esta página usa o endpoint <code>/api/examples/pagination/items</code>, que gera
                dados mockados no backend e aplica filtro, ordenação e paginação sem depender do
                banco.
            </p>
            <p class="text-slate-400">
                O bloco offset expõe totais e navegação numerada. O bloco cursor usa cursores opacos
                e navegação bidirecional.
            </p>
        </CardContent>
    </Card>

    <div class="mt-8 grid gap-8 xl:grid-cols-2">
        <section>
            <div class="space-y-2">
                <h2 class="text-xl font-semibold text-white">Offset / Limit</h2>
                <p class="text-sm text-slate-400">
                    Mock backend com <code>pagination_mode=offset</code>, ideal para experiências
                    com total de registros e páginas numeradas.
                </p>
            </div>

            {#if offsetErrorMessage}
                <Alert
                    variant="destructive"
                    class="mt-6 border-red-500/60 bg-red-950/50 text-red-200"
                >
                    <AlertDescription>{offsetErrorMessage}</AlertDescription>
                </Alert>
            {/if}

            <MockPaginationTable
                mode="offset"
                rows={offsetRows}
                page={offsetPage}
                pageSize={offsetPageSize}
                totalItems={offsetTotalItems}
                totalPages={offsetTotalPages}
                searchValue={offsetSearchValue}
                sortField={offsetSortField}
                sortDirection={offsetSortDirection}
                isLoading={offsetIsLoading}
                sortableFields={['title', 'category', 'priority', 'created_at']}
                onSearchChange={handleOffsetSearchChange}
                onPageChange={handleOffsetPageChange}
                onSortChange={handleOffsetSortChange}
            />
        </section>

        <section>
            <div class="space-y-2">
                <h2 class="text-xl font-semibold text-white">Cursor-Based</h2>
                <p class="text-sm text-slate-400">
                    Mock backend com <code>pagination_mode=cursor</code>, sem totais e com sorts
                    estáveis limitados a <code>title</code>, <code>category</code> e
                    <code>created_at</code>.
                </p>
            </div>

            {#if cursorErrorMessage}
                <Alert
                    variant="destructive"
                    class="mt-6 border-red-500/60 bg-red-950/50 text-red-200"
                >
                    <AlertDescription>{cursorErrorMessage}</AlertDescription>
                </Alert>
            {/if}

            <MockPaginationTable
                mode="cursor"
                rows={cursorRows}
                pageSize={cursorPageSize}
                searchValue={cursorSearchValue}
                sortField={cursorSortField}
                sortDirection={cursorSortDirection}
                hasNext={cursorHasNext}
                hasPrev={cursorHasPrev}
                nextCursor={cursorNextCursor}
                prevCursor={cursorPrevCursor}
                isLoading={cursorIsLoading}
                sortableFields={['title', 'category', 'created_at']}
                statusText="Cursores opacos vindos do backend mockado, sem contagem total."
                onSearchChange={handleCursorSearchChange}
                onSortChange={handleCursorSortChange}
                onCursorChange={handleCursorChange}
            />
        </section>
    </div>
</section>
