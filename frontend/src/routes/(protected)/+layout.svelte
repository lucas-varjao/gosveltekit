<!-- frontend/src/routes/(protected)/+layout.svelte -->

<script lang="ts">
    import { browser } from '$app/environment'
    import { goto } from '$app/navigation'
    import { resolve } from '$app/paths'
    import { authStore } from '$lib/stores/auth'

    let { children } = $props()

    // Redirect to login if not authenticated
    $effect(() => {
        if (browser && !$authStore.isLoading && !$authStore.isAuthenticated) {
            goto(resolve('/login'))
        }
    })
</script>

{#if $authStore.isLoading}
    <div class="flex h-screen items-center justify-center">
        <div class="text-center">
            <div
                class="h-12 w-12 animate-spin rounded-full border-4 border-slate-700 border-t-slate-300"
            ></div>
            <p class="mt-4 text-slate-300">Loading...</p>
        </div>
    </div>
{:else if $authStore.isAuthenticated}
    {@render children?.()}
{/if}
