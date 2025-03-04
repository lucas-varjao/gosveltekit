<script lang="ts">
	import { slide } from 'svelte/transition';

	// State declaration using Svelte 5 runes
	let username = $state('');
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let displayName = $state('');
	let submitted = $state(false);
	let errors = $state<Record<string, string>>({});
	let isLoading = $state(false);

	// Validation functions
	function validateUsername(value: string): string | null {
		if (!value) return 'Username is required';
		if (value.length < 3) return 'Username must be at least 3 characters long';
		if (value.length > 50) return 'Username must be less than 50 characters';
		if (!/^[a-zA-Z0-9._-]+$/.test(value)) return 'Username can only contain letters, numbers, dots, hyphens, and underscores';
		return null;
	}

	function validateEmail(value: string): string | null {
		if (!value) return 'Email is required';
		if (!/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(value)) return 'Email format is invalid';
		return null;
	}

	function validatePassword(value: string): string | null {
		if (!value) return 'Password is required';
		if (value.length < 8) return 'Password must be at least 8 characters long';
		if (!/(?=.*[a-z])/.test(value)) return 'Password must contain at least one lowercase letter';
		if (!/(?=.*[A-Z])/.test(value)) return 'Password must contain at least one uppercase letter';
		if (!/(?=.*\d)/.test(value)) return 'Password must contain at least one number';
		if (!/(?=.*[!@#$%^&*])/.test(value)) return 'Password must contain at least one special character';
		return null;
	}

	function validateConfirmPassword(value: string): string | null {
		if (!value) return 'Please confirm your password';
		if (value !== password) return 'Passwords do not match';
		return null;
	}

	function validateDisplayName(value: string): string | null {
		if (!value) return 'Display name is required';
		if (value.length < 2) return 'Display name must be at least 2 characters long';
		if (value.length > 50) return 'Display name must be less than 50 characters';
		return null;
	}

	// Reactive validation
	$effect(() => {
		if (submitted) {
			errors.username = validateUsername(username) || '';
			errors.email = validateEmail(email) || '';
			errors.password = validatePassword(password) || '';
			errors.confirmPassword = validateConfirmPassword(confirmPassword) || '';
			errors.displayName = validateDisplayName(displayName) || '';
		}
	});

	// Form submission handler
	async function handleSubmit(event: Event) {
        event.preventDefault();
		submitted = true;
		
		// Perform validation
		const usernameError = validateUsername(username);
		const emailError = validateEmail(email);
		const passwordError = validatePassword(password);
		const confirmPasswordError = validateConfirmPassword(confirmPassword);
		const displayNameError = validateDisplayName(displayName);
		
		// Update errors state
		errors = {
			username: usernameError || '',
			email: emailError || '',
			password: passwordError || '',
			confirmPassword: confirmPasswordError || '',
			displayName: displayNameError || ''
		};
		
		// Check if there are any validation errors
		const hasErrors = Object.values(errors).some(error => error !== '');
		
		if (!hasErrors) {
			try {
				isLoading = true;
				
				// Simulate API call for now (will be connected to the backend later)
				await new Promise(resolve => setTimeout(resolve, 1500));
				
				// Reset form after successful submission
				username = '';
				email = '';
				password = '';
				confirmPassword = '';
				displayName = '';
				submitted = false;
				
				// Show success message or redirect
				alert('Registration successful!');
			} catch (error) {
				console.error('Registration error:', error);
			} finally {
				isLoading = false;
			}
		}
	}
</script>

<section class="py-12">
	<div class="max-w-lg mx-auto bg-slate-900 rounded border border-slate-800 shadow-lg p-8">
		<h1 class="text-3xl font-bold mb-6 text-center text-white">Create an account</h1>

		<!-- Login Link -->
			<div class="text-center m-4 text-slate-400">
				Already have an account?
				<a href="/login" class="text-blue-500 hover:text-blue-400 font-medium">
					Login
				</a>
			</div>
		
		<form onsubmit={handleSubmit} class="space-y-4">
			<!-- Username Field -->
			<div class="space-y-2">
				<label for="username" class="block text-sm font-medium text-slate-200">
					Username
				</label>
				<input
					type="text"
					id="username"
					bind:value={username}
					placeholder="Enter username"
					class="w-full px-3 py-2 bg-slate-800 text-white border-2 rounded {errors.username ? 'border-red-500' : 'border-slate-700'} focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
				{#if errors.username && submitted}
					<p transition:slide class="text-sm text-red-500 mt-1">{errors.username}</p>
				{/if}
			</div>
			
			<!-- Email Field -->
			<div class="space-y-2">
				<label for="email" class="block text-sm font-medium text-slate-200">
					Email
				</label>
				<input
					type="email"
					id="email"
					bind:value={email}
					placeholder="Enter email address"
					class="w-full px-3 py-2 bg-slate-800 text-white border-2 rounded {errors.email ? 'border-red-500' : 'border-slate-700'} focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
				{#if errors.email && submitted}
					<p transition:slide class="text-sm text-red-500 mt-1">{errors.email}</p>
				{/if}
			</div>

			<!-- Display Name Field -->
			<div class="space-y-2">
				<label for="displayName" class="block text-sm font-medium text-slate-200">
					Display Name
				</label>
				<input
					type="text"
					id="displayName"
					bind:value={displayName}
					placeholder="Enter your full name"
					class="w-full px-3 py-2 bg-slate-800 text-white border-2 rounded {errors.displayName ? 'border-red-500' : 'border-slate-700'} focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
				{#if errors.displayName && submitted}
					<p transition:slide class="text-sm text-red-500 mt-1">{errors.displayName}</p>
				{/if}
			</div>
			
			<!-- Password Field -->
			<div class="space-y-2">
				<label for="password" class="block text-sm font-medium text-slate-200">
					Password
				</label>
				<input
					type="password"
					id="password"
					bind:value={password}
					placeholder="Enter password"
					class="w-full px-3 py-2 bg-slate-800 text-white border-2 rounded {errors.password ? 'border-red-500' : 'border-slate-700'} focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
				{#if errors.password && submitted}
					<p transition:slide class="text-sm text-red-500 mt-1">{errors.password}</p>
				{/if}
			</div>
			
			<!-- Confirm Password Field -->
			<div class="space-y-2">
				<label for="confirmPassword" class="block text-sm font-medium text-slate-200">
					Confirm Password
				</label>
				<input
					type="password"
					id="confirmPassword"
					bind:value={confirmPassword}
					placeholder="Confirm your password"
					class="w-full px-3 py-2 bg-slate-800 text-white border-2 rounded {errors.confirmPassword ? 'border-red-500' : 'border-slate-700'} focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
				{#if errors.confirmPassword && submitted}
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
						<span class="inline-block animate-pulse">Creating Account...</span>
					{:else}
						Create account
					{/if}
				</button>
			</div>			
			
		</form>
	</div>
</section>