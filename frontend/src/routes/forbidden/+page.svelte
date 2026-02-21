<script lang="ts">
    import { resolve } from '$app/paths'
</script>

<section class="mx-auto max-w-2xl py-12 text-center">
    <p class="text-sm font-semibold text-red-400">403</p>
    <h1 class="mt-2 text-3xl font-bold text-white">Access denied</h1>
    <p class="mt-4 text-slate-300">
        You are authenticated, but you do not have permission to access this resource.
    </p>

    <div class="mt-8 flex flex-wrap items-center justify-center gap-3">
        <a
            href={resolve('/profile')}
            class="rounded bg-blue-600 px-4 py-2 font-semibold text-white hover:bg-blue-700"
        >
            Go to Profile
        </a>
        <a
            href={resolve('/')}
            class="rounded border border-slate-700 px-4 py-2 font-semibold text-slate-200 hover:bg-slate-800"
        >
            Back Home
        </a>
    </div>
</section>
