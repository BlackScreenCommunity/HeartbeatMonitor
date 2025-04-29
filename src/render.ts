/* eslint @typescript-eslint/no-explicit-any: 0 */

import * as Utils from './utils';
import {makeNewWidget} from './gridStackHandler'


/**
 * Render simple widget data, like text, number, etc
 */
function renderPrimitive(data: object): string {
    return `<div class="widget-data">${data}</div>`;
}

/**
 Render set of widget indicators
 */
function renderArray(data: object[]): string {
    const itemsHtml = data
        .map((item) => `<div class="data_array_element">${renderWidgetData(item, 99)}</div>`)
        .join("");
    return `<div class="data_array">${itemsHtml}</div>`;
}

/**
 * Build html for agent title
 */
function renderAgentTitleForWidget(data: any): { html: string; data: any } {
    if ("agent_name" in data) {
        const agentName: string = data.agent_name;

        const innerData = data.data;
        const html = `
                    <div class="widget-content" data-group="${agentName}">
                        <div class="built-in-agent-name">
                            ${agentName}
                        </div>
                        {0}
                    </div>
                    `;
        return { html, data: innerData };
    }
    return { html: "{0}", data };
}

function addAgentToSwitcher(data: any): void {
    if ("agent_name" in data) {
        const agentName: string = data.agent_name;
        const agentSwitcher = document.getElementById("agent-switcher");

        if (agentSwitcher) {
            const id = `agent-switcher-${crypto.randomUUID()}`;
            agentSwitcher.insertAdjacentHTML(
                "beforeend",
                `<div class="agent-toggle btn" id="${id}">
                    ${agentName}
                </div>`
            );

            const agentSwitcherElement = document.getElementById(id);
            if (agentSwitcherElement) {
                agentSwitcherElement.addEventListener("click", () => {
                    // console.log("Clicked on:", agentName);
                    agentSwitcherElement.classList.toggle("active");
                });
            }
        }
    }
}

/**
 * Build html for plugin header
 */
function renderPluginHeader(pluginName: string): string {
        const html = `<div class="plugin-name">${pluginName}</div>{0}`;
        return html;

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

            html += renderWidgetData(item, levelClass + 1, renderPluginHeader(Object.keys(data)[0]));
        } else {
            html += `<div class="${levelClass} ${widgetClass}">
                    <div class="widget-title">${key}:</div>
                    ${renderWidgetData(item, levelClass + 1, template)}
                 </div>`;
        }
    }
    return Utils.FormatString(template, html);
}

/**
 * Render plugins data recieved from server
 */
export function renderList(data: any, levelClass: number = 0, template: string = "{0}"): void {
    if (data == null) return;

    let currentData = { ...data };
    const isWarning: boolean = currentData.isWarning || false;
    delete currentData.isWarning;

    let html = "";

    addAgentToSwitcher(currentData);
    
    const { html: agentTitleHtml, data: afterAgentData } = renderAgentTitleForWidget(currentData);
    currentData = afterAgentData;
    html += Utils.FormatString(template, agentTitleHtml);

    renderPluginData(html, currentData, levelClass + 1, isWarning);
}

/**
 * Render widget data based on its type (Array, nested widget data or primitive value)
 * Return rendered HTML or an empty string if a new widget is created.
 */
function renderWidgetData(data: any, levelClass: number = 0, template: string = "{0}"): string {
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

    let pluginDataHtml = renderPluginData(template, currentData, levelClass + 1, isWarning)

    if (pluginDataHtml.length > 0 && levelClass <= 2) {
        makeNewWidget(`${pluginDataHtml}`);
        return ""
    }

    return pluginDataHtml;
}

/**
 * A "Show Only Warnings" toggle switch element
 */
const toggleSwitch = document.getElementById("show_only_warnings_checkbox") as HTMLInputElement | null;

/**
 * Add an event listener to the toggle switch when changed
 */
toggleSwitch?.addEventListener("change", hideNonWarningWidgets);

/**
 * This function hides/shows widgets with warning class
 */
export function hideNonWarningWidgets(): void {
    if (!toggleSwitch) return;

    const widgets = document.querySelectorAll<HTMLElement>(".grid-stack-item-content");

    widgets.forEach((widget) => {
        if (!widget.classList.contains("parent-warning") && widget.parentElement) {
            widget.parentElement.style.display = toggleSwitch.checked ? "none" : "";
        }
    });
}
