<!-- frontend/src/routes/+page.svelte -->

<script lang="ts">
    import { resolve } from '$app/paths'
</script>

<section class="py-12">
    <div class="mx-auto mb-16 max-w-3xl text-center">
        <h1 class="mb-6 inline-block bg-clip-text text-4xl font-bold">Home Page Example</h1>
        <p class="mb-8 text-slate-400">
            Base template with authentication, account pages, and system states.
        </p>

        <div class="flex flex-wrap items-center justify-center gap-3">
            <a
                href={resolve('/login')}
                class="rounded bg-blue-600 px-4 py-2 font-semibold text-white hover:bg-blue-700"
            >
                Sign In
            </a>
            <a
                href={resolve('/profile')}
                class="rounded border border-slate-700 px-4 py-2 font-semibold text-slate-200 hover:bg-slate-800"
            >
                Profile
            </a>
            <a
                href={resolve('/status')}
                class="rounded border border-slate-700 px-4 py-2 font-semibold text-slate-200 hover:bg-slate-800"
            >
                Status
            </a>
        </div>
    </div>
</section>
