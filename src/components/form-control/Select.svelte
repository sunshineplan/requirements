<script lang="ts">
  let {
    id,
    label,
    value = $bindable(),
    options,
    required,
    disabled,
  }: {
    id: string;
    label: string;
    value: string;
    options: string[];
    required?: boolean;
    disabled?: boolean;
  } = $props();
</script>

<div class="form-floating">
  {#if disabled}
    <input class="form-control" {id} {value} disabled />
  {:else}
    <select class="form-select" {id} bind:value {required}>
      {#each options as option (option)}
        <option value={option}>{option}</option>
      {/each}
    </select>
  {/if}
  <label for={id}>{label}</label>
  {#if required}
    <div class="invalid-feedback">必填字段</div>
  {/if}
</div>
