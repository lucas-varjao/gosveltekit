<!-- frontend/src/routes/login/+page.svelte -->

<script lang="ts">
    import { slide } from 'svelte/transition'
    import { goto } from '$app/navigation'
    import { authStore } from '$lib/stores/auth'

    // State declaration using Svelte 5 runes
    let username = $state('')
    let password = $state('')
    // let submitted = $state(false);
    let errors = $state<Record<string, string>>({})
    let isLoading = $state(false)
    let showPassword = $state(false)
    let touched = $state<Record<string, boolean>>({
        username: false,
        password: false
    })
    let authError = $state('')

    // Validation functions
    function validateUsername(value: string): string | null {
        if (!value) return 'Username is required'
        return null
    }

    function validatePassword(value: string): string | null {
        if (!value) return 'Password is required'
        return null
    }

    // Toggle password visibility
    function togglePasswordVisibility() {
        showPassword = !showPassword
    }

    // Reactive validation using $effect
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

    // Handle input blur to mark field as touched
    function handleBlur(field: keyof typeof touched) {
        touched[field] = true
    }

    // Form submission handler
    async function handleSubmit(event: Event) {
        event.preventDefault()
        // submitted = true;
        authError = ''

        // Mark all fields as touched
        Object.keys(touched).forEach((key) => {
            touched[key as keyof typeof touched] = true
        })

        // Perform validation
        const usernameError = validateUsername(username)
        const passwordError = validatePassword(password)

        // Update errors state
        errors = {
            username: usernameError || '',
            password: passwordError || ''
        }

        // Check if there are any validation errors
        const hasErrors = Object.values(errors).some((error) => error !== '')

        if (!hasErrors) {
            try {
                isLoading = true

                // Use the auth store to login
                await authStore.login(username, password)

                // Redirect to home page after successful login
                goto('/logged')
            } catch (error) {
                console.error('Login error:', error)
                authError =
                    error instanceof Error ? error.message : 'Login failed. Please try again.'
            } finally {
                isLoading = false
            }
        }
    }
</script>

<!-- Using flexbox for main page layout -->
<section class="flex min-h-[calc(100vh-6rem)] items-center justify-center px-4 py-12">
    <div class="w-full max-w-md rounded border border-slate-800 bg-slate-900 p-8 shadow-lg">
        <!-- Using flexbox for vertical content alignment -->
        <div class="flex flex-col gap-6">
            <h1 class="text-center text-3xl font-bold text-white">Sign In</h1>

            <!-- Register Link -->
            <div class="text-center text-slate-400">
                Don't have an account?
                <a href="/register" class="font-medium text-blue-500 hover:text-blue-400">
                    Create an account
                </a>
            </div>

            {#if authError}
                <div
                    transition:slide
                    class="rounded border border-red-500 bg-red-900/50 px-4 py-3 text-red-300"
                    role="alert"
                >
                    <p>{authError}</p>
                </div>
            {/if}

            <!-- Using flexbox for the form instead of grid -->
            <form onsubmit={handleSubmit} class="flex flex-col gap-4">
                <!-- Username Field -->
                <div class="flex flex-col gap-2">
                    <label for="username" class="text-sm font-medium text-slate-200">
                        Username
                    </label>
                    <input
                        type="text"
                        id="username"
                        bind:value={username}
                        onblur={() => handleBlur('username')}
                        placeholder="Enter your username"
                        class="w-full rounded border-2 bg-slate-800 px-3 py-2 text-white {errors.username &&
                        touched.username
                            ? 'border-red-500'
                            : 'border-slate-700'} focus:ring-2 focus:ring-blue-500 focus:outline-none"
                    />
                    {#if errors.username && touched.username}
                        <p transition:slide class="text-sm text-red-500">{errors.username}</p>
                    {/if}
                </div>

                <!-- Password Field -->
                <div class="flex flex-col gap-2">
                    <label for="password" class="text-sm font-medium text-slate-200">
                        Password
                    </label>
                    <div class="relative">
                        <input
                            type={showPassword ? 'text' : 'password'}
                            id="password"
                            bind:value={password}
                            onblur={() => handleBlur('password')}
                            placeholder="Enter password"
                            class="w-full rounded border-2 bg-slate-800 px-3 py-2 text-white {errors.password &&
                            touched.password
                                ? 'border-red-500'
                                : 'border-slate-700'} pr-10 focus:ring-2 focus:ring-blue-500 focus:outline-none"
                        />
                        <button
                            type="button"
                            class="absolute inset-y-0 right-0 flex items-center pr-3 text-slate-400 hover:text-slate-200 focus:outline-none"
                            onclick={togglePasswordVisibility}
                        >
                            {#if showPassword}
                                <!-- Eye slash icon for hide password -->
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-5 w-5"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                                    />
                                </svg>
                            {:else}
                                <!-- Eye icon for show password -->
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-5 w-5"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                                    />
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                                    />
                                </svg>
                            {/if}
                        </button>
                    </div>
                    {#if errors.password && touched.password}
                        <p transition:slide class="text-sm text-red-500">{errors.password}</p>
                    {/if}
                </div>

                <!-- Forgot Password Link -->
                <div class="flex justify-end">
                    <a href="/forgot-password" class="text-sm text-blue-500 hover:text-blue-400">
                        Forgot your password?
                    </a>
                </div>

                <!-- Submit Button -->
                <div class="mt-2">
                    <button
                        type="submit"
                        disabled={isLoading}
                        class="w-full rounded bg-blue-600 px-4 py-2 font-semibold text-white transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:outline-none disabled:cursor-not-allowed disabled:opacity-70"
                    >
                        {#if isLoading}
                            <span class="inline-block animate-pulse">Signing in...</span>
                        {:else}
                            Sign in
                        {/if}
                    </button>
                </div>
            </form>
        </div>
    </div>
</section>
