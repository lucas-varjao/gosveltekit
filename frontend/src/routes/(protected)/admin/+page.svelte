<script lang="ts">
    import { resolve } from '$app/paths'
    import { ArrowRight, Database, TimerReset } from '@lucide/svelte'
    import PageHeader from '$lib/components/layout/page-header.svelte'
    import { buttonVariants } from '$lib/components/ui/button'
    import {
        Card,
        CardContent,
        CardDescription,
        CardHeader,
        CardTitle
    } from '$lib/components/ui/card'
    import { cn } from '$lib/utils'
</script>

<section class="page-shell">
    <PageHeader
        title="Admin Pagination"
        description="Hub oficial com duas estratégias de paginação backend-driven para usar como referência neste template."
        eyebrow="Stack Direction"
    />

    <div class="mt-8 grid gap-4 lg:grid-cols-2">
        <Card class="surface-card">
            <CardHeader class="space-y-3">
                <div class="flex items-center gap-3">
                    <Database class="size-5 text-cyan-300" />
                    <CardTitle>Offset / Limit</CardTitle>
                </div>
                <CardDescription>
                    Melhor para tabelas administrativas com total de registros, total de páginas e
                    navegação numérica.
                </CardDescription>
            </CardHeader>
            <CardContent class="space-y-4 text-sm text-slate-300">
                <p>
                    Exemplo com <code>pagination_mode=offset</code>, busca server-side, sorting
                    amplo e metadados completos para tabelas tradicionais.
                </p>
                <a
                    href={resolve('/admin/pagination/offset')}
                    class={cn(buttonVariants({ variant: 'outline' }), 'w-fit gap-2')}
                >
                    Abrir exemplo offset
                    <ArrowRight class="size-4" />
                </a>
            </CardContent>
        </Card>

        <Card class="surface-card">
            <CardHeader class="space-y-3">
                <div class="flex items-center gap-3">
                    <TimerReset class="size-5 text-blue-300" />
                    <CardTitle>Cursor-Based</CardTitle>
                </div>
                <CardDescription>
                    Melhor para volumes grandes, com cursores opacos, navegação eficiente e
                    ordenação estável.
                </CardDescription>
            </CardHeader>
            <CardContent class="space-y-4 text-sm text-slate-300">
                <p>
                    Exemplo com <code>pagination_mode=cursor</code>, navegação bidirecional por
                    cursor e sorts limitados a campos estáveis.
                </p>
                <a
                    href={resolve('/admin/pagination/cursor')}
                    class={cn(buttonVariants({ variant: 'outline' }), 'w-fit gap-2')}
                >
                    Abrir exemplo cursor
                    <ArrowRight class="size-4" />
                </a>
            </CardContent>
        </Card>
    </div>
</section>
