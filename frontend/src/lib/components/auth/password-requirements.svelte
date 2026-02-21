<script lang="ts">
    import { Circle, CircleCheckBig } from '@lucide/svelte'

    type PasswordRequirements = {
        length: boolean
        lowercase: boolean
        uppercase: boolean
        number: boolean
        special: boolean
    }

    let { requirements }: { requirements: PasswordRequirements } = $props()

    const requirementItems: Array<{
        key: keyof PasswordRequirements
        label: string
    }> = [
        { key: 'length', label: 'At least 8 characters' },
        { key: 'lowercase', label: 'One lowercase letter' },
        { key: 'uppercase', label: 'One uppercase letter' },
        { key: 'number', label: 'One number' },
        { key: 'special', label: 'One special character (!@#$%^&*)' }
    ]
</script>

<div class="flex flex-col gap-2 text-sm">
    <p class="text-slate-300">Your password must contain:</p>
    <ul class="ml-2 flex flex-col gap-1">
        {#each requirementItems as item (item.key)}
            {@const met = requirements[item.key]}
            <li class="flex items-center gap-2">
                {#if met}
                    <CircleCheckBig class="size-4 text-emerald-400" />
                {:else}
                    <Circle class="size-4 text-slate-500" />
                {/if}
                <span class={met ? 'text-slate-200' : 'text-slate-400'}>{item.label}</span>
            </li>
        {/each}
    </ul>
</div>
