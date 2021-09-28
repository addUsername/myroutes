//https://materialdesignicons.com/

import $ from 'jquery';
import datepickerFactory from 'jquery-datepicker';
//import datepickerJAFactory from 'jquery-datepicker/i18n/jquery.ui.datepicker-ja';
 
import 'core-js/stable';
import 'promise-polyfill/src/polyfill';
const runtime = require('@wailsapp/runtime');

// This variable keeps track of modified madedate ids
var updates = []

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
	//sel.className = "selector"
	sel.id = "route"

	for (var k = 0; k < rou.length; k++) {

		var routContent = document.createElement("div")
		routContent.className = "route"
		var op = document.createElement("option")
		op.value = rou[k].name
		op.innerText = rou[k].name
		sel.appendChild(op)

		// Add title to each route tables
		var title = document.createElement("div")
		title.className = "title-route"
		title.innerHTML=rou[k].name
		routContent.appendChild(title)

		var clients = rou[k].clients

		for (var i = 0; i < clients.length; i++) {

			var element = clients[i]

			var table = document.createElement("table")
			table.className="styled-table"
			
			var thead = document.createElement("thead")

			
			var row = document.createElement("tr")
			
			// adapt for comments event too
			table.onclick = function(ev) {
				console.log(ev.target.localName)
				if (ev.target.localName =="th" || ev.target.className =="tooltip"){
					ev.currentTarget.childNodes[1].hidden = !ev.currentTarget.childNodes[1].hidden
				}				
			}
			var th = document.createElement("th")
			th.span = 2

			// Add comments as tooltip			
			var tooltip = document.createElement("div")
			tooltip.className = "tooltip"
			tooltip.innerText = element.name

			var span = document.createElement("span")
			span.id = rou[k].routeid+"-"+element.clientid
			span.className = "tooltiptext"
			span.contentEditable = true
			span.innerText = element.comments

			tooltip.onclick = function (ev){

				if ( ev.target.className =="tooltip"){
					return false
				}
				console.log ("edit comment")
				pushUpdate(ev.target.id)


				console.log(ev.target)
				return true
			}

			tooltip.appendChild(span)
			th.appendChild(tooltip)

			var th2 = document.createElement("th")
			
			// Add comments about the company
			var th3 = document.createElement("th")
			/*
			th3.innerText = element.routename
			th3.className = "align-right"
			*/	

			var tbody = document.createElement("tbody")
			tbody.hidden = true
			var notCompleted = 0
			var timeline = 0
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

				// set colors to enddate
				switch(activity.colorwarning){
					case 114: //r
						td2.className = "date-child not-completed"
						notCompleted = notCompleted + 1
						break;
					case 121: //y
						td2.className = "date-child timeline"
						timeline = timeline+1
						break;
					case 98: //b
						td2.className = "date-child pending"
						break;
					case 103: //g
					default:
						td2.className = "date-child"
				}
				
				// last column of each activity will have an edit, and that is the whole point of the app

				var td3 = document.createElement("td")
				td3.id = rou[k].routeid+"-"+element.clientid+"-"+activity.activityid
				//console.log(td3.id) //1-206-213

				if(activity.madedate.length > 1){
					td3.innerText = activity.madedate
				}else{
					td3.innerHTML = `<svg style="width:24px;height:24px" viewBox="0 0 24 24">
										<path fill="black" d="M20.71,7.04C21.1,6.65 21.1,6 20.71,5.63L18.37,3.29C18,2.9 17.35,2.9 16.96,3.29L15.12,5.12L18.87,8.87M3,17.25V21H6.75L17.81,9.93L14.06,6.18L3,17.25Z" />
									</svg>`
				}				
				td3.className = "date-child"
				td3.onclick = editMadedate

				function editMadedate(ev){
					// Here, we just anotate which ids has been clicked and load and a datepicker inside the cell
					console.log("EDIT MADE DATTEEE")

					if(ev.target.localName == "path" || ev.target.localName == "svg"){
						this.innerHTML = `<input type="date" name="date2" id="`+ this.id+"-A"+`"" value=""/>`
						//add datepicker
						$("#"+this.id+"-A").datepicker();
					}
					pushUpdate(this.id+"-A")
				}

				row2.appendChild(td)
				row2.appendChild(td2)
				row2.appendChild(td3)

				tbody.appendChild(row2)

			};
			row.appendChild(th)
			row.appendChild(th2)
			// Add badges to thead 3rd column			
			th3.innerHTML = `<div class="wrapper">
				<span class="badge not-completed">`+notCompleted+`</span>
				<span class="badge timeline">`+timeline+`</span>
			</div>`			
			row.appendChild(th3)
			thead.appendChild(row)

			table.appendChild(thead)
			table.appendChild(tbody)

			routContent.appendChild(table)
		}
		store.appendChild(routContent)
	}
	routePicker.appendChild(sel)
}

function pushUpdate(id){

	// no duplicates
	var bol = true
	for (let index = 0; index < updates.length; index++) {
		
		if (updates[index] == id){
			bol = false
			break
		}
	}
	if(bol)	{
		updates.push(id)
	}
}

// Main entry point
function start(){

	var mystore = runtime.Store.New('Clients');
	
	
	
	// Ensure the default app div is 100% wide/high
	var app = document.getElementById('app');
	app.style.width = '100%';
	app.style.height = '100%';

	app.innerHTML = `<div class='container'>
	<div id='header' class='header'>
		<ul class="selector">			
			<li id="home-icon">
				<div class="tooltip-bottom">
					<svg style="width:34px;height:34px" viewBox="0 0 24 24">
						<path fill="white" d="M10,20V14H14V20H19V12H22L12,3L2,12H5V20H10Z" />
					</svg>
				<span class="tooltiptext-bottom"> Inicio </span>
				</div>
			</li>
			<li id="world">
				<div class="tooltip-bottom">
					<svg style="width:34px;height:34px" viewBox="0 0 24 24">
						<path fill="white" d="M17.9,17.39C17.64,16.59 16.89,16 16,16H15V13A1,1 0 0,0 14,12H8V10H10A1,1 0 0,0 11,9V7H13A2,2 0 0,0 15,5V4.59C17.93,5.77 20,8.64 20,12C20,14.08 19.2,15.97 17.9,17.39M11,19.93C7.05,19.44 4,16.08 4,12C4,11.38 4.08,10.78 4.21,10.21L9,15V16A2,2 0 0,0 11,18M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z" />
					</svg>
					<span class="tooltiptext-bottom"> no </span>
				</div>
			</li>
			<li class="no-hover">   |   </li>
			<li class="search" id="route-picker"></li>
			<li class="no-hover">
				<svg style="width:40px;height:40px" viewBox="0 0 24 24">
					<path fill="black" d="M2,3H8V5H4V19H8V21H2V3M7,17V15H9V17H7M11,17V15H13V17H11M15,17V15H17V17H15Z" />
				</svg>
			</li>
			<li class="search">			
				<div class="tooltip-bottom">	
					<input type="date" name="date" id="enddate" value=""/>
					<span class="tooltiptext-bottom"> Desde </span>
				</div>			
			</li>
			<li class="search">
				<div class="tooltip-bottom">
					<input type="date" name="date2" id="madedate" value=""/>
					<span class="tooltiptext-bottom"> Hasta </span>
				</div>
			</li>
			<li class="no-hover">
				<svg style="width:40px;height:40px" viewBox="0 0 24 24">
					<path fill="black" d="M7,17V15H9V17H7M11,17V15H13V17H11M15,17V15H17V17H15M22,3V21H16V19H20V5H16V3H22Z" />
				</svg>
			</li>
			<li class="search">
				<div class="tooltip-bottom">
					<svg id="search"  style="width:34px;height:30px;" viewBox="0 0 24 24">
						<path fill="black" style="background:white;" d="M18.68,12.32C16.92,10.56 14.07,10.57 12.32,12.33C10.56,14.09 10.56,16.94 12.32,18.69C13.81,20.17 16.11,20.43 17.89,19.32L21,22.39L22.39,21L19.3,17.89C20.43,16.12 20.17,13.8 18.68,12.32M17.27,17.27C16.29,18.25 14.71,18.24 13.73,17.27C12.76,16.29 12.76,14.71 13.74,13.73C14.71,12.76 16.29,12.76 17.27,13.73C18.24,14.71 18.24,16.29 17.27,17.27M10.9,20.1C10.25,19.44 9.74,18.65 9.42,17.78C6.27,17.25 4,15.76 4,14V17C4,19.21 7.58,21 12,21V21C11.6,20.74 11.23,20.44 10.9,20.1M4,9V12C4,13.68 6.07,15.12 9,15.7C9,15.63 9,15.57 9,15.5C9,14.57 9.2,13.65 9.58,12.81C6.34,12.3 4,10.79 4,9M12,3C7.58,3 4,4.79 4,7C4,9 7,10.68 10.85,11H10.9C12.1,9.74 13.76,9 15.5,9C16.41,9 17.31,9.19 18.14,9.56C19.17,9.09 19.87,8.12 20,7C20,4.79 16.42,3 12,3Z" />
					</svg>
					<span class="tooltiptext-bottom"> Buscar </span>
				</div>
			</li>
			<li class="no-hover">   |   </li>			
			<li id="save">
				<div class="tooltip-bottom">
					<svg style="width:34px;height:34px" viewBox="0 0 24 24">
						<path id="save-color" fill="white" d="M15,9H5V5H15M12,19A3,3 0 0,1 9,16A3,3 0 0,1 12,13A3,3 0 0,1 15,16A3,3 0 0,1 12,19M17,3H5C3.89,3 3,3.9 3,5V19A2,2 0 0,0 5,21H19A2,2 0 0,0 21,19V7L17,3Z" />
					</svg>
					<span class="tooltiptext-bottom"> Guardar </span>
				</div>
			</li>
			<li id="export">
			<div class="tooltip-bottom">
					<svg style="width:34px;height:34px" viewBox="0 0 24 24">
						<path fill="white" d="M17.86 18L18.9 19C17.5 20.2 14.94 21 12 21C7.59 21 4 19.21 4 17V7C4 4.79 7.58 3 12 3C14.95 3 17.5 3.8 18.9 5L17.86 6L17.5 6.4C16.65 5.77 14.78 5 12 5C8.13 5 6 6.5 6 7S8.13 9 12 9C13.37 9 14.5 8.81 15.42 8.54L16.38 9.5H13.5V10.92C13 10.97 12.5 11 12 11C9.61 11 7.47 10.47 6 9.64V12.45C7.3 13.4 9.58 14 12 14C12.5 14 13 13.97 13.5 13.92V14.5H16.38L15.38 15.5L15.5 15.61C14.41 15.86 13.24 16 12 16C9.72 16 7.61 15.55 6 14.77V17C6 17.5 8.13 19 12 19C14.78 19 16.65 18.23 17.5 17.61L17.86 18M18.92 7.08L17.5 8.5L20 11H15V13H20L17.5 15.5L18.92 16.92L23.84 12L18.92 7.08Z" />
					</svg>
					<span class="tooltiptext-bottom"> no </span>
				</div>
			</li>
			<li id="settings">
				<div class="tooltip-bottom">
					<svg style="width:34px;height:34px" viewBox="0 0 24 24">
						<path fill="white" d="M12,15.5A3.5,3.5 0 0,1 8.5,12A3.5,3.5 0 0,1 12,8.5A3.5,3.5 0 0,1 15.5,12A3.5,3.5 0 0,1 12,15.5M19.43,12.97C19.47,12.65 19.5,12.33 19.5,12C19.5,11.67 19.47,11.34 19.43,11L21.54,9.37C21.73,9.22 21.78,8.95 21.66,8.73L19.66,5.27C19.54,5.05 19.27,4.96 19.05,5.05L16.56,6.05C16.04,5.66 15.5,5.32 14.87,5.07L14.5,2.42C14.46,2.18 14.25,2 14,2H10C9.75,2 9.54,2.18 9.5,2.42L9.13,5.07C8.5,5.32 7.96,5.66 7.44,6.05L4.95,5.05C4.73,4.96 4.46,5.05 4.34,5.27L2.34,8.73C2.21,8.95 2.27,9.22 2.46,9.37L4.57,11C4.53,11.34 4.5,11.67 4.5,12C4.5,12.33 4.53,12.65 4.57,12.97L2.46,14.63C2.27,14.78 2.21,15.05 2.34,15.27L4.34,18.73C4.46,18.95 4.73,19.03 4.95,18.95L7.44,17.94C7.96,18.34 8.5,18.68 9.13,18.93L9.5,21.58C9.54,21.82 9.75,22 10,22H14C14.25,22 14.46,21.82 14.5,21.58L14.87,18.93C15.5,18.67 16.04,18.34 16.56,17.94L19.05,18.95C19.27,19.03 19.54,18.95 19.66,18.73L21.66,15.27C21.78,15.05 21.73,14.78 21.54,14.63L19.43,12.97Z" />
					</svg>
					<span class="tooltiptext-bottom">no</span>
				</div>
			</li>
		</ul>
	</div>
	<div id='store' class='table-conainer small-text'></div>	
	</div>`;
		 
	datepickerFactory($);
	$('#enddate').datepicker();
	$('#madedate').datepicker();	

	document.getElementById('home-icon').onclick = getAll;
	document.getElementById('search').onclick = doSearch;
	document.getElementById('save').onclick = doSave;
	document.getElementById('export').onclick = doExport;
	document.getElementById('settings').onclick = doSettings;
	document.getElementById('world').onclick = doWorld;

	
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
			drawClients(json)
		});
	}
	function doSave(e){

		console.log("doSave")
		console.log(updates)
		if(updates.length < 1){
			alert("no updated activities were found")
			return
		}
		var ids = []
		for (var i = 0; i< updates.length; i++){
			var day = document.getElementById(updates[i]).value
			
			var x = updates[i].split("-")
			console.log(x)

			// dude trust me
			if(x.length == 2){
				var up = {
					"routeid":parseInt(x[0],10),
					"clientid":parseInt(x[1],10),
					"activityid":-1,
					"comments":document.getElementById(updates[i]).innerText
				}	
				ids.push(up)
			}
			if(x.length == 4){
				var up = {
					"routeid":parseInt(x[0],10),
					"clientid":parseInt(x[1],10),
					"activityid":parseInt(x[2],10),
					"madedate": day
				}	
				ids.push(up)
			}
		}
		console.log("x")
		console.log(ids)
		
		//TODO
		var prom = new Promise(function(resolve, reject) {
			//get and send modified clients, then set them to 0
			//each client and acivity needs its id

			var json = window.backend.Handler.Save(JSON.stringify(ids))
			if (json) {
			  	console.log("stuff worked")
				resolve(json);
			}  else {
				console.log("fuuck")
			  	reject(new Error("It broke"));
			}
		});
		   
		prom.then(function(json) {			
			
			//console.log(json)
			if(json == ""){
				window.alert("something went wrong")
				return
			}
			
			updates = []
			json = JSON.parse(json)
			drawClients(json)
		});
		
	}
	function doExport(){

		console.log("doExport")
		var prom = new Promise(function(resolve, reject) {
			var json = window.backend.Handler.DoExport()
			if (json) {
			  	console.log("stuff worked")
				resolve(json);
			}  else {
				console.log("fuuck")
			  	reject(new Error("It broke"));
			}
		});
		   
		prom.then(function(json) {
			console.log(json)
			window.alert(json)
		});


	}
	function doSettings(){

		console.log("doSettings")
	}
	function doWorld(){

		console.log("doWorld")
	}	
	//mystore.set(0);
	getAll()
};
// We provide our entrypoint as a callback for runtime.Init
runtime.Init(start);

//EVENTS works but.. :(
	/*

	window.wails.Events.On("debug", function( message) {
		console.log(message)
		
		//alert(message)		
		// Get the snackbar DIV
		var x = document.getElementById("snackbar");
		console.log(x.hidden)

		// Add the "show" class to DIV
		x.innerText = "<h1>"+message+"</h1>"
		x.hidden = true
		console.log(x.hidden)
		// After 3 seconds, remove the show class from DIV
		
		setTimeout(function(){ 
			x.hidden = false
			x.innerText = ""
			console.log("timeouuutt") 
		}, 3000);	
			
	});
	*/	
	// STORE I see no difference, give it a try
	/*
	mystore.subscribe( function(state) {
		console.log("MY store ")
		drawClients(JSON.parse(state));
	});
	*/
	//console.log(mystore);