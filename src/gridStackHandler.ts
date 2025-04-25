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

export function makeNewWidget(innerHtml: string) {
    let n = getNode();
    let doc = document.implementation.createHTMLDocument();

    doc.body.innerHTML = `<div class="item"  gs-w="${2}" gs-h="${2}"><div class="grid-stack-item-content">${innerHtml}</div></div>`;
    let el = doc.body.children[0] as HTMLElement;
    grid.el.appendChild(el);
    let w = grid.makeWidget(el);
};
