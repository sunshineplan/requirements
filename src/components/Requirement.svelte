<script lang="ts">
  import { onMount, createEventDispatcher } from "svelte";
  import { valid, confirm } from "../misc";
  import { component } from "../stores";
  import { requirement, requirements } from "../requirement";

  const dispatch = createEventDispatcher();

  const types = [
    "内容策划",
    "宣传推广",
    "用户培训",
    "宣传品相关",
    "平台相关",
    "中心业务",
    "馆所业务",
  ];
  const statuses = ["进行中", "已完成", "已关闭"];
  const mode = window.location.pathname == "/add" ? "新增" : "编辑";

  let type = $requirement.type || "";
  let desc = $requirement.desc || "";
  let date = $requirement.date || "";
  let deadline = $requirement.deadline || "";
  let submitter = $requirement.submitter || "";
  let recipient = $requirement.recipient || "";
  let acceptor = $requirement.acceptor || "";
  let status = $requirement.status || "";
  let note = $requirement.note || "";
  let participating = $requirement.participating || "";
  let validated = false;

  let submitters: string[] = [];
  let recipients: string[] = [];
  let acceptors: string[] = [];

  onMount(async () => {
    submitters = await requirements.submitters();
    recipients = await requirements.recipients();
    acceptors = await requirements.acceptors();
    document.getElementById("desc").scrollTop = 0;
    document.getElementById("note").scrollTop = 0;
  });

  const save = async () => {
    if (valid()) {
      validated = false;
      const r = <Requirement>{
        type,
        desc,
        date,
        deadline,
        submitter,
        recipient,
        acceptor,
        status,
        note,
        participating,
      };
      if (mode == "编辑") r.id = $requirement.id;
      try {
        const res = await requirements.save(r);
        if (res === 0) goback();
      } catch {
        dispatch("reload");
        goback();
      }
    } else validated = true;
  };

  const del = async () => {
    if (await confirm("业务")) {
      try {
        await requirements.delete($requirement);
      } catch {
        dispatch("reload");
      }
      goback();
    }
  };

  const goback = () => {
    window.history.pushState({}, "", "/");
    $component = "show";
  };
</script>

<svelte:head><title>{mode}业务 - 业务系统</title></svelte:head>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div style="height: 100%;">
  <header>
    <h3>{mode}业务</h3>
  </header>
  <div class="row g-3" class:was-validated={validated}>
    <div class="col-md-8 col-sm-12">
      <label class="form-label" for="desc">描述</label>
      <!-- svelte-ignore a11y-autofocus -->
      <textarea
        class="form-control"
        id="desc"
        bind:value={desc}
        rows="3"
        autofocus
        required
      />
      <div class="invalid-feedback">必填字段</div>
    </div>
    <div class="w-100 m-0" />
    <div class="col-md-3 col-sm-4">
      <label class="form-label" for="type">类型</label>
      <select class="form-select" id="type" bind:value={type} required>
        {#each types as type (type)}
          <option value={type}>{type}</option>
        {/each}
      </select>
      <div class="invalid-feedback">必填字段</div>
    </div>
    <div class="col-md-3 col-sm-4">
      <label class="form-label" for="status">状态</label>
      <select class="form-select" id="status" bind:value={status} required>
        {#each statuses as status (status)}
          <option value={status}>{status}</option>
        {/each}
      </select>
      <div class="invalid-feedback">必填字段</div>
    </div>
    <div class="w-100 m-0" />
    <div class="col-md-3 col-sm-4">
      <label class="form-label" for="date">提请日期</label>
      <input
        class="form-control"
        id="date"
        type="date"
        bind:value={date}
        required
      />
      <div class="invalid-feedback">必填字段</div>
    </div>
    <div class="col-md-3 col-sm-4">
      <label class="form-label" for="deadline">期限日期</label>
      <input
        class="form-control"
        id="deadline"
        type="date"
        bind:value={deadline}
        required
      />
      <div class="invalid-feedback">必填字段</div>
    </div>
    <div class="w-100 m-0" />
    <div class="col-md-3 col-sm-4">
      <label class="form-label" for="submitter">提交人</label>
      <input
        class="form-control"
        id="submitter"
        list="submitter-list"
        bind:value={submitter}
        required
      />
      <datalist id="submitter-list">
        {#each submitters as submitter (submitter)}
          <option>{submitter}</option>
        {/each}
      </datalist>
      <div class="invalid-feedback">必填字段</div>
    </div>
    <div class="col-md-3 col-sm-4">
      <label class="form-label" for="recipient">承接人</label>
      <input
        class="form-control"
        id="recipient"
        list="recipient-list"
        bind:value={recipient}
        required
      />
      <datalist id="recipient-list">
        {#each recipients as recipient (recipient)}
          <option>{recipient}</option>
        {/each}
      </datalist>
      <div class="invalid-feedback">必填字段</div>
    </div>
    <div class="col-md-3 col-sm-4">
      <label class="form-label" for="acceptor">受理人</label>
      <input
        class="form-control"
        id="acceptor"
        list="acceptor-list"
        bind:value={acceptor}
        required
      />
      <datalist id="acceptor-list">
        {#each acceptors as acceptor (acceptor)}
          <option>{acceptor}</option>
        {/each}
      </datalist>
      <div class="invalid-feedback">必填字段</div>
    </div>
    <div class="col-md-6">
      <label class="form-label" for="participating">参与班组</label>
      <input
        class="form-control"
        id="participating"
        bind:value={participating}
        required
      />
      <div class="invalid-feedback">必填字段</div>
    </div>
    <div class="col-md-8 col-sm-12">
      <label class="form-label" for="note">备注</label>
      <textarea class="form-control" id="note" bind:value={note} />
    </div>
    <div class="col-12">
      <button class="btn btn-primary" on:click={save}>保存</button>
      <button class="btn btn-primary" on:click={goback}>取消</button>
    </div>
    {#if mode == "编辑"}
      <div class="col-12">
        <button class="btn btn-danger" on:click={del}>删除</button>
      </div>
    {/if}
  </div>
</div>

<style>
  header {
    padding-left: 20px;
  }

  .row {
    padding: 0 20px;
    overflow: auto;
    max-height: calc(100% - 60px);
  }
</style>
