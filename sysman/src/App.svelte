<script lang="ts">
  import Governors from "./lib/Governors.svelte";

  let governorsResponse = [];
  let frequencies = {};
  let message = "";

  function hertz(d) {
    if (d < 1000) {
      return d.toFixed(2) + " Hz";
    }
    if (d < 1000_000) {
      return (d / 1000).toFixed(2) + " kHz";
    }
    if (d < 1000_000_000) {
      return (d / 1000 / 1000).toFixed(2) + " MHz";
    }
    if (d < 1000_000_000_000) {
      return (d / 1000 / 1000 / 1000).toFixed(2) + " GHz";
    }
  }

  async function getFrequencies() {
    message = "Fetching Frequencies...";
    let fr = {};
    await fetch("/api/v1/sysman/frequencies")
      .then((r) => r.json())
      .then((d) => {
        Object.keys(d).forEach((cpuId) => {
          fr[cpuId] = hertz(d[cpuId][0]);
        });
        frequencies = fr;
      });
    message = "";
    // clearMessage();
  }
  async function getGovernors() {
    message = "Reloading...";
    let gr = [];
    await fetch("/api/v1/sysman/governors")
      .then((r) => r.json())
      .then((d) => {
        Object.keys(d["availableGovernors"]).forEach((cpuId) => {
          gr.push({
            cpuId: cpuId,
            availableGovernors: d["availableGovernors"][cpuId],
            selectedGovernor: d["selectedGovernor"][cpuId],
          });
        });
        governorsResponse = gr;
        message = "";
      });
  }

  async function setGovernors(e) {
    // console.log(e);
    if (e.type == "cpuGovernorChange") {
      for (let i = 0; i < governorsResponse.length; i++) {
        if (governorsResponse[i].cpuId == e.detail.cpuId) {
          if (e.detail.cpuGovernor != "") {
            governorsResponse[i].selectedGovernor = e.detail.cpuGovernor;
          }
        }
      }
      return;
    }

    console.log(governorsResponse);
    message = "Updating...";
    await fetch("/api/v1/sysman/governors", { method: "POST", body: JSON.stringify(governorsResponse) }).then((d) => {
      console.log(d);
      message = "Updated";
      clearMessage();
    });
  }

  function clearMessage() {
    setTimeout(() => {
      message = "";
    }, 2000);
  }

  getGovernors();
  getFrequencies();

  setInterval(() => {
    getFrequencies();
  }, 10000);
</script>

<main>
  <h1>SysManager</h1>

  <div class="controls">
    <button
      on:click={() => {
        getGovernors();
        getFrequencies();
      }}>Reload</button
    >
    <button on:click={setGovernors}>Submit</button>

  <div class="message">
    &nbsp; {message}
  </div>
  </div>

  {#each governorsResponse as governorData}
    <div class="cpu">
      <Governors
        freq={frequencies[governorData.cpuId]}
        cpuId={governorData.cpuId}
        availableGovernors={governorData.availableGovernors}
        bind:cpuGovernor={governorData.selectedGovernor}
        on:cpuGovernorChange={setGovernors}
      />
    </div>
  {:else}
    <div>No Data Found</div>
  {/each}

</main>

<style>
  .controls {
    padding: 2em;
  }
  .message {
    padding-top: 1em;
  }
</style>
