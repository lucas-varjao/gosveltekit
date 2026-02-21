<script lang="ts">
    import { resolve } from '$app/paths'
    import { buttonVariants } from '$lib/components/ui/button'
    import { cn } from '$lib/utils'
</script>

<section class="mx-auto max-w-2xl py-12 text-center">
    <p class="text-sm font-semibold text-amber-400">503</p>
    <h1 class="mt-2 text-3xl font-bold text-white">Scheduled maintenance</h1>
    <p class="mt-4 text-slate-300">
        The system is temporarily unavailable while maintenance is in progress.
    </p>

    <div class="mt-8 flex flex-wrap items-center justify-center gap-3">
        <a href={resolve('/status')} class={buttonVariants({ variant: 'default' })}>Check Status</a>
        <a
            href={resolve('/')}
            class={cn(buttonVariants({ variant: 'outline' }), 'border-slate-700')}
        >
            Back Home
        </a>
    </div>
</section>
