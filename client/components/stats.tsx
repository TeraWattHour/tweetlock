import { useEffect, useState } from "react";
import type { Root } from "react-dom/client";

import { useVotesStore } from "~store";
import { numberWithCommas, waitForElement, type User } from "~utils";

export const Component = ({ articleContainer, user, root }: { user: User; root: Root; articleContainer: HTMLElement }) => {
    const [status, setStatus] = useState<"loading" | "error" | "ok">("loading");
    const { votes } = useVotesStore();
    const [results, setResults] = useState<{ votes: number; hasVoted: boolean } | null>(null);

    useEffect(() => {
        articleContainer.onmouseenter = () => {
            if (status === "ok" && results) {
                return;
            }
            window.postMessage({ source: "tweetlock", type: "GET_VOTES", payload: { userId: user.id } });
        };
    }, [articleContainer, status, results, user.id]);

    useEffect(() => {
        const forUser = votes[user.id];
        if (forUser && !forUser.hasError) {
            setStatus("ok");
            setResults(forUser);
        } else if (forUser && forUser.hasError) {
            setStatus("error");
            setResults(null);
        }
    }, [votes, user.id, votes[user.id]]);

    const onBlock = async () => {
        const more = await waitForElement(articleContainer, '[aria-label="More"]');
        if (!more) return;
        more.click();

        const blockButton = await waitForElement(document.body, '[data-testid="block"]');
        if (!blockButton) return;
        blockButton.click();

        const confirmation = await waitForElement(document.body, '[data-testid="confirmationSheetConfirm"]');
        if (!confirmation) return;
        confirmation.click();

        root.unmount();
    };

    const onShame = () =>
        window.postMessage({ source: "tweetlock", type: !results?.hasVoted ? "SHAME_USER" : "UNSHAME_USER", payload: { userId: user.id } });

    return (
        <>
            <button onClick={onBlock} className="tl__action tl__action--block" title={`Block ${user.name}`}>
                <div className="tl__icon-container">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor">
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
                        />
                    </svg>
                </div>
                <span>Block</span>
            </button>
            <button
                onClick={() => {
                    if (!results || status !== "ok") {
                        return;
                    }
                    onShame();
                }}
                className={`tl__action tl__action--shame ${results?.hasVoted && "active"}`}
                data-tlstatus={status}
                title={`Shame ${user.name}`}
                disabled={status !== "ok"}>
                <div className="tl__icon-container">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor">
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M7.5 15h2.25m8.024-9.75c.011.05.028.1.052.148.591 1.2.924 2.55.924 3.977a8.96 8.96 0 01-.999 4.125m.023-8.25c-.076-.365.183-.75.575-.75h.908c.889 0 1.713.518 1.972 1.368.339 1.11.521 2.287.521 3.507 0 1.553-.295 3.036-.831 4.398C20.613 14.547 19.833 15 19 15h-1.053c-.472 0-.745-.556-.5-.96a8.95 8.95 0 00.303-.54m.023-8.25H16.48a4.5 4.5 0 01-1.423-.23l-3.114-1.04a4.5 4.5 0 00-1.423-.23H6.504c-.618 0-1.217.247-1.605.729A11.95 11.95 0 002.25 12c0 .434.023.863.068 1.285C2.427 14.306 3.346 15 4.372 15h3.126c.618 0 .991.724.725 1.282A7.471 7.471 0 007.5 19.5a2.25 2.25 0 002.25 2.25.75.75 0 00.75-.75v-.633c0-.573.11-1.14.322-1.672.304-.76.93-1.33 1.653-1.715a9.04 9.04 0 002.86-2.4c.498-.634 1.226-1.08 2.032-1.08h.384"
                        />
                    </svg>
                </div>
                <span>
                    {status === "loading" && "Loading"}
                    {status === "ok" && numberWithCommas(results?.votes || 0)}
                    {status === "error" && "Something went wrong"}
                </span>
            </button>
        </>
    );
};
