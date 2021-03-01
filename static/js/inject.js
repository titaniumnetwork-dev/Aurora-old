const rewrites = {
    url: (url) => {
        return null;
    }
};

const historyHandler = {
    apply: (target, prop, reciever) => {
        return Reflect.apply(target, prop,[reciever[1]+rewrites.url(reciever[2])]);
    }
};

window.History.prototype.pushState = new Proxy(self, historyHandler)
window.History.prototype.replaceState = new Proxy(self, historyHandler)