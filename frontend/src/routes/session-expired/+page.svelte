<script lang="ts">
    import { resolve } from '$app/paths'
    import { buttonVariants } from '$lib/components/ui/button'
    import { cn } from '$lib/utils'
</script>

<section class="mx-auto max-w-2xl py-12 text-center">
    <p class="text-sm font-semibold text-amber-400">401</p>
    <h1 class="mt-2 text-3xl font-bold text-white">Session expired</h1>
    <p class="mt-4 text-slate-300">
        Your session is no longer valid. Please sign in again to continue.
    </p>

    <div class="mt-8 flex flex-wrap items-center justify-center gap-3">
        <a href={resolve('/login')} class={buttonVariants({ variant: 'default' })}>Sign In Again</a>
        <a
            href={resolve('/')}
            class={cn(buttonVariants({ variant: 'outline' }), 'border-slate-700')}
        >
            Back Home
        </a>
    </div>
</section>
