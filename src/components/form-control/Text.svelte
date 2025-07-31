<script lang="ts">
  import Input from "./Input.svelte";
  let {
    id,
    label,
    value = $bindable(),
    list,
    required,
    disabled,
  }: {
    id: string;
    label: string;
    value: string;
    list?: string[];
    required?: boolean;
    disabled?: boolean;
  } = $props();

  const showList = $derived(!disabled && list?.length);
</script>

<Input
  type="text"
  {id}
  {label}
  bind:value
  {required}
  {disabled}
  props={{ list: showList ? id + "-list" : "" }}
  {html}
/>

{#snippet html()}
  {#if showList}
    <datalist id={id + "-list"}>
      {#each list! as option (option)}
        <option>{option}</option>
      {/each}
    </datalist>
  {/if}
{/snippet}
