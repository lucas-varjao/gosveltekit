<script lang="ts">
    import { page } from '$app/state'
    import { goto } from '$app/navigation'
    import { resolve } from '$app/paths'
    import PasswordField from '$lib/components/auth/password-field.svelte'
    import PasswordRequirements from '$lib/components/auth/password-requirements.svelte'
    import { authApi } from '$lib/api/auth'
    import { Alert, AlertDescription, AlertTitle } from '$lib/components/ui/alert'
    import { buttonVariants } from '$lib/components/ui/button'
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card'
    import { cn } from '$lib/utils'

    let password = $state('')
    let confirmPassword = $state('')
    let token = $state('')
    let errors = $state<Record<string, string>>({})
    let isLoading = $state(false)
    let touched = $state<Record<string, boolean>>({
        password: false,
        confirmPassword: false
    })
    let resetError = $state('')
    let resetSuccess = $state(false)

    let passwordRequirements = $state({
        length: false,
        uppercase: false,
        lowercase: false,
        number: false,
        special: false
    })

    $effect(() => {
        token = page.url.searchParams.get('token') || ''

        if (!token) {
            resetError = 'Invalid or missing password reset token'
        }
    })

    $effect(() => {
        passwordRequirements.length = password.length >= 8
        passwordRequirements.uppercase = /[A-Z]/.test(password)
        passwordRequirements.lowercase = /[a-z]/.test(password)
        passwordRequirements.number = /[0-9]/.test(password)
        passwordRequirements.special = /[!@#$%^&*(),.?":{}|<>]/.test(password)
    })

    function validatePassword(value: string): string | null {
        if (!value) return 'Password is required'
        if (value.length < 8) return 'Password must be at least 8 characters long'
        if (!/[A-Z]/.test(value)) return 'Password must contain at least one uppercase letter'
        if (!/[a-z]/.test(value)) return 'Password must contain at least one lowercase letter'
        if (!/[0-9]/.test(value)) return 'Password must contain at least one number'
        if (!/[!@#$%^&*(),.?":{}|<>]/.test(value)) {
            return 'Password must contain at least one special character'
        }
        return null
    }

    function validateConfirmPassword(value: string): string | null {
        if (!value) return 'Please confirm your password'
        if (value !== password) return 'Passwords do not match'
        return null
    }

    $effect(() => {
        if (touched.password) {
            errors.password = validatePassword(password) || ''
        }
    })

    $effect(() => {
        if (touched.confirmPassword) {
            errors.confirmPassword = validateConfirmPassword(confirmPassword) || ''
        }
    })

    function handleBlur(field: keyof typeof touched) {
        touched[field] = true
    }

    async function handleSubmit(event: Event) {
        event.preventDefault()
        resetError = ''

        Object.keys(touched).forEach((key) => {
            touched[key as keyof typeof touched] = true
        })

        errors = {
            password: validatePassword(password) || '',
            confirmPassword: validateConfirmPassword(confirmPassword) || ''
        }

        const hasErrors = Object.values(errors).some((error) => error !== '')
        if (hasErrors || !token) return

        try {
            isLoading = true

            await authApi.resetPassword({
                token,
                new_password: password,
                confirm_password: confirmPassword
            })

            resetSuccess = true
            setTimeout(() => {
                goto(resolve('/login'))
            }, 3000)
        } catch (error) {
            resetError =
                error instanceof Error ? error.message : 'Password reset failed. Please try again.'
        } finally {
            isLoading = false
        }
    }
</script>

<section class="page-shell flex min-h-[calc(100vh-6rem)] items-center justify-center">
    <Card class="surface-card w-full max-w-md shadow-lg">
        <CardHeader class="gap-3 text-center">
            <CardTitle class="text-3xl">Reset Your Password</CardTitle>
        </CardHeader>

        <CardContent>
            {#if resetSuccess}
                <Alert class="border-emerald-500/60 bg-emerald-950/50 text-emerald-200">
                    <AlertTitle>Password reset successful!</AlertTitle>
                    <AlertDescription>
                        You will be redirected to the login page in a moment.
                    </AlertDescription>
                </Alert>
            {:else if resetError}
                <div class="flex flex-col gap-4">
                    <Alert
                        variant="destructive"
                        class="border-red-500/60 bg-red-950/50 text-red-200"
                    >
                        <AlertDescription>{resetError}</AlertDescription>
                    </Alert>

                    <div class="flex justify-center">
                        <a
                            href={resolve('/forgot-password')}
                            class={buttonVariants({ variant: 'link', size: 'sm' })}
                        >
                            Request a new reset link
                        </a>
                    </div>
                </div>
            {:else}
                <div class="flex flex-col gap-6">
                    <p class="text-center text-slate-400">
                        Please enter and confirm your new password below.
                    </p>

                    <form onsubmit={handleSubmit} class="flex flex-col gap-4">
                        <PasswordField
                            id="password"
                            label="New Password"
                            bind:value={password}
                            placeholder="Enter your new password"
                            touched={touched.password}
                            error={errors.password}
                            onblur={() => handleBlur('password')}
                        />

                        <PasswordRequirements requirements={passwordRequirements} />

                        <PasswordField
                            id="confirmPassword"
                            label="Confirm Password"
                            bind:value={confirmPassword}
                            placeholder="Confirm your new password"
                            touched={touched.confirmPassword}
                            error={errors.confirmPassword}
                            onblur={() => handleBlur('confirmPassword')}
                        />

                        <button
                            type="submit"
                            disabled={isLoading}
                            class={cn(buttonVariants({ variant: 'default' }), 'mt-2 w-full')}
                        >
                            {isLoading ? 'Resetting Password...' : 'Reset Password'}
                        </button>
                    </form>

                    <div class="flex justify-center">
                        <a
                            href={resolve('/login')}
                            class={buttonVariants({ variant: 'link', size: 'sm' })}
                        >
                            Back to Login
                        </a>
                    </div>
                </div>
            {/if}
        </CardContent>
    </Card>
</section>
