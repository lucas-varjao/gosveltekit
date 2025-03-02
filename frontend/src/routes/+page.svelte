<script lang="ts">
	// Dados dos projetos (você pode adicionar mais projetos conforme for desenvolvendo)
	const projects = [
		{
			id: 1,
			title: "Projeto 1",
			description: "Descrição do projeto ainda será adicionada.",
			status: "Em breve",
			tags: ["Go", "SvelteKit", "Tailwind"],
			link: "#",
			image: "https://placehold.co/600x400/0f172a/e2e8f0?text=Em+Breve"
		},
		{
			id: 2,
			title: "Projeto 2",
			description: "Descrição do projeto ainda será adicionada.",
			status: "Em breve",
			tags: ["Go", "SvelteKit", "Tailwind"],
			link: "#",
			image: "https://placehold.co/600x400/0f172a/e2e8f0?text=Em+Breve"
		},
		{
			id: 3,
			title: "Projeto 3",
			description: "Descrição do projeto ainda será adicionada.",
			status: "Em breve",
			tags: ["Go", "SvelteKit", "Tailwind"],
			link: "#",
			image: "https://placehold.co/600x400/0f172a/e2e8f0?text=Em+Breve"
		}
	];
</script>

<section class="py-12">
	<div class="max-w-3xl mx-auto text-center mb-16">
		<h1 class="text-4xl font-bold mb-6 bg-gradient-to-r from-blue-500 to-purple-500 inline-block text-transparent bg-clip-text">WebApps Collection</h1>
		<p class="text-slate-300 text-lg mb-8">Um hub de aplicações web desenvolvidas com Go, SvelteKit e Tailwind para aprimorar habilidades e construir um portfolio.</p>
		<div class="flex justify-center gap-3">
			<div class="px-3 py-1 rounded-full bg-slate-800 text-xs font-medium text-slate-300 flex items-center gap-1">
				<span class="w-2 h-2 rounded-full bg-blue-500"></span>
				Go
			</div>
			<div class="px-3 py-1 rounded-full bg-slate-800 text-xs font-medium text-slate-300 flex items-center gap-1">
				<span class="w-2 h-2 rounded-full bg-orange-500"></span>
				SvelteKit
			</div>
			<div class="px-3 py-1 rounded-full bg-slate-800 text-xs font-medium text-slate-300 flex items-center gap-1">
				<span class="w-2 h-2 rounded-full bg-cyan-500"></span>
				Tailwind
			</div>
		</div>
	</div>
	
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
		{#each projects as project}
			<div class="bg-slate-900 border border-slate-800 rounded-lg overflow-hidden hover:shadow-lg hover:shadow-slate-900/50 transition-all duration-300 hover:-translate-y-1">
				<div class="h-48 overflow-hidden">
					<img src={project.image} alt={project.title} class="w-full h-full object-cover" />
				</div>
				<div class="p-6">
					<div class="flex justify-between items-center mb-3">
						<h3 class="text-lg font-bold text-slate-100">{project.title}</h3>
						<span class="text-xs font-medium px-2 py-1 rounded-full bg-slate-800 text-slate-300">{project.status}</span>
					</div>
					<p class="text-slate-400 mb-4 text-sm">{project.description}</p>
					<div class="flex flex-wrap gap-2 mb-4">
						{#each project.tags as tag}
							<span class="text-xs px-2 py-1 bg-slate-800 text-slate-400 rounded">{tag}</span>
						{/each}
					</div>
					<a href={project.link} class="inline-block bg-slate-800 hover:bg-slate-700 text-slate-200 px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200">
						Ver projeto
					</a>
				</div>
			</div>
		{/each}
	</div>

	<!-- Seção "Em Breve" com efeito de glassmorphism -->
	<div class="mt-20 py-12 px-6 bg-slate-900/50 backdrop-blur-sm rounded-xl border border-slate-800">
		<div class="text-center mb-8">
			<h2 class="text-2xl font-bold mb-4">Próximos Projetos</h2>
			<p class="text-slate-400 max-w-2xl mx-auto">Mais aplicações estão sendo desenvolvidas e serão adicionadas a esta coleção em breve. Fique atento!</p>
		</div>
		
		<div class="flex flex-wrap gap-4 justify-center">
			{#each ["API RESTful", "Auth System", "Dashboard", "E-commerce", "Real-time Chat"] as idea}
				<span class="bg-slate-800/70 px-4 py-2 rounded-md text-slate-300 text-sm">{idea}</span>
			{/each}
		</div>
	</div>
</section>
