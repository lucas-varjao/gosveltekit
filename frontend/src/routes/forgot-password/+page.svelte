<!-- frontend/src/routes/forgot-password/+page.svelte -->

<script lang="ts">
    import { slide } from 'svelte/transition'
    import { authApi } from '$lib/api/auth'

    // State declaration using Svelte 5 runes
    let email = $state('')
    let errors = $state<Record<string, string>>({})
    let isLoading = $state(false)
    let touched = $state<Record<string, boolean>>({ email: false })
    let submitted = $state(false)
    let success = $state(false)

    // Validation function
    function validateEmail(value: string): string | null {
        if (!value) return 'Email is required'

        // Basic email validation regex
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        if (!emailRegex.test(value)) {
            return 'Please enter a valid email address'
        }

        return null
    }

    // Reactive validation using $effect
    $effect(() => {
        if (touched.email) {
            errors.email = validateEmail(email) || ''
        }
    })

    // Handle input blur to mark field as touched
    function handleBlur() {
        touched.email = true
    }

    // Form submission handler
    async function handleSubmit(event: Event) {
        event.preventDefault()

        // Mark email as touched for validation
        touched.email = true

        // Update error state
        errors.email = validateEmail(email) || ''

        // Check if there are any validation errors
        if (!errors.email) {
            try {
                isLoading = true

                await authApi.requestPasswordReset({ email })

                // Show success message
                success = true
                submitted = true
            } catch (error) {
                console.error('Password reset request error:', error)
                errors.email =
                    error instanceof Error ? error.message : 'Failed to request password reset'
            } finally {
                isLoading = false
            }
        }
    }

    // Function to start over
    function resetForm() {
        email = ''
        errors = {}
        touched = { email: false }
        submitted = false
        success = false
    }
</script>

<!-- Using flexbox for main page layout -->
<section class="flex min-h-[calc(100vh-6rem)] items-center justify-center px-4 py-12">
    <div class="w-full max-w-md rounded border border-slate-800 bg-slate-900 p-8 shadow-lg">
        <!-- Using flexbox for vertical content alignment -->
        <div class="flex flex-col gap-6">
            <h1 class="text-center text-3xl font-bold text-white">Reset Password</h1>

            {#if submitted && success}
                <div class="flex flex-col gap-6">
                    <div
                        transition:slide
                        class="rounded border border-blue-500 bg-blue-900/50 px-4 py-3 text-blue-300"
                        role="alert"
                    >
                        <p class="font-medium">Check your email</p>
                        <p class="mt-1">
                            If an account exists with the email {email}, you will receive a password
                            reset link shortly.
                        </p>
                    </div>

                    <div class="flex flex-col items-center gap-4">
                        <button
                            onclick={resetForm}
                            class="font-medium text-blue-500 hover:text-blue-400"
                        >
                            Request another reset link
                        </button>

                        <a href="/login" class="text-slate-400 hover:text-slate-300">
                            Return to login
                        </a>
                    </div>
                </div>
            {:else}
                <div class="flex flex-col gap-6">
                    <p class="text-center text-slate-400">
                        Enter your email address and we'll send you a link to reset your password.
                    </p>

                    <!-- Using flexbox for the form -->
                    <form onsubmit={handleSubmit} class="flex flex-col gap-4">
                        <!-- Email Field -->
                        <div class="flex flex-col gap-2">
                            <label for="email" class="text-sm font-medium text-slate-200">
                                Email Address
                            </label>
                            <input
                                type="email"
                                id="email"
                                bind:value={email}
                                onblur={handleBlur}
                                placeholder="Enter your email address"
                                class="w-full rounded border-2 bg-slate-800 px-3 py-2 text-white {errors.email &&
                                touched.email
                                    ? 'border-red-500'
                                    : 'border-slate-700'} focus:ring-2 focus:ring-blue-500 focus:outline-none"
                            />
                            {#if errors.email && touched.email}
                                <p transition:slide class="text-sm text-red-500">{errors.email}</p>
                            {/if}
                        </div>

                        <!-- Submit Button -->
                        <div>
                            <button
                                type="submit"
                                disabled={isLoading}
                                class="w-full rounded bg-blue-600 px-4 py-2 font-semibold text-white transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:outline-none disabled:cursor-not-allowed disabled:opacity-70"
                            >
                                {#if isLoading}
                                    <span class="inline-block animate-pulse"
                                        >Sending Reset Link...</span
                                    >
                                {:else}
                                    Send Reset Link
                                {/if}
                            </button>
                        </div>
                    </form>

                    <!-- Back to Login -->
                    <div class="flex justify-center">
                        <a href="/login" class="font-medium text-blue-500 hover:text-blue-400">
                            Back to Login
                        </a>
                    </div>
                </div>
            {/if}
        </div>
    </div>
</section>
