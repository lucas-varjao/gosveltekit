<script lang="ts">
    import { onMount } from 'svelte'
    import { resolve } from '$app/paths'
    import { adminApi, type AdminUserRow } from '$lib/api/admin'
    import { isCursorPaginatedResponse, type SortDirection } from '$lib/api/pagination'
    import AdminUsersTable from '$lib/components/data-table/admin-users-table.svelte'
    import PageHeader from '$lib/components/layout/page-header.svelte'
    import { Alert, AlertDescription } from '$lib/components/ui/alert'
    import { buttonVariants } from '$lib/components/ui/button'
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card'
    import { cn } from '$lib/utils'

    const pageSize = 10

    let users = $state<AdminUserRow[]>([])
    let searchValue = $state('')
    let debouncedSearchValue = $state('')
    let sortField = $state('created_at')
    let sortDirection = $state<SortDirection>('desc')
    let nextCursor = $state<string | undefined>(undefined)
    let prevCursor = $state<string | undefined>(undefined)
    let hasNext = $state(false)
    let hasPrev = $state(false)
    let requestedCursorDirection = $state<'next' | 'prev' | null>(null)
    let requestedCursor = $state('')
    let isLoading = $state(false)
    let errorMessage = $state('')
    let initialized = $state(false)
    let requestCounter = 0

    async function loadUsers(params: {
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
            const response = await adminApi.listUsers({
                pagination_mode: 'cursor',
                page_size: pageSize,
                search: params.search,
                sort: params.sortField,
                order: params.sortDirection,
                after: params.cursorDirection === 'next' ? params.cursor : undefined,
                before: params.cursorDirection === 'prev' ? params.cursor : undefined
            })

            if (!isCursorPaginatedResponse(response)) {
                throw new Error('Resposta de paginação inesperada para o modo cursor')
            }

            if (currentRequest !== requestCounter) {
                return
            }

            users = response.items
            nextCursor = response.pagination.next_cursor
            prevCursor = response.pagination.prev_cursor
            hasNext = response.pagination.has_next
            hasPrev = response.pagination.has_prev
        } catch (error) {
            if (currentRequest !== requestCounter) {
                return
            }

            errorMessage = error instanceof Error ? error.message : 'Failed to load admin data'
            users = []
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
        requestedCursorDirection = null
        requestedCursor = ''
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
        requestedCursorDirection = direction
        requestedCursor = cursor
    }

    function resetFilters() {
        searchValue = ''
        debouncedSearchValue = ''
        sortField = 'created_at'
        sortDirection = 'desc'
        resetCursorWindow()
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

        void loadUsers({
            search: debouncedSearchValue,
            sortField,
            sortDirection,
            cursorDirection: requestedCursorDirection,
            cursor: requestedCursor
        })
    })
</script>

<section class="page-shell">
    <PageHeader
        title="Cursor Pagination"
        description="Exemplo com cursores opacos, navegação bidirecional e sorts estáveis para cenários com maior volume de dados."
        eyebrow="Admin Reference"
    />

    <div class="mt-4 flex flex-wrap gap-3">
        <a href={resolve('/admin')} class={cn(buttonVariants({ variant: 'outline' }), 'w-fit')}>
            Voltar ao hub
        </a>
    </div>

    <Card class="surface-card mt-8">
        <CardHeader>
            <CardTitle>Contrato</CardTitle>
        </CardHeader>
        <CardContent class="space-y-3 text-sm text-slate-300">
            <p>
                Endpoint usado: <code>/api/admin/users</code> com
                <code>pagination_mode=cursor</code>, <code>page_size</code>,
                <code>search</code>, <code>sort</code>, <code>order</code>,
                <code>after</code> e <code>before</code>.
            </p>
            <p class="text-slate-400">
                Este exemplo limita o sorting a <code>created_at</code>, <code>email</code> e
                <code>identifier</code>, porque o cursor depende de uma ordenação estável com
                tie-breaker por <code>id</code>.
            </p>
        </CardContent>
    </Card>

    {#if errorMessage}
        <Alert variant="destructive" class="mt-8 border-red-500/60 bg-red-950/50 text-red-200">
            <AlertDescription>{errorMessage}</AlertDescription>
        </Alert>
    {/if}

    <AdminUsersTable
        mode="cursor"
        rows={users}
        {pageSize}
        {searchValue}
        {sortField}
        {sortDirection}
        {hasNext}
        {hasPrev}
        {nextCursor}
        {prevCursor}
        {isLoading}
        sortableFields={['identifier', 'email', 'created_at']}
        primaryIdentitySortField="identifier"
        statusText="Sem contagem total: o backend navega por cursores opacos."
        onSearchChange={handleSearchChange}
        onSortChange={handleSortChange}
        onCursorChange={handleCursorChange}
    />

    <div class="mt-4 flex justify-end">
        <button
            type="button"
            class={cn(buttonVariants({ variant: 'ghost', size: 'sm' }))}
            onclick={resetFilters}
        >
            Limpar filtros
        </button>
    </div>
</section>
