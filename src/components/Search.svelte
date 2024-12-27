<script lang="ts">
  import { slide } from "svelte/transition";
  import { requirements } from "../requirement.svelte";

  let hover = $state(false);
  let showOption = $state(false);
  let tune: HTMLElement;
  let option: HTMLElement;
  const filterField = $derived(
    requirements.fields
      .filterable()
      .find((field) => field.key == requirements.search.filter),
  );

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
      ? requirements.fields.name(requirements.search.field) + "搜索"
      : "搜索"}
    oninput={() => requirements.scroll()}
  />
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <div
    bind:this={tune}
    class="icon tune"
    class:show={showOption}
    style:color={requirements.search.field ||
    requirements.search.filter ||
    requirements.search.value
      ? "#1a73e8"
      : ""}
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
    (requirements.search.search ||
      requirements.search.field ||
      requirements.search.filter)
      ? "flex"
      : "none"}
  >
    <span
      class="material-symbols-outlined"
      onclick={() => requirements.clearSearch()}
    >
      close_small
    </span>
  </div>
</div>
<div class="option" bind:this={option} style:display={showOption ? "" : "none"}>
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <span
    class="material-symbols-outlined close"
    onclick={() => (showOption = false)}
  >
    collapse_content
  </span>
  <div class="input-group px-5 py-3" transition:slide={{ duration: 200 }}>
    <div class="d-flex w-100">
      <div class="input-group">
        <label class="input-group-text" for="search">检索字段</label>
        <select
          class="form-select"
          id="search"
          bind:value={requirements.search.field}
          onchange={() => (showOption = false)}
        >
          <option value="">所有</option>
          {#each requirements.fields.searchable() as field (field.key)}
            <option value={field.key}>{field.name || field.key}</option>
          {/each}
        </select>
      </div>
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        class="remove text-danger"
        style:display={requirements.search.field ? "" : "none"}
        onclick={() => {
          requirements.search.field = "";
          showOption = false;
        }}
      >
        <span class="material-symbols-outlined">do_not_disturb_on</span>
      </div>
    </div>
    <div class="d-flex w-100">
      <div class="input-group">
        <label class="input-group-text" for="filter">筛选字段</label>
        <select
          class="form-select"
          id="filter"
          bind:value={requirements.search.filter}
          onchange={() => (requirements.search.value = "")}
        >
          <option value="">无</option>
          {#each requirements.fields.filterable() as field (field.key)}
            <option value={field.key}>{field.name || field.key}</option>
          {/each}
        </select>
      </div>
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        class="remove text-danger"
        style:display={requirements.search.filter ? "" : "none"}
        onclick={() => {
          requirements.search.filter = "";
          requirements.search.value = "";
          showOption = false;
          requirements.scrollTop = 0;
          requirements.scroll(true);
        }}
      >
        <span class="material-symbols-outlined">do_not_disturb_on</span>
      </div>
    </div>
    <div
      class="d-flex w-100"
      style:display|important={requirements.search.filter ? "" : "none"}
    >
      <div class="input-group">
        <label class="input-group-text" for="value">筛选内容</label>
        <select
          class="form-select"
          id="value"
          bind:value={requirements.search.value}
          onchange={() => {
            showOption = false;
            requirements.scrollTop = 0;
            requirements.scroll(true);
          }}
        >
          <option value="">所有</option>
          {#if requirements.search.filter === "status"}
            {#each requirements.statuses as status (status.value)}
              <option value={status.value}>{status.value}</option>
            {/each}
          {:else if filterField}
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
        style:display={requirements.search.value ? "" : "none"}
        onclick={() => {
          requirements.search.value = "";
          showOption = false;
          requirements.scrollTop = 0;
          requirements.scroll(true);
        }}
      >
        <span class="material-symbols-outlined">do_not_disturb_on</span>
      </div>
    </div>
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
    width: 320px;
  }
  .option > .input-group {
    background-color: white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    border: 1px solid rgba(0, 0, 0, 0.2);
    border-radius: 0.375rem;
    z-index: 1;
  }

  .close {
    position: absolute;
    cursor: pointer;
    right: 0;
    z-index: 100;
  }

  .d-flex + .d-flex {
    margin-top: 1rem;
  }

  .remove {
    display: flex;
    position: absolute;
    padding: 7px;
    right: 12px;
    cursor: pointer;
  }
</style>
