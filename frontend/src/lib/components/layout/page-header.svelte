<script lang="ts">
    import { cn } from '$lib/utils'

    type Props = {
        title: string
        description?: string
        eyebrow?: string
        align?: 'left' | 'center'
        class?: string
    }

    let {
        title,
        description = '',
        eyebrow = '',
        align = 'left',
        class: className
    }: Props = $props()
</script>

<header class={cn('page-header', align === 'center' && 'items-center text-center', className)}>
    {#if eyebrow}
        <p
            class="w-fit rounded-full border border-cyan-400/30 bg-cyan-500/10 px-3 py-1 text-xs font-semibold tracking-[0.08em] text-cyan-200 uppercase"
        >
            {eyebrow}
        </p>
    {/if}
    <h1 class="page-title">{title}</h1>
    {#if description}
        <p class="page-subtitle">{description}</p>
    {/if}
</header>
