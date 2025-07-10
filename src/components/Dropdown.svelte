<script lang="ts">
  import { requirements } from "../requirement.svelte";

  let {
    field,
    value = $bindable(),
    showOption = $bindable(),
  }: {
    field: string;
    value: string[];
    showOption: boolean;
  } = $props();

  let show = $state(false);
  let all = $state(true);
  let dropdown: HTMLElement;

  const options = requirements.fields
    .filterable()
    .find((f) => f.key === field)!.enum;

  const handleClickOutside = (event: MouseEvent) => {
    if (show && !dropdown.contains(event.target as Node)) show = false;
  };

  $effect(() => {
    all = value.length === 0 || value.length === options.length;
    if (value.length === options.length) value = [];
    show = false;
  });
</script>

<svelte:window onclick={handleClickOutside} />

<div bind:this={dropdown} class="dropdown w-100">
  <button
    class="btn btn-outline-primary dropdown-toggle w-100"
    type="button"
    onclick={() => (show = !show)}
  >
    {requirements.fields.name(field)}
  </button>
  <ul class="dropdown-menu" class:show>
    <li>
      <div class="form-check form-check-reverse">
        <input
          id={field}
          type="checkbox"
          class="form-check-input"
          checked={all}
          disabled={all}
          onclick={() => {
            value = [];
            showOption = false;
          }}
        />
        <label class="form-check-label" for={field}>所有</label>
      </div>
    </li>
    {#if field === "status"}
      {#each requirements.statuses as status, index (status.value)}
        <li>
          <div class="form-check form-check-reverse">
            <input
              id={`${field}${index}`}
              type="checkbox"
              class="form-check-input"
              value={status.value}
              bind:group={value}
              onclick={() => {
                showOption = false;
              }}
            />
            <label class="form-check-label" for={`${field}${index}`}>
              {status.value}
            </label>
          </div>
        </li>
      {/each}
    {:else}
      {#each options as option, index (index)}
        <li>
          <div class="form-check form-check-reverse">
            <input
              id={`${field}${index}`}
              type="checkbox"
              class="form-check-input"
              value={option}
              bind:group={value}
              onclick={() => {
                showOption = false;
              }}
            />
            <label class="form-check-label" for={`${field}${index}`}>
              {option}
            </label>
          </div>
        </li>
      {/each}
    {/if}
  </ul>
</div>

<style>
  .dropdown-menu {
    --bs-dropdown-padding-x: 0.5rem;
  }
</style>
