<script>
  // Переменные‑связки для полей
  let repoUrl = '';
  let overrideDockerfile = false;
  let dockerfileContent = '';
  let overrideCompose = false;
  let composeContent = '';
  let runCommand = '';
</script>

<style>
  /* ---------- Темная тема ---------- */
  :global(body) {
    background: #1e1e2e;
    color: #cdd6f4;
    font-family: system-ui, sans-serif;
    margin: 0;
    padding: 2rem;
  }
  h1 {
    font-size: 1.8rem;
    margin-bottom: 1rem;
    color: #f5c2e7;
  }
  form {
    background: #111;
    border-radius: 0.5rem;
    padding: 1.5rem;
    margin-top: 1rem;
    max-width: 640px;
  }
  label {
    display: flex;
    flex-direction: column;
    margin: 0.75rem 0;
  }
  input[type='text'],
  input[type='url'],
  textarea,
  button {
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
  button {
    background: #7287fd;
    color: #fff;
    cursor: pointer;
    margin-top: 1rem;
  }
  button[disabled] {
    background: #555;
    cursor: not-allowed;
  }
  .info {
    margin-top: 0.5rem;
    font-size: 0.9rem;
    color: #a6e3a1;
  }
</style>

<h1>💡 Создать Docker‑контейнер из GitHub‑репозитория</h1>

<form>
  <label>
    <strong>URL репозитория:</strong>
    <input type="url" bind:value={repoUrl} placeholder="https://github.com/owner/repo" />
  </label>

  <div class="info">(Проверка репозитория отсутствует)</div>

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