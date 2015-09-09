var DEBUG = false;
//var DEBUG = true;

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
 * for fetching given value from /api and render it to renderer
 */
function fetchAndRender(value, renderer)
{
	debugLog("< fetching " + value + "...");

	// start spinner
	renderer.spin(SPINNER_OPTIONS);

	$.get(
		"/v1/api/" + value,
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

/*
 * fetch multiple /api values
 */
function fetchValues(values)
{
	var value, element;
	var numValues = values.length;
	for(var i=0; i<numValues; i++)
	{
		value = values[i];
		element = $("#" + value);

		// fetch and render value
		fetchAndRender(value, element);

		// bind refresh event
		$("#refresh-" + value).click(function(event){
			var value = event.currentTarget.id.replace(/^refresh-/, '');
			fetchAndRender(value, $("#" + value));
		});
	}
}
