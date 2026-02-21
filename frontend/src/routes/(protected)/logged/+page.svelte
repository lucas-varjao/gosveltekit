<script lang="ts">
    import { goto } from '$app/navigation'
    import { resolve } from '$app/paths'
    import { onMount } from 'svelte'
    import { LoaderCircle, LogOut } from '@lucide/svelte'
    import { apiRequest } from '$lib/api/client'
    import { buttonVariants } from '$lib/components/ui/button'
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card'
    import { authStore } from '$lib/stores/auth'
    import { cn } from '$lib/utils'

    let user = $derived($authStore.user)
    let isLoading = $derived($authStore.isLoading)

    interface Data {
        message: string
    }

    let data = $state<Data>()

    async function fetchData() {
        try {
            const response = await apiRequest<Data>('/api/protected', {
                method: 'GET'
            })

            data = response
        } catch (error) {
            console.error('Error fetching data:', error)
        }
    }

    let intervalId: number

    onMount(() => {
        fetchData()
        intervalId = setInterval(fetchData, 300000)

        return () => {
            if (intervalId) clearInterval(intervalId)
        }
    })

    async function handleLogout() {
        try {
            await authStore.logout()
            goto(resolve('/login'))
        } catch (error) {
            console.error('Logout failed:', error)
        }
    }
</script>

<section class="py-12">
    <div class="mx-auto mb-16 max-w-3xl text-center">
        <h1 class="mb-6 inline-block bg-clip-text text-4xl font-bold">Authentication Test Page</h1>

        {#if isLoading}
            <Card class="mx-auto mb-8 max-w-md border-slate-800 bg-slate-900">
                <CardContent class="flex h-20 items-center justify-center gap-3">
                    <LoaderCircle class="size-6 animate-spin text-slate-300" />
                    <p class="text-slate-300">Loading authentication data...</p>
                </CardContent>
            </Card>
        {:else if user}
            <Card class="mx-auto mb-8 max-w-md border-slate-800 bg-slate-900 text-left">
                <CardHeader>
                    <CardTitle class="text-center text-2xl">User Information</CardTitle>
                </CardHeader>
                <CardContent>
                    <div class="space-y-3">
                        <div class="flex justify-between">
                            <span class="text-slate-400">Username:</span>
                            <span class="font-medium">{user.identifier}</span>
                        </div>

                        <div class="flex justify-between">
                            <span class="text-slate-400">Display Name:</span>
                            <span class="font-medium">{user.display_name}</span>
                        </div>

                        <div class="flex justify-between">
                            <span class="text-slate-400">Email:</span>
                            <span class="font-medium">{user.email}</span>
                        </div>

                        <div class="flex justify-between">
                            <span class="text-slate-400">Role:</span>
                            <span class="font-medium">{user.role}</span>
                        </div>

                        <div class="flex justify-between">
                            <span class="text-slate-400">User ID:</span>
                            <span class="font-medium">{user.id}</span>
                        </div>

                        <p class="text-slate-400">Data from backend: {data?.message}</p>
                    </div>
                </CardContent>
            </Card>
        {/if}

        <button
            onclick={handleLogout}
            class={cn(buttonVariants({ variant: 'destructive' }), 'gap-2')}
        >
            <LogOut class="size-4" />
            Logout
        </button>
    </div>
</section>
