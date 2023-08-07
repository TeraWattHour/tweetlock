import { tryRefresh } from "~utils";

let refreshFailed = false;

const fetcher = async (url: string, options?: RequestInit) =>
    new Promise<Response>((resolve, reject) => {
        fetch(url, options)
            .then(async (res) => {
                if (res.ok) {
                    return resolve(res);
                }

                if (res.status === 401 && !refreshFailed) {
                    return tryRefresh()
                        .then((res) => res.ok && fetch(url, options).then(resolve).catch(reject))
                        .catch(() => resolve(res));
                }

                return reject(res);
            })
            .catch(reject);
    });

chrome.runtime.onMessage.addListener((message, _sender, sendResponse) => {
    (async () => {
        if (message.type === "FETCH") {
            try {
                const res = await fetcher(message.url, {
                    method: message.method,
                    credentials: "include"
                });
                let data = null;
                try {
                    data = await res.json();
                } catch (error) {}
                sendResponse({ ok: res.ok, code: res.status, data: data });
            } catch (error) {
                sendResponse({ ok: false, code: 500, data: null });
            }
        }
    })();

    return true;
});
