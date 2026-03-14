<script lang="ts">
    type Props = {
        value: string
    }

    let { value }: Props = $props()

    function formatDate(input: string) {
        if (!input || input.startsWith('0001-')) {
            return 'Nunca'
        }

        const parsed = new Date(input)
        if (Number.isNaN(parsed.getTime())) {
            return 'Nunca'
        }

        return parsed.toLocaleString()
    }
</script>

<span class="text-slate-300">{formatDate(value)}</span>
