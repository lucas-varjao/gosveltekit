<!-- frontend/src/routes/+page.svelte -->

<script lang="ts">
    import { resolve } from '$app/paths'
    import { buttonVariants } from '$lib/components/ui/button'
    import { Activity, LogIn, UserRound } from '@lucide/svelte'
</script>

<section class="py-12">
    <div class="mx-auto mb-16 max-w-3xl text-center">
        <h1 class="mb-6 inline-block bg-clip-text text-4xl font-bold">Home Page Example</h1>
        <p class="mb-8 text-slate-400">
            Base template with authentication, account pages, and system states.
        </p>

        <div class="flex flex-wrap items-center justify-center gap-3">
            <a href={resolve('/login')} class={buttonVariants({ variant: 'default' })}>
                <LogIn class="size-4" />
                Sign In
            </a>
            <a href={resolve('/profile')} class={buttonVariants({ variant: 'outline' })}>
                <UserRound class="size-4" />
                Profile
            </a>
            <a href={resolve('/status')} class={buttonVariants({ variant: 'outline' })}>
                <Activity class="size-4" />
                Status
            </a>
        </div>
    </div>
</section>
