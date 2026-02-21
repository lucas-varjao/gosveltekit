<script lang="ts">
    import { onMount } from 'svelte'
    import { accountApi, type AccountProfile } from '$lib/api/account'

    let profile = $state<AccountProfile | null>(null)
    let displayName = $state('')
    let firstName = $state('')
    let lastName = $state('')
    let isLoading = $state(true)
    let isSaving = $state(false)
    let successMessage = $state('')
    let errorMessage = $state('')

    async function loadProfile() {
        isLoading = true
        errorMessage = ''

        try {
            const response = await accountApi.getProfile()
            profile = response
            displayName = response.display_name
            firstName = response.first_name || ''
            lastName = response.last_name || ''
        } catch (error) {
            errorMessage = error instanceof Error ? error.message : 'Failed to load profile'
        } finally {
            isLoading = false
        }
    }

    async function handleSubmit(event: Event) {
        event.preventDefault()
        successMessage = ''
        errorMessage = ''
        isSaving = true

        try {
            const updated = await accountApi.updateProfile({
                display_name: displayName,
                first_name: firstName,
                last_name: lastName
            })

            profile = updated
            successMessage = 'Profile updated successfully.'
        } catch (error) {
            errorMessage = error instanceof Error ? error.message : 'Failed to update profile'
        } finally {
            isSaving = false
        }
    }

    onMount(() => {
        void loadProfile()
    })
</script>

<section class="mx-auto max-w-3xl py-6">
    <h1 class="text-3xl font-bold text-white">Profile</h1>
    <p class="mt-2 text-slate-400">Update your account profile information.</p>

    {#if isLoading}
        <div class="mt-8 rounded border border-slate-800 bg-slate-900 p-6 text-slate-300">
            Loading profile...
        </div>
    {:else if profile}
        <div class="mt-8 rounded border border-slate-800 bg-slate-900 p-6">
            <div class="grid gap-3 text-sm text-slate-300 md:grid-cols-2">
                <p><span class="text-slate-500">Username:</span> {profile.identifier}</p>
                <p><span class="text-slate-500">Email:</span> {profile.email}</p>
                <p><span class="text-slate-500">Role:</span> {profile.role}</p>
                <p>
                    <span class="text-slate-500">Status:</span>
                    {profile.active ? 'Active' : 'Inactive'}
                </p>
            </div>

            <form class="mt-6 flex flex-col gap-4" onsubmit={handleSubmit}>
                <div class="flex flex-col gap-2">
                    <label for="display_name" class="text-sm font-medium text-slate-200"
                        >Display Name</label
                    >
                    <input
                        id="display_name"
                        type="text"
                        bind:value={displayName}
                        class="rounded border border-slate-700 bg-slate-800 px-3 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none"
                    />
                </div>

                <div class="grid gap-4 md:grid-cols-2">
                    <div class="flex flex-col gap-2">
                        <label for="first_name" class="text-sm font-medium text-slate-200"
                            >First Name</label
                        >
                        <input
                            id="first_name"
                            type="text"
                            bind:value={firstName}
                            class="rounded border border-slate-700 bg-slate-800 px-3 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none"
                        />
                    </div>
                    <div class="flex flex-col gap-2">
                        <label for="last_name" class="text-sm font-medium text-slate-200"
                            >Last Name</label
                        >
                        <input
                            id="last_name"
                            type="text"
                            bind:value={lastName}
                            class="rounded border border-slate-700 bg-slate-800 px-3 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none"
                        />
                    </div>
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
                    disabled={isSaving}
                    class="mt-2 w-fit rounded bg-blue-600 px-4 py-2 font-semibold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-70"
                >
                    {isSaving ? 'Saving...' : 'Save Changes'}
                </button>
            </form>
        </div>
    {/if}
</section>
