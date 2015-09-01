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
	// values for monitoring
	var MONITOR_VALUES = [
		'hostname',
		'uptime',
		'uname',
		'cpu_info',
		'cpu_temperature',
		'free_spaces',
		'memory_split',
		'free_memory',
	];

	// on page load,
	$(document).ready(function(){
		fetchValues(MONITOR_VALUES);
	});
</script>
