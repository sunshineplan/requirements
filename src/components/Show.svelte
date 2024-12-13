<script lang="ts">
  import { stringify } from "csv-stringify/browser/esm/sync";
  import { onMount } from "svelte";
  import { requirements } from "../requirement.svelte";
  import Action from "./Action.svelte";
  import Search from "./Search.svelte";

  const add = () => {
    requirements.requirement = {} as Requirement;
    requirements.goto("add");
  };

  const view = async (e: MouseEvent, r: Requirement) => {
    if (window.getSelection()?.toString() !== "") return;
    if (!(e.target as HTMLElement).dataset["action"]) {
      requirements.requirement = r;
      requirements.goto("view");
    }
  };

  const download = () => {
    const link = document.createElement("a");
    const file = new Blob(
      [
        stringify(requirements.results, {
          bom: true,
          header: true,
          columns: requirements.fields.columns(true).map((key) => ({
            key,
            header: requirements.fields.name(key as keyof Requirement),
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

  const labels = (s: string) => {
    const res = s.split(",");
    if (res.length > 1) {
      return res.map((i) => i.charAt(0)).join(" | ");
    }
    return s;
  };

  onMount(() => {
    requirements.scroll(true);
    requirements.subscribe(true);
    return () => requirements.controller.abort();
  });
</script>

<svelte:head><title>{requirements.brand || "业务系统"}</title></svelte:head>

<header>
  <button class="btn btn-primary" onclick={add}>新增</button>
  <button class="btn btn-primary" onclick={download}>导出</button>
  <Search />
</header>
<div class="table-responsive">
  <table class="table table-hover table-sm">
    <thead>
      <tr>
        {#each requirements.fields.columns() as field (field)}
          {@const size = requirements.fields.size(field)}
          {#if size}
            <th
              class="sortable {requirements.search.sort == field
                ? requirements.search.desc
                  ? 'desc'
                  : 'asc'
                : 'default'}"
              class:auto={size == -1}
              style:width={size > 0 ? `${size}rem` : ""}
              onclick={() => {
                const before = requirements.search.sort;
                requirements.search.sort = field;
                if (before == requirements.search.sort)
                  requirements.search.desc = !requirements.search.desc;
                else requirements.search.desc = true;
              }}
            >
              {requirements.fields.name(field)}
            </th>
          {/if}
        {/each}
        <th></th>
      </tr>
    </thead>
    <tbody>
      {#each requirements.results as requirement (requirement.id)}
        <tr onclick={(e) => view(e, requirement)}>
          {#each requirements.fields.columns() as field (field)}
            <td
              title={requirements.fields.title(field) ? requirement[field] : ""}
            >
              {field == "label"
                ? labels(requirement[field])
                : requirement[field]}
            </td>
          {/each}
          <td style="vertical-align: middle">
            <Action {requirement} --icon="18px" --margin="2px" />
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<style>
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
