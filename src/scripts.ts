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
    let html = "";

    if (typeof data === "object" && !Array.isArray(data)) {
        let duration: number = 0;
        let isWarning: boolean = data.hasOwnProperty("isWarning") ? data.isWarning : false;
        delete data["isWarning"];

        let agentName: string = data.hasOwnProperty("agent_name") ? data.agent_name : "";
        if (agentName) {
            duration = data?.duration;
            data = data?.data;
        }

        if (agentName) {
            html += `<div class='agent-name'> ${agentName}`;
            html += `<div class='timer agent-duration'> Loading time ${duration} seconds </div>`;
            html += `</div>`;
            html += `<div class='agent-data'>`;
        }

        let pluginName: string = data.hasOwnProperty("plugin_name") ? data.plugin_name : "";

        if (pluginName) {
            data = data?.data;
        }

        if (pluginName) {
            html += `<div class='plugin-data'>`;
            html += `<div class='plugin-name'> ${pluginName} </div>`;
        } else {
            for (let key in data) {
                let widgetClass: string = isWarning ? "widget warning" : "widget";

                let pluginType: string = data[key].hasOwnProperty("Type") ? data[key].Type : "";
                if (pluginType) {
                    delete data[key]["Type"];
                }

                let widgetSize: string = "";
                if (levelClass !== "inner") {
                    widgetSize = "small";

                    let isString: boolean = typeof data[key] === "string" || data[key] instanceof String;

                    if (!isString && Object.keys(data[key]).length > 4) {
                        widgetSize = "big";
                    }

                    widgetClass += " " + pluginType;
                }

                widgetClass += " " + widgetSize;

                if (Object.keys(data).length === 1 && typeof data[Object.keys(data)[0]] !== "string") {
                    data = data[Object.keys(data)[0]];
                    html += renderList(data, "inner");
                } else {
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
    } else if (Array.isArray(data)) {
        html += "<div class='data_array'>";
        data.forEach((value) => {
            html += `<div class='data_array_element'>${renderList(value, "list")}</div>`;
        });
        html += "</div>";
    } else {
        html = `<div class='widget-data'>${data}</div>`;
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
