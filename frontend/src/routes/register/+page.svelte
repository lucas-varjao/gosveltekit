<!-- frontend/src/routes/register/+page.svelte -->

<script lang="ts">
    import { goto } from '$app/navigation'
    import { slide } from 'svelte/transition'
    import { authStore } from '$lib/stores/auth'

    // State declaration using Svelte 5 runes
    let username = $state('')
    let email = $state('')
    let password = $state('')
    let confirmPassword = $state('')
    let displayName = $state('')
    let submitted = $state(false)
    let errors = $state<Record<string, string>>({})
    let isLoading = $state(false)
    let showPassword = $state(false)
    let showConfirmPassword = $state(false)
    let touched = $state<Record<string, boolean>>({
        username: false,
        email: false,
        password: false,
        confirmPassword: false,
        displayName: false
    })

    // Password requirement states
    let passwordRequirements = $state({
        length: false,
        lowercase: false,
        uppercase: false,
        number: false,
        special: false
    })

    let registerError = $state('')

    // Validation functions
    function validateUsername(value: string): string | null {
        if (!value) return 'Username is required'
        if (value.length < 3) return 'Username must be at least 3 characters long'
        if (value.length > 50) return 'Username must be less than 50 characters'
        if (!/^[a-zA-Z0-9._-]+$/.test(value))
            return 'Username can only contain letters, numbers, dots, hyphens, and underscores'
        return null
    }

    function validateEmail(value: string): string | null {
        if (!value) return 'Email is required'
        if (!/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(value))
            return 'Email format is invalid'
        return null
    }

    function validatePassword(value: string): string | null {
        if (!value) return 'Password is required'
        if (value.length < 8) return 'Password must be at least 8 characters long'
        if (!/(?=.*[a-z])/.test(value)) return 'Password must contain at least one lowercase letter'
        if (!/(?=.*[A-Z])/.test(value)) return 'Password must contain at least one uppercase letter'
        if (!/(?=.*\d)/.test(value)) return 'Password must contain at least one number'
        if (!/(?=.*[!@#$%^&*])/.test(value))
            return 'Password must contain at least one special character'
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

    // Individual password requirement validation functions
    function checkPasswordLength(value: string): boolean {
        return value.length >= 8
    }

    function checkLowercase(value: string): boolean {
        return /[a-z]/.test(value)
    }

    function checkUppercase(value: string): boolean {
        return /[A-Z]/.test(value)
    }

    function checkNumber(value: string): boolean {
        return /\d/.test(value)
    }

    function checkSpecialChar(value: string): boolean {
        return /[!@#$%^&*]/.test(value)
    }

    // Update password requirements check
    function updatePasswordRequirements(value: string) {
        passwordRequirements = {
            length: checkPasswordLength(value),
            lowercase: checkLowercase(value),
            uppercase: checkUppercase(value),
            number: checkNumber(value),
            special: checkSpecialChar(value)
        }
    }

    // Toggle password visibility
    function togglePasswordVisibility() {
        showPassword = !showPassword
    }

    function toggleConfirmPasswordVisibility() {
        showConfirmPassword = !showConfirmPassword
    }

    // Reactive validation
    $effect(() => {
        if (submitted) {
            errors.username = validateUsername(username) || ''
            errors.email = validateEmail(email) || ''
            errors.password = validatePassword(password) || ''
            errors.confirmPassword = validateConfirmPassword(confirmPassword) || ''
            errors.displayName = validateDisplayName(displayName) || ''
        }
    })

    // Reactive validation using $effect
    $effect(() => {
        // Validate username when it changes and has been touched
        if (touched.username) {
            errors.username = validateUsername(username) || ''
        }
    })

    $effect(() => {
        // Validate email when it changes and has been touched
        if (touched.email) {
            errors.email = validateEmail(email) || ''
        }
    })

    $effect(() => {
        // Update password requirements whenever password changes
        updatePasswordRequirements(password)

        // Validate password when it changes and has been touched
        if (touched.password) {
            errors.password = validatePassword(password) || ''
        }
    })

    $effect(() => {
        // Validate confirmPassword when it or password changes and has been touched
        if (touched.confirmPassword) {
            errors.confirmPassword = validateConfirmPassword(confirmPassword) || ''
        }
    })

    $effect(() => {
        // Validate displayName when it changes and has been touched
        if (touched.displayName) {
            errors.displayName = validateDisplayName(displayName) || ''
        }
    })

    // Handle input blur to mark field as touched
    function handleBlur(field: keyof typeof touched) {
        touched[field] = true
    }

    // Form submission handler
    async function handleSubmit(event: Event) {
        event.preventDefault()
        submitted = true

        registerError = ''

        // Mark all fields as touched
        Object.keys(touched).forEach((key) => {
            touched[key as keyof typeof touched] = true
        })

        // Perform validation
        const usernameError = validateUsername(username)
        const emailError = validateEmail(email)
        const passwordError = validatePassword(password)
        const confirmPasswordError = validateConfirmPassword(confirmPassword)
        const displayNameError = validateDisplayName(displayName)

        // Update errors state
        errors = {
            username: usernameError || '',
            email: emailError || '',
            password: passwordError || '',
            confirmPassword: confirmPasswordError || '',
            displayName: displayNameError || ''
        }

        // Check if there are any validation errors
        const hasErrors = Object.values(errors).some((error) => error !== '')

        if (!hasErrors) {
            try {
                isLoading = true

                // Simulate API call for now (will be connected to the backend later)
                await authStore.register({ username, email, password, display_name: displayName })

                // Reset form after successful submission
                username = ''
                email = ''
                password = ''
                confirmPassword = ''
                displayName = ''
                submitted = false

                // Show success message or redirect
                goto('/login')
            } catch (error) {
                console.error('Registration error:', error)
                registerError =
                    error instanceof Error
                        ? error.message
                        : 'Registration failed. Please try again.'
            } finally {
                isLoading = false
            }
        }
    }
</script>

<!-- Using flexbox for main page layout -->
<section class="flex min-h-[calc(100vh-6rem)] items-center justify-center px-4 py-12">
    <div class="w-full max-w-lg rounded border border-slate-800 bg-slate-900 p-8 shadow-lg">
        <!-- Using flexbox for vertical content alignment -->
        <div class="flex flex-col gap-6">
            <h1 class="text-center text-3xl font-bold text-white">Create an account</h1>

            <!-- Login Link -->
            <div class="text-center text-slate-400">
                Already have an account?
                <a href="/login" class="font-medium text-blue-500 hover:text-blue-400"> Login </a>
            </div>

            {#if registerError}
                <div
                    transition:slide
                    class="rounded border border-red-500 bg-red-900/50 px-4 py-3 text-red-300"
                    role="alert"
                >
                    <p>{registerError}</p>
                </div>
            {/if}

            <!-- Using flexbox layout for the form -->
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
                        placeholder="Enter username"
                        class="w-full rounded border-2 bg-slate-800 px-3 py-2 text-white {errors.username &&
                        touched.username
                            ? 'border-red-500'
                            : 'border-slate-700'} focus:ring-2 focus:ring-blue-500 focus:outline-none"
                    />
                    {#if errors.username && touched.username}
                        <p transition:slide class="text-sm text-red-500">{errors.username}</p>
                    {/if}
                </div>

                <!-- Email Field -->
                <div class="flex flex-col gap-2">
                    <label for="email" class="text-sm font-medium text-slate-200"> Email </label>
                    <input
                        type="email"
                        id="email"
                        bind:value={email}
                        onblur={() => handleBlur('email')}
                        placeholder="Enter email address"
                        class="w-full rounded border-2 bg-slate-800 px-3 py-2 text-white {errors.email &&
                        touched.email
                            ? 'border-red-500'
                            : 'border-slate-700'} focus:ring-2 focus:ring-blue-500 focus:outline-none"
                    />
                    {#if errors.email && touched.email}
                        <p transition:slide class="text-sm text-red-500">{errors.email}</p>
                    {/if}
                </div>

                <!-- Display Name Field -->
                <div class="flex flex-col gap-2">
                    <label for="displayName" class="text-sm font-medium text-slate-200">
                        Display Name
                    </label>
                    <input
                        type="text"
                        id="displayName"
                        bind:value={displayName}
                        onblur={() => handleBlur('displayName')}
                        placeholder="Enter your full name"
                        class="w-full rounded border-2 bg-slate-800 px-3 py-2 text-white {errors.displayName &&
                        touched.displayName
                            ? 'border-red-500'
                            : 'border-slate-700'} focus:ring-2 focus:ring-blue-500 focus:outline-none"
                    />
                    {#if errors.displayName && touched.displayName}
                        <p transition:slide class="text-sm text-red-500">{errors.displayName}</p>
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

                    <!-- Password requirements list -->
                    <div class="flex flex-col gap-1">
                        <p class="text-sm text-slate-300">Your password must contain:</p>
                        <ul class="ml-2 flex flex-col gap-1">
                            <li class="flex items-center gap-1">
                                <span
                                    class={passwordRequirements.length
                                        ? 'text-green-500'
                                        : 'text-slate-400'}
                                >
                                    {#if passwordRequirements.length}
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-4 w-4"
                                            viewBox="0 0 20 20"
                                            fill="currentColor"
                                        >
                                            <path
                                                fill-rule="evenodd"
                                                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                                clip-rule="evenodd"
                                            />
                                        </svg>
                                    {:else}
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-4 w-4"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                            />
                                        </svg>
                                    {/if}
                                </span>
                                <span
                                    class={passwordRequirements.length
                                        ? 'text-green-500'
                                        : 'text-slate-400'}>No mínimo 8 caracteres</span
                                >
                            </li>
                            <li class="flex items-center gap-1">
                                <span
                                    class={passwordRequirements.lowercase
                                        ? 'text-green-500'
                                        : 'text-slate-400'}
                                >
                                    {#if passwordRequirements.lowercase}
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-4 w-4"
                                            viewBox="0 0 20 20"
                                            fill="currentColor"
                                        >
                                            <path
                                                fill-rule="evenodd"
                                                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                                clip-rule="evenodd"
                                            />
                                        </svg>
                                    {:else}
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-4 w-4"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                            />
                                        </svg>
                                    {/if}
                                </span>
                                <span
                                    class={passwordRequirements.lowercase
                                        ? 'text-green-500'
                                        : 'text-slate-400'}>Uma letra minúscula</span
                                >
                            </li>
                            <li class="flex items-center gap-1">
                                <span
                                    class={passwordRequirements.uppercase
                                        ? 'text-green-500'
                                        : 'text-slate-400'}
                                >
                                    {#if passwordRequirements.uppercase}
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-4 w-4"
                                            viewBox="0 0 20 20"
                                            fill="currentColor"
                                        >
                                            <path
                                                fill-rule="evenodd"
                                                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                                clip-rule="evenodd"
                                            />
                                        </svg>
                                    {:else}
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-4 w-4"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                            />
                                        </svg>
                                    {/if}
                                </span>
                                <span
                                    class={passwordRequirements.uppercase
                                        ? 'text-green-500'
                                        : 'text-slate-400'}>Uma letra maiúscula</span
                                >
                            </li>
                            <li class="flex items-center gap-1">
                                <span
                                    class={passwordRequirements.number
                                        ? 'text-green-500'
                                        : 'text-slate-400'}
                                >
                                    {#if passwordRequirements.number}
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-4 w-4"
                                            viewBox="0 0 20 20"
                                            fill="currentColor"
                                        >
                                            <path
                                                fill-rule="evenodd"
                                                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                                clip-rule="evenodd"
                                            />
                                        </svg>
                                    {:else}
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-4 w-4"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                            />
                                        </svg>
                                    {/if}
                                </span>
                                <span
                                    class={passwordRequirements.number
                                        ? 'text-green-500'
                                        : 'text-slate-400'}>Um número</span
                                >
                            </li>
                            <li class="flex items-center gap-1">
                                <span
                                    class={passwordRequirements.special
                                        ? 'text-green-500'
                                        : 'text-slate-400'}
                                >
                                    {#if passwordRequirements.special}
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-4 w-4"
                                            viewBox="0 0 20 20"
                                            fill="currentColor"
                                        >
                                            <path
                                                fill-rule="evenodd"
                                                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                                clip-rule="evenodd"
                                            />
                                        </svg>
                                    {:else}
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-4 w-4"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                            />
                                        </svg>
                                    {/if}
                                </span>
                                <span
                                    class={passwordRequirements.special
                                        ? 'text-green-500'
                                        : 'text-slate-400'}>Um caractere especial (!@#$%^&*)</span
                                >
                            </li>
                        </ul>
                    </div>

                    {#if errors.password && touched.password}
                        <p transition:slide class="text-sm text-red-500">{errors.password}</p>
                    {/if}
                </div>

                <!-- Confirm Password Field -->
                <div class="flex flex-col gap-2">
                    <label for="confirmPassword" class="text-sm font-medium text-slate-200">
                        Confirm Password
                    </label>
                    <div class="relative">
                        <input
                            type={showConfirmPassword ? 'text' : 'password'}
                            id="confirmPassword"
                            bind:value={confirmPassword}
                            onblur={() => handleBlur('confirmPassword')}
                            placeholder="Confirm your password"
                            class="w-full rounded border-2 bg-slate-800 px-3 py-2 text-white {errors.confirmPassword &&
                            touched.confirmPassword
                                ? 'border-red-500'
                                : 'border-slate-700'} pr-10 focus:ring-2 focus:ring-blue-500 focus:outline-none"
                        />
                        <button
                            type="button"
                            class="absolute inset-y-0 right-0 flex items-center pr-3 text-slate-400 hover:text-slate-200 focus:outline-none"
                            onclick={toggleConfirmPasswordVisibility}
                        >
                            {#if showConfirmPassword}
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
                    {#if errors.confirmPassword && touched.confirmPassword}
                        <p transition:slide class="text-sm text-red-500">
                            {errors.confirmPassword}
                        </p>
                    {/if}
                </div>

                <!-- Submit Button -->
                <div class="mt-2">
                    <button
                        type="submit"
                        disabled={isLoading}
                        class="w-full rounded bg-blue-600 px-4 py-2 font-semibold text-white transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:outline-none disabled:cursor-not-allowed disabled:opacity-70"
                    >
                        {#if isLoading}
                            <span class="inline-block animate-pulse">Creating Account...</span>
                        {:else}
                            Create account
                        {/if}
                    </button>
                </div>
            </form>
        </div>
    </div>
</section>
