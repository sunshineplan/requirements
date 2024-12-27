<script lang="ts">
  import Swal from "sweetalert2";
  import { confirm } from "../misc.svelte";
  import { requirements } from "../requirement.svelte";

  let {
    requirement,
  }: {
    requirement: ExtendedRequirement;
  } = $props();
  const control = $derived.by(() => {
    if (requirements.doneValue.length == 1)
      return `<input class="form-control" id="done-value" disabled />`;
    const select = document.createElement("select");
    select.id = "done-value";
    select.className = "form-select";
    requirements.doneValue.forEach((done) => {
      const option = document.createElement("option");
      option.value = done;
      option.text = done;
      select.appendChild(option);
    });
    return select.outerHTML;
  });

  const done = async (r: ExtendedRequirement) => {
    const today = new Date().toISOString().split("T")[0];
    const { value } = await Swal.fire({
      title: `选择${requirements.fields.name("done")}`,
      html: `
<div class="container text-start">
  <div class="row row-cols-2">
    <div class="col-12">该条记录将被标记为:</div>
    <div class="col-12 pt-1">
      <div class="form-floating">
        ${control}
        <label for="done-value">${requirements.fields.name("status")}</label>
      </div>
    </div>
    <div class="col-5 pt-3">${requirements.fields.name("date")}:</div>
    <div class="col pt-3">${r.date}</div>
    <div class="col-5">${requirements.fields.name("deadline")}:</div>
    <div class="col">${r.deadline}</div>
    <div class="col-12 pt-1">
      <div class="form-floating">
        <input type="date" class="form-control" id="done-date">
        <label for="done-date">${requirements.fields.name("done")}</label>
      </div>
    </div>
  </div>
</div>`,
      width: "20em",
      didOpen: () => {
        document.querySelector<HTMLInputElement>("#done-value")!.value =
          requirements.doneValue[0];
        const dateInput =
          document.querySelector<HTMLInputElement>("#done-date")!;
        dateInput.value = today;
        dateInput.max = today;
      },
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      showCancelButton: true,
      buttonsStyling: false,
      customClass: {
        confirmButton: "swal btn btn-primary",
        cancelButton: "swal btn btn-primary",
      },
      preConfirm: () => {
        // need check vaild
        return {
          date: document.querySelector<HTMLInputElement>("#done-date")!.value,
          value: document.querySelector<HTMLInputElement>("#done-value")!.value,
        };
      },
    });
    if (value) {
      requirements.requirement = r;
      if (requirements.component == "show") requirements.controller.abort();
      try {
        const res = await requirements.done({ ...r }, value.date, value.status);
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

  const edit = (r: ExtendedRequirement) => {
    requirements.requirement = r;
    requirements.goto("edit");
  };

  const del = async (r: ExtendedRequirement) => {
    if (await confirm("该条记录将被永久删除。", true)) {
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
  {#if !requirements.isClosed(requirement.status) && !(requirements.component == "requirement" && requirements.mode == "edit") && requirements.doneValue.length}
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
