<script lang="ts">
    import type { Snippet } from 'svelte'
    import type { HTMLAttributes } from 'svelte/elements'
    import { cn } from '$lib/utils'

    type Props = HTMLAttributes<HTMLTableRowElement> & {
        class?: string
        children?: Snippet
    }

    let { class: className, children, ...restProps }: Props = $props()
</script>

<tr
    data-slot="table-row"
    class={cn(
        'border-b border-slate-800/60 transition-colors hover:bg-slate-900/50 data-[state=selected]:bg-slate-900/70',
        className
    )}
    {...restProps}
>
    {@render children?.()}
</tr>
