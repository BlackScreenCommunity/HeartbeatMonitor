/**
 * Create a MutationObserver that watches for changes in the DOM
 **/
const observer = new MutationObserver((mutationsList: MutationRecord[]) => {
    mutationsList.forEach((mutation: MutationRecord) => {
        if (mutation.type === "childList") {
            processWarnings();
        }
    });
});

/**
 * Select the document body as the target to observe for changes
 */
const targetNode: Node = document.body;

/**
 * Define the observer settings to watch 
 * for changes in child elements and all of its descendants
 */
const config: MutationObserverInit = {
    childList: true,
    subtree: true
};

/** 
 * Start observing
 **/
observer.observe(targetNode, config);


/** 
 * Find divs with 'warning' class 
 * and marking their parent elements with warning
 **/
function processWarnings(): void {
    document.querySelectorAll("div.warning").forEach((el) => {
        el.parentElement?.classList.add("parent-warning");
    });
}

/**
 * Create an EventSource to receive 
 * server-sent events to the "/events" endpoint.
 */
const eventSource = new EventSource("/events");

/**
 * DataObject with plugions metrics data
 */
interface PluginData {
    [key: string]: any;
}

/**
 * Sorts plugin data by the length of each [key, value] pair's JSON string, from longest to shortest.
 *
 * Takes a PluginData object, turns it into an array of [key, value] pairs,
 * Sorts data in PluginData object so plugins with longer data come first.
 */
function sortAgentsData(pluginsData: PluginData): PluginData {
    if (pluginsData) {
        let sortable: [string, any][] = Object.entries(pluginsData);

        sortable.sort((a, b) => JSON.stringify(b).length - JSON.stringify(a).length);

        let sortedData: PluginData = {};
        sortable.forEach(([key, value]) => {
            sortedData[key] = value;
        });

        return sortedData;
    }
    return {};
}
function renderList(data: any, levelClass: string = ""): string {
    if (data == null) return "";

    if (Array.isArray(data)) {
        const itemsHtml = data
            .map((item) => `<div class="data_array_element">${renderList(item, "list")}</div>`)
            .join("");
        return `<div class="data_array">${itemsHtml}</div>`;
    }

    if (typeof data !== "object") {
        return `<div class="widget-data">${data}</div>`;
    }

    let currentData = { ...data };

    const isWarning: boolean = currentData.isWarning || false;
    delete currentData.isWarning;

    let html = "";

    const hasAgent = "agent_name" in currentData;
    if (hasAgent) {
        const agentName: string = currentData.agent_name;
        const duration: number = currentData.duration;
        currentData = currentData.data;
        html += `<div class="agent-name">
                 ${agentName}
                 <div class="timer agent-duration">Loading time ${duration} seconds</div>
               </div>
               <div class="agent-data">`;
    }

    const hasPlugin = currentData && "plugin_name" in currentData;
    if (hasPlugin) {
        const pluginName: string = currentData.plugin_name;
        currentData = currentData.data;
        html += `<div class="plugin-data">
                 <div class="plugin-name">${pluginName}</div>`;
    } else {
        for (const key in currentData) {
            if (!currentData.hasOwnProperty(key)) continue;
            let item = currentData[key];
            let widgetClass = `widget${isWarning ? " warning" : ""}`;
            let pluginType = "";

            if (typeof item === "object" && item && "Type" in item) {
                pluginType = item.Type;
                item = { ...item };
                delete item.Type;
            }

            let widgetSize = "";
            if (levelClass !== "inner") {
                widgetSize = "small";
                if (typeof item === "object" && item && !Array.isArray(item) && Object.keys(item).length > 4) {
                    widgetSize = "big";
                }
                widgetClass += ` ${pluginType}`;
            }
            widgetClass += ` ${widgetSize}`;

            if (Object.keys(currentData).length === 1 && typeof item !== "string") {
                html += renderList(item, "inner");
            } else {
                html += `<div class="${levelClass} ${widgetClass}">
                     <div class="widget-title">${key}:</div>
                     ${renderList(item, "inner")}
                   </div>`;
            }
        }
    }

    if (hasPlugin) {
        html += `</div>`;
    }
    if (hasAgent) {
        html += `</div>`;
    }

    return html;
}


/**
 * Container where data from the server 
 * will be displayed 
 **/
const container = document.getElementById("data-container") as HTMLElement;

/**
 * Subscribe to Server-Sent Events (SSE)
 * Parse the received agent data as JSON and renders it on the page
 **/
eventSource.onmessage = function (event: MessageEvent) {
    const jsonData = JSON.parse(event.data);
    container.insertAdjacentHTML("beforeend", renderList(jsonData, "outer"));
};

/**
 * Closes the SSE connection
 **/
eventSource.onerror = function (event: Event) {
    eventSource.close();
    hideNonWarningWidgets();
};

const toggleSwitch = document.getElementById("toggleSwitch") as HTMLInputElement | null;

// Add an event listener to the toggle switch when changed
toggleSwitch?.addEventListener("change", hideNonWarningWidgets);

/**
 * This function hides/shows widgets with warning class
 */
function hideNonWarningWidgets(): void {
    if (!toggleSwitch) return;

    const widgets = document.querySelectorAll<HTMLElement>(".outer.widget");

    widgets.forEach((widget) => {
        if (!widget.classList.contains("parent-warning")) {
            widget.style.display = toggleSwitch.checked ? "none" : "";
        }
    });
}
