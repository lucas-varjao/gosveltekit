<!-- frontend/src/routes/(protected)/logged/+page.svelte -->

<script lang="ts">
    import { authStore } from '$lib/stores/auth'
    import { goto } from '$app/navigation'
    import { onMount } from 'svelte'
    import { apiRequest } from '$lib/api/client'

    // User data from auth store
    let user = $derived($authStore.user)
    let isLoading = $derived($authStore.isLoading)

    interface Data {
        message: string
    }

    let data = $state<Data>()

    async function fetchData() {
        try {
            // Session ID is automatically added by apiRequest
            const response = await apiRequest<Data>('/api/protected', {
                method: 'GET'
            })

            console.log(response)

            data = response
        } catch (error) {
            console.error('Error fetching data:', error)
        }
    }

    let intervalId: number

    onMount(() => {
        // Buscar dados imediatamente na montagem
        fetchData()

        // Configurar intervalo para buscar a cada 5 minutos (300000 ms)
        intervalId = setInterval(fetchData, 300000)

        // Limpar o intervalo quando o componente for destruído
        return () => {
            if (intervalId) clearInterval(intervalId)
        }
    })

    // Function to handle logout
    async function handleLogout() {
        try {
            await authStore.logout()

            goto('/login')
        } catch (error) {
            console.error('Logout failed:', error)
        }
    }
</script>

<section class="py-12">
    <div class="mx-auto mb-16 max-w-3xl text-center">
        <h1 class="mb-6 inline-block bg-clip-text text-4xl font-bold">Authentication Test Page</h1>

        {#if isLoading}
            <div class="mx-auto mb-8 max-w-md rounded-lg border border-slate-800 bg-slate-900 p-6">
                <div class="flex h-20 items-center justify-center">
                    <div
                        class="h-8 w-8 animate-spin rounded-full border-4 border-slate-700 border-t-slate-300"
                    ></div>
                    <p class="ml-3 text-slate-300">Loading authentication data...</p>
                </div>
            </div>
        {:else if user}
            <div
                class="mx-auto mb-8 max-w-md rounded-lg border border-slate-800 bg-slate-900 p-6 text-left"
            >
                <h2 class="mb-4 text-center text-2xl font-semibold">User Information</h2>

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
            </div>
        {/if}

        <button
            onclick={handleLogout}
            class="rounded-lg bg-red-600 px-6 py-3 font-medium text-white transition-colors duration-200 hover:bg-red-700 focus:ring-2 focus:ring-red-500 focus:ring-offset-2 focus:ring-offset-slate-950 focus:outline-none"
        >
            Logout
        </button>
    </div>
</section>
