'use strict';
var server = require("./server");

function TrieNode(alpha) {
    this.alpha = alpha;
    this.children = [];
}

function getAllWordsInDB(docs) {
    console.log("Getting all the words in the db");
    let docsRet = [];
    let docsWords = [];
    //Get all the docs in the db
    //Loop through all the docs and get the movie title string
    for (let i of docs) {
        //String split the movie title into different words
        docsWords = (i.title).split(" ");
        for (let j of docsWords)
            docsRet.push((j.toLowerCase()).replace(/\W/g, ''));
    }
    console.log("Got all the words in the db");
    return docsRet;
}

function printChildren(children, parentAlpha) {
    let currentChildren = "";
    console.log("Children for:", parentAlpha);
    for (let child of children)
        currentChildren += child.alpha + ", ";
    console.log(currentChildren);
    for (let child of children)
        if (child.children.length > 0)
            printChildren(child.children, child.alpha);
}
function printTrie(trieHeads) {
    for (let head of trieHeads) {
        console.log("For head", head.alpha + ": ");
        printChildren(head.children, head.alpha);
    }
}

function searchThroughChildren(children, foundStrs, results) {
    const prevFoundStr = foundStrs;
    if (children.length == 0) {
        // console.log(foundStrs);
        return;
    }
    for (let child of children) {
        if(child.alpha == "`"){
            // console.log(foundStrs);
            results.push(foundStrs);
        }else{
            foundStrs = prevFoundStr + child.alpha;
            searchThroughChildren(child.children, foundStrs, results);
        }
            
    }
}

function searchFor(query, heads) {
    let theHead = null;
    let results = [];
    for (let i = 0; i < query.length; i++) {
        // console.log("Looking for:", query[i]);
        for (let head of heads) {
            if ((query[i]).toLowerCase() == head.alpha) {
                theHead = head;
                heads = theHead.children;
                // console.log("Found:", query[i]);
                break;
            }
        }
    }
    if (theHead == null) {
        console.log("No results found");
        return;
    }
    let foundStrs = query.toLowerCase() + "";
    searchThroughChildren(theHead.children, foundStrs, results);
    return results;
}
function createTrie() {
    let heads = [];
    let currentNode = null;
    let isInHead = 0;
    let tmpCurrentNode = null;
    let counter = 0;
    console.log("Finding all the elements in the db");
    server.movies.find({}, function (err, docs) {
        if (err) {
            console.log(err);
            return [];
        }
        else {
            console.log("Creating tree");
            let allWordsForTrie = getAllWordsInDB(docs);
            for (let i of allWordsForTrie) {
                i += "`";
                for (let char of i) {
                    if (currentNode == null) {
                        for (let head of heads) {
                            if (head.alpha == char) {
                                isInHead = 1;
                                currentNode = head;
                                break;
                            }
                        }
                        if (isInHead == 0) {
                            currentNode = new TrieNode(char);
                            heads.push(currentNode);
                        }
                        isInHead = 0;
                    } else {
                        for (let child of currentNode.children) {
                            if (child.alpha == char) {
                                isInHead = 1;
                                currentNode = child;
                                break;
                            }
                        }
                        if (isInHead == 0) {
                            tmpCurrentNode = currentNode;
                            currentNode = new TrieNode(char);
                            tmpCurrentNode.children.push(currentNode);
                        }
                        isInHead = 0;
                    }
                }
                currentNode = null;
                if (++counter % 100000 == 0)
                    console.log(":>", counter);
            }
        }
        module.exports.heads = heads;
        console.log("Done with creating tree");
        console.log("Number of heads", heads.length);
        // printTrie(heads);
    });
}

module.exports.create = createTrie;
module.exports.searchFor = searchFor;