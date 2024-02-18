<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { fire, post } from "../misc";
  import { goto } from "../stores";

  const dispatch = createEventDispatcher();

  let username = localStorage.getItem("username");
  let password = "";
  let rememberme = localStorage.getItem("rememberme") === "true";

  const login = async () => {
    if (!document.querySelector<HTMLSelectElement>("#username").checkValidity())
      await fire("错误", "用户名不能为空", "error");
    else if (
      !document.querySelector<HTMLSelectElement>("#password").checkValidity()
    )
      await fire("错误", "密码不能为空", "error");
    else {
      const resp = await post("/login", { username, password, rememberme });
      if (resp.ok) {
        const json = await resp.json();
        if (json.status == 1) {
          localStorage.setItem("username", username);
          if (rememberme) localStorage.setItem("rememberme", "true");
          else localStorage.removeItem("rememberme");
          dispatch("info");
          goto("show");
        } else await fire("错误", json.message, "error");
      } else await fire("错误", await resp.text(), "error");
    }
  };

  const handleEnter = async (event: KeyboardEvent) => {
    if (event.key === "Enter") await login();
  };
</script>

<svelte:head><title>登录</title></svelte:head>

<header>
  <h3
    class="d-flex justify-content-center align-items-center"
    style="height: 100%"
  >
    登录
  </h3>
</header>
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="login" on:keyup={handleEnter}>
  <div class="form-floating mb-3">
    <!-- svelte-ignore a11y-autofocus -->
    <input
      class="form-control"
      bind:value={username}
      id="username"
      maxlength="20"
      placeholder="Username"
      autofocus
      required
    />
    <label for="username">用户名</label>
  </div>
  <div class="form-floating mb-3">
    <input
      class="form-control"
      type="password"
      bind:value={password}
      id="password"
      maxlength="20"
      placeholder="Password"
      required
    />
    <label for="password">密码</label>
  </div>
  <div class="mb-3 form-check">
    <input
      type="checkbox"
      class="form-check-input"
      bind:checked={rememberme}
      id="rememberme"
    />
    <label class="form-check-label" for="rememberme">记住我</label>
  </div>
  <hr />
  <button class="btn btn-primary login" on:click={login}>登录</button>
</div>

<style>
  .login {
    width: 250px;
    margin: 0 auto;
  }

  .form-control {
    width: 250px;
  }
</style>
