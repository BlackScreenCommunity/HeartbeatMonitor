"use strict";
const observer = new MutationObserver((mutationsList) => {
    mutationsList.forEach((mutation) => {
        if (mutation.type === "childList") {
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
    document.querySelectorAll("div.warning").forEach((el) => {
        var _a;
        (_a = el.parentElement) === null || _a === void 0 ? void 0 : _a.classList.add("parent-warning");
    });
}
const eventSource = new EventSource("/events");
function sortAgentsData(pluginsData) {
    if (pluginsData) {
        let sortable = Object.entries(pluginsData);
        sortable.sort((a, b) => JSON.stringify(b).length - JSON.stringify(a).length);
        let sortedData = {};
        sortable.forEach(([key, value]) => {
            sortedData[key] = value;
        });
        return sortedData;
    }
    return {};
}
function renderList(data, levelClass = "") {
    let html = "";
    if (typeof data === "object" && !Array.isArray(data)) {
        let duration = 0;
        let isWarning = data.hasOwnProperty("isWarning") ? data.isWarning : false;
        delete data["isWarning"];
        let agentName = data.hasOwnProperty("agent_name") ? data.agent_name : "";
        if (agentName) {
            duration = data === null || data === void 0 ? void 0 : data.duration;
            data = data === null || data === void 0 ? void 0 : data.data;
        }
        if (agentName) {
            html += `<div class='agent-name'> ${agentName}`;
            html += `<div class='timer agent-duration'> Loading time ${duration} seconds </div>`;
            html += `</div>`;
            html += `<div class='agent-data'>`;
        }
        let pluginName = data.hasOwnProperty("plugin_name") ? data.plugin_name : "";
        if (pluginName) {
            data = data === null || data === void 0 ? void 0 : data.data;
        }
        if (pluginName) {
            html += `<div class='plugin-data'>`;
            html += `<div class='plugin-name'> ${pluginName} </div>`;
        }
        else {
            for (let key in data) {
                let widgetClass = isWarning ? "widget warning" : "widget";
                let pluginType = data[key].hasOwnProperty("Type") ? data[key].Type : "";
                if (pluginType) {
                    delete data[key]["Type"];
                }
                let widgetSize = "";
                if (levelClass !== "inner") {
                    widgetSize = "small";
                    let isString = typeof data[key] === "string" || data[key] instanceof String;
                    if (!isString && Object.keys(data[key]).length > 4) {
                        widgetSize = "big";
                    }
                    widgetClass += " " + pluginType;
                }
                widgetClass += " " + widgetSize;
                if (Object.keys(data).length === 1 && typeof data[Object.keys(data)[0]] !== "string") {
                    data = data[Object.keys(data)[0]];
                    html += renderList(data, "inner");
                }
                else {
                    html += `<div class='${levelClass} ${widgetClass}'>`;
                    html += `<div class='widget-title'>${key}:</div>`;
                    html += renderList(data[key], "inner");
                    html += `</div>`;
                }
            }
        }
        if (pluginName) {
            html += `</div>`;
        }
        if (agentName) {
            html += `</div>`;
        }
    }
    else if (Array.isArray(data)) {
        html += "<div class='data_array'>";
        data.forEach((value) => {
            html += `<div class='data_array_element'>${renderList(value, "list")}</div>`;
        });
        html += "</div>";
    }
    else {
        html = `<div class='widget-data'>${data}</div>`;
    }
    return html;
}
const container = document.getElementById("data-container");
/**
 * Subscribe to Server-Sent Events (SSE)
 * Parse the received agent data as JSON and renders it on the page
 **/
eventSource.onmessage = function (event) {
    const jsonData = JSON.parse(event.data);
    container.insertAdjacentHTML("beforeend", renderList(jsonData, "outer"));
};
/**
 * Closes the SSE connection
 **/
eventSource.onerror = function (event) {
    eventSource.close();
    hideNonWarningWidgets();
};
const toggleSwitch = document.getElementById("toggleSwitch");
// Add an event listener to the toggle switch when changed
toggleSwitch === null || toggleSwitch === void 0 ? void 0 : toggleSwitch.addEventListener("change", hideNonWarningWidgets);
/**
 * This function hides/shows widgets with warning class
 */
function hideNonWarningWidgets() {
    if (!toggleSwitch)
        return;
    const widgets = document.querySelectorAll(".outer.widget");
    widgets.forEach((widget) => {
        if (!widget.classList.contains("parent-warning")) {
            widget.style.display = toggleSwitch.checked ? "none" : "";
        }
    });
}
