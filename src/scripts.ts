import {renderList, hideNonWarningWidgets} from './render'


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


