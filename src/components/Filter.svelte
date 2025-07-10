<script lang="ts">
  import { requirements } from "../requirement.svelte";
  import Dropdown from "./Dropdown.svelte";

  let {
    showOption = $bindable(),
  }: {
    showOption: boolean;
  } = $props();

  let field = $state("");

  const filterable = $derived(
    requirements.fields
      .filterable()
      .filter(
        (i) => !requirements.search.filter.some((f) => f.field === i.key),
      ),
  );
</script>

<div class="d-flex w-100">
  <div class="input-group">
    <label class="input-group-text" for="filter">筛选字段</label>
    <select class="form-select" id="filter" bind:value={field}>
      <option value="">请选择</option>
      {#each filterable as field (field.key)}
        <option value={field.key}>{field.name || field.key}</option>
      {/each}
    </select>
  </div>
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="add text-success"
    style:display={field ? "" : "none"}
    onclick={() => {
      requirements.search.filter.push({ field, value: [] });
      field = "";
    }}
  >
    <span class="material-symbols-outlined">add_circle</span>
  </div>
</div>
{#each requirements.search.filter as filter, index (index)}
  <div class="d-flex w-100">
    <Dropdown field={filter.field} bind:value={filter.value} bind:showOption />
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
      class="remove text-danger"
      onclick={() => {
        requirements.search.filter.splice(index, 1);
        field = "";
        showOption = false;
        requirements.scrollTop = 0;
        requirements.scroll(true);
      }}
    >
      <span class="material-symbols-outlined">do_not_disturb_on</span>
    </div>
  </div>
{/each}
