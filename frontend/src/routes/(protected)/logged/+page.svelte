<!-- frontend/src/routes/(protected)/logged/+page.svelte -->

<script lang="ts">
	import { authStore } from '$lib/stores/auth';
	import { goto } from '$app/navigation';	
	import { onMount } from 'svelte';
	import { apiRequest, getAccessToken } from '$lib/api/client';

	// User data from auth store
	let user = $derived($authStore.user);
	let isLoading = $derived($authStore.isLoading);

	interface Data {
		message: string;
	}

	let data = $state<Data>()

	async function fetchData() {
		try {
			const response = await apiRequest<Data>('/api/protected', {
				method: 'GET',
				headers: {
					'Authorization': `Bearer ${getAccessToken()}`
				}
			})

			console.log(response);

			data = response;
		} catch (error) {
			console.error('Error fetching data:', error);
		}
	}

	let intervalId: number;

	onMount(() => {
		// Buscar dados imediatamente na montagem
		fetchData();
		
		// Configurar intervalo para buscar a cada 5 minutos (300000 ms)
		intervalId = setInterval(fetchData, 300000);
		
		// Limpar o intervalo quando o componente for destruído
		return () => {
			if (intervalId) clearInterval(intervalId);
		};
	});


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

<section class="py-12">
	<div class="max-w-3xl mx-auto text-center mb-16">
		<h1 class="text-4xl font-bold mb-6 inline-block bg-clip-text">Authentication Test Page</h1>
		
		{#if isLoading}
			<div class="bg-slate-900 border border-slate-800 rounded-lg p-6 mb-8 max-w-md mx-auto">
				<div class="flex justify-center items-center h-20">
					<div class="h-8 w-8 animate-spin rounded-full border-4 border-slate-700 border-t-slate-300"></div>
					<p class="ml-3 text-slate-300">Loading authentication data...</p>
				</div>
			</div>
		{:else if user}
			<div class="bg-slate-900 border border-slate-800 rounded-lg p-6 mb-8 max-w-md mx-auto text-left">
				<h2 class="text-2xl font-semibold mb-4 text-center">User Information</h2>
				
				<div class="space-y-3">
					<div class="flex justify-between">
						<span class="text-slate-400">Username:</span>
						<span class="font-medium">{user.Username}</span>
					</div>
					
					<div class="flex justify-between">
						<span class="text-slate-400">Display Name:</span>
						<span class="font-medium">{user.DisplayName}</span>
					</div>
					
					<div class="flex justify-between">
						<span class="text-slate-400">Email:</span>
						<span class="font-medium">{user.Email}</span>
					</div>
					
					<div class="flex justify-between">
						<span class="text-slate-400">Role:</span>
						<span class="font-medium">{user.Role}</span>
					</div>
					
					<div class="flex justify-between">
						<span class="text-slate-400">User ID:</span>
						<span class="font-medium">{user.ID}</span>
					</div>
					<p class="text-slate-400">Data from backend: {data?.message}</p>
				</div>
			</div>		
		{/if}
		
		<button 
			onclick={handleLogout}
			class="px-6 py-3 bg-red-600 hover:bg-red-700 text-white font-medium rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 focus:ring-offset-slate-950"
		>
			Logout
		</button>
	</div>
</section>
