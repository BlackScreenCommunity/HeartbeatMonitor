
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
