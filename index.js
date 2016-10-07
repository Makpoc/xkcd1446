(function() {
	var SERVERS = [
		'//c0.xkcd.com',
		'//c1.xkcd.com',
		'//c2.xkcd.com',
		'//c3.xkcd.com',
		'//c4.xkcd.com',
		'//c5.xkcd.com',
		'//c6.xkcd.com',
		'//c7.xkcd.com'
	]

	var IMGPREFIX= 'data:image/png;base64,'

	//	function record(name) {
	//		new Image().src = location.protocol + '//xkcd.com/events/' + name
	//	}

	function log() {
		if (location.hash == '#verbose') {
			console.log.apply(console, arguments)
		}
	}

	try {
		//		var server = location.protocol + SERVERS[Math.floor(Math.random() * SERVERS.length)],
		//			esURL = server + '/stream/comic/landing?method=EventSource',
		//			source = new EventSource(esURL)
		//
		//		log('connecting to event source:', esURL)
		//		source.addEventListener('open', function(ev) {
		//			log('open');
		//			record('connect_start');
		//		}, false)
		//
		//		source.addEventListener('error', function(ev) {
		//			log('connection error', ev)
		//			record('connect_error')
		//		}, false)
		//
		//		source.addEventListener('comic/landing', function(ev) {
		//			log('comming/langing');
		//		}, false)
		//		var currImage = 0;
		//
		//		var image = document.getElementById('landing').getElementsByTagName('img')[0];
		//		if (image === undefined) {
		//			image = document.createElement('img');
		//			document.getElementById('landing').appendChild(image);
		//		};
		//
		//		setTimeout(function() {
		//			var newimg = IMGPREFIX;
		//			if (image.src != newimg) { log('changing image', image.src, newimg); }
		//			image.src = newimg;
		//			currImage++;
		//		}, 5000);

		var source = new EventSource('/landing');
		source.addEventListener('landing', function(ev){
			var newimage = IMGPREFIX + ev.data;
			var image = document.getElementById('landing').getElementsByTagName('img')[0];
			if (image === undefined ) {
				log('creating image');
				image = document.createElement('img');
				document.getElementById('landing').appendChild(image);
			}
			image.src = newimage;
		}, false);

			//	source.addEventListener('comic/landing/reload', function(ev) {
			//		var delay = Math.round(Math.random() * 55)
			//		log('reloading in', delay + 5, 'seconds')
			//		setTimeout(function() {
			//			record('reloading')
			//
			//			// give the record a little time to be sent
			//			setTimeout(function() {
			//				location.reload()
			//			}, 5 * 1000)
			//		}, delay * 1000)
			//	}, false)
		} catch (e) {
			log('js_error', e)
		}
	})()
