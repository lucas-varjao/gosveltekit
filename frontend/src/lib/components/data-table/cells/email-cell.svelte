<script lang="ts">
    type Props = {
        email: string
    }

    let { email }: Props = $props()
</script>

<a class="text-slate-200 transition-colors hover:text-cyan-300" href={`mailto:${email}`}>
    {email}
</a>
