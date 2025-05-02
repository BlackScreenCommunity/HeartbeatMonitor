import 'gridstack/dist/gridstack.min.css';
import { GridStack } from 'gridstack';

var grid = GridStack.init(
    {
        float: false,
        layout: 'move',
        staticGrid: false,
        minRow: 2
    }
);

createPlaceholdersWidgets();


export function makeNewWidget(innerHtml: string) {
    let doc = document.implementation.createHTMLDocument();

    doc.body.innerHTML = `<div class="item"  gs-w="${2}" gs-h="${2}"><div class="grid-stack-item-content">${innerHtml}</div></div>`;
    let el = doc.body.children[0] as HTMLElement;
    grid.el.appendChild(el);
    grid.makeWidget(el);
};

export function createPlaceholdersWidgets() {
    grid.on('added', function (e, items) {
        for (let i = 0; i < items.length; i++) {
            if (items[i] && items[i].el && items[i].el?.firstElementChild) {
                items[i].el?.firstElementChild?.classList?.add("placeholder");
            }
        }
    });

    var items = [
        { w: 2, h: 2 },
        { w: 2, h: 2 },
        { w: 2, h: 2 },
        { w: 2, h: 2 },
        { w: 2, h: 2 },
        { w: 2, h: 2 }
    ];

    grid.load(items);
    grid.off('added');
}

export function getGrid() {
    return grid;
}