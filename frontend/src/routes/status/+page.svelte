<script lang="ts">
    import { onMount } from 'svelte'
    import { RefreshCw } from '@lucide/svelte'
    import { apiRequest } from '$lib/api/client'
    import { buttonVariants } from '$lib/components/ui/button'
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card'
    import { cn } from '$lib/utils'

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

    <Card class="mt-8 border-slate-800 bg-slate-900">
        <CardHeader>
            <CardTitle class="text-sm font-medium text-slate-400">Backend health</CardTitle>
        </CardHeader>

        <CardContent>
            <p class="text-2xl font-semibold text-white">
                {isLoading ? 'Checking...' : healthStatus}
            </p>

            {#if errorMessage}
                <p class="mt-3 text-sm text-red-300">{errorMessage}</p>
            {/if}

            <button
                type="button"
                onclick={checkStatus}
                disabled={isLoading}
                class={cn(buttonVariants({ variant: 'default' }), 'mt-6')}
            >
                <RefreshCw class={cn('size-4', isLoading && 'animate-spin')} />
                Refresh
            </button>
        </CardContent>
    </Card>
</section>
