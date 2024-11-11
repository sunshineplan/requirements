<script lang="ts">
  import { onMount } from "svelte";
  import { confirm, valid } from "../misc.svelte";
  import { requirements } from "../requirement.svelte";
  import Action from "./Action.svelte";

  const modeList: { [key: string]: string } = {
    add: "新增",
    edit: "编辑",
    view: "查看",
  };

  let type = $state(requirements.requirement.type || "");
  let desc = $state(requirements.requirement.desc || "");
  let date = $state(requirements.requirement.date || "");
  let deadline = $state(requirements.requirement.deadline || "");
  let done = $state(requirements.requirement.done || "");
  let submitter = $state(requirements.requirement.submitter || "");
  let recipient = $state(requirements.requirement.recipient || "");
  let acceptor = $state(requirements.requirement.acceptor || "");
  let status = $state(requirements.requirement.status || "");
  let note = $state(requirements.requirement.note || "");
  let participating = $state(
    requirements.requirement.participating
      ? requirements.requirement.participating.split(",")
      : [],
  );
  let validated = $state(false);

  let doneValue = $state("");
  let participants: string[] = $state([]);

  let submitters: string[] = $state([]);
  let recipients: string[] = $state([]);
  let acceptors: string[] = $state([]);

  let descElement: HTMLElement;
  let noteElement: HTMLElement;

  onMount(async () => {
    const res = await requirements.init();
    participants = res.participants;
    doneValue = res.done;
    submitters = await requirements.submitters();
    recipients = await requirements.recipients();
    acceptors = await requirements.acceptors();
    descElement.scrollTop = 0;
    noteElement.scrollTop = 0;
  });

  const current = () => {
    return {
      type,
      desc,
      date,
      deadline,
      done,
      submitter,
      recipient,
      acceptor,
      status,
      note,
      participating: participating.join(","),
    } as Requirement;
  };

  const save = async () => {
    if (valid() && (!participants.length || participating.length > 0)) {
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
    {modeList[requirements.mode]}业务 - {requirements.brand || "业务系统"}
  </title>
</svelte:head>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div style="height: 100%;">
  <header>
    <div class="back">
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <span class="material-symbols-outlined" onclick={back}>arrow_back</span>
    </div>
    <h3>{modeList[requirements.mode]}业务</h3>
    {#if requirements.mode != "add"}
      <Action
        requirement={requirements.requirement}
        --icon="22px"
        --margin="10px"
      />
    {/if}
  </header>
  <div class="row g-3" class:was-validated={validated}>
    <div class="col-md-8 col-sm-12">
      <label for="desc" class="form-label">描述</label>
      <textarea
        class="form-control"
        id="desc"
        bind:this={descElement}
        bind:value={desc}
        required
        disabled={requirements.mode == "view"}
      ></textarea>
      <div class="invalid-feedback">必填字段</div>
    </div>
    <div class="w-100 m-0"></div>
    <div class="col-md-3 col-sm-4">
      <div class="form-floating">
        {#if requirements.mode == "view"}
          <input class="form-control" id="type" value={type} disabled />
        {:else}
          <select class="form-select" id="type" bind:value={type} required>
            {#each requirements.types as type (type)}
              <option value={type}>{type}</option>
            {/each}
          </select>
        {/if}
        <label for="type">类型</label>
        <div class="invalid-feedback">必填字段</div>
      </div>
    </div>
    <div class="col-md-3 col-sm-4">
      <div class="form-floating">
        {#if requirements.mode == "view"}
          <input class="form-control" id="status" value={status} disabled />
        {:else}
          <select class="form-select" id="status" bind:value={status} required>
            {#each requirements.statuses as status (status.value)}
              <option value={status.value}>{status.value}</option>
            {/each}
          </select>
        {/if}
        <label for="status">状态</label>
        <div class="invalid-feedback">必填字段</div>
      </div>
    </div>
    <div class="w-100 m-0"></div>
    <div class="col-md-3 col-sm-4">
      <div class="form-floating">
        <input
          class="form-control"
          id="date"
          type="date"
          bind:value={date}
          required
          disabled={requirements.mode == "view"}
        />
        <label for="date">提请日期</label>
        <div class="invalid-feedback">必填字段</div>
      </div>
    </div>
    <div class="col-md-3 col-sm-4">
      <div class="form-floating">
        <input
          class="form-control"
          id="deadline"
          type="date"
          bind:value={deadline}
          disabled={requirements.mode == "view"}
        />
        <label for="deadline">期限日期</label>
      </div>
    </div>
    {#if status === doneValue}
      <div class="col-md-3 col-sm-4">
        <div class="form-floating">
          <input
            class="form-control"
            id="done"
            type="date"
            bind:value={done}
            required
            disabled={requirements.mode == "view"}
          />
          <label for="deadline">完成日期</label>
          <div class="invalid-feedback">必填字段</div>
        </div>
      </div>
    {/if}
    <div class="w-100 m-0"></div>
    <div class="col-md-3 col-sm-4">
      <div class="form-floating">
        <input
          class="form-control"
          id="submitter"
          list="submitter-list"
          bind:value={submitter}
          placeholder="submitter"
          required
          disabled={requirements.mode == "view"}
        />
        <datalist id="submitter-list">
          {#each submitters as submitter (submitter)}
            <option>{submitter}</option>
          {/each}
        </datalist>
        <label for="submitter">提交人</label>
        <div class="invalid-feedback">必填字段</div>
      </div>
    </div>
    <div class="col-md-3 col-sm-4">
      <div class="form-floating">
        <input
          class="form-control"
          id="recipient"
          list="recipient-list"
          bind:value={recipient}
          placeholder="recipient"
          disabled={requirements.mode == "view"}
        />
        <datalist id="recipient-list">
          {#each recipients as recipient (recipient)}
            <option>{recipient}</option>
          {/each}
        </datalist>
        <label for="recipient">承接人</label>
      </div>
    </div>
    <div class="col-md-3 col-sm-4">
      <div class="form-floating">
        <input
          class="form-control"
          id="acceptor"
          list="acceptor-list"
          bind:value={acceptor}
          placeholder="acceptor"
          required
          disabled={requirements.mode == "view"}
        />
        <datalist id="acceptor-list">
          {#each acceptors as acceptor (acceptor)}
            <option>{acceptor}</option>
          {/each}
        </datalist>
        <label for="acceptor">受理人</label>
        <div class="invalid-feedback">必填字段</div>
      </div>
    </div>
    <div class="col-md-6">
      <label class="form-label" for="participating">参与班组</label>
      <div id="participating">
        {#if doneValue}
          {#if participants.length == 0}
            <div class="form-check form-check-inline">
              <input
                type="checkbox"
                class="form-check-input"
                id={"noparticipant"}
                disabled
              />
              <label class="form-check-label" for={"noparticipant"}>无</label>
            </div>
          {/if}
          {#each participants as participant, index (participant)}
            <div class="form-check form-check-inline">
              <input
                type="checkbox"
                class="form-check-input"
                class:invalid={validated && participating.length == 0}
                id={"participant" + index}
                bind:group={participating}
                value={participant}
                disabled={requirements.mode == "view"}
              />
              <label
                class="form-check-label"
                class:invalid={validated && participating.length == 0}
                for={"participant" + index}
              >
                {participant}
              </label>
            </div>
          {/each}
        {:else}
          <div class="form-check"></div>
        {/if}
      </div>
      <div
        class="invalid-feedback"
        class:invalid={validated &&
          participants.length &&
          participating.length == 0}
      >
        必选字段
      </div>
    </div>
    <div class="col-md-8 col-sm-12">
      <label for="note" class="form-label">备注</label>
      <textarea
        class="form-control"
        id="note"
        bind:this={noteElement}
        bind:value={note}
        disabled={requirements.mode == "view"}
      ></textarea>
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

  #desc {
    height: 5rem;
  }

  #participating {
    display: flex;
  }

  #note {
    height: 7rem;
  }

  .invalid {
    display: block;
    color: var(--bs-form-invalid-color) !important;
    border-color: var(--bs-form-invalid-border-color) !important;
  }
</style>
