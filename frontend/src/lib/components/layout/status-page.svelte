<script lang="ts">
    import { resolve } from '$app/paths'
    import { buttonVariants } from '$lib/components/ui/button'
    import { cn } from '$lib/utils'

    type Tone = 'info' | 'warning' | 'danger'
    type AppPath = Parameters<typeof resolve>[0]

    type Props = {
        code: string | number
        title: string
        description: string
        primaryPath: AppPath
        primaryLabel: string
        secondaryPath?: AppPath
        secondaryLabel?: string
        tone?: Tone
    }

    let {
        code,
        title,
        description,
        primaryPath,
        primaryLabel,
        secondaryPath,
        secondaryLabel = '',
        tone = 'warning'
    }: Props = $props()

    const toneClass = $derived.by(() => {
        if (tone === 'danger') return 'text-red-300'
        if (tone === 'info') return 'text-blue-300'
        return 'text-amber-300'
    })
</script>

<section class="page-shell">
    <div class="mx-auto max-w-2xl text-center">
        <p class={cn('text-sm font-semibold', toneClass)}>{code}</p>
        <h1 class="mt-2 text-3xl font-semibold text-white md:text-4xl">{title}</h1>
        <p class="mx-auto mt-4 max-w-xl text-slate-300">{description}</p>

        <div class="mt-8 flex flex-wrap items-center justify-center gap-3">
            <a href={resolve(primaryPath)} class={buttonVariants({ variant: 'default' })}
                >{primaryLabel}</a
            >
            {#if secondaryPath && secondaryLabel}
                <a
                    href={resolve(secondaryPath)}
                    class={cn(buttonVariants({ variant: 'outline' }), 'border-slate-700')}
                >
                    {secondaryLabel}
                </a>
            {/if}
        </div>
    </div>
</section>
