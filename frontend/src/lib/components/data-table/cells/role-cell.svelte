<script lang="ts">
    import { cn } from '$lib/utils'

    type Props = {
        role: string
    }

    let { role }: Props = $props()
</script>

<span
    class={cn(
        'inline-flex w-fit rounded-full border px-2.5 py-1 text-xs font-semibold tracking-[0.08em] uppercase',
        role === 'admin'
            ? 'border-cyan-400/30 bg-cyan-500/10 text-cyan-200'
            : 'border-slate-700 bg-slate-800/70 text-slate-300'
    )}
>
    {role}
</span>
