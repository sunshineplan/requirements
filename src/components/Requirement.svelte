<script lang="ts">
  import { confirm, valid } from "../misc.svelte";
  import { requirements } from "../requirement.svelte";
  import Action from "./Action.svelte";
  import Checkbox from "./form-control/Checkbox.svelte";
  import Date from "./form-control/Date.svelte";
  import Input from "./form-control/Input.svelte";
  import Select from "./form-control/Select.svelte";
  import Textarea from "./form-control/Textarea.svelte";

  const modeList: { [key: string]: string } = {
    add: "新增",
    edit: "编辑",
    view: "查看",
  };

  let type = $state(requirements.requirement.type || "");
  let title = $state(requirements.requirement.title || "");
  let date = $state(requirements.requirement.date || "");
  let deadline = $state(requirements.requirement.deadline || "");
  let done = $state(requirements.requirement.done || "");
  let submitter = $state(requirements.requirement.submitter || "");
  let recipient = $state(requirements.requirement.recipient || "");
  let acceptor = $state(requirements.requirement.acceptor || "");
  let status = $state(requirements.requirement.status || "");
  let note = $state(requirements.requirement.note || "");
  let group = $state(
    requirements.requirement.group
      ? requirements.requirement.group.split(",")
      : [],
  );
  let validated = $state(false);

  let doneValue = $state("");
  let groups: string[] = $state([]);

  let submitters: string[] = $state([]);
  let recipients: string[] = $state([]);
  let acceptors: string[] = $state([]);

  const init = async () => {
    const res = await requirements.init();
    groups = res.groups;
    doneValue = res.done;
    submitters = await requirements.submitters();
    recipients = await requirements.recipients();
    acceptors = await requirements.acceptors();
  };
  const promise = init();

  const current = () => {
    return {
      type,
      title,
      date,
      deadline,
      done,
      submitter,
      recipient,
      acceptor,
      status,
      note,
      group: group.join(","),
    } as Requirement;
  };

  const save = async () => {
    if (valid() && (!groups.length || group.length > 0)) {
      validated = false;
      const r = current();
      if (requirements.mode == "edit") r.id = requirements.requirement.id;
      try {
        if (r.status != doneValue) r.done = "";
        const res = await requirements.save(r);
        if (res === 0) {
          if (requirements.mode == "add") requirements.clearSearch();
          requirements.goto("show");
        }
      } catch {
        await requirements.init();
        requirements.goto("show");
      }
    } else validated = true;
  };

  const back = async () => {
    const r = current();
    let edited = false;
    switch (requirements.mode) {
      case "view":
        break;
      case "add":
        for (const k in r) {
          if (r[k as keyof Requirement] != "") {
            edited = true;
            break;
          }
        }
        break;
      case "edit":
        for (const k in r) {
          const key = k as keyof Requirement;
          if (r[key] != requirements.requirement[key]) {
            edited = true;
            break;
          }
        }
    }
    if (edited && !(await confirm("数据未保存，确定将放弃保存并返回。", true)))
      return;
    requirements.goto("show");
  };
</script>

<svelte:head>
  <title>
    {modeList[requirements.mode]} - {requirements.brand || "业务系统"}
  </title>
</svelte:head>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div style="height: 100%;">
  <header>
    <div class="back">
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <span class="material-symbols-outlined" onclick={back}>arrow_back</span>
    </div>
    <h3>{modeList[requirements.mode]}</h3>
    {#if requirements.mode != "add"}
      <Action
        requirement={requirements.requirement}
        --icon="22px"
        --margin="10px"
      />
    {/if}
  </header>
  {#await promise then _}
    <div class="row g-3" class:was-validated={validated}>
      <div class="col-md-8 col-sm-12">
        <Textarea
          id="title"
          bind:value={title}
          required={true}
          disabled={requirements.mode == "view"}
          label={requirements.fields.name("title")}
        />
      </div>
      <div class="w-100 m-0"></div>
      <div class="col-md-3 col-sm-4">
        <Select
          id="type"
          bind:value={type}
          options={requirements.types}
          required={true}
          disabled={requirements.mode == "view"}
          label={requirements.fields.name("type")}
        />
      </div>
      <div class="col-md-3 col-sm-4">
        <Select
          id="status"
          bind:value={status}
          options={requirements.statuses.map((status) => status.value)}
          required={true}
          disabled={requirements.mode == "view"}
          label={requirements.fields.name("status")}
        />
      </div>
      <div class="w-100 m-0"></div>
      <div class="col-md-3 col-sm-4">
        <Date
          id="date"
          bind:value={date}
          required={true}
          disabled={requirements.mode == "view"}
          label={requirements.fields.name("date")}
        />
      </div>
      <div class="col-md-3 col-sm-4">
        <Date
          id="deadline"
          bind:value={deadline}
          min={date}
          disabled={requirements.mode == "view"}
          label={requirements.fields.name("deadline")}
        />
      </div>
      {#if status === doneValue}
        <div class="col-md-3 col-sm-4">
          <Date
            id="done"
            bind:value={done}
            min={date}
            required={true}
            disabled={requirements.mode == "view"}
            label={requirements.fields.name("done")}
          />
        </div>
      {/if}
      <div class="w-100 m-0"></div>
      <div class="col-md-3 col-sm-4">
        <Input
          id="submitter"
          bind:value={submitter}
          required={true}
          disabled={requirements.mode == "view"}
          label={requirements.fields.name("submitter")}
          list={submitters}
        />
      </div>
      <div class="col-md-3 col-sm-4">
        <Input
          id="recipient"
          bind:value={recipient}
          disabled={requirements.mode == "view"}
          label={requirements.fields.name("recipient")}
          list={recipients}
        />
      </div>
      <div class="col-md-3 col-sm-4">
        <Input
          id="acceptor"
          bind:value={acceptor}
          required={true}
          disabled={requirements.mode == "view"}
          label={requirements.fields.name("acceptor")}
          list={acceptors}
        />
      </div>
      <div class="col-md-6">
        <Checkbox
          id="group"
          bind:value={group}
          required={true}
          disabled={requirements.mode == "view"}
          label={requirements.fields.name("group")}
          options={groups}
          {validated}
        />
      </div>
      <div class="col-md-8 col-sm-12">
        <Textarea
          id="note"
          height="7rem"
          bind:value={note}
          required={true}
          disabled={requirements.mode == "view"}
          label={requirements.fields.name("note")}
        />
      </div>
      <div class="col-md-8 col-sm-12">
        {#if requirements.mode == "view"}
          <button class="btn btn-primary float-end mb-2" onclick={back}>
            返回
          </button>
        {:else}
          <button class="btn btn-primary float-end mb-2" onclick={save}>
            保存
          </button>
          <button class="btn btn-primary float-end mx-2 mb-2" onclick={back}>
            取消
          </button>
        {/if}
      </div>
    </div>
  {/await}
</div>

<style>
  header {
    display: flex;
    align-items: center;
  }

  header h3 {
    margin: 0;
    cursor: default;
  }

  .row {
    padding: 0 20px;
    overflow: auto;
    margin-top: 0;
    max-height: calc(100% - 60px);
  }
</style>
