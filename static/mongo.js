// mongo.js

const MIN_SEARCH_LENGTH = 3;
const MAX_RESULTS = 10;

async function searchForStop(text) {
    if (text.length < MIN_SEARCH_LENGTH) {
        return;
    }
    let postData = {
        queryString: text
    };
    let mongoResponse = await postRequest('/search', postData);
    let stops = mongoResponse["mongo-response"];
    clearTableBody(searchResultsTable);
    // add a search result row for each stop, up to the MAX_RESULTS threshold
    for (let i = 0; i < Math.min(stops.length, MAX_RESULTS); i++) {
        let stop = stops[i];
        addSearchResultRow(stop);
    }
}

function addSearchResultRow(stop) {
    let tableBody = searchResultsTable.getElementsByTagName("tbody")[0];
    let row = tableBody.insertRow();
    // when the row is clicked, run getPredictions() and update the div text
    row.onclick = () => {
        getPredictions(stop.stpid);
        document.getElementById("stopDiv").innerHTML = stop.stpnm;
    };
    let idCell = row.insertCell(0);
    idCell.innerHTML = "<b>" + stop.stpid + "</b>";
    let descCell = row.insertCell(1);
    descCell.innerHTML = stop.StopDesc;
}