<script lang="ts">
  import Swal from "sweetalert2";
  import { confirm, fire, post, valid } from "../misc.svelte";
  import { requirements } from "../requirement.svelte";

  let participants = $state("");
  let types = $state("");
  let users: string[] = $state([]);
  let validated = $state(false);

  const load = async () => {
    const res = await requirements.init();
    participants = res.participants.join("\n");
    types = res.types.join("\n");
    users = res.users;
  };

  const updateTypes = async () => {
    if (valid()) {
      validated = false;
      const restult = types.split("\n").filter(Boolean);
      const resp = await post("/types", restult);
      if (resp.ok) {
        await fire("成功", "保存成功", "success");
      } else await fire("错误", await resp.text(), "error");
    } else validated = true;
  };

  const updateParticipants = async () => {
    if (valid()) {
      validated = false;
      const restult = participants.split("\n").filter(Boolean);
      const resp = await post("/participants", restult);
      if (resp.ok) {
        await fire("成功", "保存成功", "success");
      } else await fire("错误", await resp.text(), "error");
    } else validated = true;
  };

  const addUser = async () => {
    const { value: user } = await Swal.fire({
      title: "添加用户",
      html: `
<div id="addUser">
  <div class="mx-5 mb-3">
    <label for="username" class="form-label">用户名</label>
    <input class="form-control" id="newUser" placeholder="Username" required autofocus />
    <div class="invalid-feedback">不能为空</div>
  </div>
  <div class="mx-5 mb-1">
    <label for="password" class="form-label">密码</label>
    <input class="form-control" type="password" id="newPwd" placeholder="Password" required />
    <div class="invalid-feedback">不能为空</div>
  </div>
</div>`,
      focusConfirm: false,
      confirmButtonText: "添加",
      cancelButtonText: "取消",
      showCancelButton: true,
      allowOutsideClick: false,
      buttonsStyling: false,
      customClass: {
        confirmButton: "swal btn btn-primary",
        cancelButton: "swal btn btn-primary",
      },
      preConfirm: () => {
        if (valid())
          return {
            username:
              document.querySelector<HTMLInputElement>("#newUser")!.value,
            password:
              document.querySelector<HTMLInputElement>("#newPwd")!.value,
          };
        document.getElementById("addUser")!.classList.add("was-validated");
        return false;
      },
    });
    if (user) {
      const resp = await post("/addUser", user);
      if (resp.ok) {
        users = [...users, user.username];
        await fire("成功", "添加成功", "success");
      } else await fire("错误", await resp.text(), "error");
    }
  };

  const reset = async (username: string) => {
    const { value: password } = await Swal.fire({
      title: `重置密码(${username})`,
      html: `
<div id="reset">
  <div class="mx-5 mb-1">
    <label for="password" class="form-label">密码</label>
    <input class="form-control" type="password" id="chgpwd" placeholder="Password" required />
    <div class="invalid-feedback">不能为空</div>
  </div>
</div>`,
      focusConfirm: false,
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      showCancelButton: true,
      allowOutsideClick: false,
      buttonsStyling: false,
      customClass: {
        confirmButton: "swal btn btn-primary",
        cancelButton: "swal btn btn-primary",
      },
      preConfirm: () => {
        if (valid())
          return document.querySelector<HTMLInputElement>("#chgpwd")!.value;
        document.getElementById("reset")!.classList.add("was-validated");
        return false;
      },
    });
    if (password) {
      const resp = await post("/chgpwd", { username, password });
      if (resp.ok) {
        await fire("成功", "重置成功", "success");
      } else await fire("错误", await resp.text(), "error");
    }
  };

  const del = async (username: string) => {
    if (await confirm("该用户将被永久删除。", true)) {
      const resp = await post("/deleteUser", { username });
      if (resp.ok) {
        users = users.filter((user) => user != username);
        await fire("成功", "删除成功", "success");
      } else await fire("错误", await resp.text(), "error");
    }
  };
</script>

<svelte:head>
  <title>设置 - {requirements.brand || "业务系统"}</title>
</svelte:head>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div style="height: 100%;">
  <header>
    <div class="back">
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <span
        class="material-symbols-outlined"
        onclick={() => requirements.goto("show")}
      >
        arrow_back
      </span>
    </div>
    <h3>设置</h3>
  </header>
  {#await load() then _}
    <div class="row g-3" class:was-validated={validated}>
      <div class="col-md-6 col-sm-12">
        <label for="types" class="form-label">类型</label>
        <textarea class="form-control" id="types" bind:value={types} required
        ></textarea>
        <div class="invalid-feedback">必填字段</div>
        <button class="btn btn-primary float-end mt-2" onclick={updateTypes}>
          保存类型
        </button>
      </div>
      <div class="col-md-6 col-sm-12">
        <label for="participants" class="form-label">班组</label>
        <textarea
          class="form-control"
          id="participants"
          bind:value={participants}
        ></textarea>
        <button
          class="btn btn-primary float-end mt-2"
          onclick={updateParticipants}
        >
          保存班组
        </button>
      </div>
      <hr />
      <div class="col-md-4 col-sm-12">
        <label class="form-label" for="users">用户</label>
        <ol class="list-group list-group-numbered" id="users">
          {#each users as user (user)}
            <li class="list-group-item d-flex justify-content-between">
              {user}
              <div class="d-flex">
                <!-- svelte-ignore a11y_click_events_have_key_events -->
                <!-- svelte-ignore a11y_no_static_element_interactions -->
                <span
                  title="重置"
                  class="material-symbols-outlined link-primary"
                  onclick={async () => reset(user)}
                >
                  lock_reset
                </span>
                {#if user != "admin"}
                  <!-- svelte-ignore a11y_click_events_have_key_events -->
                  <!-- svelte-ignore a11y_no_static_element_interactions -->
                  <span
                    title="删除"
                    class="material-symbols-outlined link-danger"
                    onclick={async () => del(user)}
                  >
                    delete
                  </span>
                {:else}
                  <span class="material-symbols-outlined hidden">delete</span>
                {/if}
              </div>
            </li>
          {/each}
        </ol>
        <button class="btn btn-primary float-end my-2" onclick={addUser}>
          添加用户
        </button>
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
  }

  #participants,
  #types {
    height: 9rem;
  }

  .row {
    padding: 0 20px;
    overflow: auto;
    margin-top: 0;
    max-height: calc(100% - 60px);
  }

  .invalid-feedback {
    position: absolute;
  }
</style>
