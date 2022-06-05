// index.js

const busTable = document.getElementById("busTable");
const searchResultsTable = document.getElementById("searchResultsTable");

async function load() {
    await getPredictions("1517");
    await getRoutes();
}

// Generic POST request maker via REST API
const postRequest = async (reqURL, req) => {
    return await fetch(reqURL, {
        headers: {
            'Content-Type': 'application/json'
        },
        method: 'POST',
        body: JSON.stringify(req)
    })
        .then(response => response.json())
        .then(data =>{
            return data;
        });
}


/*
    Routes
*/

async function getRoutes() {
    let postData = {
        reqType: "getroutes"
    };
    let ctaResponse = await postRequest('/getCTAData', postData);
    let routes = ctaResponse["bustime-response"].routes;
    let routeSelect = document.getElementById("routeSelect");
    // remove the first option
    routeSelect.removeChild(routeSelect.firstChild);
    routes.forEach((route) => {
        let option = document.createElement("option");
        option.text = route.rt;
        option.value = route.rt;
        routeSelect.add(option);
    });
}


/*
    Predictions
*/
async function getPredictions(stopID) {
    let postData = {
        reqType: "getpredictions",
        stpid: stopID
    };
    let ctaResponse = await postRequest('/getCTAData', postData);
    let predictions = ctaResponse["bustime-response"].prd;
    console.log(predictions);
    clearTableBody(busTable);
    predictions.forEach((prd) => {
        let prdTime = parseCTATime(prd.prdtm);
        let currentTime = new Date();
        let diff = Math.ceil((prdTime - currentTime)/(1000*60));
        console.log(prdTime);
        addPredictionTableRow(busTable, prd.rt, prd.rtdir, diff + " minutes");
    });
}

// CTA time format is yyyymmdd hh:mm (24 hour)
function parseCTATime(datetm) {
    let dateArray = datetm.split(" ");
    let yr = dateArray[0].slice(0,4);
    let month = dateArray[0].slice(4,6);
    let day = dateArray[0].slice(6,8);
    let tmArray = dateArray[1].split(":");
    let hr = tmArray[0];
    let min = tmArray[1];
    console.log(yr, month, day, hr, min);
    let dateString = yr + "-" + month + "-" + day + "T" + hr + ":" + min + ":" + "00";
    date = new Date(dateString);
    return date;
}

/*
    Directions
*/

async function getDirections(routeID) {
    let postData = {
        reqType: "getdirections",
        rt: routeID
    };
    let ctaResponse = await postRequest('/getCTAData', postData);
    let directions = ctaResponse["bustime-response"].directions;
    console.log(directions);
    updateDirections(directions);
}

// update directionSelect with the results of getDirections
function updateDirections(directions) {
    let select = document.getElementById("directionSelect");
    clearSelect(select);
    let dotOption = document.createElement("option");
    dotOption.text = "...";
    dotOption.value = "";
    select.add(dotOption);
    directions.forEach((dir) => {
        let option = document.createElement("option");
        option.text = dir.dir;
        option.value = dir.dir;
        select.add(option);
    });
}


/*
    Stops
*/

async function getStops(routeID, direction) {
    let postData = {
        reqType: "getstops",
        rt: routeID,
        dir: direction
    };
    let ctaResponse = await postRequest('/getCTAData', postData);
    let stops = ctaResponse["bustime-response"].stops;
    console.log(stops);
    updateStops(stops);
}

// update stopSelect with the results of getStops
function updateStops(stops) {
    let select = document.getElementById("stopSelect");
    clearSelect(select);
    let dotOption = document.createElement("option");
    dotOption.text = "...";
    dotOption.value = "";
    select.add(dotOption);
    stops.forEach((stop) => {
        let option = document.createElement("option");
        option.text = stop.stpnm;
        option.value = stop.stpid;
        select.add(option);
    });
}


/*
    DOM updates
*/

function changeRoute() {
    let routeSelect = document.getElementById("routeSelect");
    let routeID = routeSelect.options[routeSelect.selectedIndex].value;
    getDirections(routeID);
}

function changeDirection() {
    let directionSelect = document.getElementById("directionSelect");
    let direction = directionSelect.options[directionSelect.selectedIndex].value;
    let routeSelect = document.getElementById("routeSelect");
    let routeID = routeSelect.options[routeSelect.selectedIndex].value;
    getStops(routeID, direction);
}

function changeStops() {
    let stopSelect = document.getElementById("stopSelect");
    let stopID = stopSelect.options[stopSelect.selectedIndex].value;
    getPredictions(stopID);
}

// Delete all options in a select element
function clearSelect(select) {
    while (select.firstChild) {
        select.removeChild(select.firstChild);
    }
}

// Delete all non-header rows in busTable
function clearTableBody(table) {
    let tableBody = table.getElementsByTagName('tbody')[0];
    while (tableBody.firstChild) {
        tableBody.removeChild(tableBody.firstChild);
    }
}

function addPredictionTableRow(table, routeNumber, direction, eta) {
    let tableBody = table.getElementsByTagName('tbody')[0];
    let row = tableBody.insertRow();
    let cell = row.insertCell(0);
    cell.innerHTML = "<b>" + routeNumber + "</b>";
    let dir = row.insertCell(1);
    dir.innerHTML = direction;
    let theEta = row.insertCell(2);
    theEta.innerHTML = eta;
}
