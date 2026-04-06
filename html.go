package main

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>rmail</title>
<style>
  @import url('https://fonts.googleapis.com/css2?family=DM+Mono:wght@300;400;500&family=DM+Sans:wght@300;400;500&display=swap');

  :root {
    --cream: #f7f5f0;
    --cream2: #eeebe4;
    --gray-light: #e8e5de;
    --gray-mid: #c8c4bb;
    --gray-dark: #6b6760;
    --text: #2a2825;
    --text-soft: #7a7772;
    --border: #dedad2;
    --border-focus: #a8a49c;
    --accent: #3d3b37;
    --white: #fafaf8;
    --success: #4a7c59;
    --error: #8b3a3a;
    --warn: #7a5c1e;
    --shadow: 0 1px 3px rgba(42,40,37,0.08);
    --shadow-md: 0 4px 16px rgba(42,40,37,0.12);
  }

  * { box-sizing: border-box; margin: 0; padding: 0; }

  body {
    font-family: 'DM Sans', sans-serif;
    background: var(--cream);
    color: var(--text);
    min-height: 100vh;
    font-size: 14px;
    line-height: 1.5;
  }

  .app { display: flex; flex-direction: column; min-height: 100vh; }

  header {
    padding: 18px 32px;
    border-bottom: 1px solid var(--border);
    background: var(--white);
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .logo {
    font-family: 'DM Mono', monospace;
    font-weight: 500;
    font-size: 15px;
    letter-spacing: 0.05em;
    color: var(--text);
  }

  .logo span { color: var(--gray-dark); }

  main { flex: 1; padding: 32px; max-width: 860px; margin: 0 auto; width: 100%; }

  .token-screen {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding-top: 80px;
  }

  .token-card {
    background: var(--white);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 36px 40px;
    width: 100%;
    max-width: 440px;
    box-shadow: var(--shadow);
  }

  .token-title {
    font-size: 15px;
    font-weight: 500;
    margin-bottom: 6px;
    color: var(--text);
  }

  .token-sub {
    font-size: 12px;
    color: var(--text-soft);
    margin-bottom: 24px;
    font-family: 'DM Mono', monospace;
  }

  label {
    display: block;
    font-size: 11px;
    font-weight: 500;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--text-soft);
    margin-bottom: 6px;
  }

  input[type="text"], input[type="password"], input[type="email"], select {
    width: 100%;
    background: var(--cream);
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 9px 12px;
    font-family: 'DM Mono', monospace;
    font-size: 12px;
    color: var(--text);
    outline: none;
    transition: border-color 0.15s;
  }

  input:focus, select:focus {
    border-color: var(--border-focus);
    background: var(--white);
  }

  select {
    appearance: none;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='10' height='6' viewBox='0 0 10 6'%3E%3Cpath d='M1 1l4 4 4-4' stroke='%23a8a49c' stroke-width='1.5' fill='none' stroke-linecap='round'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 10px center;
    cursor: pointer;
  }

  .input-group { margin-bottom: 16px; }

  .save-token-row {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 16px;
    font-size: 12px;
    color: var(--text-soft);
  }

  .save-token-row input[type="checkbox"] {
    width: auto;
    margin: 0;
  }

  .btn {
    font-family: 'DM Sans', sans-serif;
    font-size: 12px;
    font-weight: 500;
    letter-spacing: 0.04em;
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 8px 18px;
    cursor: pointer;
    transition: all 0.15s;
    background: var(--white);
    color: var(--text);
    display: inline-flex;
    align-items: center;
    gap: 6px;
  }

  .btn:hover { background: var(--cream2); border-color: var(--gray-mid); }
  .btn:active { transform: translateY(1px); }

  .btn-primary {
    background: var(--accent);
    color: var(--cream);
    border-color: var(--accent);
  }
  .btn-primary:hover { background: #2a2825; border-color: #2a2825; }

  .btn-danger {
    color: var(--error);
    border-color: transparent;
  }
  .btn-danger:hover { background: #f5eded; border-color: #e8d4d4; }

  .btn-sm { padding: 5px 12px; font-size: 11px; }
  .btn-full { width: 100%; justify-content: center; }

  .zone-bar {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 28px;
    padding-bottom: 20px;
    border-bottom: 1px solid var(--border);
  }

  .zone-bar label { margin: 0; white-space: nowrap; }
  .zone-bar select { max-width: 260px; }

  .routing-badge {
    font-family: 'DM Mono', monospace;
    font-size: 10px;
    padding: 3px 8px;
    border-radius: 3px;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    white-space: nowrap;
  }

  .badge-on { background: #e8f2eb; color: var(--success); border: 1px solid #c8e0ce; }
  .badge-off { background: var(--cream2); color: var(--text-soft); border: 1px solid var(--border); }

  .warn-banner {
    background: #fdf6e3;
    border: 1px solid #e8d99a;
    border-radius: 4px;
    padding: 10px 14px;
    font-size: 12px;
    color: var(--warn);
    font-family: 'DM Mono', monospace;
    margin-bottom: 20px;
    line-height: 1.6;
  }

  .section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
  }

  .section-title {
    font-size: 11px;
    font-weight: 500;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--text-soft);
  }

  .rules-table {
    width: 100%;
    border: 1px solid var(--border);
    border-radius: 6px;
    overflow: hidden;
    background: var(--white);
  }

  .rules-head {
    display: grid;
    grid-template-columns: 1fr 140px 120px 80px;
    padding: 10px 16px;
    background: var(--cream2);
    border-bottom: 1px solid var(--border);
    font-size: 10px;
    font-weight: 500;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--text-soft);
  }

  .rule-row {
    display: grid;
    grid-template-columns: 1fr 140px 120px 80px;
    padding: 13px 16px;
    border-bottom: 1px solid var(--border);
    align-items: center;
    transition: background 0.1s;
  }

  .rule-row:last-child { border-bottom: none; }
  .rule-row:hover { background: var(--cream); }

  .rule-email {
    font-family: 'DM Mono', monospace;
    font-size: 12px;
    color: var(--text);
  }

  .rule-action-type {
    font-size: 11px;
    color: var(--text-soft);
    text-transform: capitalize;
  }

  .rule-dest {
    font-family: 'DM Mono', monospace;
    font-size: 11px;
    color: var(--text-soft);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .rule-actions { display: flex; gap: 4px; justify-content: flex-end; }

  .empty-state {
    padding: 48px 24px;
    text-align: center;
    color: var(--text-soft);
    font-size: 12px;
    font-family: 'DM Mono', monospace;
  }

  .overlay {
    position: fixed; inset: 0;
    background: rgba(42,40,37,0.35);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 100;
    backdrop-filter: blur(2px);
  }

  .overlay.hidden { display: none; }

  .modal {
    background: var(--white);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 28px 32px;
    width: 100%;
    max-width: 480px;
    box-shadow: var(--shadow-md);
  }

  .modal-title {
    font-size: 14px;
    font-weight: 500;
    margin-bottom: 4px;
  }

  .modal-sub {
    font-size: 11px;
    color: var(--text-soft);
    margin-bottom: 24px;
    font-family: 'DM Mono', monospace;
  }

  .modal-footer {
    display: flex;
    gap: 8px;
    justify-content: flex-end;
    margin-top: 24px;
    padding-top: 16px;
    border-top: 1px solid var(--border);
  }

  .email-split {
    display: flex;
    align-items: center;
    gap: 0;
    border: 1px solid var(--border);
    border-radius: 4px;
    overflow: hidden;
    background: var(--cream);
    transition: border-color 0.15s;
  }

  .email-split:focus-within {
    border-color: var(--border-focus);
    background: var(--white);
  }

  .email-split input {
    border: none;
    background: transparent;
    border-radius: 0;
    flex: 1;
  }

  .email-split input:focus { background: transparent; }

  .at-sep {
    padding: 9px 6px;
    font-family: 'DM Mono', monospace;
    font-size: 12px;
    color: var(--text-soft);
    background: transparent;
    flex-shrink: 0;
  }

  .domain-label {
    padding: 9px 12px 9px 4px;
    font-family: 'DM Mono', monospace;
    font-size: 12px;
    color: var(--text-soft);
    background: transparent;
    flex-shrink: 0;
  }

  #toast {
    position: fixed;
    bottom: 28px;
    right: 28px;
    background: var(--accent);
    color: var(--cream);
    padding: 11px 18px;
    border-radius: 4px;
    font-size: 12px;
    font-family: 'DM Mono', monospace;
    letter-spacing: 0.03em;
    opacity: 0;
    transform: translateY(8px);
    transition: all 0.2s;
    pointer-events: none;
    z-index: 200;
    max-width: 300px;
  }

  #toast.show { opacity: 1; transform: translateY(0); }
  #toast.error { background: var(--error); }

  .spinner {
    display: inline-block;
    width: 12px; height: 12px;
    border: 1.5px solid var(--gray-mid);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
  }

  @keyframes spin { to { transform: rotate(360deg); } }

  .hidden { display: none !important; }

  .toggle-wrap { display: flex; align-items: center; gap: 10px; }
  .toggle {
    position: relative;
    width: 34px; height: 18px;
    cursor: pointer;
  }
  .toggle input { opacity: 0; width: 0; height: 0; }
  .toggle-track {
    position: absolute; inset: 0;
    background: var(--gray-light);
    border-radius: 9px;
    border: 1px solid var(--border);
    transition: background 0.2s;
  }
  .toggle input:checked + .toggle-track { background: var(--accent); border-color: var(--accent); }
  .toggle-thumb {
    position: absolute;
    top: 2px; left: 2px;
    width: 12px; height: 12px;
    background: var(--white);
    border-radius: 50%;
    transition: transform 0.2s;
    pointer-events: none;
  }
  .toggle input:checked ~ .toggle-thumb { transform: translateX(16px); }
</style>
</head>
<body>
<div class="app">
  <header>
    <div class="logo">r<span>mail</span></div>
  </header>

  <main>
    <!-- Token screen -->
    <div id="tokenScreen" class="token-screen">
      <div class="token-card">
        <div class="token-title">Cloudflare API Token</div>
        <div class="token-sub">cloudflare.com &gt; profile &gt; api tokens</div>
        <div class="input-group">
          <label>Token</label>
          <input type="password" id="tokenInput" placeholder="paste token here" autocomplete="off">
        </div>
        <div class="save-token-row">
          <input type="checkbox" id="saveTokenCheck" checked>
          <label for="saveTokenCheck" style="text-transform:none;letter-spacing:0;font-size:12px;margin:0;color:var(--text-soft)">Remember token (saved to config.yml)</label>
        </div>
        <button class="btn btn-primary btn-full" id="checkBtn" onclick="checkToken()">
          <span id="checkLabel">Verify &amp; Continue</span>
          <span id="checkSpinner" class="spinner hidden"></span>
        </button>
      </div>
    </div>

    <!-- Dashboard screen -->
    <div id="dashboard" class="hidden">
      <div class="zone-bar">
        <label>Domain</label>
        <select id="zoneSelect" onchange="selectZone()"></select>
        <span id="routingBadge" class="routing-badge badge-off">routing off</span>
        <div style="flex:1"></div>
        <button class="btn btn-sm" onclick="logout()">change token</button>
      </div>

      <div id="warnBanner" class="warn-banner hidden">
        Token is missing <strong>Email Routing:Edit</strong> permission. Go to Cloudflare &gt; My Profile &gt; API Tokens &gt; edit token &gt; add <strong>Zone / Email Routing / Edit</strong> for All zones.
      </div>

      <div class="section-header">
        <div class="section-title">Custom Addresses</div>
        <button class="btn btn-sm btn-primary" onclick="openCreate()">New address</button>
      </div>

      <div class="rules-table" id="rulesTable">
        <div class="rules-head">
          <div>Address</div>
          <div>Action</div>
          <div>Destination</div>
          <div></div>
        </div>
        <div id="rulesBody"></div>
      </div>
    </div>
  </main>
</div>

<!-- Create/Edit Modal -->
<div class="overlay hidden" id="ruleModal">
  <div class="modal">
    <div class="modal-title" id="modalTitle">New custom address</div>
    <div class="modal-sub" id="modalSub">email routing / custom address</div>

    <div class="input-group">
      <label>Rule name</label>
      <input type="text" id="rName" placeholder="e.g. work alias">
    </div>

    <div class="input-group">
      <label>Custom address</label>
      <div class="email-split">
        <input type="text" id="rLocal" placeholder="username">
        <span class="at-sep">@</span>
        <span class="domain-label" id="rDomainLabel">domain.com</span>
      </div>
    </div>

    <div class="input-group">
      <label>Action</label>
      <select id="rAction" onchange="toggleDestField()">
        <option value="forward">Send to an email</option>
        <option value="worker">Send to a Worker</option>
        <option value="drop">Drop</option>
      </select>
    </div>

    <div class="input-group" id="destGroup">
      <label>Destination</label>
      <input type="text" id="rDest" placeholder="destination@email.com">
    </div>

    <div class="input-group" id="enabledGroup">
      <label>Enabled</label>
      <label class="toggle">
        <input type="checkbox" id="rEnabled" checked>
        <span class="toggle-track"></span>
        <span class="toggle-thumb"></span>
      </label>
    </div>

    <div class="modal-footer">
      <button class="btn" onclick="closeModal()">Cancel</button>
      <button class="btn btn-primary" id="saveBtn" onclick="saveRule()">
        <span id="saveLabel">Create</span>
        <span id="saveSpinner" class="spinner hidden"></span>
      </button>
    </div>
  </div>
</div>

<div id="toast"></div>

<script>
let state = {
  token: '',
  zones: [],
  currentZone: null,
  rules: [],
  editingTag: null
};

window.addEventListener('DOMContentLoaded', async () => {
  try {
    const res = await fetch('/api/config');
    const data = await res.json();
    if (data.token) {
      document.getElementById('tokenInput').value = data.token;
      await checkToken(true);
    }
  } catch(e) {}
});
async function checkToken(silent = false) {
  const token = document.getElementById('tokenInput').value.trim();
  if (!token) return;

  setLoading('check', true);
  try {
    const res = await fetch('/api/check?token=' + encodeURIComponent(token));
    const data = await res.json();
    if (!data.ok) {
      if (!silent) toast(data.error, true);
      return;
    }
    if (!data.zones || data.zones.length === 0) {
      if (!silent) toast('No zones found', true);
      return;
    }

    state.token = token;
    state.zones = data.zones;

    const shouldSave = silent || document.getElementById('saveTokenCheck').checked;
    if (shouldSave) {
      fetch('/api/config/save', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token })
      });
    }

    showDashboard(data.warn_routing_perm);
    if (!silent) toast('Token verified');
  } catch(e) {
    if (!silent) toast('Connection error', true);
  } finally {
    setLoading('check', false);
  }
}

function logout() {
  state = { token: '', zones: [], currentZone: null, rules: [], editingTag: null };
  document.getElementById('tokenInput').value = '';
  document.getElementById('tokenScreen').classList.remove('hidden');
  document.getElementById('dashboard').classList.add('hidden');
  fetch('/api/config/save', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ token: '' })
  });
}

function showDashboard(warnRoutingPerm) {
  document.getElementById('tokenScreen').classList.add('hidden');
  document.getElementById('dashboard').classList.remove('hidden');

  const banner = document.getElementById('warnBanner');
  if (warnRoutingPerm) {
    banner.classList.remove('hidden');
  } else {
    banner.classList.add('hidden');
  }

  const sel = document.getElementById('zoneSelect');
  sel.innerHTML = '';
  state.zones.forEach(z => {
    const opt = document.createElement('option');
    opt.value = z.id;
    opt.textContent = z.name;
    opt.dataset.routing = z.routing;
    sel.appendChild(opt);
  });

  selectZone();
}

function selectZone() {
  const sel = document.getElementById('zoneSelect');
  const opt = sel.options[sel.selectedIndex];
  state.currentZone = state.zones.find(z => z.id === sel.value);

  const badge = document.getElementById('routingBadge');
  if (opt.dataset.routing === 'true') {
    badge.textContent = 'routing on';
    badge.className = 'routing-badge badge-on';
  } else {
    badge.textContent = 'routing off';
    badge.className = 'routing-badge badge-off';
  }

  loadRules();
}
async function loadRules() {
  if (!state.currentZone) return;
  document.getElementById('rulesBody').innerHTML = '<div class="empty-state"><span class="spinner"></span></div>';

  try {
    const res = await fetch('/api/rules?token=' + encodeURIComponent(state.token) + '&zone=' + state.currentZone.id);
    const data = await res.json();
    if (!data.ok) { toast(data.error, true); return; }
    state.rules = data.rules || [];
    renderRules();
  } catch(e) { toast('Failed to load rules', true); }
}

function renderRules() {
  const body = document.getElementById('rulesBody');
  const rules = state.rules.filter(r => {
    return r.matchers && r.matchers.some(m => m.type === 'literal');
  });

  if (rules.length === 0) {
    body.innerHTML = '<div class="empty-state">no custom addresses yet</div>';
    return;
  }

  body.innerHTML = rules.map(r => {
    const email = r.matchers.find(m => m.field === 'to')?.value || '';
    const action = r.actions[0] || {};
    const actionLabel = { forward: 'email', worker: 'worker', drop: 'drop' }[action.type] || action.type;
    const dest = action.value || '';
    return '<div class="rule-row">' +
      '<div class="rule-email">' + email + '</div>' +
      '<div class="rule-action-type">' + actionLabel + '</div>' +
      '<div class="rule-dest" title="' + dest + '">' + (dest || '\u2014') + '</div>' +
      '<div class="rule-actions">' +
        '<button class="btn btn-sm" onclick="openEdit(\'' + r.tag + '\')">edit</button>' +
        '<button class="btn btn-sm btn-danger" onclick="deleteRule(\'' + r.tag + '\')">del</button>' +
      '</div>' +
    '</div>';
  }).join('');
}
function openCreate() {
  state.editingTag = null;
  document.getElementById('modalTitle').textContent = 'New custom address';
  document.getElementById('modalSub').textContent = 'email routing / custom address';
  document.getElementById('saveLabel').textContent = 'Create';
  document.getElementById('rName').value = '';
  document.getElementById('rLocal').value = '';
  document.getElementById('rDest').value = '';
  document.getElementById('rAction').value = 'forward';
  document.getElementById('rEnabled').checked = true;
  document.getElementById('rDomainLabel').textContent = state.currentZone?.name || '';
  document.getElementById('enabledGroup').classList.add('hidden');
  toggleDestField();
  document.getElementById('ruleModal').classList.remove('hidden');
  setTimeout(() => document.getElementById('rLocal').focus(), 50);
}

function openEdit(tag) {
  const rule = state.rules.find(r => r.tag === tag);
  if (!rule) return;

  state.editingTag = tag;
  const email = rule.matchers.find(m => m.field === 'to')?.value || '';
  const [local] = email.split('@');
  const action = rule.actions[0] || {};

  document.getElementById('modalTitle').textContent = 'Edit custom address';
  document.getElementById('modalSub').textContent = 'email routing / custom address';
  document.getElementById('saveLabel').textContent = 'Save';
  document.getElementById('rName').value = rule.name || '';
  document.getElementById('rLocal').value = local;
  document.getElementById('rDest').value = action.value || '';
  document.getElementById('rAction').value = action.type || 'forward';
  document.getElementById('rEnabled').checked = rule.enabled;
  document.getElementById('rDomainLabel').textContent = state.currentZone?.name || '';
  document.getElementById('enabledGroup').classList.remove('hidden');
  toggleDestField();
  document.getElementById('ruleModal').classList.remove('hidden');
}

function closeModal() {
  document.getElementById('ruleModal').classList.add('hidden');
  state.editingTag = null;
}

function toggleDestField() {
  const action = document.getElementById('rAction').value;
  document.getElementById('destGroup').classList.toggle('hidden', action === 'drop');
}

async function saveRule() {
  const name = document.getElementById('rName').value.trim();
  const local = document.getElementById('rLocal').value.trim();
  const actionType = document.getElementById('rAction').value;
  const dest = document.getElementById('rDest').value.trim();
  const enabled = document.getElementById('rEnabled').checked;
  const domain = state.currentZone?.name || '';

  if (!local) { toast('Enter address username', true); return; }
  if (actionType !== 'drop' && !dest) { toast('Enter destination', true); return; }

  setLoading('save', true);

  try {
    let res, data;
    if (state.editingTag) {
      res = await fetch('/api/rules/update', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          token: state.token,
          zone_id: state.currentZone.id,
          tag: state.editingTag,
          name: name || local + ' rule',
          local_part: local,
          domain,
          action_type: actionType,
          destination: dest,
          enabled
        })
      });
    } else {
      res = await fetch('/api/rules/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          token: state.token,
          zone_id: state.currentZone.id,
          name: name || local + ' rule',
          local_part: local,
          domain,
          action_type: actionType,
          destination: dest
        })
      });
    }

    data = await res.json();
    if (!data.ok) { toast(data.error, true); return; }

    toast(state.editingTag ? 'Rule updated' : 'Rule created');
    closeModal();
    loadRules();
  } catch(e) { toast('Request failed', true); }
  finally { setLoading('save', false); }
}

async function deleteRule(tag) {
  if (!confirm('Delete this rule?')) return;

  try {
    const res = await fetch('/api/rules/delete', {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ token: state.token, zone_id: state.currentZone.id, tag })
    });
    const data = await res.json();
    if (!data.ok) { toast(data.error, true); return; }
    toast('Rule deleted');
    loadRules();
  } catch(e) { toast('Failed to delete', true); }
}
function setLoading(ctx, on) {
  if (ctx === 'check') {
    document.getElementById('checkLabel').textContent = on ? 'Verifying...' : 'Verify & Continue';
    document.getElementById('checkSpinner').classList.toggle('hidden', !on);
    document.getElementById('checkBtn').disabled = on;
  } else if (ctx === 'save') {
    document.getElementById('saveLabel').textContent = on ? 'Saving...' : (state.editingTag ? 'Save' : 'Create');
    document.getElementById('saveSpinner').classList.toggle('hidden', !on);
    document.getElementById('saveBtn').disabled = on;
  }
}

let toastTimer;
function toast(msg, isError = false) {
  const el = document.getElementById('toast');
  el.textContent = msg;
  el.className = 'show' + (isError ? ' error' : '');
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => { el.className = ''; }, 3500);
}

document.getElementById('tokenInput').addEventListener('keydown', e => {
  if (e.key === 'Enter') checkToken();
});

document.getElementById('ruleModal').addEventListener('click', e => {
  if (e.target === document.getElementById('ruleModal')) closeModal();
});
</script>
</body>
</html>
`
