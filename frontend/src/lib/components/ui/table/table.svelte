<script lang="ts">
    import type { Snippet } from 'svelte'
    import type { HTMLTableAttributes } from 'svelte/elements'
    import { cn } from '$lib/utils'

    type Props = HTMLTableAttributes & {
        class?: string
        children?: Snippet
    }

    let { class: className, children, ...restProps }: Props = $props()
</script>

<div data-slot="table-container" class="relative w-full overflow-x-auto">
    <table
        data-slot="table"
        class={cn('w-full caption-bottom text-sm text-slate-100', className)}
        {...restProps}
    >
        {@render children?.()}
    </table>
</div>
