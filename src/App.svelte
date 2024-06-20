<script lang="ts">
  import type { ComponentType } from "svelte";
  import Nav from "./components/Nav.svelte";
  import Login from "./components/Login.svelte";
  import Setting from "./components/Setting.svelte";
  import Show from "./components/Show.svelte";
  import Requirement from "./components/Requirement.svelte";
  import { mode, component, clear, loading } from "./stores";
  import { info } from "./requirement";

  let username = "";

  const load = async () => {
    clear();
    loading.start();
    const res = await info(true);
    loading.end();
    username = res.username;
  };

  const components: { [component: string]: ComponentType } = {
    setting: Setting,
    show: Show,
    requirement: Requirement,
  };

  const handlePopstate = () => {
    if (username) {
      switch (window.location.pathname) {
        case "/":
          $component = "show";
          return;
        case "/add":
          $mode = "add";
          break;
        case "/edit":
          $mode = "edit";
          break;
        default:
          $mode = "view";
      }
      $component = "requirement";
    }
  };
</script>

<svelte:window on:popstate={handlePopstate} />

<Nav bind:username on:reload={load} />
{#await load() then _}
  <div class="content" style="opacity: {$loading ? 0.5 : 1}">
    {#if !username}
      {#if !$loading}
        <Login on:info={load} />
      {/if}
    {:else}
      <svelte:component this={components[$component]} on:reload={load} />
    {/if}
  </div>
{/await}
<div class={username ? "loading" : "initializing"} hidden={!$loading}>
  <div class="sk-wave sk-center">
    <div class="sk-wave-rect" />
    <div class="sk-wave-rect" />
    <div class="sk-wave-rect" />
    <div class="sk-wave-rect" />
    <div class="sk-wave-rect" />
  </div>
</div>

<style>
  .content {
    position: fixed;
    top: 0;
    padding-top: 60px;
    width: 100%;
    height: 100%;
  }
</style>
