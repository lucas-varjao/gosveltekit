<!-- frontend/src/routes/+layout.svelte -->

<script lang="ts">
    import '../app.css'
    import { LoaderCircle, LogIn, LogOut, ShieldCheck } from '@lucide/svelte'
    import { page } from '$app/state'
    import { authStore } from '$lib/stores/auth'
    import { buttonVariants } from '$lib/components/ui/button'
    import { goto } from '$app/navigation'
    import { resolve } from '$app/paths'
    import { cn } from '$lib/utils'

    let { children } = $props()

    let isAuthenticated = $derived($authStore.isAuthenticated)
    let isLoading = $derived($authStore.isLoading)
    let user = $derived($authStore.user)
    let pathname = $derived(page.url.pathname)

    type AppPath = Parameters<typeof resolve>[0]

    const navLinks: Array<{ path: AppPath; label: string }> = [
        { path: '/', label: 'Home' },
        { path: '/status', label: 'Status' }
    ]

    const protectedNavLinks: Array<{ path: AppPath; label: string }> = [
        { path: '/profile', label: 'Profile' },
        { path: '/settings', label: 'Settings' },
        { path: '/admin', label: 'Admin' }
    ]

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

<div class="flex min-h-screen flex-col text-slate-100">
    <header class="sticky top-0 z-40 border-b border-slate-800/80 bg-slate-950/75 backdrop-blur-xl">
        <div class="mx-auto w-full max-w-6xl px-4 py-4">
            <nav class="flex items-center justify-between">
                <div class="flex items-center gap-6">
                    <a href={resolve('/')} class="flex items-center gap-2 text-xl font-bold">
                        <ShieldCheck class="size-5 text-blue-300" />
                        <span class="brand-gradient">GoSvelteKit</span>
                    </a>

                    <div class="hidden items-center gap-2 md:flex">
                        {#each navLinks as nav (nav.path)}
                            <a
                                href={resolve(nav.path)}
                                class={cn(
                                    'top-nav-link',
                                    pathname === resolve(nav.path) && 'top-nav-link-active'
                                )}
                            >
                                {nav.label}
                            </a>
                        {/each}

                        {#if isAuthenticated}
                            {#each protectedNavLinks as nav (nav.path)}
                                <a
                                    href={resolve(nav.path)}
                                    class={cn(
                                        'top-nav-link',
                                        pathname.startsWith(resolve(nav.path)) &&
                                            'top-nav-link-active'
                                    )}
                                >
                                    {nav.label}
                                </a>
                            {/each}
                        {/if}
                    </div>
                </div>

                <div class="text-sm font-semibold">
                    {#if isLoading}
                        <LoaderCircle class="size-4 animate-spin text-slate-300" />
                    {:else if isAuthenticated}
                        {#if user}
                            <span>{user.display_name} | </span>
                        {/if}
                        <button
                            onclick={handleLogout}
                            class={cn(buttonVariants({ variant: 'ghost', size: 'sm' }), 'px-2')}
                        >
                            <LogOut class="size-4" />
                            Sign Out
                        </button>
                    {:else}
                        <a
                            href={resolve('/login')}
                            class={cn(buttonVariants({ variant: 'ghost', size: 'sm' }), 'px-2')}
                        >
                            <LogIn class="size-4" />
                            Sign In</a
                        >
                    {/if}
                </div>
            </nav>
        </div>
    </header>

    <main class="flex-grow">
        {@render children()}
    </main>

    <footer class="flex h-24 items-center justify-center border-t border-slate-800/80">
        <div class="mx-auto w-full max-w-6xl px-4 py-6 text-center text-sm text-slate-400">
            <span>&copy; {new Date().getFullYear()} Lucas Varjão - Built with SvelteKit and Go</span
            >
        </div>
    </footer>
</div>
