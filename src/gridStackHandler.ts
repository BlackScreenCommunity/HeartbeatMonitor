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

export function makeNewWidget(innerHtml: string) {
    let doc = document.implementation.createHTMLDocument();

    doc.body.innerHTML = `<div class="item"  gs-w="${2}" gs-h="${2}"><div class="grid-stack-item-content">${innerHtml}</div></div>`;
    let el = doc.body.children[0] as HTMLElement;
    grid.el.appendChild(el);
    grid.makeWidget(el);
};
