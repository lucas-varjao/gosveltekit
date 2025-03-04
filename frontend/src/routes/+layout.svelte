<script lang="ts">
	import '../app.css';
	let { children } = $props();
</script>

<div class="min-h-screen bg-slate-950 text-slate-100">
	<header class="border-b border-slate-800">
		<div class="container mx-auto px-4 py-4">
			<nav class="flex items-center justify-between">
				<div class="font-bold text-xl">GoSvelteKit</div>
				<div>
					<a href="https://github.com/lucas-varjao" target="_blank" rel="noopener noreferrer" class="text-slate-400 hover:text-white">
						Full Stack
					</a>
				</div>
			</nav>
		</div>
	</header>
	
	<main class="container mx-auto px-4 py-8">
		{@render children()}
	</main>
	
	<footer class="border-t border-slate-800 mt-auto">
		<div class="container mx-auto px-4 py-6 text-center text-slate-400 text-sm">
			&copy; {new Date().getFullYear()} Lucas Varjão - Desenvolvido com SvelteKit e Go
		</div>
	</footer>
</div>
