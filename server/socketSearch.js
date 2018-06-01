var server = require("./server");
function onIMDBQuery(query)
{
    console.log("executing query for:", query);
    // var nameRegex = new RegExp(query);
    server.movies.find({
        title: new RegExp(query, 'i')
    }, function (err, docs) {
        if (err)
            console.log(err);
        else {
            console.log("did query");
            console.log(docs);
        }
    });
}
module.exports = {
    onIMDBQuery: onIMDBQuery
}