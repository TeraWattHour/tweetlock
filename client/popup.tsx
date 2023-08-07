import "./styles/popup.css";

import icon from "data-base64:~assets/icon.png";
import { useEffect, useState } from "react";

import { BACKEND_URL } from "~consts";
import { decodeJWT, tryRefresh } from "~utils";

function IndexPopup() {
    const [isLoading, setIsLoading] = useState(true);
    const [user, setUser] = useState<any>(null);

    const getUser = (firstTry = true) => {
        chrome.cookies.get({ url: BACKEND_URL, name: "x-access" }, async (cookie) => {
            if (!cookie && firstTry) {
                try {
                    await tryRefresh();
                    return getUser(false);
                } catch (error) {
                    setUser(null);
                    setIsLoading(false);
                }
            }
            if (cookie) {
                try {
                    const decoded = decodeJWT(cookie.value)?.user;
                    setUser(decoded || null);
                } catch (error) {
                    setUser(null);
                } finally {
                    setIsLoading(false);
                }
            }
        });
    };

    useEffect(() => {
        getUser();
    }, []);

    return (
        <div className="tl__wrapper">
            <div className="tl__header">
                <img src={icon} width="20" height="20" style={{ marginRight: 8 }} />
                <header style={{ fontSize: "14px", fontWeight: 600 }}>TweetLock</header>
            </div>

            <div className="tl__container">
                {isLoading ? (
                    <p style={{ color: "lightgray", textAlign: "center" }}>Loading...</p>
                ) : user ? (
                    <div style={{ width: "100%", textAlign: "center" }}>
                        <p style={{ color: "lightgray" }}>Welcome</p>
                        <p
                            style={{
                                marginTop: 2,
                                fontWeight: 600,
                                fontSize: "15px"
                            }}>
                            {user.name}
                        </p>
                    </div>
                ) : (
                    <a id="login-with-twitter" href={`${BACKEND_URL}/twitter-redirect`} target="_blank" rel="noreferrer">
                        Login with Twitter
                    </a>
                )}
            </div>
        </div>
    );
}

export default IndexPopup;
