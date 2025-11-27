(function() {
    const editor = document.getElementById('editor');
    const highlight = document.getElementById('highlight');
    const status = document.getElementById('status');
    const saveBtn = document.getElementById('save-btn');
    const apiUrl = '/' + CHANNEL + '/' + DOCUMENT;

    function setStatus(msg, isError) {
        status.textContent = msg;
        status.className = isError ? 'error' : 'success';
        if (!isError) setTimeout(() => status.textContent = '', 3000);
    }

    // Syntax highlighting
    function highlightJSON(text) {
        if (!text) return '';
        return text
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/("(?:\\.|[^"\\])*")\s*:/g, '<span class="json-key">$1</span>:')
            .replace(/:(\s*)("(?:\\.|[^"\\])*")/g, ':$1<span class="json-string">$2</span>')
            .replace(/:\s*(-?\d+\.?\d*)/g, ': <span class="json-number">$1</span>')
            .replace(/:\s*(true|false)/g, ': <span class="json-boolean">$1</span>')
            .replace(/:\s*(null)/g, ': <span class="json-null">$1</span>');
    }

    function updateHighlight() {
        highlight.innerHTML = highlightJSON(editor.value) + '\n';
    }

    function syncScroll() {
        highlight.scrollTop = editor.scrollTop;
        highlight.scrollLeft = editor.scrollLeft;
    }

    // Load document on page load
    async function loadDocument() {
        try {
            const res = await fetch(apiUrl);
            if (res.ok) {
                const data = await res.json();
                editor.value = JSON.stringify(data, null, 2);
                updateHighlight();
            } else if (res.status === 404) {
                // New document - show empty textarea
                editor.value = '';
                updateHighlight();
            } else {
                const err = await res.json();
                setStatus('Error loading: ' + err.message, true);
            }
        } catch (e) {
            setStatus('Failed to load document', true);
        }
    }

    // Save document
    async function saveDocument() {
        const content = editor.value.trim();

        // Validate JSON
        try {
            if (content) JSON.parse(content);
        } catch (e) {
            setStatus('Invalid JSON: ' + e.message, true);
            return;
        }

        saveBtn.disabled = true;
        try {
            const res = await fetch(apiUrl, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: content || '{}'
            });
            if (res.ok) {
                setStatus('Saved successfully', false);
            } else {
                const err = await res.json();
                setStatus('Error: ' + err.message, true);
            }
        } catch (e) {
            setStatus('Failed to save', true);
        }
        saveBtn.disabled = false;
    }

    // Event listeners
    saveBtn.addEventListener('click', saveDocument);
    editor.addEventListener('input', updateHighlight);
    editor.addEventListener('scroll', syncScroll);

    // Ctrl+S / Cmd+S to save
    document.addEventListener('keydown', function(e) {
        if ((e.ctrlKey || e.metaKey) && e.key === 's') {
            e.preventDefault();
            saveDocument();
        }
    });

    // Initialize
    loadDocument();
})();
