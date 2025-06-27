;(function hideAllResponseSectionsViaCSS(){
    const css = `
    .schemes-server-container .servers table {
      display: none !important;
    }
    .responses-wrapper .responses-table:not(.live-responses-table) {
      display: none !important;
    }
    .responses-wrapper .curl-command,.request-url {
      display: none !important;
    }
  `;
    const style = document.createElement('style');
    style.type = 'text/css';
    style.appendChild(document.createTextNode(css));
    document.head.appendChild(style);
})();
window.onload = function() {
    const forbidden = ['Server variables', 'Responses'];
    const ui = SwaggerUIBundle({
        url: "openapi.yaml",
        dom_id: "#swagger-ui",
        deepLinking: true,
        docExpansion: "none",
        defaultModelsExpandDepth: 1,
        defaultModelExpandDepth: 1,
        displayRequestDuration: true,
        showExtensions: false,
        showCommonExtensions: false,
        displayOperationId: true,
        filter: false,
        tryItOutEnabled: true,
        presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIStandalonePreset
        ],
        plugins: [
            SwaggerUIBundle.plugins.DownloadUrl
        ],
    });

    function purgeH4s(root = document) {
        root.querySelectorAll('h4').forEach(h4 => {
            if (forbidden.includes(h4.textContent.trim())) {
                h4.remove();
            }
        });
    }

    purgeH4s();

    const container = document.getElementById('swagger-ui') || document.body;
    const mo = new MutationObserver(mutations => {
        for (let m of mutations) {
            m.addedNodes.forEach(node => {
                if (node.nodeType === 1) purgeH4s(node);
            });
        }
    });
    mo.observe(container, { childList: true, subtree: true });

    window.ui = ui;
};
