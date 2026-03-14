<script lang="ts">
    import MockCursorPaginationDemo from '$lib/components/pagination-playground/mock-cursor-pagination-demo.svelte'
    import MockOffsetPaginationDemo from '$lib/components/pagination-playground/mock-offset-pagination-demo.svelte'
    import PageHeader from '$lib/components/layout/page-header.svelte'
    import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card'
</script>

<section class="page-shell">
    <PageHeader
        title="Pagination Playground"
        description="Exemplo independente de admin, com dados mockados no backend e as duas estratégias de paginação na mesma tela."
        eyebrow="Examples"
    />

    <Card class="surface-card mt-8">
        <CardHeader>
            <CardTitle>Backend mockado</CardTitle>
        </CardHeader>
        <CardContent class="space-y-3 text-sm text-slate-300">
            <p>
                Esta página usa o endpoint <code>/api/examples/pagination/items</code>, que gera
                dados mockados no backend e aplica filtro, ordenação e paginação sem depender do
                banco.
            </p>
            <p class="text-slate-400">
                O bloco offset expõe totais e navegação numerada. O bloco cursor usa cursores opacos
                e navegação bidirecional.
            </p>
        </CardContent>
    </Card>

    <div class="mt-8 space-y-10">
        <MockOffsetPaginationDemo />
        <MockCursorPaginationDemo />
    </div>
</section>
