package main

import "net/http"

// serveIndex отдает основной HTML-интерфейс
func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(indexHTML))
}

const indexHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>lazyPostman</title>
    <style>
        body {
            font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
            margin: 0;
            padding: 0;
            background: #0f172a;
            color: #e5e7eb;
        }
        header {
            padding: 12px 24px;
            background: #020617;
            border-bottom: 1px solid #1f2937;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }
        header h1 {
            font-size: 18px;
            margin: 0;
        }
        main {
            padding: 24px;
            max-width: 1200px;
            margin: 0 auto;
        }
        .layout {
            display: flex;
            gap: 16px;
            align-items: flex-start;
        }
        .sidebar {
            width: 260px;
            background: #020617;
            border-radius: 12px;
            border: 1px solid #1f2937;
            padding: 12px 12px 10px 12px;
            box-sizing: border-box;
        }
        .sidebar-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 8px;
        }
        .sidebar-title {
            font-size: 13px;
            font-weight: 500;
        }
        .sidebar-btn {
            border-radius: 999px;
            border: 1px solid #374151;
            background: #020617;
            color: #e5e7eb;
            font-size: 11px;
            padding: 4px 8px;
            cursor: pointer;
        }
        .sidebar-list {
            max-height: 520px;
            overflow-y: auto;
            padding-right: 4px;
        }
        .sidebar-item {
            border-radius: 8px;
            padding: 6px 8px;
            border: 1px solid transparent;
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 6px;
            cursor: pointer;
            font-size: 12px;
            color: #e5e7eb;
        }
        .sidebar-item:hover {
            border-color: #374151;
            background: #020617;
        }
        .sidebar-item-main {
            display: flex;
            flex-direction: column;
            gap: 2px;
            flex: 1;
            min-width: 0;
        }
        .sidebar-item-name {
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
        }
        .sidebar-item-meta {
            font-size: 10px;
            color: #9ca3af;
        }
        .sidebar-item-actions {
            display: flex;
            gap: 4px;
        }
        .sidebar-mini-btn {
            border-radius: 999px;
            border: 1px solid #374151;
            background: #020617;
            color: #9ca3af;
            font-size: 10px;
            padding: 2px 6px;
            cursor: pointer;
        }
        .sidebar-mini-btn:hover {
            border-color: #4b5563;
            color: #e5e7eb;
        }
        .content {
            flex: 1;
        }
        .card {
            background: #020617;
            border-radius: 12px;
            border: 1px solid #1f2937;
            padding: 16px 20px;
            margin-bottom: 16px;
        }
        .card h2 {
            font-size: 16px;
            margin: 0 0 8px 0;
        }
        .muted {
            color: #9ca3af;
            font-size: 13px;
        }
        .row {
            display: flex;
            gap: 8px;
            margin-bottom: 8px;
        }
        select, input, textarea, button {
            font-family: inherit;
            font-size: 13px;
        }
        select, input {
            border-radius: 8px;
            border: 1px solid #374151;
            background: #020617;
            color: #e5e7eb;
            padding: 6px 8px;
        }
        input, select {
            height: 32px;
        }
        textarea {
            width: 100%;
            min-height: 120px;
            border-radius: 8px;
            border: 1px solid #374151;
            background: #020617;
            color: #e5e7eb;
            padding: 8px 10px;
            resize: vertical;
            box-sizing: border-box;
        }
        button {
            border-radius: 999px;
            border: none;
            padding: 8px 16px;
            cursor: pointer;
            background: linear-gradient(to right, #22c55e, #3b82f6);
            color: #020617;
            font-weight: 500;
            display: inline-flex;
            align-items: center;
            gap: 6px;
        }
        button:disabled {
            opacity: 0.5;
            cursor: default;
        }
        pre {
            margin: 0;
            background: #020617;
            border-radius: 8px;
            padding: 10px 12px;
            border: 1px solid #111827;
            font-size: 12px;
            overflow: auto;
        }
        .pill {
            border-radius: 999px;
            border: 1px solid #374151;
            padding: 3px 8px;
            font-size: 11px;
            text-transform: uppercase;
            letter-spacing: 0.06em;
            color: #9ca3af;
        }
        .tabs {
            display: flex;
            gap: 4px;
            margin-bottom: 12px;
        }
        .tab {
            padding: 6px 10px;
            border-radius: 999px;
            border: 1px solid transparent;
            font-size: 12px;
            cursor: pointer;
            color: #9ca3af;
            background: transparent;
        }
        .tab.active {
            border-color: #4b5563;
            background: #020617;
            color: #e5e7eb;
        }
        .hidden {
            display: none;
        }
    </style>
</head>
<body>
<header>
    <h1>lazyPostman</h1>
    <span class="pill">MVP · HTTP + gRPC (soon)</span>
    </header>
<main>
    <div class="layout">
        <aside class="sidebar">
            <div class="sidebar-header">
                <span class="sidebar-title">Сохранённые запросы</span>
                <button id="save-current" class="sidebar-btn">Сохранить текущее</button>
            </div>
            <div id="saved-list" class="sidebar-list">
                <p class="muted">Пока нет сохранённых запросов.</p>
            </div>
        </aside>

        <div class="content card">
        <div class="tabs">
            <button class="tab active" data-tab="http">HTTP</button>
            <button class="tab" data-tab="grpc">gRPC</button>
        </div>

        <section id="http-panel">
            <h2>HTTP запрос</h2>
            <p class="muted">MVP: пока без сохранения коллекций. Введи URL и тело, отправим через backend.</p>
            <div class="row">
                <select id="http-method">
                    <option>GET</option>
                    <option>POST</option>
                    <option>PUT</option>
                    <option>DELETE</option>
                    <option>PATCH</option>
                </select>
                <input id="http-url" placeholder="https://api.example.com/resource" style="flex:1">
                <button id="http-send">Send</button>
            </div>
            <textarea id="http-body" placeholder="{&#10;  &quot;example&quot;: true&#10;}"></textarea>
            <p class="muted" style="margin-top:6px;">Ответ:</p>
            <pre id="http-response">// Здесь появится ответ HTTP</pre>
        </section>

        <section id="grpc-panel" class="hidden">
            <h2>gRPC запрос</h2>
            <p class="muted">Укажи адрес gRPC-сервера c включенной reflection, загрузим список методов и провалидируем JSON тела.</p>
            <div class="row">
                <input id="grpc-target" placeholder="localhost:50051" style="flex:1">
                <button id="grpc-load-methods">Загрузить методы</button>
            </div>
            <div class="row">
                <select id="grpc-method" style="flex:1">
                    <option value="">Сначала загрузи методы</option>
                </select>
            </div>
            <textarea id="grpc-body" placeholder="{&#10;  &quot;field&quot;: &quot;value&quot;&#10;}"></textarea>
            <div style="margin-top:8px; display:flex; gap:8px; align-items:center;">
                <button id="grpc-send" disabled>Send</button>
                <span class="muted" id="grpc-hint">Выбери метод и введи JSON-тело запроса.</span>
            </div>
            <p class="muted" style="margin-top:6px;">Ответ:</p>
            <pre id="grpc-response">// Здесь будет ответ gRPC</pre>
        </section>
        </div>
    </div>
</main>

<script>
    const tabs = document.querySelectorAll('.tab');
    const httpPanel = document.getElementById('http-panel');
    const grpcPanel = document.getElementById('grpc-panel');

    let activeTab = 'http';

    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            tabs.forEach(t => t.classList.remove('active'));
            tab.classList.add('active');
            const target = tab.dataset.tab;
            if (target === 'http') {
                httpPanel.classList.remove('hidden');
                grpcPanel.classList.add('hidden');
                activeTab = 'http';
            } else {
                httpPanel.classList.add('hidden');
                grpcPanel.classList.remove('hidden');
                activeTab = 'grpc';
            }
        });
    });

    const httpSendBtn = document.getElementById('http-send');
    const httpMethod = document.getElementById('http-method');
    const httpUrl = document.getElementById('http-url');
    const httpBody = document.getElementById('http-body');
    const httpResp = document.getElementById('http-response');

    httpSendBtn.addEventListener('click', async () => {
        const method = httpMethod.value;
        const url = httpUrl.value.trim();
        let body = httpBody.value.trim();

        if (!url) {
            httpResp.textContent = '// Укажи URL';
            return;
        }

        let parsedBody = null;
        if (body) {
            try {
                parsedBody = JSON.parse(body);
            } catch (e) {
                httpResp.textContent = '// Некорректный JSON в теле запроса: ' + e.message;
                return;
            }
        }

        httpSendBtn.disabled = true;
        httpResp.textContent = '// Отправляем запрос...';

        try {
            const res = await fetch('/api/http/request', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    method,
                    url,
                    body: parsedBody
                })
            });
            const text = await res.text();
            try {
                const json = JSON.parse(text);
                httpResp.textContent = JSON.stringify(json, null, 2);
            } catch {
                httpResp.textContent = text;
            }
        } catch (err) {
            httpResp.textContent = '// Ошибка отправки: ' + err.message;
        } finally {
            httpSendBtn.disabled = false;
        }
    });

    // gRPC
    const grpcTarget = document.getElementById('grpc-target');
    const grpcLoadBtn = document.getElementById('grpc-load-methods');
    const grpcMethod = document.getElementById('grpc-method');
    const grpcBody = document.getElementById('grpc-body');
    const grpcSendBtn = document.getElementById('grpc-send');
    const grpcResp = document.getElementById('grpc-response');
    const grpcHint = document.getElementById('grpc-hint');

    // Saved requests
    const saveCurrentBtn = document.getElementById('save-current');
    const savedListEl = document.getElementById('saved-list');
    let savedRequests = [];
    let currentRequestId = null;

    grpcLoadBtn.addEventListener('click', async () => {
        const target = grpcTarget.value.trim();
        if (!target) {
            grpcResp.textContent = '// Укажи адрес gRPC-сервера (target)';
            return;
        }
        grpcLoadBtn.disabled = true;
        grpcResp.textContent = '// Загружаем список сервисов и методов через reflection...';
        grpcMethod.innerHTML = '<option value=\"\">Загрузка...</option>';

        try {
            const res = await fetch('/api/grpc/services', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ target })
            });
            const data = await res.json();
            if (data.error) {
                grpcResp.textContent = '// Ошибка: ' + data.error;
                grpcMethod.innerHTML = '<option value=\"\">Ошибка при загрузке методов</option>';
                grpcSendBtn.disabled = true;
                return;
            }
            const options = [];
            if (data.services) {
                data.services.forEach(svc => {
                    if (svc.methods) {
                        svc.methods.forEach(m => {
                            const label = (svc.name || '') + '/' + (m.name || '');
                            options.push('<option value="' + m.fullMethod + '">' + label + '</option>');
                        });
                    }
                });
            }
            if (options.length === 0) {
                grpcMethod.innerHTML = '<option value=\"\">Методы не найдены (проверь reflection)</option>';
                grpcSendBtn.disabled = true;
                grpcResp.textContent = '// Методы не найдены. Убедись, что на сервере включен gRPC reflection.';
            } else {
                grpcMethod.innerHTML = '<option value=\"\">Выбери метод...</option>' + options.join('');
                grpcSendBtn.disabled = false;
                grpcResp.textContent = '// Методы загружены. Выбери метод и введи JSON-тело.';
            }
        } catch (e) {
            grpcResp.textContent = '// Ошибка при загрузке методов: ' + e.message;
            grpcMethod.innerHTML = '<option value=\"\">Ошибка при загрузке методов</option>';
            grpcSendBtn.disabled = true;
        } finally {
            grpcLoadBtn.disabled = false;
        }
    });

    grpcSendBtn.addEventListener('click', async () => {
        const target = grpcTarget.value.trim();
        const fullMethod = grpcMethod.value;
        let bodyText = grpcBody.value.trim();

        if (!target) {
            grpcResp.textContent = '// Укажи адрес gRPC-сервера (target)';
            return;
        }
        if (!fullMethod) {
            grpcResp.textContent = '// Выбери метод';
            return;
        }
        if (!bodyText) {
            grpcResp.textContent = '// Введи JSON-тело запроса';
            return;
        }

        let parsedBody;
        try {
            parsedBody = JSON.parse(bodyText);
        } catch (e) {
            grpcResp.textContent = '// Некорректный JSON (клиентская проверка): ' + e.message;
            return;
        }

        grpcSendBtn.disabled = true;
        grpcResp.textContent = '// Отправляем gRPC-запрос...';
        grpcHint.textContent = 'Выполняем вызов...';

        try {
            const res = await fetch('/api/grpc/invoke', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    target,
                    fullMethod,
                    jsonBody: parsedBody
                })
            });
            const text = await res.text();
            try {
                const json = JSON.parse(text);
                if (json && json.error) {
                    grpcResp.textContent = '// Ошибка (серверная проверка/вызов): ' + json.error;
                    grpcHint.textContent = 'Исправь JSON или параметры и попробуй снова.';
                } else {
                    grpcResp.textContent = JSON.stringify(json, null, 2);
                    grpcHint.textContent = 'Ответ успешно получен.';
                }
            } catch {
                grpcResp.textContent = text;
            }
        } catch (e) {
            grpcResp.textContent = '// Ошибка при вызове: ' + e.message;
            grpcHint.textContent = 'Проверь подключение и попробуй снова.';
        } finally {
            grpcSendBtn.disabled = false;
        }
    });

    async function loadSavedRequestsFromServer() {
        try {
            const res = await fetch('/api/requests');
            const data = await res.json();
            if (data && Array.isArray(data.items)) {
                savedRequests = data.items;
            } else {
                savedRequests = [];
            }
        } catch (e) {
            savedRequests = [];
        }
        renderSavedRequests();
    }

    function renderSavedRequests() {
        if (!savedRequests || savedRequests.length === 0) {
            savedListEl.innerHTML = '<p class="muted">Пока нет сохранённых запросов.</p>';
            return;
        }
        savedListEl.innerHTML = '';
        savedRequests.forEach(req => {
            const item = document.createElement('div');
            item.className = 'sidebar-item';
            item.dataset.id = req.id;

            const main = document.createElement('div');
            main.className = 'sidebar-item-main';

            const nameEl = document.createElement('div');
            nameEl.className = 'sidebar-item-name';
            const type = (req.type || '').toUpperCase();
            nameEl.textContent = req.name || '(без имени)';

            const metaEl = document.createElement('div');
            metaEl.className = 'sidebar-item-meta';
            metaEl.textContent = type === 'HTTP' ? 'HTTP' : 'gRPC';

            main.appendChild(nameEl);
            main.appendChild(metaEl);

            const actions = document.createElement('div');
            actions.className = 'sidebar-item-actions';

            const dupBtn = document.createElement('button');
            dupBtn.className = 'sidebar-mini-btn';
            dupBtn.textContent = 'Дубль';

            const delBtn = document.createElement('button');
            delBtn.className = 'sidebar-mini-btn';
            delBtn.textContent = 'Удалить';

            actions.appendChild(dupBtn);
            actions.appendChild(delBtn);

            item.appendChild(main);
            item.appendChild(actions);

            item.addEventListener('click', () => {
                loadSavedRequestIntoUI(req);
                highlightCurrentSidebarItem(req.id);
            });

            dupBtn.addEventListener('click', (ev) => {
                ev.stopPropagation();
                duplicateSavedRequest(req);
            });

            delBtn.addEventListener('click', (ev) => {
                ev.stopPropagation();
                deleteSavedRequest(req);
            });

            savedListEl.appendChild(item);
        });
    }

    function loadSavedRequestIntoUI(req) {
        if (!req || !req.type) {
            return;
        }
        currentRequestId = req.id || null;
        if (req.type === 'http' && req.http) {
            tabs.forEach(t => {
                if (t.dataset.tab === 'http') {
                    t.click();
                }
            });
            httpMethod.value = req.http.method || 'GET';
            httpUrl.value = req.http.url || '';
            httpBody.value = req.http.body || '';
            httpResp.textContent = '// Загружен сохранённый HTTP-запрос "' + (req.name || '') + '"';
        } else if (req.type === 'grpc' && req.grpc) {
            tabs.forEach(t => {
                if (t.dataset.tab === 'grpc') {
                    t.click();
                }
            });
            grpcTarget.value = req.grpc.target || '';

            // если метода ещё нет в select – добавим
            if (req.grpc.fullMethod) {
                let found = false;
                for (let i = 0; i < grpcMethod.options.length; i++) {
                    if (grpcMethod.options[i].value === req.grpc.fullMethod) {
                        found = true;
                        break;
                    }
                }
                if (!found) {
                    const opt = document.createElement('option');
                    opt.value = req.grpc.fullMethod;
                    opt.textContent = req.grpc.fullMethod;
                    grpcMethod.appendChild(opt);
                }
                grpcMethod.value = req.grpc.fullMethod;
            }

            grpcBody.value = req.grpc.body || '';
            grpcResp.textContent = '// Загружен сохранённый gRPC-запрос "' + (req.name || '') + '"';
            // позволяем сразу отправлять, без повторной загрузки методов через reflection
            grpcSendBtn.disabled = false;
        }
    }

    saveCurrentBtn.addEventListener('click', async () => {
        // если редактируем существующий запрос, подставим его имя в prompt
        let existing = null;
        if (currentRequestId) {
            existing = savedRequests.find(it => it.id === currentRequestId) || null;
        }

        let name = window.prompt('Имя для сохранённого запроса:', existing && existing.name ? existing.name : '');
        if (!name) {
            return;
        }
        name = name.trim();
        if (!name) {
            return;
        }

        let payload = { name: name };
        if (currentRequestId) {
            payload.id = currentRequestId;
        }
        if (activeTab === 'http') {
            payload.type = 'http';
            payload.http = {
                method: httpMethod.value || 'GET',
                url: httpUrl.value || '',
                body: httpBody.value || ''
            };
        } else {
            payload.type = 'grpc';
            payload.grpc = {
                target: grpcTarget.value || '',
                fullMethod: grpcMethod.value || '',
                body: grpcBody.value || ''
            };
        }

        try {
            const res = await fetch('/api/requests/save', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });
            const saved = await res.json();
            if (saved && saved.id) {
                // если это обновление – заменим элемент, иначе добавим
                const idx = savedRequests.findIndex(it => it.id === saved.id);
                if (idx >= 0) {
                    savedRequests[idx] = saved;
                } else {
                    savedRequests.push(saved);
                }
                currentRequestId = saved.id;
                renderSavedRequests();
            } else if (saved && saved.error) {
                alert('Ошибка сохранения: ' + saved.error);
            }
        } catch (e) {
            alert('Ошибка сохранения: ' + e.message);
        }
    });

    async function duplicateSavedRequest(req) {
        const baseName = req.name || '';
        const dupName = baseName ? baseName + ' (копия)' : 'Запрос (копия)';
        let payload = { name: dupName, type: req.type };
        if (req.type === 'http' && req.http) {
            payload.http = {
                method: req.http.method || 'GET',
                url: req.http.url || '',
                body: req.http.body || ''
            };
        } else if (req.type === 'grpc' && req.grpc) {
            payload.grpc = {
                target: req.grpc.target || '',
                fullMethod: req.grpc.fullMethod || '',
                body: req.grpc.body || ''
            };
        }
        try {
            const res = await fetch('/api/requests/save', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });
            const saved = await res.json();
            if (saved && saved.id) {
                savedRequests.push(saved);
                renderSavedRequests();
                currentRequestId = saved.id;
                highlightCurrentSidebarItem(saved.id);
            }
        } catch (e) {
            // ignore
        }
    }

    async function deleteSavedRequest(req) {
        if (!req || !req.id) {
            return;
        }
        if (!window.confirm('Удалить "' + (req.name || '') + '"?')) {
            return;
        }
        try {
            const res = await fetch('/api/requests/delete', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ id: req.id })
            });
            const data = await res.json();
            if (data && data.ok) {
                savedRequests = savedRequests.filter(it => it.id !== req.id);
                renderSavedRequests();
                if (currentRequestId === req.id) {
                    currentRequestId = null;
                }
            }
        } catch (e) {
            // ignore
        }
    }

    function highlightCurrentSidebarItem(id) {
        const items = savedListEl.querySelectorAll('.sidebar-item');
        items.forEach(el => {
            if (el.dataset.id === id) {
                el.style.borderColor = '#4b5563';
                el.style.backgroundColor = '#020617';
            } else {
                el.style.borderColor = 'transparent';
                el.style.backgroundColor = 'transparent';
            }
        });
    }

    // initial load
    loadSavedRequestsFromServer();
</script>
</body>
</html>`

