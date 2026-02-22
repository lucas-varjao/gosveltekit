<!-- frontend/src/routes/+page.svelte -->

<script lang="ts">
    import { resolve } from '$app/paths'
    import PageHeader from '$lib/components/layout/page-header.svelte'
    import { buttonVariants } from '$lib/components/ui/button'
    import {
        Card,
        CardContent,
        CardDescription,
        CardHeader,
        CardTitle
    } from '$lib/components/ui/card'
    import { Activity, LogIn, ShieldCheck, UserRound } from '@lucide/svelte'
</script>

<section class="page-shell">
    <PageHeader
        align="center"
        eyebrow="Production starter"
        title="Build session-first apps faster"
        description="A fullstack baseline with SvelteKit + Go, pre-wired auth flows and account management pages."
    />

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

    <div class="mt-10 grid gap-4 md:grid-cols-3">
        <Card class="surface-card">
            <CardHeader class="flex-row items-center gap-3">
                <ShieldCheck class="size-5 text-blue-300" />
                <CardTitle class="text-lg">Session Auth</CardTitle>
            </CardHeader>
            <CardContent>
                <CardDescription>
                    Database-backed sessions and protected routes already connected.
                </CardDescription>
            </CardContent>
        </Card>

        <Card class="surface-card">
            <CardHeader class="flex-row items-center gap-3">
                <UserRound class="size-5 text-cyan-300" />
                <CardTitle class="text-lg">Account Area</CardTitle>
            </CardHeader>
            <CardContent>
                <CardDescription>
                    Profile editing, password change and session revocation screens included.
                </CardDescription>
            </CardContent>
        </Card>

        <Card class="surface-card">
            <CardHeader class="flex-row items-center gap-3">
                <Activity class="size-5 text-slate-200" />
                <CardTitle class="text-lg">Status Checks</CardTitle>
            </CardHeader>
            <CardContent>
                <CardDescription>
                    Ready-to-use maintenance, forbidden, session-expired and health status pages.
                </CardDescription>
            </CardContent>
        </Card>
    </div>
</section>
