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
    import type { AdminUserRow } from '$lib/api/admin'
    import type { SortDirection } from '$lib/api/pagination'
    import DateCell from '$lib/components/data-table/cells/date-cell.svelte'
    import EmailCell from '$lib/components/data-table/cells/email-cell.svelte'
    import NameCell from '$lib/components/data-table/cells/name-cell.svelte'
    import RoleCell from '$lib/components/data-table/cells/role-cell.svelte'
    import StatusCell from '$lib/components/data-table/cells/status-cell.svelte'
    import {
        Table,
        TableBody,
        TableCell,
        TableHead,
        TableHeader,
        TableRow
    } from '$lib/components/ui/table'
    import DataTablePagination from './data-table-pagination.svelte'
    import DataTableToolbar from './data-table-toolbar.svelte'

    type Props = {
        rows: AdminUserRow[]
        page: number
        pageSize: number
        totalItems: number
        totalPages: number
        searchValue: string
        sortField: string
        sortDirection: SortDirection
        isLoading?: boolean
        onPageChange?: (page: number) => void
        onSearchChange?: (value: string) => void
        onSortChange?: (field: string, direction: SortDirection) => void
    }

    let {
        rows,
        page,
        pageSize,
        totalItems,
        totalPages,
        searchValue,
        sortField,
        sortDirection,
        isLoading = false,
        onPageChange,
        onSearchChange,
        onSortChange
    }: Props = $props()

    const columns: ColumnDef<AdminUserRow>[] = [
        {
            accessorKey: 'display_name',
            header: 'Usuário',
            enableSorting: true,
            cell: (info) =>
                renderComponent(NameCell, {
                    displayName: info.row.original.display_name,
                    identifier: info.row.original.identifier
                })
        },
        {
            accessorKey: 'email',
            header: 'Email',
            enableSorting: true,
            cell: (info) => renderComponent(EmailCell, { email: String(info.getValue() ?? '') })
        },
        {
            accessorKey: 'role',
            header: 'Papel',
            enableSorting: true,
            cell: (info) => renderComponent(RoleCell, { role: String(info.getValue() ?? 'user') })
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
            enableSorting: true,
            cell: (info) => renderComponent(DateCell, { value: String(info.getValue() ?? '') })
        },
        {
            accessorKey: 'last_login',
            header: 'Último login',
            enableSorting: true,
            cell: (info) => renderComponent(DateCell, { value: String(info.getValue() ?? '') })
        }
    ]

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

<div class="surface-card mt-8 overflow-hidden">
    <DataTableToolbar {searchValue} {totalItems} {isLoading} {onSearchChange} />

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
                            Carregando usuários do backend...
                        </TableCell>
                    </TableRow>
                {:else if rows.length === 0}
                    <TableRow class="hover:bg-transparent">
                        <TableCell
                            colspan={columns.length}
                            class="py-10 text-center text-slate-400"
                        >
                            Nenhum usuário encontrado para os filtros atuais.
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

    <DataTablePagination {page} {pageSize} {totalItems} {totalPages} {onPageChange} />
</div>
