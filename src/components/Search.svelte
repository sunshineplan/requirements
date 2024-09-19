<script lang="ts">
  import { slide } from "svelte/transition";
  import { fields } from "../requirement";
  import { clear, scroll, search, searchField } from "../stores";

  let hover = false;
  let showOption = false;
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

<svelte:window on:click={handleClickOutside} />

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
  class="search"
  on:mouseenter={() => (hover = true)}
  on:mouseleave={() => (hover = false)}
>
  <div class="icon">
    <span class="material-symbols-outlined">search</span>
  </div>
  <input
    bind:value={$search}
    placeholder={$searchField ? fields.name($searchField) + "搜索" : "搜索"}
    on:input={() => scroll()}
  />
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div
    bind:this={tune}
    class="icon tune"
    class:show={showOption}
    style:color={$searchField ? "#1a73e8" : ""}
    on:click={() => {
      showOption = !showOption;
    }}
  >
    <span class="material-symbols-outlined">tune</span>
  </div>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div
    class="icon reset"
    style:display={hover && ($search || $searchField) ? "flex" : "none"}
  >
    <span class="material-symbols-outlined" on:click={clear}>close_small</span>
  </div>
</div>
{#if showOption}
  <div class="option" bind:this={option}>
    <div class="input-group px-5 py-3" transition:slide={{ duration: 50 }}>
      <label class="input-group-text" for="option">检索字段</label>
      <select class="form-select" id="option" bind:value={$searchField}>
        <option value="">所有</option>
        {#each fields.searchable() as field (field)}
          <option value={field}>{fields.name(field)}</option>
        {/each}
      </select>
    </div>
  </div>
{/if}

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
    user-select: none;
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
