<script lang="ts">
  import Nav from "./components/Nav.svelte";
  import Login from "./components/Login.svelte";
  import Show from "./components/Show.svelte";
  import Requirement from "./components/Requirement.svelte";
  import { component, loading } from "./stores";
  import { init } from "./requirement";
  import { onMount } from "svelte";

  let username: string = "";

  const load = async () => {
    loading.start();
    username = await init();
    loading.end();
  };
  const promise = load();

  const components: {
    [component: string]: typeof Show | typeof Requirement;
  } = {
    show: Show,
    requirement: Requirement,
  };

  const handlePopstate = () => {
    if (username)
      switch (window.location.pathname) {
        case "/":
          $component = "show";
          break;
        case "/add":
        case "/edit":
          $component = "requirement";
      }
  };

  onMount(() => {
    window.addEventListener("popstate", handlePopstate);
    return () => window.removeEventListener("popstate", handlePopstate);
  });
</script>

<Nav bind:username on:reload={load} />
{#await promise then _}
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
