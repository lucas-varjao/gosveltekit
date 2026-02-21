<script lang="ts">
    import { Eye, EyeOff } from '@lucide/svelte'
    import { buttonVariants } from '$lib/components/ui/button'
    import { Input } from '$lib/components/ui/input'
    import { Label } from '$lib/components/ui/label'
    import { cn } from '$lib/utils'
    import type { HTMLInputAttributes } from 'svelte/elements'

    type Props = {
        id: string
        label: string
        value?: string
        placeholder?: string
        autocomplete?: HTMLInputAttributes['autocomplete']
        error?: string
        touched?: boolean
        onblur?: () => void
    }

    let {
        id,
        label,
        value = $bindable(''),
        placeholder = '',
        autocomplete = 'current-password',
        error = '',
        touched = false,
        onblur
    }: Props = $props()

    let showPassword = $state(false)
    let hasError = $derived(Boolean(touched && error))

    function togglePasswordVisibility() {
        showPassword = !showPassword
    }
</script>

<div class="flex flex-col gap-2">
    <Label for={id}>{label}</Label>
    <div class="relative">
        <Input
            {id}
            type={showPassword ? 'text' : 'password'}
            bind:value
            {placeholder}
            {autocomplete}
            aria-invalid={hasError}
            class={cn('pr-10', hasError && 'border-red-500 focus-visible:ring-red-500/30')}
            {onblur}
        />

        <button
            type="button"
            class={cn(
                buttonVariants({ variant: 'ghost', size: 'icon' }),
                'absolute top-0 right-0 h-9 w-9 rounded-l-none text-slate-400 hover:text-slate-200'
            )}
            onclick={togglePasswordVisibility}
            aria-label={showPassword ? 'Hide password' : 'Show password'}
        >
            {#if showPassword}
                <EyeOff class="size-4" />
            {:else}
                <Eye class="size-4" />
            {/if}
        </button>
    </div>

    {#if hasError}
        <p class="text-sm text-red-500">{error}</p>
    {/if}
</div>
