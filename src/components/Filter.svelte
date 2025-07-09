<script lang="ts">
  import { requirements } from "../requirement.svelte";

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

  const getField = (key: string) => {
    return requirements.fields.filterable().find((field) => field.key === key);
  };
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
      requirements.search.filter.push({ field, value: "" });
      field = "";
    }}
  >
    <span class="material-symbols-outlined">add_circle</span>
  </div>
</div>
{#each requirements.search.filter as filter, index (index)}
  <div class="d-flex w-100">
    <div class="input-group">
      <label class="input-group-text" for={`value${index}`}>
        {requirements.fields.name(filter.field)}
      </label>
      <select
        class="form-select"
        id={`value${index}`}
        bind:value={filter.value}
        onchange={() => {
          showOption = false;
          requirements.scrollTop = 0;
          requirements.scroll(true);
        }}
      >
        <option value="">所有</option>
        {#if filter.field === "status"}
          {#each requirements.statuses as status (status.value)}
            <option value={status.value}>{status.value}</option>
          {/each}
        {:else if getField(filter.field)}
          {@const filterField = getField(filter.field)!}
          {#if !filterField.enum || filterField.enum.length === 0}
            <option value="">无选项</option>
          {/if}
          {#each filterField.enum as value (value)}
            <option {value}>{value}</option>
          {/each}
        {/if}
      </select>
    </div>
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
