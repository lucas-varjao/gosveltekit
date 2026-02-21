<script lang="ts">
    import { onMount } from 'svelte'
    import { accountApi, type AccountProfile } from '$lib/api/account'
    import { Alert, AlertDescription } from '$lib/components/ui/alert'
    import { buttonVariants } from '$lib/components/ui/button'
    import { Card, CardContent } from '$lib/components/ui/card'
    import { Input } from '$lib/components/ui/input'
    import { Label } from '$lib/components/ui/label'
    import { cn } from '$lib/utils'

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
        <Card class="mt-8 border-slate-800 bg-slate-900">
            <CardContent class="text-slate-300">Loading profile...</CardContent>
        </Card>
    {:else if profile}
        <Card class="mt-8 border-slate-800 bg-slate-900">
            <CardContent>
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
                        <Label for="display_name">Display Name</Label>
                        <Input id="display_name" type="text" bind:value={displayName} />
                    </div>

                    <div class="grid gap-4 md:grid-cols-2">
                        <div class="flex flex-col gap-2">
                            <Label for="first_name">First Name</Label>
                            <Input id="first_name" type="text" bind:value={firstName} />
                        </div>

                        <div class="flex flex-col gap-2">
                            <Label for="last_name">Last Name</Label>
                            <Input id="last_name" type="text" bind:value={lastName} />
                        </div>
                    </div>

                    {#if errorMessage}
                        <Alert
                            variant="destructive"
                            class="border-red-500/60 bg-red-950/50 text-red-200"
                        >
                            <AlertDescription>{errorMessage}</AlertDescription>
                        </Alert>
                    {/if}

                    {#if successMessage}
                        <Alert class="border-emerald-500/60 bg-emerald-950/50 text-emerald-200">
                            <AlertDescription>{successMessage}</AlertDescription>
                        </Alert>
                    {/if}

                    <button
                        type="submit"
                        disabled={isSaving}
                        class={cn(buttonVariants({ variant: 'default' }), 'mt-2 w-fit')}
                    >
                        {isSaving ? 'Saving...' : 'Save Changes'}
                    </button>
                </form>
            </CardContent>
        </Card>
    {/if}
</section>
