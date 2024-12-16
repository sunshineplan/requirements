<script lang="ts">
  let {
    id,
    label,
    value = $bindable(),
    height = "5rem",
    required,
    disabled,
    absolute,
  }: {
    id: string;
    label: string;
    value: string;
    height?: string;
    required?: boolean;
    disabled?: boolean;
    absolute?: boolean;
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
  <div class="invalid-feedback" style:position={absolute ? "absolute" : ""}>
    必填字段
  </div>
{/if}
