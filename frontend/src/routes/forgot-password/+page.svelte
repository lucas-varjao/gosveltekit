<script lang="ts">
    import { resolve } from '$app/paths'
    import { authApi } from '$lib/api/auth'
    import { Alert, AlertDescription, AlertTitle } from '$lib/components/ui/alert'
    import { buttonVariants } from '$lib/components/ui/button'
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card'
    import { Input } from '$lib/components/ui/input'
    import { Label } from '$lib/components/ui/label'
    import { cn } from '$lib/utils'

    let email = $state('')
    let errors = $state<Record<string, string>>({})
    let isLoading = $state(false)
    let touched = $state<Record<string, boolean>>({ email: false })
    let submitted = $state(false)
    let success = $state(false)

    function validateEmail(value: string): string | null {
        if (!value) return 'Email is required'

        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        if (!emailRegex.test(value)) {
            return 'Please enter a valid email address'
        }

        return null
    }

    $effect(() => {
        if (touched.email) {
            errors.email = validateEmail(email) || ''
        }
    })

    function handleBlur() {
        touched.email = true
    }

    async function handleSubmit(event: Event) {
        event.preventDefault()
        touched.email = true
        errors.email = validateEmail(email) || ''

        if (errors.email) return

        try {
            isLoading = true
            await authApi.requestPasswordReset({ email })
            success = true
            submitted = true
        } catch (error) {
            errors.email =
                error instanceof Error ? error.message : 'Failed to request password reset'
        } finally {
            isLoading = false
        }
    }

    function resetForm() {
        email = ''
        errors = {}
        touched = { email: false }
        submitted = false
        success = false
    }
</script>

<section class="flex min-h-[calc(100vh-6rem)] items-center justify-center px-4 py-12">
    <Card class="w-full max-w-md border-slate-800 bg-slate-900 shadow-lg">
        <CardHeader class="gap-3 text-center">
            <CardTitle class="text-3xl">Reset Password</CardTitle>
        </CardHeader>

        <CardContent>
            {#if submitted && success}
                <div class="flex flex-col gap-6">
                    <Alert class="border-blue-500/60 bg-blue-950/50 text-blue-200">
                        <AlertTitle>Check your email</AlertTitle>
                        <AlertDescription>
                            If an account exists with the email {email}, you will receive a password
                            reset link shortly.
                        </AlertDescription>
                    </Alert>

                    <div class="flex flex-col items-center gap-3">
                        <button
                            type="button"
                            onclick={resetForm}
                            class={buttonVariants({ variant: 'link', size: 'sm' })}
                        >
                            Request another reset link
                        </button>

                        <a
                            href={resolve('/login')}
                            class={buttonVariants({ variant: 'ghost', size: 'sm' })}
                        >
                            Return to login
                        </a>
                    </div>
                </div>
            {:else}
                <div class="flex flex-col gap-6">
                    <p class="text-center text-slate-400">
                        Enter your email address and we'll send you a link to reset your password.
                    </p>

                    <form onsubmit={handleSubmit} class="flex flex-col gap-4">
                        <div class="flex flex-col gap-2">
                            <Label for="email">Email Address</Label>
                            <Input
                                id="email"
                                type="email"
                                bind:value={email}
                                onblur={handleBlur}
                                placeholder="Enter your email address"
                                aria-invalid={Boolean(errors.email && touched.email)}
                                class={cn(
                                    errors.email &&
                                        touched.email &&
                                        'border-red-500 focus-visible:ring-red-500/30'
                                )}
                            />
                            {#if errors.email && touched.email}
                                <p class="text-sm text-red-500">{errors.email}</p>
                            {/if}
                        </div>

                        <button
                            type="submit"
                            disabled={isLoading}
                            class={cn(buttonVariants({ variant: 'default' }), 'w-full')}
                        >
                            {isLoading ? 'Sending Reset Link...' : 'Send Reset Link'}
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
