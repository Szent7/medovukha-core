<script lang="ts">
	import FrameElement from '@templates/frameElement.svelte';
	import Sidebar from '@templates/sidebar.svelte';
	import Frame from '@templates/frame.svelte';
	import logo from '@assets/logo_small.svg';

    // Переменные‑связки для полей
  let repoUrl = '';
  let overrideDockerfile = false;
  let dockerfileContent = '';
  let overrideCompose = false;
  let composeContent = '';
  let runCommand = '';
</script>

<svelte:head>
	<title>Create from git</title>
	<link rel="icon" type="image/svg+xml" href={logo} />
</svelte:head>

<Sidebar />
<Frame {content} />

{#snippet content()}
	<FrameElement
		--display="flex"
		--justify-content="space-between"
		--align-items="center"
		content={header}
	/>
	<FrameElement content={createFromGitForm} />
{/snippet}

{#snippet header()}
	<img class="logoHeader" src={logo} alt="logo" />
	<h2>Create from git</h2>
	<div class="releaseLabel">Alpha v0.0.1</div>
{/snippet}

{#snippet createFromGitForm()}
<form>
  <label>
    <strong>URL репозитория:</strong>
    <input type="url" bind:value={repoUrl} placeholder="https://github.com/owner/repo" />
  </label>

  <label>
    <input type="checkbox" bind:checked={overrideDockerfile} />
    Переопределить Dockerfile
  </label>
  {#if overrideDockerfile}
    <textarea
      rows="6"
      bind:value={dockerfileContent}
      placeholder="Новый Dockerfile…"
    ></textarea>
  {/if}

  <label>
    <input type="checkbox" bind:checked={overrideCompose} />
    Переопределить docker‑compose.yml
  </label>
  {#if overrideCompose}
    <textarea
      rows="6"
      bind:value={composeContent}
      placeholder="Новый docker‑compose.yml…"
    ></textarea>
  {/if}

  <label>
    <strong>Или команда docker run:</strong>
    <input type="text" bind:value={runCommand} placeholder="docker run -d …" />
  </label>

  <button type="submit" disabled={!repoUrl}>
    Создать контейнер
  </button>
</form>
{/snippet}

<style>
	h2 {
		font-size: 24pt;
	}

	.releaseLabel {
		background-color: #ffffff;
		color: #1c1c1c;
		padding: 5px;
		border-radius: 5px;
		font-size: 14pt;
		font-weight: bold;
	}

	.logoHeader {
		display: block;
		max-height: 100px;
	}

  label {
    display: flex;
    flex-direction: column;
    margin: 0.75rem 0;
  }

  input[type='text'],
  input[type='url'],
  textarea {
    background: #292c3c;
    color: #cdd6f4;
    border: 1px solid #444;
    padding: 0.5rem;
    border-radius: 0.25rem;
    width: 100%;
    font-size: 1rem;
  }

  input[type='checkbox'] {
    width: auto;
    margin-right: 0.5rem;
  }
</style>
