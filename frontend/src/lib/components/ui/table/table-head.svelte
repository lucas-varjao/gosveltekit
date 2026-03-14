<script lang="ts">
    import type { Snippet } from 'svelte'
    import type { HTMLThAttributes } from 'svelte/elements'
    import { cn } from '$lib/utils'

    type Props = HTMLThAttributes & {
        class?: string
        children?: Snippet
    }

    let { class: className, children, ...restProps }: Props = $props()
</script>

<th
    data-slot="table-head"
    class={cn(
        'h-11 px-4 text-left align-middle text-xs font-semibold tracking-[0.08em] text-slate-400 uppercase',
        className
    )}
    {...restProps}
>
    {@render children?.()}
</th>
