<script lang="ts">
    import StatusPage from '$lib/components/layout/status-page.svelte'
</script>

<StatusPage
    code="401"
    title="Session expired"
    description="Your session is no longer valid. Please sign in again to continue."
    primaryPath="/login"
    primaryLabel="Sign In Again"
    secondaryPath="/"
    secondaryLabel="Back Home"
    tone="warning"
/>
