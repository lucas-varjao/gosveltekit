<script lang="ts">
    import { goto } from '$app/navigation'
    import { resolve } from '$app/paths'
    import PasswordField from '$lib/components/auth/password-field.svelte'
    import { accountApi } from '$lib/api/account'
    import { authStore } from '$lib/stores/auth'
    import { Alert, AlertDescription } from '$lib/components/ui/alert'
    import { buttonVariants } from '$lib/components/ui/button'
    import { Card, CardContent } from '$lib/components/ui/card'
    import { cn } from '$lib/utils'

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

    <Card class="mt-8 border-slate-800 bg-slate-900">
        <CardContent>
            <form onsubmit={handleSubmit} class="flex flex-col gap-4">
                <PasswordField
                    id="current_password"
                    label="Current Password"
                    bind:value={currentPassword}
                    placeholder="Current password"
                />

                <PasswordField
                    id="new_password"
                    label="New Password"
                    bind:value={newPassword}
                    placeholder="New password"
                />

                <PasswordField
                    id="confirm_password"
                    label="Confirm Password"
                    bind:value={confirmPassword}
                    placeholder="Confirm new password"
                />

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
                    disabled={isLoading}
                    class={cn(buttonVariants({ variant: 'default' }), 'mt-2 w-fit')}
                >
                    {isLoading ? 'Updating...' : 'Update Password'}
                </button>
            </form>
        </CardContent>
    </Card>
</section>
