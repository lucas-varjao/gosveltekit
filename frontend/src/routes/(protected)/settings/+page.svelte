<script lang="ts">
    import { resolve } from '$app/paths'
    import { LockKeyhole, MonitorSmartphone } from '@lucide/svelte'
    import {
        Card,
        CardContent,
        CardDescription,
        CardHeader,
        CardTitle
    } from '$lib/components/ui/card'
</script>

<section class="mx-auto max-w-4xl py-6">
    <h1 class="text-3xl font-bold text-white">Settings</h1>
    <p class="mt-2 text-slate-400">Manage account security and active sessions.</p>

    <div class="mt-8 grid gap-4 md:grid-cols-2">
        <a href={resolve('/settings/security')}>
            <Card
                class="border-slate-800 bg-slate-900 transition-colors hover:border-slate-700 hover:bg-slate-800"
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
                class="border-slate-800 bg-slate-900 transition-colors hover:border-slate-700 hover:bg-slate-800"
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
