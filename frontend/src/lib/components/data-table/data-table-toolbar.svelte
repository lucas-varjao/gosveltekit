<script lang="ts">
    import { Search } from '@lucide/svelte'
    import { Input } from '$lib/components/ui/input'

    type Props = {
        searchValue: string
        totalItems: number
        isLoading?: boolean
        onSearchChange?: (value: string) => void
    }

    let { searchValue, totalItems, isLoading = false, onSearchChange }: Props = $props()
</script>

<div class="flex flex-col gap-4 px-4 py-4 md:flex-row md:items-center md:justify-between">
    <div class="relative w-full max-w-md">
        <Search
            class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-slate-500"
        />
        <Input
            value={searchValue}
            placeholder="Buscar por nome, usuário ou email..."
            class="border-slate-800 bg-slate-950/70 pl-9 text-slate-100 placeholder:text-slate-500"
            oninput={(event) => onSearchChange?.(event.currentTarget.value)}
        />
    </div>

    <div class="text-sm text-slate-400">
        {#if isLoading}
            Atualizando resultados...
        {:else}
            {totalItems} registro{totalItems === 1 ? '' : 's'}
        {/if}
    </div>
</div>
