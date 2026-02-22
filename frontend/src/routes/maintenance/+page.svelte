<script lang="ts">
    import StatusPage from '$lib/components/layout/status-page.svelte'
</script>

<StatusPage
    code="503"
    title="Scheduled maintenance"
    description="The system is temporarily unavailable while maintenance is in progress."
    primaryPath="/status"
    primaryLabel="Check Status"
    secondaryPath="/"
    secondaryLabel="Back Home"
    tone="warning"
/>
