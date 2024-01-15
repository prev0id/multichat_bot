const notFoundRoutIndex = 0;
const checked = "checked";
const unchecked = "";

const navigateTo = url => {
    history.pushState(null, null, url);
    router();
};

let previousPathID = "";

function convertPathIDToHTMLID(path) {
    return `${path}-tab-button`;
}

const router = async () => {
    const routes = [
        {
            path: "/404", id: "", view: () => {
                console.log("viewing 404")
            }
        },
        {
            path: "/", id: "stream-settings", view: () => {
                console.log("viewing stream settings")
            }
        },
        {
            path: "/bot-settings", id: "bot-settings", view: () => {
                console.log("viewing bot settings")
            }
        },
        {
            path: "/commands", id: "commands", view: () => {
                console.log("viewing commands")
            }
        },
    ];

    let match = routes.find(route => {
        return location.pathname === route.path;
    });

    if (!match) {
        match = routes[notFoundRoutIndex];
    }


    if (previousPathID) {
        const button = document.getElementById(convertPathIDToHTMLID(previousPathID));
        button.dataset.value = unchecked;
    }

    const button = document.getElementById(convertPathIDToHTMLID(match.id));
    button.dataset.value = checked;
    previousPathID = match.id;

    match.view()
};

window.addEventListener("popstate", router);

document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", e => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.href);
            return;
        }

        if (e.target.parentNode.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.parentNode.href);
        }
    });

    router();
});