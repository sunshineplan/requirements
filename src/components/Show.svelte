<script lang="ts">
  import Cookies from "js-cookie";
  import Action from "./Action.svelte";
  import { onMount, createEventDispatcher } from "svelte";
  import { stringify } from "csv-stringify/browser/esm/sync";
  import { search, sort, desc, goto, scroll, loading } from "../stores";
  import { requirement, requirements, info } from "../requirement";
  import { poll } from "../misc";

  const dispatch = createEventDispatcher();

  const headers: { [key: string]: string } = {
    id: "编号",
    type: "类型",
    desc: "描述",
    date: "提请日期",
    deadline: "期限日期",
    done: "完成日期",
    submitter: "提交人",
    recipient: "承接人",
    acceptor: "受理人",
    status: "状态",
    note: "备注",
    participating: "参与班组",
  };

  const columns: { [key in keyof Requirement]: boolean } = {
    id: true,
    type: true,
    desc: true,
    date: true,
    deadline: true,
    done: true,
    submitter: true,
    recipient: false,
    acceptor: true,
    status: true,
    note: true,
    participating: true,
  };

  let output: Requirement[] = [];

  $: $search, $sort, $desc, filter();

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
    else
      array = $requirements.filter((i) => {
        return (
          i.type.includes($search) ||
          i.desc.includes($search) ||
          i.submitter.includes($search) ||
          i.recipient.includes($search) ||
          i.acceptor.includes($search) ||
          i.note.includes($search)
        );
      });
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

<svelte:head><title>业务系统</title></svelte:head>

<header>
  <button class="btn btn-primary" on:click={add}>新增业务</button>
  <button class="btn btn-primary" on:click={download}>导出</button>
  <div class="search">
    <div class="icon">
      <span class="material-symbols-outlined">search</span>
    </div>
    <input
      bind:value={$search}
      type="search"
      placeholder="搜索"
      on:input={() => scroll()}
    />
  </div>
</header>
<div class="table-responsive">
  <table class="table table-hover table-sm">
    <thead>
      <tr>
        {#each Object.entries(columns) as [key, show] (key)}
          {#if show}
            <th
              class="sortable {$sort == key
                ? $desc
                  ? 'desc'
                  : 'asc'
                : 'default'}"
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

  .icon {
    flex-direction: column;
    display: flex;
    justify-content: center;
    padding-left: 20px;
  }

  .search {
    position: relative;
    width: 250px;
    display: flex;
    float: right;
    margin-bottom: 10px;
    margin-right: 0;
    background-color: #e6ecf0;
    border-radius: 9999px;
  }
  .search:hover {
    box-shadow: 0 1px 6px 0 rgba(32, 33, 36, 0.28);
  }

  .search > input {
    background-color: transparent;
    padding: 10px;
    border: 0;
    width: 100%;
  }
  .search > input:focus {
    outline: none;
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

  th:nth-of-type(1),
  th:nth-of-type(2),
  th:nth-of-type(11) {
    width: 6rem;
  }
  th:nth-of-type(7),
  th:nth-of-type(8),
  th:nth-of-type(9) {
    width: 5rem;
  }
  @media (max-width: 1200px) {
    th:nth-of-type(3) {
      width: 10rem;
    }
  }
  th:nth-of-type(4),
  th:nth-of-type(5),
  th:nth-of-type(6) {
    width: 8rem;
  }
  th:nth-of-type(10) {
    width: 9rem;
  }
  th:nth-of-type(12) {
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
