<script lang="ts">
  import { fire, post } from "../misc.svelte";
  import { requirements } from "../requirement.svelte";

  let username = $state(localStorage.getItem("username") || "");
  let password = $state("");
  let rememberme = $state(localStorage.getItem("rememberme") === "true");
  let usernameInput: HTMLInputElement;
  let passwordInput: HTMLInputElement;

  const login = async () => {
    if (!usernameInput!.checkValidity())
      await fire("错误", "用户名不能为空", "error");
    else if (!passwordInput.checkValidity())
      await fire("错误", "密码不能为空", "error");
    else {
      const resp = await post("/login", { username, password, rememberme });
      if (resp.ok) {
        const json = await resp.json();
        if (json.status == 1) {
          localStorage.setItem("username", username);
          if (rememberme) localStorage.setItem("rememberme", "true");
          else localStorage.removeItem("rememberme");
          await requirements.init(true);
          requirements.goto("show");
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
  <h3 class="d-flex h-100 justify-content-center align-items-center">登录</h3>
</header>
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="login" onkeyup={handleEnter}>
  <div class="form-floating mb-3">
    <!-- svelte-ignore a11y_autofocus -->
    <input
      class="form-control"
      bind:this={usernameInput}
      bind:value={username}
      id="username"
      name="username"
      autocomplete="username"
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
      bind:this={passwordInput}
      bind:value={password}
      id="password"
      name="password"
      autocomplete="current-password"
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
  <button class="btn btn-primary login" onclick={login}>登录</button>
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
