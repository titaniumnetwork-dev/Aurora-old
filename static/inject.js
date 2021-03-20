// TODO: Remove space indent

const config = {httpprefix: '{{.HTTPPrefix}}', wsprefix: '{{.WSPrefix}}', url: new URL('{{.URL}}'), proxyurl: new URL('{{.ProxyURL}}')}

// TODO: Rewrite html and css
const rewrites = {
	isUrl: url => {
		try {
			return Boolean(new URL(url));
		} catch (err) {
			return false
		}
	},
	url: url => {
		if (rewrites.isUrl(url) == false) {
			pURL = config.proxyurl.toString()
			pathSplit = pURL.split('/')

			var split
		
			switch (true) {
			case url.split(':').length>=2:
			case url.split('../').length=2:
				split = url.split('../')
				url = config.url.origin + config.httpprefix+btoa(pURL.splice(0, len(split).join('/'))+split.pop());
			case url.startsWith('//'):
				split = url.split('/')
				url = config.url.origin + config.httpprefix+btoa(split.pop());
			case url.startsWith('/'):
				url = config.url.origin + config.httpprefix+btoa(pURL.origin + url);
			default:
				url = config.url.origin + config.httpprefix+btoa(pURL + '/' + url);
			}
		} else if (rewrites.isUrl(url) == true) {
			url = config.url.origin + config.httpprefix + btoa(url);
		}

		return url;
	},
	html: html => {
		// TODO: Avoid selecting id
		var dom = new DOMParser().parseFromString(html, 'text/html'), sel = dom.querySelector('*');

		sel.querySelectorAll('*').forEach(node => {
			switch(node.tagName) {
			case 'SCRIPT':
				node.textContent = "{let document=audocument;" + node.textContent + "}"
				break;
			case 'STYLE':
				node.textContent = node.textContent.replace(/(?<=url\((?<a>["']?)).*?(?=\k<a>\))/gi, rewrites.url)
				break;
			}
			node.getAttributeNames().forEach(attr => {
				switch (attr) {
				case 'script':
					node.setAttribute("{let document=audocument;" + node.getAttribute(attr) + "}")
				case 'style':
					node.setAttribute(node.getAttribute(attr).replace(/(?<=url\((?<a>["']?)).*?(?=\k<a>\))/gi, rewrites.url))
				// TODO: Handle other attributes see server side rewrites for reference
				}
			});
		});

		return sel.innerHTML
	},
};

audocument = new Proxy(document, {
	get: (target, prop) => {
		switch (prop) {
		case 'location':
			return config.proxyurl
		default:
			Reflect.get(target, prop);
		}

		return typeof(prop=Reflect.get(target,prop))=='function'?prop.bind(target):prop;
	},
	set: (target, prop, value, reciever) => {
		target[prop] = value;
	}
});

document.write = new Proxy(document.write, {
	apply: (target, thisArg, args) => {
        html = rewrites.html(args[0])

		return Reflect.apply(target, thisArg, args);
    }
});

const historyHandler = {
	apply: (target, thisArg, args) => {
		args[2] = rewrites.url(args[2]);
		return Reflect.apply(target, thisArg, args);
	}
};

window.History.prototype.pushState = new Proxy(window.History.prototype.pushState, historyHandler);
window.History.prototype.replaceState = new Proxy(window.History.prototype.replaceState, historyHandler);

window.open = new Proxy(window.open, {
    apply: (target, thisArg, args) => {
		args[0] = rewrites.url(args[0]);

		return Reflect.apply(target, thisArg, args);
    }
});

window.fetch = new Proxy(window.fetch, {
	apply: (target, thisArg, args) => {
		args[0] = rewrites.url(args[0]);

		return Reflect.apply(target, thisArg, args);
    }
});

window.XMLHttpRequest.prototype.open = new Proxy(window.XMLHttpRequest.prototype.open, {
	apply: (target, thisArg, args) => {
		args[1] = rewrites.url(args[1])
		
		return Reflect.apply(target, thisArg, args);
    }
});

window.Navigator.prototype.sendBeacon = new Proxy(window.Navigator.prototype.sendBeacon, {
    apply: (target, thisArg, args) => {
		args[0] = rewrites.url(args[0]);
		
		return Reflect.apply(target, thisArg, args);
    }
});

/*
window.Websocket = new Proxy(window.Websocket, {
    construct: (target, args) => {
		// TODO: rewrite
		Reflect.construct(target, args)
    }
});
*/

// Delete non-proxified objects so requests don't escape the proxy

// WebSocket
delete window.WebSocket;

// WebRTC
delete window.MediaStreamTrack; 
delete window.RTCPeerConnection;
delete window.RTCSessionDescription;
delete window.mozMediaStreamTrack;
delete window.mozRTCPeerConnection;
delete window.mozRTCSessionDescription;
delete window.navigator.getUserMedia;
delete window.navigator.mozGetUserMedia;
delete window.navigator.webkitGetUserMedia;
delete window.webkitMediaStreamTrack;
delete window.webkitRTCPeerConnection;
delete window.webkitRTCSessionDescription;
