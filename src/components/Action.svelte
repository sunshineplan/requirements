<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { confirm } from "../misc";
  import { mode, component, goHome } from "../stores";
  import { requirement as current, requirements } from "../requirement";

  const dispatch = createEventDispatcher();

  export let requirement: Requirement;

  const done = async (r: Requirement) => {
    if (await confirm("这条业务将被标记为已完成。")) {
      $current = r;
      try {
        const res = await requirements.done({ ...r });
        if (res === 0)
          if ($component == "requirement") goHome();
          else dispatch("refresh");
      } catch {
        dispatch("reload");
        goHome();
      }
    }
  };

  const edit = (r: Requirement) => {
    $current = r;
    window.history.pushState({}, "", "/edit");
    $mode = "edit";
    $component = "requirement";
  };
</script>

<div>
  {#if requirement.status != "已完成"}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <span
      data-action="done"
      title="完成"
      class="material-symbols-outlined done"
      on:click={() => done(requirement)}
    >
      done_outline
    </span>
  {:else}
    <span class="material-symbols-outlined hidden">done_outline</span>
  {/if}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <span
    title="编辑"
    class="material-symbols-outlined edit"
    on:click={() => edit(requirement)}
  >
    edit
  </span>
</div>

<style>
  div {
    display: flex;
  }

  span {
    font-size: var(--icon);
    cursor: pointer;
    margin-left: var(--margin);
  }

  .hidden {
    visibility: hidden;
  }

  .done {
    color: #198754 !important;
  }

  .done:hover {
    color: #157347 !important;
  }

  .edit {
    color: #007bff !important;
  }

  .edit:hover {
    color: #0056b3 !important;
  }
</style>
