<script lang="ts">
    import { onMount } from 'svelte'
    import { goto } from '$app/navigation'
    import { resolve } from '$app/paths'
    import { accountApi, type AccountSession } from '$lib/api/account'

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
        <p class="mt-6 rounded border border-red-500 bg-red-900/50 px-4 py-3 text-sm text-red-300">
            {errorMessage}
        </p>
    {/if}

    {#if isLoading}
        <div class="mt-8 rounded border border-slate-800 bg-slate-900 p-6 text-slate-300">
            Loading sessions...
        </div>
    {:else if sessions.length === 0}
        <div class="mt-8 rounded border border-slate-800 bg-slate-900 p-6 text-slate-300">
            No active sessions found.
        </div>
    {:else}
        <div class="mt-8 space-y-4">
            {#each sessions as session (session.id)}
                <article class="rounded border border-slate-800 bg-slate-900 p-5">
                    <div class="flex flex-wrap items-start justify-between gap-3">
                        <div>
                            <h2 class="font-semibold text-white">
                                {session.is_current ? 'Current session' : 'Active session'}
                            </h2>
                            <p class="mt-1 text-sm text-slate-400">
                                {session.user_agent || 'Unknown device'}
                            </p>
                            <p class="mt-1 text-sm text-slate-500">IP: {session.ip || 'N/A'}</p>
                            <p class="mt-1 text-sm text-slate-500">
                                Created: {formatDate(session.created_at)}
                            </p>
                            <p class="mt-1 text-sm text-slate-500">
                                Expires: {formatDate(session.expires_at)}
                            </p>
                        </div>

                        <button
                            type="button"
                            disabled={isRevoking}
                            onclick={() => revokeSession(session)}
                            class="rounded border border-red-700 px-3 py-2 text-sm font-semibold text-red-300 transition-colors hover:bg-red-900/40 disabled:cursor-not-allowed disabled:opacity-60"
                        >
                            {session.is_current ? 'Revoke Current Session' : 'Revoke'}
                        </button>
                    </div>
                </article>
            {/each}
        </div>
    {/if}
</section>
