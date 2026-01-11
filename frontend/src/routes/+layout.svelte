<!-- frontend/src/routes/+layout.svelte -->

<script lang="ts">
	import '../app.css';
	import { authStore } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	let { children } = $props();

	let isAuthenticated = $derived($authStore.isAuthenticated);
	let isLoading = $derived($authStore.isLoading);
	let user = $derived($authStore.user);

	// Function to handle logout
	async function handleLogout() {
		try {
			await authStore.logout();

			goto('/login');
		} catch (error) {
			console.error('Logout failed:', error);
		}
	}
</script>

<div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col">
	<header class="border-b border-slate-800">
		<div class="container mx-auto px-4 py-4">
			<nav class="flex items-center justify-between">
				<div class="font-bold text-xl"><a href="/">GoSvelteKit</a></div>
				<div class="font-bold text-xl">
					{#if isLoading}	
						<div class="h-4 w-4 animate-spin rounded-full border-4 border-slate-700 border-t-slate-300"></div>
					{:else if isAuthenticated}
						{#if user}
							<span>{user.display_name} | </span>							
						{/if}
						<button onclick={handleLogout} class="text-slate-400 hover:text-white text-base">Sign Out</button>
					{:else}
						<a href="/login" class="text-slate-400 hover:text-white">Sign In</a>
					{/if}
				</div>
			</nav>
		</div>
	</header>
	
	<main class="container mx-auto px-4 py-8 flex-grow">
		{@render children()}
	</main>
	
	<footer class="border-t border-slate-800 h-24 flex items-center justify-center">
		<div class="container mx-auto px-4 py-6 text-center text-slate-400 text-sm">
			<span>&copy; {new Date().getFullYear()} Lucas Varjão - Desenvolvido com SvelteKit e Go</span>
		</div>
	</footer>
</div>
