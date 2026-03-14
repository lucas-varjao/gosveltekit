<script lang="ts">
    import { Button } from '$lib/components/ui/button'

    type Props = {
        pageSize: number
        hasNext: boolean
        hasPrev: boolean
        nextCursor?: string
        prevCursor?: string
        onCursorChange?: (direction: 'next' | 'prev', cursor: string) => void
    }

    let { pageSize, hasNext, hasPrev, nextCursor, prevCursor, onCursorChange }: Props = $props()

    function goToCursor(direction: 'next' | 'prev', cursor?: string) {
        if (!cursor) {
            return
        }

        onCursorChange?.(direction, cursor)
    }
</script>

<div class="flex flex-col gap-4 px-4 py-4 md:flex-row md:items-center md:justify-between">
    <p class="text-sm text-slate-400">
        Paginação por cursor com lotes de {pageSize} registro{pageSize === 1 ? '' : 's'}
    </p>

    <div class="flex flex-wrap items-center gap-2">
        <Button
            variant="outline"
            size="sm"
            disabled={!hasPrev || !prevCursor}
            onclick={() => goToCursor('prev', prevCursor)}
        >
            Anterior
        </Button>

        <Button
            variant="outline"
            size="sm"
            disabled={!hasNext || !nextCursor}
            onclick={() => goToCursor('next', nextCursor)}
        >
            Próxima
        </Button>
    </div>
</div>
