const config = JSON.parse(atob(document.currentScript.getAttribute('data-config')));
config.url = new URL(config.url);

const rewrites = {url: url => (url = new URL(url) ? config.url + btoa(url) : config.url + btoa(atob(window.location.href) + url))};

let document = new Proxy(document, {
	get: (target, prop) => {
		switch (prop) {
		case 'location':
			return rewrites.url(prop);
		default:
			Reflect.get(target, prop);
		}
	}
});

let window = new Proxy(document, {
	get: (target, prop) => {
		switch (prop) {
		case 'document':
			return document;
		case 'window':
			return window;
		default:
			Reflect.get(target, prop);
		}
	}
});

document.prototype.write = new Proxy(document.prototype.write, {
    apply: (target, thisArg, args) => {
        var doc = domparser.parseFromString(args[0], 'text/html');
        // TODO: Rewrite and send back data
    }
});

const historyHandler = {
    apply: (target, thisArg, args) => {
        args[2] = rewrites.url();
        return Reflect.apply(target, thisArg, args);
    }
};

window.History.prototype.pushState = new Proxy(window.History.prototype.pushState, historyHandler);
window.History.prototype.replaceState = new Proxy(window.History.prototype.replaceState, historyHandler);

window.open = new Proxy(window.open, {
    apply: (target, thisArg, args) => {
        args[0] = rewrites.url();
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

// Fetc and XMLHttpRequest
delete window.fetch;
delete window.XMLHttpRequest;

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
