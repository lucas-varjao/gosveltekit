<script lang="ts">
    import { resolve } from '$app/paths'
    import { LockKeyhole, MonitorSmartphone } from '@lucide/svelte'
    import PageHeader from '$lib/components/layout/page-header.svelte'
    import {
        Card,
        CardContent,
        CardDescription,
        CardHeader,
        CardTitle
    } from '$lib/components/ui/card'
</script>

<section class="page-shell">
    <PageHeader title="Settings" description="Manage account security and active sessions." />

    <div class="mt-8 grid gap-4 md:grid-cols-2">
        <a href={resolve('/settings/security')}>
            <Card
                class="surface-card transition-colors hover:border-slate-700 hover:bg-slate-800/80"
            >
                <CardHeader class="flex-row items-center gap-3">
                    <LockKeyhole class="size-5 text-slate-300" />
                    <CardTitle class="text-xl">Security</CardTitle>
                </CardHeader>
                <CardContent>
                    <CardDescription>
                        Change your password and review account security.
                    </CardDescription>
                </CardContent>
            </Card>
        </a>

        <a href={resolve('/settings/sessions')}>
            <Card
                class="surface-card transition-colors hover:border-slate-700 hover:bg-slate-800/80"
            >
                <CardHeader class="flex-row items-center gap-3">
                    <MonitorSmartphone class="size-5 text-slate-300" />
                    <CardTitle class="text-xl">Sessions</CardTitle>
                </CardHeader>
                <CardContent>
                    <CardDescription>Inspect active sessions and revoke access.</CardDescription>
                </CardContent>
            </Card>
        </a>
    </div>
</section>
