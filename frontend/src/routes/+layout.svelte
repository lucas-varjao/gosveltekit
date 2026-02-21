<!-- frontend/src/routes/+layout.svelte -->

<script lang="ts">
    import '../app.css'
    import { authStore } from '$lib/stores/auth'
    import { goto } from '$app/navigation'
    import { resolve } from '$app/paths'

    let { children } = $props()

    let isAuthenticated = $derived($authStore.isAuthenticated)
    let isLoading = $derived($authStore.isLoading)
    let user = $derived($authStore.user)

    // Function to handle logout
    async function handleLogout() {
        try {
            await authStore.logout()

            goto(resolve('/login'))
        } catch (error) {
            console.error('Logout failed:', error)
        }
    }
</script>

<div class="flex min-h-screen flex-col bg-slate-950 text-slate-100">
    <header class="border-b border-slate-800">
        <div class="container mx-auto px-4 py-4">
            <nav class="flex items-center justify-between">
                <div class="text-xl font-bold"><a href={resolve('/')}>GoSvelteKit</a></div>
                <div class="text-xl font-bold">
                    {#if isLoading}
                        <div
                            class="h-4 w-4 animate-spin rounded-full border-4 border-slate-700 border-t-slate-300"
                        ></div>
                    {:else if isAuthenticated}
                        {#if user}
                            <span>{user.display_name} | </span>
                        {/if}
                        <button
                            onclick={handleLogout}
                            class="text-base text-slate-400 hover:text-white">Sign Out</button
                        >
                    {:else}
                        <a href={resolve('/login')} class="text-slate-400 hover:text-white">
                            Sign In
                        </a>
                    {/if}
                </div>
            </nav>
        </div>
    </header>

    <main class="container mx-auto flex-grow px-4 py-8">
        {@render children()}
    </main>

    <footer class="flex h-24 items-center justify-center border-t border-slate-800">
        <div class="container mx-auto px-4 py-6 text-center text-sm text-slate-400">
            <span
                >&copy; {new Date().getFullYear()} Lucas Varjão - Desenvolvido com SvelteKit e Go</span
            >
        </div>
    </footer>
</div>
