<script lang="ts">
    import { slide } from 'svelte/transition';
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';

    // State declaration using Svelte 5 runes
    let password = $state('');
    let confirmPassword = $state('');
    let token = $state('');
    let errors = $state<Record<string, string>>({});
    let isLoading = $state(false);
    let showPassword = $state(false);
    let showConfirmPassword = $state(false);
    let touched = $state<Record<string, boolean>>({
        password: false,
        confirmPassword: false
    });
    let resetError = $state('');
    let resetSuccess = $state(false);

    // Password requirements state
    let passwordRequirements = $state({
        length: false,
        uppercase: false,
        lowercase: false,
        number: false,
        special: false
    });

    // Get token from URL when page loads
    $effect(() => {
        token = $page.url.searchParams.get('token') || '';
        
        // Redirect if no token is provided
        if (!token) {
            resetError = 'Invalid or missing password reset token';
        }
    });

    // Toggle password visibility
    function togglePasswordVisibility() {
        showPassword = !showPassword;
    }

    // Toggle confirm password visibility
    function toggleConfirmPasswordVisibility() {
        showConfirmPassword = !showConfirmPassword;
    }

    // Update password requirements as user types
    $effect(() => {
        passwordRequirements.length = password.length >= 8;
        passwordRequirements.uppercase = /[A-Z]/.test(password);
        passwordRequirements.lowercase = /[a-z]/.test(password);
        passwordRequirements.number = /[0-9]/.test(password);
        passwordRequirements.special = /[!@#$%^&*(),.?":{}|<>]/.test(password);
    });

    // Validation functions
    function validatePassword(value: string): string | null {
        if (!value) return 'Password is required';
        
        if (value.length < 8) {
            return 'Password must be at least 8 characters long';
        }
        
        if (!/[A-Z]/.test(value)) {
            return 'Password must contain at least one uppercase letter';
        }
        
        if (!/[a-z]/.test(value)) {
            return 'Password must contain at least one lowercase letter';
        }
        
        if (!/[0-9]/.test(value)) {
            return 'Password must contain at least one number';
        }
        
        if (!/[!@#$%^&*(),.?":{}|<>]/.test(value)) {
            return 'Password must contain at least one special character';
        }
        
        return null;
    }

    function validateConfirmPassword(value: string): string | null {
        if (!value) return 'Please confirm your password';
        
        if (value !== password) {
            return 'Passwords do not match';
        }
        
        return null;
    }

    // Reactive validation using $effect
    $effect(() => {
        if (touched.password) {
            errors.password = validatePassword(password) || '';
        }
    });

    $effect(() => {
        if (touched.confirmPassword) {
            errors.confirmPassword = validateConfirmPassword(confirmPassword) || '';
        }
    });

    // Handle input blur to mark field as touched
    function handleBlur(field: keyof typeof touched) {
        touched[field] = true;
    }

    // Form submission handler
    async function handleSubmit(event: Event) {
        event.preventDefault();
        resetError = '';
        
        // Mark all fields as touched
        Object.keys(touched).forEach(key => {
            touched[key as keyof typeof touched] = true;
        });
        
        // Perform validation
        const passwordError = validatePassword(password);
        const confirmPasswordError = validateConfirmPassword(confirmPassword);
        
        // Update errors state
        errors = {
            password: passwordError || '',
            confirmPassword: confirmPasswordError || ''
        };
        
        // Check if there are any validation errors
        const hasErrors = Object.values(errors).some(error => error !== '');
        
        if (!hasErrors && token) {
            try {
                isLoading = true;
                
                // Simulate API call for now (will be connected to the backend later)
                // The actual implementation would look like this:
                //
                // const response = await fetch('/api/auth/password-reset', {
                //   method: 'POST',
                //   headers: { 'Content-Type': 'application/json' },
                //   body: JSON.stringify({ 
                //     token,
                //     new_password: password,
                //     confirm_password: confirmPassword
                //   })
                // });
                // 
                // if (!response.ok) {
                //   const data = await response.json();
                //   throw new Error(data.error || 'Password reset failed');
                // }
                
                await new Promise(resolve => setTimeout(resolve, 1500));
                
                // Show success message
                resetSuccess = true;
                
                // Will redirect after a delay
                setTimeout(() => {
                    goto('/login');
                }, 3000);
                
            } catch (error) {
                console.error('Password reset error:', error);
                resetError = error instanceof Error ? error.message : 'Password reset failed. Please try again.';
            } finally {
                isLoading = false;
            }
        }
    }
</script>

<section class="py-12">
    <div class="max-w-md mx-auto bg-slate-900 rounded border border-slate-800 shadow-lg p-8">
        <h1 class="text-3xl font-bold mb-6 text-center text-white">Reset Your Password</h1>
        
        {#if resetSuccess}
            <div transition:slide class="bg-green-900/50 border border-green-500 text-green-300 px-4 py-3 rounded" role="alert">
                <p class="font-medium">Password reset successful!</p>
                <p class="mt-1">You will be redirected to the login page in a moment.</p>
            </div>
        {:else if resetError}
            <div transition:slide class="bg-red-900/50 border border-red-500 text-red-300 px-4 py-3 rounded mb-4" role="alert">
                <p>{resetError}</p>
            </div>
            
            <div class="text-center mt-6">
                <a href="/forgot-password" class="text-blue-500 hover:text-blue-400 font-medium">
                    Request a new reset link
                </a>
            </div>
        {:else}
            <div class="space-y-5">
                <p class="text-slate-400 text-center mb-6">
                    Please enter and confirm your new password below.
                </p>
                
                <form onsubmit={handleSubmit} class="space-y-6">
                    <!-- Password Field -->
                    <div class="space-y-2">
                        <label for="password" class="block text-sm font-medium text-slate-200">
                            New Password
                        </label>
                        <div class="relative">
                            <input
                                type={showPassword ? "text" : "password"}
                                id="password"
                                bind:value={password}
                                onblur={() => handleBlur('password')}
                                placeholder="Enter your new password"
                                class="w-full px-3 py-2 bg-slate-800 text-white border-2 rounded {errors.password && touched.password ? 'border-red-500' : 'border-slate-700'} focus:outline-none focus:ring-2 focus:ring-blue-500 pr-10"
                            />
                            <button 
                                type="button" 
                                class="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-200 focus:outline-none" 
                                onclick={togglePasswordVisibility}
                            >
                                {#if showPassword}
                                    <!-- Eye slash icon for hide password -->
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                                    </svg>
                                {:else}
                                    <!-- Eye icon for show password -->
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                                    </svg>
                                {/if}
                            </button>
                        </div>
                        {#if errors.password && touched.password}
                            <p transition:slide class="text-sm text-red-500 mt-1">{errors.password}</p>
                        {/if}
                    </div>
                    
                    <!-- Password requirements list -->
                    <div class="mt-2 space-y- text-sm">
                        <p class="text-slate-300">Your password must contain:</p>
                        <ul class="ml-2 space-y-1">
                            <li class="flex items-center gap-1">
                                <span class={passwordRequirements.length ? "text-green-500" : "text-slate-400"}>
                                    {#if passwordRequirements.length}
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                                            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                                        </svg>
                                    {:else}
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                        </svg>
                                    {/if}
                                </span>
                                <span class={passwordRequirements.length ? "text-slate-200" : "text-slate-400"}>At least 8 characters</span>
                            </li>
                            <li class="flex items-center gap-1">
                                <span class={passwordRequirements.uppercase ? "text-green-500" : "text-slate-400"}>
                                    {#if passwordRequirements.uppercase}
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                                            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                                        </svg>
                                    {:else}
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                        </svg>
                                    {/if}
                                </span>
                                <span class={passwordRequirements.uppercase ? "text-slate-200" : "text-slate-400"}>One uppercase letter</span>
                            </li>
                            <li class="flex items-center gap-1">
                                <span class={passwordRequirements.lowercase ? "text-green-500" : "text-slate-400"}>
                                    {#if passwordRequirements.lowercase}
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                                            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                                        </svg>
                                    {:else}
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                        </svg>
                                    {/if}
                                </span>
                                <span class={passwordRequirements.lowercase ? "text-slate-200" : "text-slate-400"}>One lowercase letter</span>
                            </li>
                            <li class="flex items-center gap-1">
                                <span class={passwordRequirements.number ? "text-green-500" : "text-slate-400"}>
                                    {#if passwordRequirements.number}
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                                            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                                        </svg>
                                    {:else}
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                        </svg>
                                    {/if}
                                </span>
                                <span class={passwordRequirements.number ? "text-slate-200" : "text-slate-400"}>One number</span>
                            </li>
                            <li class="flex items-center gap-1">
                                <span class={passwordRequirements.special ? "text-green-500" : "text-slate-400"}>
                                    {#if passwordRequirements.special}
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                                            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                                        </svg>
                                    {:else}
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                        </svg>
                                    {/if}
                                </span>
                                <span class={passwordRequirements.special ? "text-slate-200" : "text-slate-400"}>One special character</span>
                            </li>
                        </ul>
                    </div>
                    
                    <!-- Confirm Password Field -->
                    <div class="space-y-2">
                        <label for="confirmPassword" class="block text-sm font-medium text-slate-200">
                            Confirm Password
                        </label>
                        <div class="relative">
                            <input
                                type={showConfirmPassword ? "text" : "password"}
                                id="confirmPassword"
                                bind:value={confirmPassword}
                                onblur={() => handleBlur('confirmPassword')}
                                placeholder="Confirm your new password"
                                class="w-full px-3 py-2 bg-slate-800 text-white border-2 rounded {errors.confirmPassword && touched.confirmPassword ? 'border-red-500' : 'border-slate-700'} focus:outline-none focus:ring-2 focus:ring-blue-500 pr-10"
                            />
                            <button 
                                type="button" 
                                class="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-200 focus:outline-none" 
                                onclick={toggleConfirmPasswordVisibility}
                            >
                                {#if showConfirmPassword}
                                    <!-- Eye slash icon for hide password -->
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                                    </svg>
                                {:else}
                                    <!-- Eye icon for show password -->
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                                    </svg>
                                {/if}
                            </button>
                        </div>
                        {#if errors.confirmPassword && touched.confirmPassword}
                            <p transition:slide class="text-sm text-red-500 mt-1">{errors.confirmPassword}</p>
                        {/if}
                    </div>
                    
                    <!-- Submit Button -->
                    <div class="pt-2">
                        <button
                            type="submit"
                            disabled={isLoading}
                            class="w-full py-2 px-4 bg-blue-600 hover:bg-blue-700 text-white rounded font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-70 disabled:cursor-not-allowed"
                        >
                            {#if isLoading}
                                <span class="inline-block animate-pulse">Resetting Password...</span>
                            {:else}
                                Reset Password
                            {/if}
                        </button>
                    </div>
                </form>
                
                <!-- Back to Login -->
                <div class="text-center mt-4">
                    <a href="/login" class="text-blue-500 hover:text-blue-400 font-medium">
                        Back to Login
                    </a>
                </div>
            </div>
        {/if}
    </div>
</section>