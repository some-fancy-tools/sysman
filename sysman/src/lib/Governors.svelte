<script>
    import { createEventDispatcher } from "svelte";

    export let cpuGovernor, cpuId, availableGovernors;
    export let freq;
    let dispatch = createEventDispatcher();

    function emitEvent(e) {
        console.log(e.target.id);
        // let els = e.target.id.split(":")
        let ev = {
            cpuId: e.target.id,
            cpuGovernor: cpuGovernor,
        };
        dispatch("cpuGovernorChange", ev);
    }
</script>

<div>
    <span class="label">{cpuId}</span>

    <datalist id={cpuId + "-availableGovernors"}>
        {#each availableGovernors as governor}
            <option>{governor}</option>
        {/each}
    </datalist>

    <input
        class="input-governor"
        type="text"
        list={cpuId + "-availableGovernors"}
        on:focus={() => {
            cpuGovernor = "";
        }}
        bind:value={cpuGovernor}
        on:change={emitEvent}
        id={cpuId}
    />
    <!-- <select on:change={emitEvent} id={cpuId} bind:value={cpuGovernor}>
        {#each availableGovernors as governor}
            <option value={governor}>-{governor}</option>
        {/each}
    </select> -->
    <span class="freq">{freq}</span>
</div>

<style>
    .label {
        display: inline-block;
        min-width: 3em;
    }
    .freq {
        padding-left: 2em;
    }

    .input-governor {
        max-width: 8em;
    }
</style>
