<script lang="ts">
    import { onMount } from 'svelte'
    import { resolve } from '$app/paths'
    import { adminApi, type AdminUserRow } from '$lib/api/admin'
    import { isOffsetPaginatedResponse, type SortDirection } from '$lib/api/pagination'
    import AdminUsersTable from '$lib/components/data-table/admin-users-table.svelte'
    import PageHeader from '$lib/components/layout/page-header.svelte'
    import { Alert, AlertDescription } from '$lib/components/ui/alert'
    import { buttonVariants } from '$lib/components/ui/button'
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card'
    import { cn } from '$lib/utils'

    const pageSize = 10

    let users = $state<AdminUserRow[]>([])
    let page = $state(1)
    let totalItems = $state(0)
    let totalPages = $state(1)
    let searchValue = $state('')
    let debouncedSearchValue = $state('')
    let sortField = $state('created_at')
    let sortDirection = $state<SortDirection>('desc')
    let isLoading = $state(false)
    let errorMessage = $state('')
    let initialized = $state(false)
    let requestCounter = 0

    async function loadUsers(params: {
        page: number
        search: string
        sortField: string
        sortDirection: SortDirection
    }) {
        const currentRequest = ++requestCounter
        isLoading = true
        errorMessage = ''

        try {
            const response = await adminApi.listUsers({
                pagination_mode: 'offset',
                page: params.page,
                page_size: pageSize,
                search: params.search,
                sort: params.sortField,
                order: params.sortDirection
            })

            if (!isOffsetPaginatedResponse(response)) {
                throw new Error('Resposta de paginação inesperada para o modo offset')
            }

            if (currentRequest !== requestCounter) {
                return
            }

            users = response.items
            totalItems = response.pagination.total_items
            totalPages = response.pagination.total_pages
        } catch (error) {
            if (currentRequest !== requestCounter) {
                return
            }

            errorMessage = error instanceof Error ? error.message : 'Failed to load admin data'
            users = []
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

    function resetFilters() {
        searchValue = ''
        debouncedSearchValue = ''
        sortField = 'created_at'
        sortDirection = 'desc'
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

        void loadUsers({
            page,
            search: debouncedSearchValue,
            sortField,
            sortDirection
        })
    })
</script>

<section class="page-shell">
    <PageHeader
        title="Offset Pagination"
        description="Exemplo com total de páginas, total de registros e navegação numérica controlados pelo backend Go."
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
                <code>pagination_mode=offset</code>, <code>page</code>,
                <code>page_size</code>, <code>search</code>, <code>sort</code> e
                <code>order</code>.
            </p>
            <p class="text-slate-400">
                Use esta estratégia quando a experiência depende de contagem total e navegação por
                páginas numeradas.
            </p>
        </CardContent>
    </Card>

    {#if errorMessage}
        <Alert variant="destructive" class="mt-8 border-red-500/60 bg-red-950/50 text-red-200">
            <AlertDescription>{errorMessage}</AlertDescription>
        </Alert>
    {/if}

    <AdminUsersTable
        mode="offset"
        rows={users}
        {page}
        {pageSize}
        {totalItems}
        {totalPages}
        {searchValue}
        {sortField}
        {sortDirection}
        {isLoading}
        sortableFields={['display_name', 'email', 'role', 'created_at', 'last_login']}
        primaryIdentitySortField="display_name"
        onSearchChange={handleSearchChange}
        onPageChange={handlePageChange}
        onSortChange={handleSortChange}
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
