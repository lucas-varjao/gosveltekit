<script lang="ts">
    import { onMount } from 'svelte'
    import { apiRequest } from '$lib/api/client'
    import { Alert, AlertDescription } from '$lib/components/ui/alert'
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card'

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
        <Card class="mt-8 border-slate-800 bg-slate-900">
            <CardContent class="text-slate-300">Loading admin dashboard...</CardContent>
        </Card>
    {:else if errorMessage}
        <Alert variant="destructive" class="mt-8 border-red-500/60 bg-red-950/50 text-red-200">
            <AlertDescription>{errorMessage}</AlertDescription>
        </Alert>
    {:else if data}
        <Card class="mt-8 border-slate-800 bg-slate-900">
            <CardHeader>
                <CardTitle>Dashboard</CardTitle>
            </CardHeader>
            <CardContent>
                <p class="text-slate-300">{data.message}</p>
            </CardContent>
        </Card>
    {/if}
</section>
