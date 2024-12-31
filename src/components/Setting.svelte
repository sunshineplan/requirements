<script lang="ts">
  import Swal from "sweetalert2";
  import { confirm, fire, post, valid } from "../misc.svelte";
  import { requirements } from "../requirement.svelte";
  import Textarea from "./form-control/Textarea.svelte";

  let fields = $state("");
  let custom = $state("");
  let users: string[] = $state([]);
  let validated = $state(false);

  const height = "18rem";

  const load = async () => {
    users = await requirements.init();
    fields = JSON.stringify(requirements.fields.raw, null, 2);
    custom = JSON.stringify(requirements.fields.custom, null, 2);
  };
  const promise = load();

  const updateFields = async () => {
    if (valid()) {
      try {
        validated = false;
        const resp = await post("/fields", JSON.parse(fields));
        if (resp.ok) {
          await fire("成功", "字段保存成功", "success");
          await requirements.init();
        } else await fire("错误", await resp.text(), "error");
      } catch (e) {
        await fire("错误", String(e), "error");
      }
    } else validated = true;
  };

  const updateCustom = async () => {
    try {
      validated = false;
      const resp = await post("/custom", JSON.parse(custom));
      if (resp.ok) {
        await fire("成功", "自定义保存成功", "success");
        await requirements.init();
      } else await fire("错误", await resp.text(), "error");
    } catch (e) {
      await fire("错误", String(e), "error");
    }
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
<div class="h-100">
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
  {#await promise then _}
    <div class="row g-3" class:was-validated={validated}>
      <div class="col-md-6 col-sm-12 mt-0">
        <Textarea
          id="fields"
          label="字段"
          {height}
          bind:value={fields}
          required={true}
          absolute={true}
        />
        <button class="btn btn-primary float-end mt-2" onclick={updateFields}>
          保存字段
        </button>
      </div>
      <div class="col-md-6 col-sm-12 mt-0">
        <Textarea id="custom" label="自定义" {height} bind:value={custom} />
        <button class="btn btn-primary float-end mt-2" onclick={updateCustom}>
          保存自定义
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

  .row {
    padding: 0 20px;
    overflow: auto;
    margin-top: 0;
    max-height: calc(100% - 60px);
  }
</style>
