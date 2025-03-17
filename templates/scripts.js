
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

function sortAgentsData(pluginsData) {
    if(pluginsData) {
        let sortable = [];
        for (var plugin in pluginsData) {
            sortable.push([plugin, pluginsData[plugin]]);
        }
    
        sortable.sort(function(a, b) {
            return JSON.stringify(b).length - JSON.stringify(a).length;
        });


        pluginsData = {};
        sortable.forEach(function(item){
            pluginsData[item[0]]=item[1]
        })

        return pluginsData;
    }
}

function renderList(data, levelClass) {

    levelClass = levelClass ? levelClass : "";

    if(data?.data) {
        data.data = sortAgentsData(data.data);
    }
    


    let html = "";

    if (typeof data === "object" && !Array.isArray(data)) { 
        let duration = 0;
        let isWarning = data.hasOwnProperty("isWarning") ? data.isWarning : false;
        delete data["isWarning"];

        
        let agentName = data.hasOwnProperty("agent_name") ? data.agent_name : "";
        if(agentName) {
            duration = data?.duration;
            data = data?.data;
        }

        if(agentName) {
            html += `<div class='agent-name'> ${agentName}`;
            html += `<div class='timer agent-duration'> Loading time ${duration} seconds </div>`;
            html += `</div>`
            html += `<div class='agent-data'>`;

        }

        let pluginName = data.hasOwnProperty("plugin_name") ? data.plugin_name : "";

        if(pluginName) {
            data = data?.data;
        }

        if(pluginName) {
            html += `<div class='plugin-data'>`;
            html += `<div class='plugin-name'> ${pluginName} </div>`;
        } 
        else 
        {
            for (let key in data) {
                let widgetClass = isWarning ? "widget warning" : "widget";
                

                let pluginType = data[key].hasOwnProperty("Type") ? data[key].Type : "";

                if(pluginType) {
                    delete data[key]["Type"];
                }


                widgetSize = "";
                if(levelClass != "inner") {
                    widgetSize = "small";
                    if(Object.keys(data[key]).length > 4 ) {
                        widgetSize = "big"
                    }

                    widgetClass += " " + pluginType;
                }   
                
                
                
                widgetClass += " " + widgetSize
                
                if (Object.keys(data).length == 1 && typeof(data[Object.keys(data)[0]]) != "string")
                {
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

        if(pluginName) {
            html += `</div>`;
        }
        if(agentName) {
            html += `</div>`;
        }

        // if(agentName) {
        //     html += `</div>`;
        // }

    } else if (Array.isArray(data)) { 
        html += "<div class='data_array'>";
        data.forEach(value => {
            html += `<div class='data_array_element'>${renderList(value, "list")}</div>`;
        });
        html += "</div>";
    } else { 
        html = `<div class='widget-data'>${data}</div>`;
    }
    return html;
}
const container = $("#data-container");

 /**
 * Subscribe to Server-Sent Events (SSE)
 * Parse the received agent data as JSON and renders it on the page
 **/
eventSource.onmessage = function(event) {
    const jsonData = JSON.parse(event.data);
    container.append(renderList(jsonData, "outer"));
};

/**
 * Closes the SSE connection when an error occurs.
 * Logs the error details to the console.
 **/
eventSource.onerror = function(event) {
    console.error("SSE connection error:", event);
    eventSource.close();
};
