;(function hideAllResponseSectionsViaCSS(){
    const css = `
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
    const ui = SwaggerUIBundle({
        url: "openapi.yaml",
        dom_id: "#swagger-ui",
        deepLinking: true,
        docExpansion: "none",
        defaultModelsExpandDepth: -1,
        defaultModelExpandDepth: 0,
        displayRequestDuration: true,
        showExtensions: false,
        showCommonExtensions: false,
        displayOperationId: false,
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

    window.ui = ui;
};
