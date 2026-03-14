<script lang="ts">
    type Props = {
        displayName: string
        identifier: string
    }

    let { displayName, identifier }: Props = $props()
</script>

<div class="flex flex-col">
    <span class="font-medium text-slate-100">{displayName}</span>
    <span class="text-xs text-slate-500">@{identifier}</span>
</div>
