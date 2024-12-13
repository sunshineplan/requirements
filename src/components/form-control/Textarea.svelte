<script lang="ts">
  let {
    id,
    label,
    value = $bindable(),
    height = "5rem",
    required,
    disabled,
  }: {
    id: string;
    label: string;
    value: string;
    height?: string;
    required?: boolean;
    disabled?: boolean;
  } = $props();
</script>

<label for={id} class="form-label">
  {label}
</label>
<textarea
  class="form-control"
  style:height
  {id}
  bind:value
  {required}
  {disabled}
></textarea>
{#if required}
  <div class="invalid-feedback">必填字段</div>
{/if}
