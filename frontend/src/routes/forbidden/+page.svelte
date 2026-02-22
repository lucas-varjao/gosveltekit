<script lang="ts">
    import StatusPage from '$lib/components/layout/status-page.svelte'
</script>

<StatusPage
    code="403"
    title="Access denied"
    description="You are authenticated, but you do not have permission to access this resource."
    primaryPath="/profile"
    primaryLabel="Go to Profile"
    secondaryPath="/"
    secondaryLabel="Back Home"
    tone="danger"
/>
