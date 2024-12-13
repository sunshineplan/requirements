<script lang="ts">
  let {
    id,
    label,
    value = $bindable(),
    options,
    required,
    validated,
    disabled,
  }: {
    id: string;
    label: string;
    options: string[];
    value: string[];
    required?: boolean;
    validated?: boolean;
    disabled?: boolean;
  } = $props();

  const invalid = $derived(
    required && validated && options.length && value.length == 0,
  );
</script>

<label class="form-label" for={id}>
  {label}
</label>
<div class="checkbox" {id}>
  {#if options.length}
    {#each options as v, index (v)}
      <div class="form-check form-check-inline">
        <input
          type="checkbox"
          class="form-check-input"
          class:invalid
          id={id + index}
          bind:group={value}
          value={v}
          {disabled}
        />
        <label class="form-check-label" class:invalid for={id + index}>
          {v}
        </label>
      </div>
    {/each}
  {:else}
    <div class="form-check form-check-inline">
      <input type="checkbox" class="form-check-input" id={"no" + id} disabled />
      <label class="form-check-label" for={"no" + id}>无</label>
    </div>
  {/if}
</div>
{#if required}
  <div class="invalid-feedback" class:invalid>必选字段</div>
{/if}

<style>
  .checkbox {
    display: flex;
  }

  .invalid {
    display: block;
    color: var(--bs-form-invalid-color) !important;
    border-color: var(--bs-form-invalid-border-color) !important;
  }
</style>
