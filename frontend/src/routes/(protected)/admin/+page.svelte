<script lang="ts">
    import { onMount } from 'svelte'
    import { apiRequest } from '$lib/api/client'

    interface AdminDashboardResponse {
        message: string
    }

    let data = $state<AdminDashboardResponse | null>(null)
    let isLoading = $state(true)
    let errorMessage = $state('')

    async function loadDashboard() {
        isLoading = true
        errorMessage = ''

        try {
            data = await apiRequest<AdminDashboardResponse>('/api/admin/dashboard', {
                method: 'GET',
                requiresAuth: true
            })
        } catch (error) {
            errorMessage = error instanceof Error ? error.message : 'Failed to load admin data'
        } finally {
            isLoading = false
        }
    }

    onMount(() => {
        void loadDashboard()
    })
</script>

<section class="mx-auto max-w-3xl py-6">
    <h1 class="text-3xl font-bold text-white">Admin</h1>
    <p class="mt-2 text-slate-400">Administrative area.</p>

    {#if isLoading}
        <div class="mt-8 rounded border border-slate-800 bg-slate-900 p-6 text-slate-300">
            Loading admin dashboard...
        </div>
    {:else if errorMessage}
        <div class="mt-8 rounded border border-red-500 bg-red-900/50 p-6 text-red-300">
            {errorMessage}
        </div>
    {:else if data}
        <div class="mt-8 rounded border border-slate-800 bg-slate-900 p-6">
            <h2 class="text-xl font-semibold text-white">Dashboard</h2>
            <p class="mt-2 text-slate-300">{data.message}</p>
        </div>
    {/if}
</section>
