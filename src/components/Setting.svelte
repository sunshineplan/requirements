<script lang="ts">
  import Swal from "sweetalert2";
  import { valid, fire, post, confirm } from "../misc";
  import { loading, goto } from "../stores";
  import { info } from "../requirement";

  let teams = "";
  let users: string[] = [];
  let validated = false;

  const load = async () => {
    loading.start();
    const res = await info();
    loading.end();
    teams = res.participants.join("\n");
    users = res.users;
  };

  const updateParticipants = async () => {
    if (valid()) {
      validated = false;
      const restult = teams.split("\n").filter(Boolean);
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
    <input class="form-control" id="username" placeholder="Username" required autofocus />
    <div class="invalid-feedback">不能为空</div>
  </div>
  <div class="mx-5 mb-1">
    <label for="password" class="form-label">密码</label>
    <input class="form-control" type="password" id="password" placeholder="Password" required />
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
            username: (document.getElementById("username") as HTMLInputElement)
              .value,
            password: (document.getElementById("password") as HTMLInputElement)
              .value,
          };
        document.getElementById("addUser").classList.add("was-validated");
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
          return (document.getElementById("chgpwd") as HTMLInputElement).value;
        document.getElementById("reset").classList.add("was-validated");
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

<svelte:head><title>设置 - 业务系统</title></svelte:head>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div style="height: 100%;">
  <header>
    <div class="back">
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <span class="material-symbols-outlined" on:click={() => goto("show")}>
        arrow_back
      </span>
    </div>
    <h3>设置</h3>
  </header>
  {#await load() then _}
    <div class="row g-3" class:was-validated={validated}>
      <div class="col-md-6 col-sm-12">
        <label class="form-label" for="participants">班组</label>
        <textarea
          class="form-control"
          id="participants"
          bind:value={teams}
          rows="6"
          required
        />
        <div class="invalid-feedback">必填字段</div>
        <button
          class="btn btn-primary float-end mt-2"
          on:click={updateParticipants}
        >
          保存
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
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <!-- svelte-ignore a11y-no-static-element-interactions -->
                <span
                  title="重置"
                  class="material-symbols-outlined link-primary"
                  on:click={async () => reset(user)}
                >
                  lock_reset
                </span>
                {#if user != "admin"}
                  <!-- svelte-ignore a11y-click-events-have-key-events -->
                  <!-- svelte-ignore a11y-no-static-element-interactions -->
                  <span
                    title="删除"
                    class="material-symbols-outlined link-danger"
                    on:click={async () => del(user)}
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
        <button class="btn btn-primary float-end mt-2" on:click={addUser}>
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
    height: 70px;
  }

  header h3 {
    margin: 0;
  }

  .back {
    height: 50px;
    width: 50px;
    margin-right: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .back:hover {
    background-color: rgba(15, 20, 25, 0.1);
    border-radius: 50%;
  }

  .back span {
    font-size: 30px;
    cursor: default;
  }

  .row {
    padding: 0 20px;
    overflow: auto;
    max-height: calc(100% - 60px);
  }
</style>
