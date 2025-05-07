import 'gridstack/dist/gridstack.min.css';
import { GridStack } from 'gridstack';

let grid: GridStack;

/**
 * Return grid instance
 */
export function getGrid() {
    return grid;
}

/**
 * Initialize grid with placeholder widgets
 */
export function initGrid () {
    grid = GridStack.init(
        {
            float: false,
            layout: 'move',
            staticGrid: false,
            minRow: 2
        }
    );

    createPlaceholdersWidgets();
}

/**
 * Add new widget on a grid with html content
 */
export function makeNewWidget(innerHtml: string) {
    createWidget(innerHtml, 2, 2);
};

/**
 * Add new doubled size widget on a grid with html content
 */
export function makeNewBigWidget(innerHtml: string) {
    createWidget(innerHtml, 4, 4);
}

/**
 * Add new widget on a grid with html content
 * with specified parameters
 */
function createWidget(innerHtml: string, height: number, width: number) {
    let doc = document.implementation.createHTMLDocument();
    doc.body.innerHTML = `<div class="item"  gs-w="${width}" gs-h="${height}"><div class="grid-stack-item-content">${innerHtml}</div></div>`;
    let el = doc.body.children[0] as HTMLElement;
    grid.el.appendChild(el);
    grid.makeWidget(el);
};

/**
 * Create placeholder widgets without content
 * but with loading animation
 */
export function createPlaceholdersWidgets() {
    grid.on('added', function (e, items) {
        items.map(item => 
            item.el?.firstElementChild?.classList?.add("placeholder")
        );
    });

    const placeHolderWidgets = Array.from({ length: 6 }, () => ({ w: 2, h: 2 }));

    grid.load(placeHolderWidgets);
    grid.off('added');
}