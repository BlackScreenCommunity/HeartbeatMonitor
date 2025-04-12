import 'gridstack/dist/gridstack.min.css';
import { GridStack } from 'gridstack';

var grid = GridStack.init(
    {
        float: false,
        layout: 'move'
    }
);


export function FormatString(str: string, ...val: string[]) {
    for (let index = 0; index < val.length; index++) {
      str = str.replace(`{${index}}`, val[index]);
    }
    return str;
  }

let items = [
    { x: 1, y: 1, w: 1, h: 1 }, //, locked:true, content:"locked"},
  ];
  let count = 0;

function getNode() {
    let n = items[count] || {
      x: 2,
      y: 2,
      w: 2,
      h: 2
    };
    count++;
    return n;
  };

function addNewWidget() {
let w = grid.addWidget(getNode());
};

function makeNewWidget(innerHtml: string) {
let n = getNode();
let doc = document.implementation.createHTMLDocument();

doc.body.innerHTML = `<div class="item"  gs-w="${2}" gs-h="${2}"><div class="grid-stack-item-content">${innerHtml}</div></div>`;
let el = doc.body.children[0] as HTMLElement; 
grid.el.appendChild(el);
let w = grid.makeWidget(el);
};


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
 * Start observing chancges in DOM
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

/**
 * Render simple widget data, like text, number, etc
 */
function renderPrimitive(data: any): string {
    return `<div class="widget-data">${data}</div>`;
}

/**
 Render set of widget indicators
 */
function renderArray(data: any[]): string {
    const itemsHtml = data
        .map((item) => `<div class="data_array_element">${renderList(item, 99)}</div>`)
        .join("");
    return `<div class="data_array">${itemsHtml}</div>`;
}

/**
 * Render agent's plugins data
 */
function renderAgentSection(data: any): { html: string; data: any } {
    if ("agent_name" in data) {
        const agentName: string = data.agent_name;
        const duration: number = data.duration;

        const innerData = data.data;
        const html = `<div class="agent-name">
                      ${agentName}
                      <div class="timer agent-duration">Loading time ${duration} seconds</div>
                    </div>
                    <div class="agent-data">`;
        return { html, data: innerData };
    }
    return { html: "", data };
}

/**
 * Build html for agent title
 */
function renderAgentTitleForWidget(data: any): { html: string; data: any } {
    if ("agent_name" in data) {
        const agentName: string = data.agent_name;
        const duration: number = data.duration;

        const innerData = data.data;
        const html = `<div class="built-in-agent-name">
                      ${agentName}
                    </div>
                    {0}
                    `;
        return { html, data: innerData };
    }
    return { html: "{0}", data };
}

/**
 * Build html for plugin header
 */
function renderPluginHeader(data: any): { html: string; data: any } {
    if (data && "plugin_name" in data) {
        const pluginName: string = data.plugin_name;

        const innerData = data.data;
        const html = `<div class="plugin-data">
                      <div class="plugin-name">${pluginName}</div>`;
        return { html, data: innerData };
    }
    return { html: "", data };
}

/**
 * Build html for plugin data
 */
function renderPluginData(
    template: string, 
    data: any,
    levelClass: number,
    isWarning: boolean
): string {
    let html = "";
    for (const key in data) {
        if (!data.hasOwnProperty(key)) continue;

        let item = data[key];
        let widgetClass = `widget${isWarning ? " warning" : ""}`;
        let pluginType = "";

        if (typeof item === "object" && item && "Type" in item) {
            pluginType = item.Type;
            item = { ...item };
            delete item.Type;
        }

        let widgetSize = "";
        if (levelClass > 1) {
            widgetSize = "small";
            if (
                typeof item === "object" &&
                item &&
                !Array.isArray(item) &&
                Object.keys(item).length > 4
            ) {
                widgetSize = "big";
            }
            widgetClass += ` ${pluginType}`;
        }
        widgetClass += ` ${widgetSize}`;

        if (Object.keys(data).length === 1 && typeof item !== "string") {
            html += renderList(item, levelClass + 1, template);
        } else {
            html += `<div class="${levelClass} ${widgetClass}">
                    <div class="widget-title">${key}:</div>
                    ${renderList(item, levelClass + 1, template)}
                 </div>`;
        }
    }
    return FormatString(template, html);
}

/**
 * Render plugins data recieved from server
 */
function renderList(data: any, levelClass: number = 0, template: string = "{0}"): string {
    if (data == null) return "";

    if (Array.isArray(data)) {
        return renderArray(data);
    }

    if (typeof data !== "object") {
        return renderPrimitive(data);
    }

    let currentData = { ...data };
    const isWarning: boolean = currentData.isWarning || false;
    delete currentData.isWarning;

    let html = "";

    const { html: agentTitleHtml, data: afterAgentData } = renderAgentTitleForWidget(currentData);    // html += agentHtml;
    currentData = afterAgentData;

    const { html: pluginHtml, data: afterPluginData } = renderPluginHeader(currentData);
    html += pluginHtml;
    currentData = afterPluginData;

    let pluginDataHtml = FormatString(template, agentTitleHtml);

    // if (!pluginHtml) {
        pluginDataHtml = FormatString(pluginDataHtml, renderPluginData(agentTitleHtml, currentData, levelClass+1, isWarning));
        
    // }

    if(pluginDataHtml.length > 0 && levelClass >= 1 && levelClass <= 2) {
        makeNewWidget(`${pluginDataHtml}`);
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
    container.insertAdjacentHTML("beforeend", renderList(jsonData, 0));
};

/**
 * Closes the SSE connection
 **/
eventSource.onerror = function (event: Event) {
    eventSource.close();
    hideNonWarningWidgets();
};

const toggleSwitch = document.getElementById("toggleSwitch") as HTMLInputElement | null;

/**
 * Add an event listener to the toggle switch when changed
 */
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

