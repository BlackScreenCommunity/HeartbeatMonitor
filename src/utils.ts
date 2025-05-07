
export function FormatString(str: string, ...val: string[]) {
    for (let index = 0; index < val.length; index++) {
        str = str.replace(`{${index}}`, val[index]);
    }
    
    return str;
}

export function GetOccurencesCount(text: string, search: string) { 
    let trimmedText = text.replace(/ /g, '');
    let trimmedSearch = search.replace(/ /g, '');

    return (trimmedText.match(new RegExp(trimmedSearch, 'g')) || []).length;
}