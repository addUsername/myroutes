import $ from 'jquery';
import datepickerFactory from 'jquery-datepicker';
import datepickerJAFactory from 'jquery-datepicker/i18n/jquery.ui.datepicker-ja';
 
import 'core-js/stable';
import 'promise-polyfill/src/polyfill';
const runtime = require('@wailsapp/runtime');

function getAll(){
	var prom = new Promise(function(resolve, reject) {
		var json = window.backend.Handler.GetAllClients()
		if (json) {
		  //console.log("stuff worked")
			resolve(json);
		}  else {
			//console.log("fuuck")
		  reject(new Error("It broke"));
		}
	  });
	   
	  prom.then(function(json) {
		json = JSON.parse(json)
		drawClients(json)
	  });
}
// OLD getAll() function working fine but not ie11 compatible
/*
function getAll(){
	window.backend.Handler.GetAllClients().then( (json) => {
		console.log("GET ALL CLIENTS ")
		json = JSON.parse(json)
		drawClients(json)
	});
}
*/

function drawClients(rou){
	console.log("DRAAAAAWWWWW")
	
	var store = document.getElementById('store');
	store.innerHTML = ""

	var routePicker = document.getElementById("route-picker")
	routePicker.innerHTML = ""
	var sel = document.createElement("select")
	sel.id = "route"

	for (var k = 0; k < rou.length; k++) {

		var routContent = document.createElement("div")
		routContent.className = "route"
		var op = document.createElement("option")
		op.value = rou[k].name
		op.innerText = rou[k].name
		sel.appendChild(op)


		var clients = rou[k].clients

		for (var i = 0; i < clients.length; i++) {
			//console.log("client number "+ i)
			var element = clients[i]

			var table = document.createElement("table")
			table.className="styled-table"
			var thead = document.createElement("thead")
			var row = document.createElement("tr")
			//table.id = i
			table.onclick = function(ev) {			
				ev.currentTarget.childNodes[1].hidden = !ev.currentTarget.childNodes[1].hidden
			}

			var th = document.createElement("th")
			th.innerText = element.name
			var th2 = document.createElement("th")
			
			var th3 = document.createElement("th")

			//  = element.warnings, a num that indicates how many visits the client has in 30, 60 days
			th3.innerText = element.routename
			th3.className = "align-right"

			row.appendChild(th)
			row.appendChild(th2)
			row.appendChild(th3)
			thead.appendChild(row)
			table.appendChild(thead)

			var tbody = document.createElement("tbody")
			tbody.hidden = true
			for (var j = 0; j < element.activities.length; j++) {
				
				var activity = element.activities[j]
				var row2 = document.createElement("tr")
				var th2 = document.createElement("th")
				th2.className = "row"
				
				var td = document.createElement("td")
				td.innerText = activity.name
				td.className = "align-left"
				var td2 = document.createElement("td")
				td2.innerText = activity.enddate
				td2.className = "date-child"
				var td3 = document.createElement("td")
				td3.innerText = activity.madedate
				td3.className = "date-child"

				row2.appendChild(td)
				row2.appendChild(td2)
				row2.appendChild(td3)

				tbody.appendChild(row2)
			};
			table.appendChild(tbody)			
			routContent.appendChild(table)
		}
		store.appendChild(routContent)
	}
	routePicker.appendChild(sel)
}

// Main entry point
function start(){

	var mystore = runtime.Store.New('Clients');
	//console.log(mystore);

	
	// Ensure the default app div is 100% wide/high
	var app = document.getElementById('app');
	app.style.width = '100%';
	app.style.height = '100%';

	app.innerHTML = `<div class='container'>
	<div id='header' class='header'>
		<ul>
			<li id="route-picker"></li>
			<li><input type="date" name="date" id="enddate" value=""/></li>
			<li><input type="date" name="date2" id="madedate" value=""/></li>
			<li><button id="search">search</button></li>
		</ul>
	</div>
	<div id='store' class='table-conainer small-text'></div>
	</div>`;
		 
	datepickerFactory($);
	$('#enddate').datepicker(); 
	$('#madedate').datepicker(); 
	
	document.getElementById('search').onclick = doSearch;
	
	function doSearch(e) {

		console.log("doSearch")
		var prom = new Promise(function(resolve, reject) {
			var json = window.backend.Handler.DoSearch(
					document.getElementById("route").value,
					document.getElementById('enddate').value,
					document.getElementById('madedate').value
				)
			if (json) {
			  	//console.log("stuff worked")
				resolve(json);
			}  else {
				//console.log("fuuck")
			  	reject(new Error("It broke"));
			}
		  });
		   
		  prom.then(function(json) {
			
			json = JSON.parse(json)
			console.log(json)
			drawClients(json)
		  });	

	}
	/*
	// I see no difference, give it a try

	mystore.subscribe( function(state) {
		console.log("MY store ")
		drawClients(JSON.parse(state));
	});
	*/	
	//mystore.set(0);
	getAll()
};

// We provide our entrypoint as a callback for runtime.Init
runtime.Init(start);