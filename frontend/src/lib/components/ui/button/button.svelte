<script lang="ts" module>
    import { tv, type VariantProps } from 'tailwind-variants'

    export const buttonVariants = tv({
        base: 'inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors outline-none disabled:pointer-events-none disabled:opacity-50 focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-950',
        variants: {
            variant: {
                default: 'bg-blue-600 text-white hover:bg-blue-700 focus-visible:ring-blue-500',
                destructive: 'bg-red-600 text-white hover:bg-red-700 focus-visible:ring-red-500',
                secondary:
                    'bg-slate-800 text-slate-100 hover:bg-slate-700 focus-visible:ring-slate-500',
                outline:
                    'border border-slate-700 bg-transparent text-slate-200 hover:bg-slate-800 focus-visible:ring-slate-500',
                ghost: 'text-slate-200 hover:bg-slate-800 focus-visible:ring-slate-500',
                link: 'h-auto p-0 text-blue-400 underline-offset-4 hover:text-blue-300 hover:underline focus-visible:ring-blue-500'
            },
            size: {
                default: 'h-10 px-4 py-2',
                sm: 'h-9 rounded-md px-3',
                lg: 'h-11 rounded-md px-8',
                icon: 'size-10'
            }
        },
        defaultVariants: {
            variant: 'default',
            size: 'default'
        }
    })

    export type ButtonVariant = VariantProps<typeof buttonVariants>['variant']
    export type ButtonSize = VariantProps<typeof buttonVariants>['size']
</script>

<script lang="ts">
    import type { Snippet } from 'svelte'
    import type { HTMLButtonAttributes } from 'svelte/elements'
    import { cn } from '$lib/utils'

    type Props = Omit<HTMLButtonAttributes, 'children'> & {
        variant?: ButtonVariant
        size?: ButtonSize
        children?: Snippet
    }

    let {
        class: className,
        variant = 'default',
        size = 'default',
        children,
        ...restProps
    }: Props = $props()
</script>

<button class={cn(buttonVariants({ variant, size }), className)} {...restProps}>
    {@render children?.()}
</button>
