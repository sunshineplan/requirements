<script lang="ts">
  import Swal from "sweetalert2";
  import { fire, post } from "../misc.svelte";
  import { requirements } from "../requirement.svelte";

  const statistics = async () => {
    let url = "/statistics";
    const select = await Swal.fire({
      title: "请选择统计方法",
      confirmButtonText: `${requirements.fields.name("date")}至${requirements.fields.name("deadline")}`,
      denyButtonText: `仅${requirements.fields.name("date")}`,
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
      const link = document.createElement("a");
      link.href = URL.createObjectURL(blob);
      link.download = "统计.csv";
      link.click();
      URL.revokeObjectURL(link.href);
    } else await fire("错误", "内部错误", "error");
  };

  const logout = async () => {
    requirements.abort();
    const resp = await post("/logout");
    if (resp.ok) {
      await requirements.init();
      requirements.goto("show");
    } else await fire("错误", "未知错误", "error");
  };
</script>

<nav class="navbar navbar-light topbar">
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <span class="brand" onclick={() => requirements.goto("show")}>
    {requirements.brand || "业务系统"}
  </span>
  <div class="navbar-nav flex-row">
    {#if requirements.username}
      <span class="nav-link">{requirements.username}</span>
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <span class="nav-link link" onclick={statistics}>统计</span>
      {#if requirements.username == "admin"}
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <span class="nav-link link" onclick={() => requirements.goto("setting")}
          >设置</span
        >
      {/if}
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <span class="nav-link link" onclick={logout}>退出</span>
    {:else}
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
