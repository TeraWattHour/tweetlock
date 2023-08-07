// @ts-ignore
import styleText from "data-text:~styles/content.css";
import type { PlasmoCSConfig } from "plasmo";
import { createRoot } from "react-dom/client";

import { Component } from "~components/stats";
import { useVotesStore } from "~store";
import { findUserData } from "~utils";

const style = document.createElement("style");
style.textContent = styleText;
(document.head || document.documentElement).appendChild(style);

const observer = new MutationObserver(() => {
    const articles = [...document.querySelectorAll('article:where(:not([data-tlfound="true"]))')] as HTMLElement[];
    for (const article of articles) {
        prepareArticle(article);
    }
});

observer.observe(document.body, {
    subtree: true,
    childList: true
});

const prepareArticle = (article: HTMLElement) => {
    article.dataset.tlfound = "true";

    const user = findUserData(article);
    if (!user) {
        return;
    }

    const outlet = article.firstChild?.firstChild?.lastChild?.lastChild as HTMLElement | undefined;
    if (!outlet) {
        return;
    }

    window.postMessage({ source: "tweetlock", type: "GET_VOTES", payload: { userId: user.id } });

    outlet.style.position = "relative";
    outlet.style.paddingBottom = "36px";

    const wrapper = document.createElement("div");
    wrapper.className = "tl__wrapper";
    outlet.appendChild(wrapper);

    const container = document.createElement("div");
    container.className = "tl__container";
    wrapper.appendChild(container);

    const root = createRoot(container);
    root.render(<Component user={user} root={root} articleContainer={article} />);
};

window.addEventListener("message", ({ data }) => {
    if (data.source !== "tweetlock") {
        return;
    }

    if (data.type === "UPDATE_VOTES") {
        const { userId, votes, hasVoted, hasError } = data.payload;

        const snap = useVotesStore.getState().setVotesForUser;
        return snap(userId, votes, hasVoted, hasError);
    }
});

export const config: PlasmoCSConfig = {
    matches: ["https://twitter.com/*"],
    world: "MAIN"
};

export default () => {};
