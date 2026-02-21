<script lang="ts">
    import { onMount } from 'svelte'
    import { apiRequest } from '$lib/api/client'

    interface HealthResponse {
        status: string
    }

    let isLoading = $state(true)
    let healthStatus = $state('unknown')
    let errorMessage = $state('')

    async function checkStatus() {
        isLoading = true
        errorMessage = ''

        try {
            const response = await apiRequest<HealthResponse>('/health', {
                method: 'GET'
            })
            healthStatus = response.status
        } catch (error) {
            healthStatus = 'unavailable'
            errorMessage = error instanceof Error ? error.message : 'Status check failed'
        } finally {
            isLoading = false
        }
    }

    onMount(() => {
        void checkStatus()
    })
</script>

<section class="mx-auto max-w-2xl py-6">
    <h1 class="text-3xl font-bold text-white">System Status</h1>
    <p class="mt-2 text-slate-400">Live status from the backend health endpoint.</p>

    <div class="mt-8 rounded border border-slate-800 bg-slate-900 p-6">
        <p class="text-sm text-slate-400">Backend health</p>
        <p class="mt-2 text-2xl font-semibold text-white">
            {isLoading ? 'Checking...' : healthStatus}
        </p>

        {#if errorMessage}
            <p class="mt-3 text-sm text-red-300">{errorMessage}</p>
        {/if}

        <button
            type="button"
            onclick={checkStatus}
            disabled={isLoading}
            class="mt-6 rounded bg-blue-600 px-4 py-2 font-semibold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-70"
        >
            Refresh
        </button>
    </div>
</section>
