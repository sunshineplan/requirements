<script lang="ts">
  import { stringify } from "csv-stringify/browser/esm/sync";
  import Cookies from "js-cookie";
  import { createEventDispatcher, onMount } from "svelte";
  import { poll } from "../misc";
  import {
    columns,
    headers,
    info,
    requirement,
    requirements,
    searchable,
  } from "../requirement";
  import {
    desc,
    goto,
    loading,
    name,
    scroll,
    search,
    searchField,
    sort,
  } from "../stores";
  import Action from "./Action.svelte";
  import Search from "./Search.svelte";

  const dispatch = createEventDispatcher();

  let output: Requirement[] = [];

  $: $search, $searchField, $sort, $desc, filter();

  const getField = (r: Requirement, key: string) => {
    return r[key as keyof Requirement];
  };

  const add = () => {
    $requirement = <Requirement>{};
    goto("add");
  };

  const view = async (e: MouseEvent, r: Requirement) => {
    await new Promise((sleep) => setTimeout(sleep, 50));
    if (window.getSelection()?.toString() !== "") return;
    if ((e.target as HTMLElement).dataset["action"] != "done") {
      $requirement = r;
      goto("view");
    }
  };

  const download = () => {
    const link = document.createElement("a");
    const file = new Blob(
      [
        stringify(output, {
          bom: true,
          header: true,
          columns: Object.keys(headers).map((key) => ({
            key,
            header: headers[key],
          })),
        }),
      ],
      {
        type: "text/csv",
      },
    );
    link.href = URL.createObjectURL(file);
    link.download = "download.csv";
    link.click();
    URL.revokeObjectURL(link.href);
  };

  const filter = () => {
    let array: Requirement[] = [];
    if (!$search) array = $requirements;
    else if ($searchField)
      array = $requirements.filter((i) => i[$searchField].includes($search));
    else
      array = $requirements.filter((i) =>
        searchable.some((field) => i[field].includes($search)),
      );
    if (!$sort) output = array.sort();
    else
      output = array.toSorted((a, b) => {
        const v1 = getField(a, $sort),
          v2 = getField(b, $sort);
        let res = 0;
        if (v1 < v2) res = 1;
        else if (v1 > v2) res = -1;
        if ($desc) return res;
        else return -res;
      });
  };

  const restore = () => {
    filter();
    scroll(true);
  };

  const participants = (s: string) => {
    const res = s.split(",");
    if (res.length > 1) {
      return res.map((i) => i.charAt(0)).join(" | ");
    }
    return s;
  };

  const subscribe = async (signal: AbortSignal) => {
    const resp = await poll(signal);
    if (resp.ok) {
      if (Cookies.get("last") != (await resp.text())) {
        loading.start();
        await info(true);
        filter();
        loading.end();
      }
      await subscribe(signal);
    } else if (resp.status == 401) {
      dispatch("reload");
    } else {
      await new Promise((sleep) => setTimeout(sleep, 30000));
      await subscribe(signal);
    }
  };
  onMount(() => {
    const controller = new AbortController();
    subscribe(controller.signal);
    return () => controller.abort();
  });

  onMount(() => scroll(true));
</script>

<svelte:head><title>{$name || "业务系统"}</title></svelte:head>

<header>
  <button class="btn btn-primary" on:click={add}>新增业务</button>
  <button class="btn btn-primary" on:click={download}>导出</button>
  <Search />
</header>
<div class="table-responsive">
  <table class="table table-hover table-sm">
    <thead>
      <tr>
        {#each Object.entries(columns) as [key, width] (key)}
          {#if width}
            <th
              class="sortable {$sort == key
                ? $desc
                  ? 'desc'
                  : 'asc'
                : 'default'}"
              class:auto={width == -1}
              style:width={width > 0 ? `${width}rem` : ""}
              on:click={() => {
                const before = $sort;
                $sort = key;
                if (before == $sort) $desc = !$desc;
                else $desc = true;
              }}
            >
              {headers[key]}
            </th>
          {/if}
        {/each}
        <th />
      </tr>
    </thead>
    <tbody>
      {#each output as requirement (requirement.id)}
        <tr on:click={(e) => view(e, requirement)}>
          {#each Object.entries(columns) as [key, show] (key)}
            {#if show}
              <td
                title={/编号|类型|日期|班组/i.test(headers[key])
                  ? ""
                  : getField(requirement, key)}
              >
                {key == "participating"
                  ? participants(requirement[key])
                  : getField(requirement, key)}
              </td>
            {/if}
          {/each}
          <td style="vertical-align: middle">
            <Action
              {requirement}
              --icon="18px"
              --margin="2px"
              on:reload
              on:refresh={restore}
            />
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<style>
  header {
    height: 60px;
  }

  table {
    table-layout: fixed;
  }

  .table-responsive {
    height: calc(100% - 60px);
    padding: 0px 10px;
    cursor: default;
    width: 100%;
  }

  th {
    position: sticky;
    top: 0;
    background-color: white;
    user-select: none;
  }

  @media (max-width: 1200px) {
    .auto {
      width: 10rem;
    }
  }

  th:nth-last-of-type(1) {
    width: 3.2rem;
  }

  td {
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow: hidden;
  }

  .sortable {
    cursor: pointer;
    background-position: right;
    background-repeat: no-repeat;
    padding-right: 30px !important;
  }

  .default {
    background-image: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABMAAAATCAQAAADYWf5HAAAAkElEQVQoz7X QMQ5AQBCF4dWQSJxC5wwax1Cq1e7BAdxD5SL+Tq/QCM1oNiJidwox0355mXnG/DrEtIQ6azioNZQxI0ykPhTQIwhCR+BmBYtlK7kLJYwWCcJA9M4qdrZrd8pPjZWPtOqdRQy320YSV17OatFC4euts6z39GYMKRPCTKY9UnPQ6P+GtMRfGtPnBCiqhAeJPmkqAAAAAElFTkSuQmCC");
  }

  .asc {
    background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABMAAAATCAYAAAByUDbMAAAAZ0lEQVQ4y2NgGLKgquEuFxBPAGI2ahhWCsS/gDibUoO0gPgxEP8H4ttArEyuQYxAPBdqEAxPBImTY5gjEL9DM+wTENuQahAvEO9DMwiGdwAxOymGJQLxTyD+jgWDxCMZRsEoGAVoAADeemwtPcZI2wAAAABJRU5ErkJggg==);
  }

  .desc {
    background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABMAAAATCAYAAAByUDbMAAAAZUlEQVQ4y2NgGAWjYBSggaqGu5FA/BOIv2PBIPFEUgxjB+IdQPwfC94HxLykus4GiD+hGfQOiB3J8SojEE9EM2wuSJzcsFMG4ttQgx4DsRalkZENxL+AuJQaMcsGxBOAmGvopk8AVz1sLZgg0bsAAAAASUVORK5CYII=);
  }
</style>
