<script lang="ts">
  import Action from "./Action.svelte";
  import { onMount, createEventDispatcher } from "svelte";
  import { valid, confirm } from "../misc";
  import { loading, mode, goto, clear } from "../stores";
  import { requirement, requirements, statuses, info } from "../requirement";

  const dispatch = createEventDispatcher();

  const modeList: { [key: string]: string } = {
    add: "新增",
    edit: "编辑",
    view: "查看",
  };

  let type = $requirement.type || "";
  let desc = $requirement.desc || "";
  let date = $requirement.date || "";
  let deadline = $requirement.deadline || "";
  let done = $requirement.done || "";
  let submitter = $requirement.submitter || "";
  let recipient = $requirement.recipient || "";
  let acceptor = $requirement.acceptor || "";
  let status = $requirement.status || "";
  let note = $requirement.note || "";
  let participating = $requirement.participating
    ? $requirement.participating.split(",")
    : [];
  let validated = false;

  let doneValue = "";
  let participants: string[] = [];
  let types: string[] = [];

  let submitters: string[] = [];
  let recipients: string[] = [];
  let acceptors: string[] = [];

  onMount(async () => {
    loading.start();
    const res = await info();
    participants = res.participants;
    types = res.types;
    doneValue = res.done;
    loading.end();
    submitters = await requirements.submitters();
    recipients = await requirements.recipients();
    acceptors = await requirements.acceptors();
    document.getElementById("desc")!.scrollTop = 0;
    document.getElementById("note")!.scrollTop = 0;
  });

  const current = () => {
    return <Requirement>{
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
    };
  };

  const save = async () => {
    if (valid() && (!participants.length || participating.length > 0)) {
      validated = false;
      const r = current();
      if ($mode == "edit") r.id = $requirement.id;
      try {
        if (r.status != doneValue) r.done = "";
        const res = await requirements.save(r);
        if (res === 0) {
          if ($mode == "add") clear();
          goto("show");
        }
      } catch {
        dispatch("reload");
        goto("show");
      }
    } else validated = true;
  };

  const back = async () => {
    const r = current();
    let edited = false;
    switch ($mode) {
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
          if (r[key] != $requirement[key]) {
            edited = true;
            break;
          }
        }
    }
    if (edited && !(await confirm("数据未保存，确定将放弃保存并返回。", true)))
      return;
    goto("show");
  };
</script>

<svelte:head><title>{modeList[$mode]}业务 - 业务系统</title></svelte:head>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div style="height: 100%;">
  <header>
    <div class="back">
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <span class="material-symbols-outlined" on:click={back}>arrow_back</span>
    </div>
    <h3>{modeList[$mode]}业务</h3>
    {#if $mode != "add"}
      <Action
        requirement={$requirement}
        --icon="22px"
        --margin="10px"
        on:reload
      />
    {/if}
  </header>
  <div class="row g-3" class:was-validated={validated}>
    <div class="col-md-8 col-sm-12">
      <div class="form-floating">
        <textarea
          class="form-control"
          id="desc"
          bind:value={desc}
          placeholder="desc"
          required
          disabled={$mode == "view"}
        />
        <label for="desc">描述</label>
        <div class="invalid-feedback">必填字段</div>
      </div>
    </div>
    <div class="w-100 m-0" />
    <div class="col-md-3 col-sm-4">
      <div class="form-floating">
        {#if $mode == "view"}
          <input class="form-control" id="type" value={type} disabled />
        {:else}
          <select class="form-select" id="type" bind:value={type} required>
            {#each types as type (type)}
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
        {#if $mode == "view"}
          <input class="form-control" id="status" value={status} disabled />
        {:else}
          <select class="form-select" id="status" bind:value={status} required>
            {#each $statuses as status (status.value)}
              <option value={status.value}>{status.value}</option>
            {/each}
          </select>
        {/if}
        <label for="status">状态</label>
        <div class="invalid-feedback">必填字段</div>
      </div>
    </div>
    <div class="w-100 m-0" />
    <div class="col-md-3 col-sm-4">
      <div class="form-floating">
        <input
          class="form-control"
          id="date"
          type="date"
          bind:value={date}
          required
          disabled={$mode == "view"}
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
          required
          disabled={$mode == "view"}
        />
        <label for="deadline">期限日期</label>
        <div class="invalid-feedback">必填字段</div>
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
            disabled={$mode == "view"}
          />
          <label for="deadline">完成日期</label>
          <div class="invalid-feedback">必填字段</div>
        </div>
      </div>
    {/if}
    <div class="w-100 m-0" />
    <div class="col-md-3 col-sm-4">
      <div class="form-floating">
        <input
          class="form-control"
          id="submitter"
          list="submitter-list"
          bind:value={submitter}
          placeholder="submitter"
          required
          disabled={$mode == "view"}
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
          required
          disabled={$mode == "view"}
        />
        <datalist id="recipient-list">
          {#each recipients as recipient (recipient)}
            <option>{recipient}</option>
          {/each}
        </datalist>
        <label for="recipient">承接人</label>
        <div class="invalid-feedback">必填字段</div>
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
          disabled={$mode == "view"}
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
                disabled={$mode == "view"}
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
      <div class="form-floating">
        <textarea
          class="form-control"
          id="note"
          bind:value={note}
          placeholder="note"
          disabled={$mode == "view"}
        />
        <label for="note">备注</label>
      </div>
    </div>
    <div class="col-md-8 col-sm-12">
      {#if $mode == "view"}
        <button class="btn btn-primary float-end mb-2" on:click={back}>
          返回
        </button>
      {:else}
        <button class="btn btn-primary float-end mb-2" on:click={save}>
          保存
        </button>
        <button class="btn btn-primary float-end mx-2 mb-2" on:click={back}>
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
    height: 70px;
  }

  header h3 {
    margin: 0;
  }

  .row {
    padding: 0 20px;
    overflow: auto;
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
