<script lang="ts">
    import { resolve } from '$app/paths'
</script>

<section class="mx-auto max-w-2xl py-12 text-center">
    <p class="text-sm font-semibold text-amber-400">401</p>
    <h1 class="mt-2 text-3xl font-bold text-white">Unauthorized</h1>
    <p class="mt-4 text-slate-300">You need to sign in before accessing this page.</p>

    <div class="mt-8 flex flex-wrap items-center justify-center gap-3">
        <a
            href={resolve('/login')}
            class="rounded bg-blue-600 px-4 py-2 font-semibold text-white hover:bg-blue-700"
        >
            Sign In
        </a>
        <a
            href={resolve('/')}
            class="rounded border border-slate-700 px-4 py-2 font-semibold text-slate-200 hover:bg-slate-800"
        >
            Back Home
        </a>
    </div>
</section>
