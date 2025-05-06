import {hideNonWarningWidgets} from './render'

/** 
 * Find divs with 'warning' class 
 * and marking their parent elements with warning
 **/
function processWarnings(): void {
    document.querySelectorAll("div.warning").forEach((el) => {
        el.parentElement?.parentElement?.classList.add("parent-warning");
    });

    hideNonWarningWidgets();
}
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

