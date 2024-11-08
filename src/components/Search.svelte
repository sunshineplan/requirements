<script lang="ts">
  import { slide } from "svelte/transition";
  import { fields } from "../fields";
  import { requirements } from "../requirement.svelte";

  let hover = $state(false);
  let showOption = $state(false);
  let tune: HTMLElement;
  let option: HTMLElement;

  const handleClickOutside = (event: MouseEvent) => {
    if (
      showOption &&
      !tune.contains(event.target as Node) &&
      !option.contains(event.target as Node)
    )
      showOption = false;
  };
</script>

<svelte:window onclick={handleClickOutside} />

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  class="search"
  onmouseenter={() => (hover = true)}
  onmouseleave={() => (hover = false)}
>
  <div class="icon">
    <span class="material-symbols-outlined">search</span>
  </div>
  <input
    bind:value={requirements.search.search}
    placeholder={requirements.search.field
      ? fields.name(requirements.search.field) + "搜索"
      : "搜索"}
    oninput={() => requirements.scroll()}
  />
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <div
    bind:this={tune}
    class="icon tune"
    class:show={showOption}
    style:color={requirements.search.field ? "#1a73e8" : ""}
    onclick={() => {
      showOption = !showOption;
    }}
  >
    <span class="material-symbols-outlined">tune</span>
  </div>
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <div
    class="icon reset"
    style:display={hover &&
    (requirements.search.search || requirements.search.field)
      ? "flex"
      : "none"}
  >
    <span
      class="material-symbols-outlined"
      onclick={() => requirements.clearSearch()}>close_small</span
    >
  </div>
</div>
<div class="option" bind:this={option} style:display={showOption ? "" : "none"}>
  <div class="input-group px-5 py-3" transition:slide={{ duration: 200 }}>
    <label class="input-group-text" for="option">检索字段</label>
    <select
      class="form-select"
      id="option"
      bind:value={requirements.search.field}
      onchange={() => (showOption = false)}
    >
      <option value="">所有</option>
      {#each fields.searchable() as field (field)}
        <option value={field}>{fields.name(field)}</option>
      {/each}
    </select>
  </div>
</div>

<style>
  .search {
    position: relative;
    width: 250px;
    display: flex;
    float: right;
    margin-bottom: 10px;
    background-color: #e6ecf0;
    border-radius: 9999px;
  }
  .search:hover {
    box-shadow: 0 1px 6px 0 rgba(32, 33, 36, 0.28);
  }

  .icon {
    position: absolute;
    display: flex;
    padding: 10px;
    cursor: default;
  }

  input {
    margin-left: 40px;
    margin-right: 74px;
    background-color: transparent;
    padding: 10px 0;
    border: 0;
    width: 100%;
  }
  input:focus {
    outline: none;
  }

  .tune {
    right: 30px;
    cursor: pointer;
  }
  .show,
  .tune:hover {
    background-color: #ddd;
    border-radius: 50%;
  }

  .reset {
    right: 0;
    cursor: pointer;
  }
  .reset > span:hover {
    text-shadow: 0 0 4px;
  }

  .option {
    position: absolute;
    right: 20px;
    margin-top: 10px;
  }
  .option > .input-group {
    background-color: white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    border: 1px solid rgba(0, 0, 0, 0.2);
    z-index: 1;
  }
</style>
