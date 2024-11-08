<script lang="ts">
  import { confirm } from "../misc.svelte";
  import { requirements } from "../requirement.svelte";

  let {
    requirement,
  }: {
    requirement: Requirement;
  } = $props();

  const done = async (r: Requirement) => {
    if (await confirm("该条业务将被标记为已完成。")) {
      requirements.requirement = r;
      if (requirements.component == "show") requirements.controller.abort();
      try {
        const res = await requirements.done({ ...r });
        if (res === 0)
          if (requirements.component == "requirement") {
            requirements.goto("show");
            return;
          } else {
            requirements.saveScrollTop();
            await requirements.init();
            requirements.scroll(true);
          }
        requirements.subscribe(true);
      } catch {
        await requirements.init();
        requirements.goto("show");
      }
    }
  };

  const edit = (r: Requirement) => {
    requirements.requirement = r;
    requirements.goto("edit");
  };

  const del = async (r: Requirement) => {
    if (await confirm("该条业务将被永久删除。", true)) {
      try {
        await requirements.delete(r);
      } catch {
        await requirements.init();
      }
      requirements.goto("show");
    }
  };
</script>

<div>
  {#if !requirements.isClosed(requirement)}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <span
      data-action="done"
      title="完成"
      class="material-symbols-outlined link-success"
      onclick={() => done(requirement)}
    >
      done_outline
    </span>
  {:else}
    <span class="material-symbols-outlined hidden">done_outline</span>
  {/if}
  {#if requirements.component == "requirement" && requirements.mode == "edit"}
    {#if requirements.username == "admin"}
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <span
        title="删除"
        class="material-symbols-outlined link-danger"
        onclick={() => del(requirement)}
      >
        delete_outline
      </span>
    {/if}
  {:else}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <span
      data-action="edit"
      title="编辑"
      class="material-symbols-outlined link-primary"
      onclick={() => edit(requirement)}
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
    cursor: pointer;
  }
</style>
