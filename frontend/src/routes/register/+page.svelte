<script lang="ts">
    import { goto } from '$app/navigation'
    import { resolve } from '$app/paths'
    import PasswordField from '$lib/components/auth/password-field.svelte'
    import PasswordRequirements from '$lib/components/auth/password-requirements.svelte'
    import { Alert, AlertDescription } from '$lib/components/ui/alert'
    import { buttonVariants } from '$lib/components/ui/button'
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card'
    import { Input } from '$lib/components/ui/input'
    import { Label } from '$lib/components/ui/label'
    import { authStore } from '$lib/stores/auth'
    import { cn } from '$lib/utils'

    let username = $state('')
    let email = $state('')
    let password = $state('')
    let confirmPassword = $state('')
    let displayName = $state('')
    let submitted = $state(false)
    let errors = $state<Record<string, string>>({})
    let isLoading = $state(false)
    let touched = $state<Record<string, boolean>>({
        username: false,
        email: false,
        password: false,
        confirmPassword: false,
        displayName: false
    })

    let passwordRequirements = $state({
        length: false,
        lowercase: false,
        uppercase: false,
        number: false,
        special: false
    })

    let registerError = $state('')

    function validateUsername(value: string): string | null {
        if (!value) return 'Username is required'
        if (value.length < 3) return 'Username must be at least 3 characters long'
        if (value.length > 50) return 'Username must be less than 50 characters'
        if (!/^[a-zA-Z0-9._-]+$/.test(value)) {
            return 'Username can only contain letters, numbers, dots, hyphens, and underscores'
        }
        return null
    }

    function validateEmail(value: string): string | null {
        if (!value) return 'Email is required'
        if (!/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(value)) {
            return 'Email format is invalid'
        }
        return null
    }

    function validatePassword(value: string): string | null {
        if (!value) return 'Password is required'
        if (value.length < 8) return 'Password must be at least 8 characters long'
        if (!/(?=.*[a-z])/.test(value)) return 'Password must contain at least one lowercase letter'
        if (!/(?=.*[A-Z])/.test(value)) return 'Password must contain at least one uppercase letter'
        if (!/(?=.*\d)/.test(value)) return 'Password must contain at least one number'
        if (!/(?=.*[!@#$%^&*])/.test(value)) {
            return 'Password must contain at least one special character'
        }
        return null
    }

    function validateConfirmPassword(value: string): string | null {
        if (!value) return 'Please confirm your password'
        if (value !== password) return 'Passwords do not match'
        return null
    }

    function validateDisplayName(value: string): string | null {
        if (!value) return 'Display name is required'
        if (value.length < 2) return 'Display name must be at least 2 characters long'
        if (value.length > 50) return 'Display name must be less than 50 characters'
        return null
    }

    function updatePasswordRequirements(value: string) {
        passwordRequirements = {
            length: value.length >= 8,
            lowercase: /[a-z]/.test(value),
            uppercase: /[A-Z]/.test(value),
            number: /\d/.test(value),
            special: /[!@#$%^&*]/.test(value)
        }
    }

    $effect(() => {
        if (submitted) {
            errors.username = validateUsername(username) || ''
            errors.email = validateEmail(email) || ''
            errors.password = validatePassword(password) || ''
            errors.confirmPassword = validateConfirmPassword(confirmPassword) || ''
            errors.displayName = validateDisplayName(displayName) || ''
        }
    })

    $effect(() => {
        if (touched.username) errors.username = validateUsername(username) || ''
    })

    $effect(() => {
        if (touched.email) errors.email = validateEmail(email) || ''
    })

    $effect(() => {
        updatePasswordRequirements(password)
        if (touched.password) errors.password = validatePassword(password) || ''
    })

    $effect(() => {
        if (touched.confirmPassword) {
            errors.confirmPassword = validateConfirmPassword(confirmPassword) || ''
        }
    })

    $effect(() => {
        if (touched.displayName) errors.displayName = validateDisplayName(displayName) || ''
    })

    function handleBlur(field: keyof typeof touched) {
        touched[field] = true
    }

    async function handleSubmit(event: Event) {
        event.preventDefault()
        submitted = true
        registerError = ''

        Object.keys(touched).forEach((key) => {
            touched[key as keyof typeof touched] = true
        })

        errors = {
            username: validateUsername(username) || '',
            email: validateEmail(email) || '',
            password: validatePassword(password) || '',
            confirmPassword: validateConfirmPassword(confirmPassword) || '',
            displayName: validateDisplayName(displayName) || ''
        }

        const hasErrors = Object.values(errors).some((error) => error !== '')
        if (hasErrors) return

        try {
            isLoading = true
            await authStore.register({ username, email, password, display_name: displayName })

            username = ''
            email = ''
            password = ''
            confirmPassword = ''
            displayName = ''
            submitted = false

            goto(resolve('/login'))
        } catch (error) {
            registerError =
                error instanceof Error ? error.message : 'Registration failed. Please try again.'
        } finally {
            isLoading = false
        }
    }
</script>

<section class="page-shell flex min-h-[calc(100vh-6rem)] items-center justify-center">
    <Card class="surface-card w-full max-w-lg shadow-lg">
        <CardHeader class="gap-3 text-center">
            <CardTitle class="text-3xl">Create an account</CardTitle>
            <p class="text-slate-400">
                Already have an account?
                <a href={resolve('/login')} class={buttonVariants({ variant: 'link', size: 'sm' })}>
                    Login
                </a>
            </p>
        </CardHeader>

        <CardContent>
            {#if registerError}
                <Alert
                    variant="destructive"
                    class="mb-4 border-red-500/60 bg-red-950/50 text-red-200"
                >
                    <AlertDescription>{registerError}</AlertDescription>
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
                        placeholder="Enter username"
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

                <div class="flex flex-col gap-2">
                    <Label for="email">Email</Label>
                    <Input
                        id="email"
                        type="email"
                        bind:value={email}
                        onblur={() => handleBlur('email')}
                        placeholder="Enter email address"
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

                <div class="flex flex-col gap-2">
                    <Label for="displayName">Display Name</Label>
                    <Input
                        id="displayName"
                        type="text"
                        bind:value={displayName}
                        onblur={() => handleBlur('displayName')}
                        placeholder="Enter your full name"
                        aria-invalid={Boolean(errors.displayName && touched.displayName)}
                        class={cn(
                            errors.displayName &&
                                touched.displayName &&
                                'border-red-500 focus-visible:ring-red-500/30'
                        )}
                    />
                    {#if errors.displayName && touched.displayName}
                        <p class="text-sm text-red-500">{errors.displayName}</p>
                    {/if}
                </div>

                <PasswordField
                    id="password"
                    label="Password"
                    bind:value={password}
                    placeholder="Enter password"
                    touched={touched.password}
                    error={errors.password}
                    onblur={() => handleBlur('password')}
                />

                <PasswordRequirements requirements={passwordRequirements} />

                <PasswordField
                    id="confirmPassword"
                    label="Confirm Password"
                    bind:value={confirmPassword}
                    placeholder="Confirm your password"
                    touched={touched.confirmPassword}
                    error={errors.confirmPassword}
                    onblur={() => handleBlur('confirmPassword')}
                />

                <button
                    type="submit"
                    disabled={isLoading}
                    class={cn(buttonVariants({ variant: 'default' }), 'mt-2 w-full')}
                >
                    {isLoading ? 'Creating Account...' : 'Create account'}
                </button>
            </form>
        </CardContent>
    </Card>
</section>
