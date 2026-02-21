<script lang="ts">
    import { goto } from '$app/navigation'
    import { resolve } from '$app/paths'
    import { accountApi } from '$lib/api/account'
    import { authStore } from '$lib/stores/auth'

    let currentPassword = $state('')
    let newPassword = $state('')
    let confirmPassword = $state('')
    let isLoading = $state(false)
    let errorMessage = $state('')
    let successMessage = $state('')

    async function handleSubmit(event: Event) {
        event.preventDefault()
        errorMessage = ''
        successMessage = ''

        if (newPassword !== confirmPassword) {
            errorMessage = 'Passwords do not match.'
            return
        }

        isLoading = true

        try {
            await accountApi.changePassword({
                current_password: currentPassword,
                new_password: newPassword,
                confirm_password: confirmPassword
            })

            successMessage = 'Password changed. Please sign in again.'
            await authStore.init()

            setTimeout(() => {
                goto(resolve('/session-expired'))
            }, 900)
        } catch (error) {
            errorMessage = error instanceof Error ? error.message : 'Failed to change password'
        } finally {
            isLoading = false
        }
    }
</script>

<section class="mx-auto max-w-2xl py-6">
    <h1 class="text-3xl font-bold text-white">Security</h1>
    <p class="mt-2 text-slate-400">Change your password.</p>

    <form onsubmit={handleSubmit} class="mt-8 rounded border border-slate-800 bg-slate-900 p-6">
        <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-2">
                <label for="current_password" class="text-sm font-medium text-slate-200"
                    >Current Password</label
                >
                <input
                    id="current_password"
                    type="password"
                    bind:value={currentPassword}
                    class="rounded border border-slate-700 bg-slate-800 px-3 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none"
                />
            </div>

            <div class="flex flex-col gap-2">
                <label for="new_password" class="text-sm font-medium text-slate-200"
                    >New Password</label
                >
                <input
                    id="new_password"
                    type="password"
                    bind:value={newPassword}
                    class="rounded border border-slate-700 bg-slate-800 px-3 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none"
                />
            </div>

            <div class="flex flex-col gap-2">
                <label for="confirm_password" class="text-sm font-medium text-slate-200"
                    >Confirm Password</label
                >
                <input
                    id="confirm_password"
                    type="password"
                    bind:value={confirmPassword}
                    class="rounded border border-slate-700 bg-slate-800 px-3 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none"
                />
            </div>

            {#if errorMessage}
                <p
                    class="rounded border border-red-500 bg-red-900/50 px-3 py-2 text-sm text-red-300"
                >
                    {errorMessage}
                </p>
            {/if}

            {#if successMessage}
                <p
                    class="rounded border border-green-500 bg-green-900/50 px-3 py-2 text-sm text-green-300"
                >
                    {successMessage}
                </p>
            {/if}

            <button
                type="submit"
                disabled={isLoading}
                class="mt-2 w-fit rounded bg-blue-600 px-4 py-2 font-semibold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-70"
            >
                {isLoading ? 'Updating...' : 'Update Password'}
            </button>
        </div>
    </form>
</section>
