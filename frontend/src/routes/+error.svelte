<script lang="ts">
    import { page } from '$app/state'
    import StatusPage from '$lib/components/layout/status-page.svelte'

    const status = $derived(page.status)
    const error = $derived(page.error as { message?: string } | null)
    const title = $derived(status === 404 ? 'Page not found' : 'Something went wrong')
    const description = $derived(
        status === 404
            ? "The page you tried to access doesn't exist."
            : error?.message || 'An unexpected error happened while loading this page.'
    )
</script>

<StatusPage
    code={`Error ${status}`}
    {title}
    {description}
    primaryPath="/"
    primaryLabel="Go Home"
    secondaryPath="/login"
    secondaryLabel="Sign In"
    tone={status === 404 ? 'warning' : 'danger'}
/>
