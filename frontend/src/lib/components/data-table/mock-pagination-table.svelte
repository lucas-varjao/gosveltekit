<script lang="ts">
    import { ArrowDown, ArrowUp, ArrowUpDown } from '@lucide/svelte'
    import {
        FlexRender,
        createTable,
        getCoreRowModel,
        renderComponent,
        type ColumnDef,
        type SortingState
    } from '@tanstack/svelte-table'
    import type { MockPaginationItem } from '$lib/api/mock-pagination'
    import type { PaginationMode, SortDirection } from '$lib/api/pagination'
    import DateCell from '$lib/components/data-table/cells/date-cell.svelte'
    import StatusCell from '$lib/components/data-table/cells/status-cell.svelte'
    import {
        Table,
        TableBody,
        TableCell,
        TableHead,
        TableHeader,
        TableRow
    } from '$lib/components/ui/table'
    import DataTableCursorPagination from './data-table-cursor-pagination.svelte'
    import DataTablePagination from './data-table-pagination.svelte'
    import DataTableToolbar from './data-table-toolbar.svelte'

    type Props = {
        mode: PaginationMode
        rows: MockPaginationItem[]
        pageSize: number
        searchValue: string
        sortField: string
        sortDirection: SortDirection
        sortableFields?: string[]
        page?: number
        totalItems?: number
        totalPages?: number
        hasNext?: boolean
        hasPrev?: boolean
        nextCursor?: string
        prevCursor?: string
        statusText?: string
        isLoading?: boolean
        onPageChange?: (page: number) => void
        onCursorChange?: (direction: 'next' | 'prev', cursor: string) => void
        onSearchChange?: (value: string) => void
        onSortChange?: (field: string, direction: SortDirection) => void
    }

    let {
        mode,
        rows,
        pageSize,
        searchValue,
        sortField,
        sortDirection,
        sortableFields = ['title', 'category', 'priority', 'created_at'],
        page = 1,
        totalItems,
        totalPages = 1,
        hasNext = false,
        hasPrev = false,
        nextCursor,
        prevCursor,
        statusText = '',
        isLoading = false,
        onPageChange,
        onCursorChange,
        onSearchChange,
        onSortChange
    }: Props = $props()

    let columns = $derived.by(() => {
        const sortableFieldSet = new Set(sortableFields)

        return [
            {
                accessorKey: 'title',
                header: 'Item',
                enableSorting: sortableFieldSet.has('title'),
                cell: (info) => String(info.getValue() ?? '')
            },
            {
                accessorKey: 'category',
                header: 'Categoria',
                enableSorting: sortableFieldSet.has('category'),
                cell: (info) => String(info.getValue() ?? '')
            },
            {
                accessorKey: 'priority',
                header: 'Prioridade',
                enableSorting: sortableFieldSet.has('priority'),
                cell: (info) => String(info.getValue() ?? '')
            },
            {
                accessorKey: 'active',
                header: 'Status',
                enableSorting: false,
                cell: (info) => renderComponent(StatusCell, { active: Boolean(info.getValue()) })
            },
            {
                accessorKey: 'created_at',
                header: 'Criado em',
                enableSorting: sortableFieldSet.has('created_at'),
                cell: (info) => renderComponent(DateCell, { value: String(info.getValue() ?? '') })
            }
        ] satisfies ColumnDef<MockPaginationItem>[]
    })

    let sorting = $derived<SortingState>(
        sortField ? [{ id: sortField, desc: sortDirection === 'desc' }] : []
    )

    let table = $derived(
        createTable({
            data: rows,
            columns,
            getCoreRowModel: getCoreRowModel(),
            manualPagination: true,
            manualSorting: true,
            pageCount: totalPages,
            renderFallbackValue: '',
            state: {
                sorting,
                pagination: {
                    pageIndex: page - 1,
                    pageSize
                }
            },
            onStateChange: () => {},
            onSortingChange: (updater) => {
                const nextSorting = typeof updater === 'function' ? updater(sorting) : updater
                const nextSort = nextSorting[0]

                if (!nextSort) {
                    return
                }

                onSortChange?.(String(nextSort.id), nextSort.desc ? 'desc' : 'asc')
            }
        })
    )

    function getSortIconState(columnId: string) {
        if (columnId !== sortField) {
            return 'none'
        }

        return sortDirection
    }
</script>

<div class="surface-card mt-6 overflow-hidden">
    <DataTableToolbar {searchValue} {totalItems} {statusText} {isLoading} {onSearchChange} />

    <div class="border-y border-slate-800/70">
        <Table>
            <TableHeader>
                {#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
                    <TableRow class="hover:bg-transparent">
                        {#each headerGroup.headers as header (header.id)}
                            <TableHead>
                                {#if header.isPlaceholder}
                                    <span></span>
                                {:else if header.column.getCanSort()}
                                    <button
                                        type="button"
                                        class="inline-flex items-center gap-2 text-left text-inherit"
                                        onclick={header.column.getToggleSortingHandler()}
                                    >
                                        <FlexRender
                                            content={header.column.columnDef.header}
                                            context={header.getContext()}
                                        />

                                        {#if getSortIconState(header.column.id) === 'asc'}
                                            <ArrowUp class="size-3.5 text-cyan-300" />
                                        {:else if getSortIconState(header.column.id) === 'desc'}
                                            <ArrowDown class="size-3.5 text-cyan-300" />
                                        {:else}
                                            <ArrowUpDown class="size-3.5 text-slate-500" />
                                        {/if}
                                    </button>
                                {:else}
                                    <FlexRender
                                        content={header.column.columnDef.header}
                                        context={header.getContext()}
                                    />
                                {/if}
                            </TableHead>
                        {/each}
                    </TableRow>
                {/each}
            </TableHeader>

            <TableBody>
                {#if isLoading}
                    <TableRow class="hover:bg-transparent">
                        <TableCell
                            colspan={columns.length}
                            class="py-10 text-center text-slate-400"
                        >
                            Carregando dados mockados do backend...
                        </TableCell>
                    </TableRow>
                {:else if rows.length === 0}
                    <TableRow class="hover:bg-transparent">
                        <TableCell
                            colspan={columns.length}
                            class="py-10 text-center text-slate-400"
                        >
                            Nenhum item encontrado para os filtros atuais.
                        </TableCell>
                    </TableRow>
                {:else}
                    {#each table.getRowModel().rows as row (row.id)}
                        <TableRow>
                            {#each row.getVisibleCells() as cell (cell.id)}
                                <TableCell>
                                    <FlexRender
                                        content={cell.column.columnDef.cell}
                                        context={cell.getContext()}
                                    />
                                </TableCell>
                            {/each}
                        </TableRow>
                    {/each}
                {/if}
            </TableBody>
        </Table>
    </div>

    {#if mode === 'offset' && totalItems !== undefined}
        <DataTablePagination {page} {pageSize} {totalItems} {totalPages} {onPageChange} />
    {:else if mode === 'cursor'}
        <DataTableCursorPagination
            {pageSize}
            {hasNext}
            {hasPrev}
            {nextCursor}
            {prevCursor}
            {onCursorChange}
        />
    {/if}
</div>
