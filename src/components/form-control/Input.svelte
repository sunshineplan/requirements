<script lang="ts">
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

  const showList = $derived(!disabled && list && list.length);
</script>

<div class="form-floating">
  <input
    class="form-control"
    {id}
    list={showList ? id + "-list" : ""}
    bind:value
    {required}
    {disabled}
  />
  {#if showList}
    <datalist id={id + "-list"}>
      {#each list! as option (option)}
        <option>{option}</option>
      {/each}
    </datalist>
  {/if}
  <label for={id}>{label}</label>
  {#if required}
    <div class="invalid-feedback">必填字段</div>
  {/if}
</div>
