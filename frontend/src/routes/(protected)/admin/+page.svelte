<script lang="ts">
    import { onMount } from 'svelte'
    import { adminApi, type AdminUserRow } from '$lib/api/admin'
    import type { SortDirection } from '$lib/api/pagination'
    import AdminUsersTable from '$lib/components/data-table/admin-users-table.svelte'
    import PageHeader from '$lib/components/layout/page-header.svelte'
    import { Alert, AlertDescription } from '$lib/components/ui/alert'
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card'
    import { Button } from '$lib/components/ui/button'

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
                page: params.page,
                page_size: pageSize,
                search: params.search,
                sort: params.sortField,
                order: params.sortDirection
            })

            if (currentRequest !== requestCounter) {
                return
            }

            users = response.items
            totalItems = response.total_items
            totalPages = response.total_pages
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
        title="Admin"
        description="Exemplo oficial de Data Table server-side: o frontend controla a experiência e o backend Go controla paginação, filtro e sorting."
        eyebrow="Stack Direction"
    />

    <Card class="surface-card">
        <CardHeader>
            <CardTitle>Referência da stack</CardTitle>
        </CardHeader>
        <CardContent class="space-y-3 text-sm text-slate-300">
            <p>
                Esta página demonstra a direção recomendada para listagens administrativas no
                template: <strong class="text-white"
                    ><code>shadcn-svelte</code> + <code>TanStack Table</code></strong
                >
                no frontend e dados paginados vindos do backend Go.
            </p>
            <p class="text-slate-400">
                Endpoint usado: <code>/api/admin/users</code> com <code>page</code>,
                <code>page_size</code>, <code>search</code>, <code>sort</code> e
                <code>order</code>.
            </p>
        </CardContent>
    </Card>

    {#if errorMessage}
        <Alert variant="destructive" class="mt-8 border-red-500/60 bg-red-950/50 text-red-200">
            <AlertDescription>{errorMessage}</AlertDescription>
        </Alert>
    {/if}

    <AdminUsersTable
        rows={users}
        {page}
        {pageSize}
        {totalItems}
        {totalPages}
        {searchValue}
        {sortField}
        {sortDirection}
        {isLoading}
        onSearchChange={handleSearchChange}
        onPageChange={handlePageChange}
        onSortChange={handleSortChange}
    />

    <div class="mt-4 flex justify-end">
        <Button variant="ghost" size="sm" onclick={resetFilters}>Limpar filtros</Button>
    </div>
</section>
