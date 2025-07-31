<script lang="ts">
  import type { Snippet } from "svelte";

  let {
    type,
    id,
    label,
    value = $bindable(),
    required,
    disabled,
    props,
    html,
  }: {
    type: string;
    id: string;
    label: string;
    value: string;
    required?: boolean;
    disabled?: boolean;
    props?: Record<string, any>;
    html?: Snippet<[]>;
  } = $props();
</script>

<div class="form-floating">
  <input
    class="form-control"
    {type}
    {id}
    bind:value
    {...props}
    {required}
    {disabled}
  />
  {@render html?.()}
  <label for={id}>{label}</label>
  {#if required}
    <div class="invalid-feedback">必填字段</div>
  {/if}
</div>
