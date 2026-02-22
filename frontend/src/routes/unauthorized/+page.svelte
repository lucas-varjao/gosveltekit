<script lang="ts">
    import StatusPage from '$lib/components/layout/status-page.svelte'
</script>

<StatusPage
    code="401"
    title="Unauthorized"
    description="You need to sign in before accessing this page."
    primaryPath="/login"
    primaryLabel="Sign In"
    secondaryPath="/"
    secondaryLabel="Back Home"
    tone="warning"
/>
