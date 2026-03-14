<script lang="ts">
    import type { Snippet } from 'svelte'
    import type { HTMLTdAttributes } from 'svelte/elements'
    import { cn } from '$lib/utils'

    type Props = HTMLTdAttributes & {
        class?: string
        children?: Snippet
    }

    let { class: className, children, ...restProps }: Props = $props()
</script>

<td data-slot="table-cell" class={cn('p-4 align-middle text-slate-200', className)} {...restProps}>
    {@render children?.()}
</td>
