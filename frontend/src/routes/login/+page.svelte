<script lang="ts">
    import { goto } from '$app/navigation'
    import { resolve } from '$app/paths'
    import PasswordField from '$lib/components/auth/password-field.svelte'
    import { Alert, AlertDescription } from '$lib/components/ui/alert'
    import { buttonVariants } from '$lib/components/ui/button'
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card'
    import { Input } from '$lib/components/ui/input'
    import { Label } from '$lib/components/ui/label'
    import { authStore } from '$lib/stores/auth'
    import { cn } from '$lib/utils'

    let username = $state('')
    let password = $state('')
    let errors = $state<Record<string, string>>({})
    let isLoading = $state(false)
    let touched = $state<Record<string, boolean>>({
        username: false,
        password: false
    })
    let authError = $state('')

    function validateUsername(value: string): string | null {
        if (!value) return 'Username is required'
        return null
    }

    function validatePassword(value: string): string | null {
        if (!value) return 'Password is required'
        return null
    }

    $effect(() => {
        if (touched.username) {
            errors.username = validateUsername(username) || ''
        }
    })

    $effect(() => {
        if (touched.password) {
            errors.password = validatePassword(password) || ''
        }
    })

    function handleBlur(field: keyof typeof touched) {
        touched[field] = true
    }

    async function handleSubmit(event: Event) {
        event.preventDefault()
        authError = ''

        Object.keys(touched).forEach((key) => {
            touched[key as keyof typeof touched] = true
        })

        errors = {
            username: validateUsername(username) || '',
            password: validatePassword(password) || ''
        }

        const hasErrors = Object.values(errors).some((error) => error !== '')
        if (hasErrors) return

        try {
            isLoading = true
            await authStore.login(username, password)
            goto(resolve('/profile'))
        } catch (error) {
            authError = error instanceof Error ? error.message : 'Login failed. Please try again.'
        } finally {
            isLoading = false
        }
    }
</script>

<section class="flex min-h-[calc(100vh-6rem)] items-center justify-center px-4 py-12">
    <Card class="w-full max-w-md border-slate-800 bg-slate-900 shadow-lg">
        <CardHeader class="gap-3 text-center">
            <CardTitle class="text-3xl">Sign In</CardTitle>
            <p class="text-slate-400">
                Don't have an account?
                <a
                    href={resolve('/register')}
                    class={buttonVariants({ variant: 'link', size: 'sm' })}
                >
                    Create an account
                </a>
            </p>
        </CardHeader>

        <CardContent>
            {#if authError}
                <Alert
                    variant="destructive"
                    class="mb-4 border-red-500/60 bg-red-950/50 text-red-200"
                >
                    <AlertDescription>{authError}</AlertDescription>
                </Alert>
            {/if}

            <form onsubmit={handleSubmit} class="flex flex-col gap-4">
                <div class="flex flex-col gap-2">
                    <Label for="username">Username</Label>
                    <Input
                        id="username"
                        type="text"
                        bind:value={username}
                        onblur={() => handleBlur('username')}
                        placeholder="Enter your username"
                        aria-invalid={Boolean(errors.username && touched.username)}
                        class={cn(
                            errors.username &&
                                touched.username &&
                                'border-red-500 focus-visible:ring-red-500/30'
                        )}
                    />
                    {#if errors.username && touched.username}
                        <p class="text-sm text-red-500">{errors.username}</p>
                    {/if}
                </div>

                <PasswordField
                    id="password"
                    label="Password"
                    bind:value={password}
                    placeholder="Enter your password"
                    touched={touched.password}
                    error={errors.password}
                    onblur={() => handleBlur('password')}
                />

                <div class="flex justify-end">
                    <a
                        href={resolve('/forgot-password')}
                        class={buttonVariants({ variant: 'link', size: 'sm' })}
                    >
                        Forgot your password?
                    </a>
                </div>

                <button
                    type="submit"
                    disabled={isLoading}
                    class={cn(buttonVariants({ variant: 'default' }), 'mt-2 w-full')}
                >
                    {isLoading ? 'Signing in...' : 'Sign in'}
                </button>
            </form>
        </CardContent>
    </Card>
</section>
