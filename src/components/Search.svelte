<script lang="ts">
  import { headers, searchable } from "../requirement";
  import { scroll, search, searchField } from "../stores";

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

<div class="search">
  <div class="icon">
    <span class="material-symbols-outlined">search</span>
  </div>
  <input
    bind:value={$search}
    type="search"
    placeholder={$searchField ? headers[$searchField] + "搜索" : "搜索"}
    on:input={() => scroll()}
  />
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div
    class="tune"
    bind:this={tune}
    style:display={showOption ? "none" : "flex"}
    on:click={() => {
      showOption = true;
    }}
  >
    <span class="material-symbols-outlined"> tune </span>
  </div>
</div>
<div
  class="option"
  bind:this={option}
  style:display={showOption ? "flex" : "none"}
>
  <div class="input-group px-5 py-3">
    <label class="input-group-text" for="option">检索字段</label>
    <select class="form-select" id="option" bind:value={$searchField}>
      <option value="">所有</option>
      {#each searchable as field (field)}
        <option value={field}>{headers[field]}</option>
      {/each}
    </select>
  </div>
</div>

<style>
  .icon {
    flex-direction: column;
    display: flex;
    justify-content: center;
    padding-left: 20px;
  }

  .search {
    position: relative;
    width: 250px;
    display: flex;
    float: right;
    margin-bottom: 10px;
    margin-right: 0;
    background-color: #e6ecf0;
    border-radius: 9999px;
  }
  .search:hover {
    box-shadow: 0 1px 6px 0 rgba(32, 33, 36, 0.28);
  }

  .search > input {
    background-color: transparent;
    padding: 10px;
    border: 0;
    width: 100%;
  }
  .search > input:focus {
    outline: none;
  }

  .tune {
    position: absolute;
    right: 30px;
    top: 3px;
    width: 2.4rem;
    height: 2.4rem;
    justify-content: center;
    align-items: center;
  }

  .tune:hover {
    background-color: #ddd;
    border-radius: 50%;
  }

  .option {
    background-color: white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    border: 1px solid rgba(0, 0, 0, 0.2);
    z-index: 1;
    position: absolute;
    right: 20px;
    margin-top: 10px;
  }
</style>
