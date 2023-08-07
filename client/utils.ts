import { BACKEND_URL } from "~consts";

export interface User {
    handle: string;
    id: string;
    name: string;
}

export function findUserData(element: HTMLElement): User | null {
    const key = Object.keys(element).find((key) => key.startsWith("__reactFiber$"));
    let fiber = element[key]?.return;

    while (fiber) {
        let tweet = fiber.stateNode?.props?.tweet;
        if (tweet) {
            if (tweet.retweeted_status) {
                tweet = tweet.retweeted_status;
            }

            const { user } = tweet;
            if (user) {
                return { handle: user.screen_name, id: user.id_str, name: user.name };
            }
        }
        fiber = fiber.return;
    }

    return null;
}

export function waitForElement(base: HTMLElement, selector: string): Promise<HTMLElement | null> {
    return new Promise((resolve, reject) => {
        let element = base.querySelector(selector);
        if (element) {
            return resolve(element as HTMLElement);
        }
        const observer = new MutationObserver(() => {
            const t = setTimeout(() => {
                observer.disconnect();
                reject(null);
            }, 1000);
            element = base.querySelector(selector);
            if (element) {
                clearTimeout(t);
                observer.disconnect();
                resolve(element as HTMLElement);
            }
        });
        observer.observe(document.body, { childList: true, subtree: true });
    });
}

export function numberWithCommas(x: number) {
    return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

export async function tryRefresh() {
    return new Promise<Response>((resolve, reject) =>
        fetch(`${BACKEND_URL}/refresh`, {
            method: "POST",
            credentials: "include"
        })
            .then((res) => {
                if (!res.ok) return reject("couldnt refresh");

                return resolve(res);
            })
            .catch((err) => reject(err))
    );
}

export function decodeJWT(jwt: string) {
    const payloadBase64 = jwt.split(".")[1];
    if (!payloadBase64) return;
    const base64 = payloadBase64.replace(/-/g, "+").replace(/_/g, "/");
    return JSON.parse(window.atob(base64));
}
