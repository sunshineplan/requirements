<script lang="ts">
  let {
    id,
    label,
    value = $bindable(),
    max,
    min,
    required,
    disabled,
  }: {
    id: string;
    label: string;
    value: string;
    max?: string;
    min?: string;
    required?: boolean;
    disabled?: boolean;
  } = $props();
</script>

<div class="form-floating">
  <input
    class="form-control"
    {id}
    type="date"
    bind:value
    {max}
    {min}
    {required}
    {disabled}
  />
  <label for={id}>{label}</label>
  {#if required}
    <div class="invalid-feedback">必填字段</div>
  {/if}
</div>
