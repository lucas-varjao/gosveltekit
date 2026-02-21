<script lang="ts">
    import { resolve } from '$app/paths'
</script>

<section class="mx-auto max-w-4xl py-6">
    <h1 class="text-3xl font-bold text-white">Settings</h1>
    <p class="mt-2 text-slate-400">Manage account security and active sessions.</p>

    <div class="mt-8 grid gap-4 md:grid-cols-2">
        <a
            href={resolve('/settings/security')}
            class="rounded border border-slate-800 bg-slate-900 p-5 transition-colors hover:border-slate-700 hover:bg-slate-800"
        >
            <h2 class="text-xl font-semibold text-white">Security</h2>
            <p class="mt-2 text-sm text-slate-400">
                Change your password and review account security.
            </p>
        </a>

        <a
            href={resolve('/settings/sessions')}
            class="rounded border border-slate-800 bg-slate-900 p-5 transition-colors hover:border-slate-700 hover:bg-slate-800"
        >
            <h2 class="text-xl font-semibold text-white">Sessions</h2>
            <p class="mt-2 text-sm text-slate-400">Inspect active sessions and revoke access.</p>
        </a>
    </div>
</section>
