<script lang="ts">
  import type { Component } from "svelte";
  import Login from "./components/Login.svelte";
  import Nav from "./components/Nav.svelte";
  import Requirement from "./components/Requirement.svelte";
  import Setting from "./components/Setting.svelte";
  import Show from "./components/Show.svelte";
  import { loading } from "./misc.svelte";
  import { requirements } from "./requirement.svelte";

  const promise = requirements.init(true);

  const components: { [component: string]: Component } = {
    setting: Setting,
    show: Show,
    requirement: Requirement,
  };

  const Content = $derived(components[requirements.component]);

  const handlePopstate = () => {
    if (requirements.username) {
      switch (window.location.pathname) {
        case "/":
          requirements.component = "show";
          return;
        case "/add":
          requirements.mode = "add";
          break;
        case "/edit":
          requirements.mode = "edit";
          break;
        default:
          requirements.mode = "view";
      }
      requirements.component = "requirement";
    }
  };
</script>

<svelte:head><title>{requirements.brand || "业务系统"}</title></svelte:head>
<svelte:window onpopstate={handlePopstate} />

<Nav />
{#await promise then _}
  <div class="content" style:opacity={loading.show ? 0.5 : 1}>
    {#if !requirements.username}
      {#if !loading.show}
        <Login />
      {/if}
    {:else}
      <Content />
    {/if}
  </div>
{/await}
<div class="loading" hidden={!loading.show}>
  <div class="sk-wave sk-center">
    <div class="sk-wave-rect"></div>
    <div class="sk-wave-rect"></div>
    <div class="sk-wave-rect"></div>
    <div class="sk-wave-rect"></div>
    <div class="sk-wave-rect"></div>
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

  .loading {
    position: fixed;
    z-index: 2;
    top: 60px;
    height: calc(100% - 60px);
    width: 100%;
    display: flex;
  }
</style>
