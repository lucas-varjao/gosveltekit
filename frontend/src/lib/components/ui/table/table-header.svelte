<script lang="ts">
    import type { Snippet } from 'svelte'
    import type { HTMLAttributes } from 'svelte/elements'
    import { cn } from '$lib/utils'

    type Props = HTMLAttributes<HTMLTableSectionElement> & {
        class?: string
        children?: Snippet
    }

    let { class: className, children, ...restProps }: Props = $props()
</script>

<thead
    data-slot="table-header"
    class={cn('[&_tr]:border-b [&_tr]:border-slate-800/70', className)}
    {...restProps}
>
    {@render children?.()}
</thead>
