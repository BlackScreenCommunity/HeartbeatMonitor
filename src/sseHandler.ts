import {renderList, hideNonWarningWidgets} from './render'

/**
 * Create an EventSource to receive 
 * server-sent events to the "/events" endpoint.
 */
const eventSource = new EventSource("/events");

/**
 * Subscribe to Server-Sent Events (SSE)
 * Parse the received agent data as JSON and renders it on the page
 **/
eventSource.onmessage = function (event: MessageEvent) {
    const jsonData = JSON.parse(event.data);
    renderList(jsonData, 0);
};

/**
 * Closes the SSE connection
 **/
eventSource.onerror = function () {
    eventSource.close();
    hideNonWarningWidgets();
};
