<script lang="ts">
    import { Button } from '$lib/components/ui/button'

    type Props = {
        page: number
        pageSize: number
        totalItems: number
        totalPages: number
        onPageChange?: (page: number) => void
    }

    let { page, pageSize, totalItems, totalPages, onPageChange }: Props = $props()

    let startItem = $derived(totalItems === 0 ? 0 : (page - 1) * pageSize + 1)
    let endItem = $derived(Math.min(page * pageSize, totalItems))
    let visiblePages = $derived.by(() => {
        const pages: number[] = []
        const firstPage = Math.max(1, page - 2)
        const lastPage = Math.min(totalPages, firstPage + 4)

        for (let current = Math.max(1, lastPage - 4); current <= lastPage; current += 1) {
            pages.push(current)
        }

        return pages
    })

    function goToPage(nextPage: number) {
        if (nextPage < 1 || nextPage > totalPages || nextPage === page) {
            return
        }

        onPageChange?.(nextPage)
    }
</script>

<div class="flex flex-col gap-4 px-4 py-4 md:flex-row md:items-center md:justify-between">
    <p class="text-sm text-slate-400">
        Exibindo {startItem}-{endItem} de {totalItems}
    </p>

    <div class="flex flex-wrap items-center gap-2">
        <Button variant="outline" size="sm" disabled={page <= 1} onclick={() => goToPage(page - 1)}>
            Anterior
        </Button>

        {#each visiblePages as visiblePage (visiblePage)}
            <Button
                variant={visiblePage === page ? 'secondary' : 'ghost'}
                size="sm"
                onclick={() => goToPage(visiblePage)}
            >
                {visiblePage}
            </Button>
        {/each}

        <Button
            variant="outline"
            size="sm"
            disabled={page >= totalPages}
            onclick={() => goToPage(page + 1)}
        >
            Próxima
        </Button>
    </div>
</div>
