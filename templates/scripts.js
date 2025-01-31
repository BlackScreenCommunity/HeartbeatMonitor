
const observer = new MutationObserver((mutationsList) => {
    mutationsList.forEach((mutation) => {
        if (mutation.type === 'childList') {
            processWarnings();
        }
    });
});

const targetNode = document.body;
const config = {
    childList: true,
    subtree: true
};
observer.observe(targetNode, config);

function processWarnings() {
    $('div.warning').each(function () {
        $(this).parent().addClass('parent-warning');
    });
}


const eventSource = new EventSource("/events");

function renderList(data) {
    debugger;
    let html = "";

    if (typeof data === "object" && !Array.isArray(data)) { 
        let isWarning = data.hasOwnProperty("isWarning") ? data.isWarning : false;
        delete data["isWarning"];

        let pluginName = data.hasOwnProperty("plugin_name") ? data.plugin_name : "";
        if(pluginName) {
            data = data?.data;
        }

        if(pluginName) {
            html += `<div class='plugin-data'>`;
            html += `<div class='plugin-name'> ${pluginName} </div>`;
        }

        for (let key in data) {
            let widgetClass = isWarning ? "widget warning" : "widget";
            html += `<div class='${widgetClass}'>`;
            html += `<div class='widget-title'>${key}:</div>`;
            html += renderList(data[key]);
            html += `</div>`;
        }
        html += `</div>`;
    } else if (Array.isArray(data)) { 
        html += "<div class='data_array'>";
        data.forEach(value => {
            html += `<div class='data_array_element'>${renderList(value)}</div>`;
        });
        html += "</div>";
    } else { 
        html = `<div class='widget-data'>${data}</div>`;
    }
    return html;
}
const container = $("#data-container");


eventSource.onmessage = function(event) {
    const jsonData = JSON.parse(event.data);
    container.append(renderList(jsonData));
};

eventSource.onerror = function() {
    console.log("Ошибка соединения с сервером");
    eventSource.close();
};