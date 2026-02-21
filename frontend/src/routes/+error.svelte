<script lang="ts">
    import { page } from '$app/state'
    import { resolve } from '$app/paths'
    import { buttonVariants } from '$lib/components/ui/button'
    import { cn } from '$lib/utils'

    const status = $derived(page.status)
    const error = $derived(page.error as { message?: string } | null)
    const title = $derived(status === 404 ? 'Page not found' : 'Something went wrong')
    const description = $derived(
        status === 404
            ? "The page you tried to access doesn't exist."
            : error?.message || 'An unexpected error happened while loading this page.'
    )
</script>

<section class="flex min-h-[calc(100vh-12rem)] items-center justify-center px-4 py-12">
    <div
        class="w-full max-w-xl rounded border border-slate-800 bg-slate-900 p-8 text-center shadow-lg"
    >
        <p class="text-sm font-semibold text-slate-400">Error {status}</p>
        <h1 class="mt-2 text-3xl font-bold text-white">{title}</h1>
        <p class="mt-4 text-slate-300">{description}</p>

        <div class="mt-8 flex flex-wrap items-center justify-center gap-4">
            <a href={resolve('/')} class={buttonVariants({ variant: 'default' })}> Go Home </a>
            <a
                href={resolve('/login')}
                class={cn(buttonVariants({ variant: 'outline' }), 'border-slate-700')}
            >
                Sign In
            </a>
        </div>
    </div>
</section>
