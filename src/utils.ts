
/**
 * Replace {0}, {1} ... {n} in a string with provided values
 */
export function FormatString(str: string, ...val: string[]) {
    for (let index = 0; index < val.length; index++) {
        str = str.replace(new RegExp(`\\{${index}\\}`, 'g'), val[index]);
    }
    return str;
}

/**
 * Escape string for regexp to prevent parsing errors
 */
function escapeRegExp(str: string) {
    return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

/**
 *  Count how many times a pattern appears in a text
 */
export function GetOccurencesCount(text: string, search: string) { 
    let trimmedText = text.replace(/ /g, '');
    let trimmedSearch = escapeRegExp(search.replace(/ /g, ''));

    return (trimmedText.match(new RegExp(trimmedSearch, 'g')) || []).length;
}