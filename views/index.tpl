<h3>
  Raspberry Pi Monitoring with Go
</h3>

<div class="server-stat">
	<h6>Hostname</h6>
	<pre class="stat" id="hostname">
		<!-- hostname here -->
	</pre>
	<br/>
	<h6>
		Uptime <span class="glyphicon glyphicon-refresh refresh" id="refresh-uptime"/>
	</h6>
	<pre class="stat" id="uptime">
		<!-- uptime here -->
	</pre>
	<br/>
	<h6>
		Machine
	</h6>
	<pre class="stat" id="uname">
		<!-- uname here -->
	</pre>
	<pre class="stat" id="cpu_info">
		<!-- cpu_info here -->
	</pre>
	<br/>
	<h6>
		CPU Temperature <span class="glyphicon glyphicon-refresh refresh" id="refresh-cpu_temperature"/>
	</h6>
	<pre class="stat" id="cpu_temperature">
		<!-- cpu_temperature here -->
	</pre>
	<br/>
	<h6>
		Free Spaces <span class="glyphicon glyphicon-refresh refresh" id="refresh-free_spaces"/>
	</h6>
	<pre class="stat" id="free_spaces">
		<!-- free_spaces here -->
	</pre>
	<br/>
	<h6>
		Memory Split
	</h6>
	<pre class="stat" id="memory_split">
		<!-- memory_split here -->
	</pre>
	<br/>
	<h6>
		Free Memory <span class="glyphicon glyphicon-refresh refresh" id="refresh-free_memory"/>
	</h6>
	<pre class="stat" id="free_memory">
		<!-- free_memory here -->
	</pre>
	<br/>
</div>
<script type="text/javascript">
	// on page load,
	$(document).ready(function(){
		var value, element;
		var numValues = MONITOR_VALUES.length;
		for(var i=0; i<numValues; i++)
		{
			value = MONITOR_VALUES[i];
			element = $("#" + value);

			// fetch value
			fetch(value, element);

			// bind refresh event
			$("#refresh-" + value).click(function(event){
				var value = event.currentTarget.id.replace(/^refresh-/, '');
				fetch(value, $("#" + value));
			});
		}
	});
</script>
