if ("serviceWorker" in navigator) {
    window.addEventListener("load", function() {
        navigator.serviceWorker
            .register("/static/service-worker.js")
            .then((res) => console.log("service worker registered"))
            .catch((err) =>
                console.log("service worker not registered", err),
            );
    });
}

self.addEventListener("fetch", (fetchEvent) => {
    fetchEvent.respondWith(
        caches.match(fetchEvent.request).then((res) => {
            return res || fetch(fetchEvent.request);
        }),
    );
});

function go_home() {
    // make POST request to /gohome, which will redirect to /
    $.post("/gohome", function(data, status) {
        window.location.href = "/";
    });
}

var menu_shown = false;

function show_menu(event) {
    var menu = document.getElementById("menu");
    menu.style.display = "block";
    menu.style.left = event.clientX + "px";
    menu.style.top = event.clientY + "px";
    menu_shown = true;
}

function hide_menu() {
    var menu = document.getElementById("menu");
    menu.style.display = "none";
    menu_shown = false;
}

function mouse_inside_menu(event) {
    var menu = document.getElementById("menu");
    var x = event.clientX;
    var y = event.clientY;
    var rect = menu.getBoundingClientRect();
    if (
        x < rect.left ||
        x > rect.right ||
        y < rect.top ||
        y > rect.bottom
    ) {
        return false;
    }
    return true;
}

function new_folder() {
    menu_item = document.getElementById("menu-item");
    menu_item.children[0].innerHTML = "Clicked";
}

function mouse_handler(event) {
    switch (event.button) {
        case 0:
            event.preventDefault();
            event.stopImmediatePropagation();
            if (menu_shown && !mouse_inside_menu(event)) {
                hide_menu();
            }
            // left click
            break;
        case 1:
            // middle click
            break;
        case 2:
            // trigger menu
            event.preventDefault();
            switch (menu_shown) {
                case true:
                    hide_menu();
                    break;
                case false:
                    show_menu(event);
                    break;
            }
            break;
        case 3:
            window.history.back();
            break;
        case 4:
            window.history.forward();
            break;
        default:
            break;
    }
}

document.addEventListener("mousedown", mouse_handler);
