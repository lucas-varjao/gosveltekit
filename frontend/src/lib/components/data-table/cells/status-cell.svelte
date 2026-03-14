<script lang="ts">
    import { cn } from '$lib/utils'

    type Props = {
        active: boolean
    }

    let { active }: Props = $props()
</script>

<span
    class={cn(
        'inline-flex w-fit rounded-full border px-2.5 py-1 text-xs font-semibold tracking-[0.08em] uppercase',
        active
            ? 'border-emerald-400/30 bg-emerald-500/10 text-emerald-200'
            : 'border-amber-400/30 bg-amber-500/10 text-amber-200'
    )}
>
    {active ? 'Ativo' : 'Inativo'}
</span>
