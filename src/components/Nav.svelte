<script lang="ts">
  import Swal from "sweetalert2";
  import { createEventDispatcher } from "svelte";
  import { fire, post } from "../misc";
  import { component } from "../stores";

  const dispatch = createEventDispatcher();

  export let username: string;

  const statistics = async () => {
    let url = "/statistics";
    const select = await Swal.fire({
      title: "请选择统计方法",
      confirmButtonText: "提请日期至期限日期",
      denyButtonText: "仅提请日期",
      showDenyButton: true,
      showCloseButton: true,
      customClass: {
        confirmButton: "swal btn btn-primary",
        denyButton: "swal btn btn-primary",
      },
      buttonsStyling: false,
    });
    if (select.isDenied) url += "?isNew=1";
    else if (select.isDismissed) return;
    const resp = await fetch(url);
    if (resp.ok) {
      const blob = await resp.blob();
      console.log(blob.type);
      const link = document.createElement("a");
      link.href = URL.createObjectURL(blob);
      link.download = "统计.csv";
      link.click();
      URL.revokeObjectURL(link.href);
    } else await fire("错误", "内部错误", "error");
  };

  const logout = async () => {
    const resp = await post("/logout", undefined);
    if (resp.ok) {
      dispatch("reload");
      window.history.pushState({}, "", "/");
      $component = "show";
    } else await fire("错误", "未知错误", "error");
  };
</script>

<nav class="navbar navbar-light topbar">
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <span
    class="brand"
    on:click={() => {
      window.history.pushState({}, "", "/");
      $component = "show";
    }}
  >
    公关推广部业务系统
  </span>
  <div class="navbar-nav flex-row">
    {#if username}
      <span class="nav-link">{username}</span>
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <!-- svelte-ignore a11y-no-static-element-interactions -->
      <span class="nav-link link" on:click={statistics}>统计</span>
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <!-- svelte-ignore a11y-no-static-element-interactions -->
      <span class="nav-link link" on:click={logout}>退出</span>
    {:else}
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <!-- svelte-ignore a11y-no-static-element-interactions -->
      <span class="nav-link">登录</span>
    {/if}
  </div>
</nav>

<style>
  .topbar {
    position: fixed;
    top: 0;
    z-index: 2;
    width: 100%;
    height: 60px;
    background-color: #1a73e8;
    padding: 0.5rem 1rem;
    user-select: none;
  }

  .brand {
    padding-left: 20px;
    font-size: 25px;
    letter-spacing: 0.3px;
    color: white;
    cursor: pointer;
  }

  .nav-link {
    padding-left: 8px;
    padding-right: 8px;
    color: white !important;
    cursor: pointer;
  }

  .link:hover {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 5px;
  }
</style>
