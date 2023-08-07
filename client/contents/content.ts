import type { PlasmoCSConfig } from "plasmo";

import { BACKEND_URL } from "~consts";

const voteCache: {
    [key: string]: { hasVoted: boolean | null; votes: number | null; hasError?: boolean };
} = {};
let queue: string[] = [];
let queueLastAt: Date | null = null;
let queueRequest: NodeJS.Timeout;

const notify = (type: string, payload: any) =>
    window.postMessage({
        source: "tweetlock",
        type,
        payload
    });

const backgroundFetch = (method: string, url: string, callback: (response: any) => void) =>
    chrome.runtime.sendMessage({ type: "FETCH", method, url }, callback);

const updateVotes = (userId: string, votes: number | null, hasVoted: boolean | null, hasError?: boolean) => {
    voteCache[userId] = {
        votes,
        hasVoted,
        hasError
    };
    sendVotes(userId);
};

const sendVotes = (userId: string) => {
    const entry = voteCache[userId];
    if (entry) {
        notify("UPDATE_VOTES", { userId, ...entry });
    }
};

const dispatchQueue = (userId: string) => {
    if (!queue.includes(userId)) {
        queue.push(userId);
    }
    if (queueLastAt && Date.now() - +queueLastAt < 200) {
        clearTimeout(queueRequest);
    } else {
        queueLastAt = new Date();
    }
    queueRequest = setTimeout(queueFetch, 200);
};

const queueFetch = () => {
    const currentQueue = queue.slice();
    const shift = currentQueue.length;
    if (shift == 0) {
        return;
    }
    backgroundFetch("GET", `${BACKEND_URL}/vote-count?targets=${currentQueue.join(",")}`, (res) => {
        if (res.ok && res.data) {
            for (const userId in res.data) {
                const entry = res.data[userId];
                if (entry) {
                    updateVotes(userId, entry.votes, entry.hasVoted);
                }
            }
            queue = queue.slice(shift);
        }
        if (!res.ok) {
            for (const userId of currentQueue) {
                updateVotes(userId, null, null, true);
            }
        }
    });
};

window.addEventListener("message", async ({ data }) => {
    if (data.source !== "tweetlock") {
        return;
    }

    if (data.type === "GET_VOTES") {
        const { userId } = data.payload;
        if (Object.keys(voteCache).includes(userId)) {
            return sendVotes(userId);
        }

        return dispatchQueue(userId);
    }

    if (data.type === "SHAME_USER" || data.type === "UNSHAME_USER") {
        const { userId } = data.payload;
        const addMode = data.type === "SHAME_USER";

        const delta = addMode ? 1 : -1;
        const currentVotes = Math.max(0, voteCache[userId].votes + delta);
        updateVotes(userId, currentVotes, addMode);

        return backgroundFetch(addMode ? "POST" : "DELETE", `${BACKEND_URL}/vote/${userId}`, (response) => {
            if (!response?.ok) {
                return updateVotes(userId, voteCache[userId].votes - delta, !addMode);
            }
            if (response.code == 200 && addMode) {
                return updateVotes(userId, voteCache[userId].votes - 1, true);
            }
        });
    }
});

export const config: PlasmoCSConfig = {
    matches: ["https://twitter.com/*"]
};
