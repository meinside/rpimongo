var DEBUG = false;
//var DEBUG = true;

var MONITOR_VALUES = [
	"hostname",
	"uptime",
	"uname",
	"cpu_info",
	"cpu_temperature",
	"free_spaces",
	"memory_split",
	"free_memory",
];

var SPINNER_OPTIONS = {
	color: "#A0A0A0",
	length: 3,
	width: 2,
	radius: 3,
};

/*
 * for printing debug messages
 */
function debugLog(log)
{
	if(!DEBUG)
		return;

	console.log(log);
}

/*
 * for fetching values from /api
 */
function fetch(value, renderer)
{
	debugLog("< fetching " + value + "...");

	// start spinner
	renderer.spin(SPINNER_OPTIONS);

	$.get(
		"/api/" + value + ".json",
		{},
		function(json) {
			var value = json["value"];
			renderer.text(value);

			debugLog("> fetched: " + value);

			// stop spinner
			renderer.spin(false);
		},
		"json");
}
