'use strict';
var server = require("./server");
var trie = require("./trie");

function getSuggestionSearchArray(query) {
    let queryWords = query.split(" ");
    let suggestions;
    let stringFiltersArray = [];
    let ret = [];
    for (let word of queryWords) {
        let results = [];
        suggestions = [];
        suggestions = trie.searchFor(word, trie.heads);
        if (suggestions && suggestions.length > 0)
            for (let suggestion of suggestions)
                results.push(suggestion);
        results.push(word);
        results = results.join(" ");
        stringFiltersArray.push(results);
    }
    return stringFiltersArray;
}

const diffBy = (pred) => (a, b) => a.filter(x => b.some(y => pred(x, y)))
const makeSymmDiffFunc = (pred) => (a, b) => diffBy(pred)(a, b).concat(diffBy(pred)(b, a))

const myDiff = makeSymmDiffFunc((x, y) => x.mId != y.mId)

function obj1HasObj2Property(aElem, bA, propName) {
    for (let i of bA)
        if (i[propName] == aElem[propName])
            return true;
    return false;
}

function retObjArrSimlairities(a, b) {
    return a.filter(e => obj1HasObj2Property(e, b, "mId"));
}

function chainQueries(query, fx) {
    let queries = getSuggestionSearchArray(query);
    let totalDocs;
    let currentDoc = null;
    for (let i = 0; i < queries.length; i++) {
        // console.log(suggestions);
        server.movies.find({
            $text: {
                $search: queries[i]
            }
        }).sort({"numVotes": -1, "rating": -1, "year": -1}).limit(50).exec(function (err, docs) {
            if (err)
                console.log(err);
            else {
                console.log("did query");
                // console.log(docs);
                if (currentDoc == null)
                    currentDoc = docs;
                else {
                    if (docs.length > 0)
                        currentDoc = retObjArrSimlairities(currentDoc, docs);
                }
                if (i == queries.length - 1)
                    fx(currentDoc);
            }
        });
    }
}
function onIMDBQuery(query, socket) {
    console.log("executing query for:", query);
    // var nameRegex = new RegExp(query);
    // console.log("suggestions: ", trie.searchFor(query, trie.heads)); 
    // let stringFilters = getSuggestionSearchString(query);
    // console.log(stringFilters);
    chainQueries(query, function (currentDoc) {
        // console.log(currentDoc);
        socket.emit("server: search result", JSON.stringify(currentDoc));
    });
}
module.exports = {
    onIMDBQuery: onIMDBQuery
}