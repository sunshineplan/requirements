<script lang="ts">
  import Swal from "sweetalert2";
  import { confirm } from "../misc.svelte";
  import { requirements } from "../requirement.svelte";

  let {
    requirement,
  }: {
    requirement: Requirement;
  } = $props();

  const done = async (r: Requirement) => {
    const today = new Date().toISOString().split("T")[0];
    const { value: date } = await Swal.fire({
      title: "选择完成日期",
      html: `
<div class="container text-start">
  <div class="row row-cols-2">
    <div class="col-12">该条业务将被标记为已完成。</div>
    <div class="col-5">提请日期:</div>
    <div class="col">${r.date}</div>
    <div class="col-5">期限日期:</div>
    <div class="col">${r.deadline}</div>
  </div>
</div>`,
      width: "20em",
      input: "date",
      inputAttributes: {
        max: today,
      },
      inputAutoFocus: false,
      didOpen: () => {
        Swal.getInput()!.value = today;
      },
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      showCancelButton: true,
      buttonsStyling: false,
      customClass: {
        input: "form-control w-auto text-center d-block",
        confirmButton: "swal btn btn-primary",
        cancelButton: "swal btn btn-primary",
      },
    });
    if (date) {
      requirements.requirement = r;
      if (requirements.component == "show") requirements.controller.abort();
      try {
        const res = await requirements.done({ ...r }, date);
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
  {#if !requirements.isClosed(requirement) && !(requirements.component == "requirement" && requirements.mode == "edit")}
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
