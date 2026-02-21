<script lang="ts">
    import { onMount } from 'svelte'
    import { goto } from '$app/navigation'
    import { resolve } from '$app/paths'
    import { accountApi, type AccountSession } from '$lib/api/account'
    import { Alert, AlertDescription } from '$lib/components/ui/alert'
    import { buttonVariants } from '$lib/components/ui/button'
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card'
    import { cn } from '$lib/utils'

    let sessions = $state<AccountSession[]>([])
    let isLoading = $state(true)
    let isRevoking = $state(false)
    let errorMessage = $state('')

    async function loadSessions() {
        isLoading = true
        errorMessage = ''

        try {
            sessions = await accountApi.listSessions()
        } catch (error) {
            errorMessage = error instanceof Error ? error.message : 'Failed to load sessions'
        } finally {
            isLoading = false
        }
    }

    async function revokeSession(session: AccountSession) {
        isRevoking = true
        errorMessage = ''

        try {
            await accountApi.revokeSession(session.id)

            if (session.is_current) {
                goto(resolve('/session-expired'))
                return
            }

            await loadSessions()
        } catch (error) {
            errorMessage = error instanceof Error ? error.message : 'Failed to revoke session'
        } finally {
            isRevoking = false
        }
    }

    function formatDate(value: string) {
        const date = new Date(value)
        return Number.isNaN(date.getTime()) ? 'N/A' : date.toLocaleString()
    }

    onMount(() => {
        void loadSessions()
    })
</script>

<section class="mx-auto max-w-4xl py-6">
    <h1 class="text-3xl font-bold text-white">Sessions</h1>
    <p class="mt-2 text-slate-400">Review active sessions and revoke access when needed.</p>

    {#if errorMessage}
        <Alert variant="destructive" class="mt-6 border-red-500/60 bg-red-950/50 text-red-200">
            <AlertDescription>{errorMessage}</AlertDescription>
        </Alert>
    {/if}

    {#if isLoading}
        <Card class="mt-8 border-slate-800 bg-slate-900">
            <CardContent class="text-slate-300">Loading sessions...</CardContent>
        </Card>
    {:else if sessions.length === 0}
        <Card class="mt-8 border-slate-800 bg-slate-900">
            <CardContent class="text-slate-300">No active sessions found.</CardContent>
        </Card>
    {:else}
        <div class="mt-8 space-y-4">
            {#each sessions as session (session.id)}
                <Card class="border-slate-800 bg-slate-900">
                    <CardHeader class="flex-row items-start justify-between gap-3">
                        <div>
                            <CardTitle class="text-lg">
                                {session.is_current ? 'Current session' : 'Active session'}
                            </CardTitle>
                            <p class="mt-1 text-sm text-slate-400">
                                {session.user_agent || 'Unknown device'}
                            </p>
                        </div>

                        <button
                            type="button"
                            disabled={isRevoking}
                            onclick={() => revokeSession(session)}
                            class={cn(
                                buttonVariants({ variant: 'destructive', size: 'sm' }),
                                'text-sm'
                            )}
                        >
                            {session.is_current ? 'Revoke Current Session' : 'Revoke'}
                        </button>
                    </CardHeader>

                    <CardContent>
                        <p class="text-sm text-slate-500">IP: {session.ip || 'N/A'}</p>
                        <p class="mt-1 text-sm text-slate-500">
                            Created: {formatDate(session.created_at)}
                        </p>
                        <p class="mt-1 text-sm text-slate-500">
                            Expires: {formatDate(session.expires_at)}
                        </p>
                    </CardContent>
                </Card>
            {/each}
        </div>
    {/if}
</section>
