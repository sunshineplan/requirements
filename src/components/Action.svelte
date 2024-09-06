<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { confirm } from "../misc";
  import { component, saveScrollTop, mode, goto } from "../stores";
  import {
    requirement as current,
    requirements,
    isClosed,
  } from "../requirement";

  const dispatch = createEventDispatcher();

  export let requirement: Requirement;

  const done = async (r: Requirement) => {
    if (await confirm("该条业务将被标记为已完成。")) {
      $current = r;
      try {
        const res = await requirements.done({ ...r });
        if (res === 0)
          if ($component == "requirement") goto("show");
          else {
            saveScrollTop();
            dispatch("refresh");
          }
      } catch {
        dispatch("reload");
        goto("show");
      }
    }
  };

  const edit = (r: Requirement) => {
    $current = r;
    goto("edit");
  };

  const del = async (r: Requirement) => {
    if (await confirm("该条业务将被永久删除。", true)) {
      try {
        await requirements.delete(r);
      } catch {
        dispatch("reload");
      }
      goto("show");
    }
  };
</script>

<div>
  {#if !isClosed(requirement)}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <span
      data-action="done"
      title="完成"
      class="material-symbols-outlined link-success"
      on:click={() => done(requirement)}
    >
      done_outline
    </span>
  {:else}
    <span class="material-symbols-outlined hidden">done_outline</span>
  {/if}
  {#if $component == "requirement" && $mode == "edit"}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <span
      title="删除"
      class="material-symbols-outlined link-danger"
      on:click={() => del(requirement)}
    >
      delete_outline
    </span>
  {:else}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <span
      title="编辑"
      class="material-symbols-outlined link-primary"
      on:click={() => edit(requirement)}
    >
      edit
    </span>
  {/if}
</div>

<style>
  div {
    display: flex;
  }

  span {
    font-size: var(--icon);
    margin-left: var(--margin);
  }
</style>
