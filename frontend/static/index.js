if ("serviceWorker" in navigator) {
    window.addEventListener("load", function() {
        navigator.serviceWorker
            .register("/static/service-worker.js")
            .then((_) => console.log("service worker registered"))
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
    $.post("/gohome", function(_, _) {
        window.location.href = "/";
    });
}

var menu_shown = false;

function get_drive_under_mouse(event) {
    var drive = document.getElementsByClassName("drive");
    for (var i = 0; i < drive.length; i++) {
        var rect = drive[i].getBoundingClientRect();
        if (
            event.clientX >= rect.left &&
            event.clientX <= rect.right &&
            event.clientY >= rect.top &&
            event.clientY <= rect.bottom
        ) {
            selected_drive = drive[i].getAttribute("hx-post");
            return drive[i];
        }
    }
    return null;
}

function show_menu(event) {
    var menu = document.getElementById("menu");
    if (get_drive_under_mouse(event) != null) {
        menu.children.namedItem("index").style.display = "block";
    } else {
        menu.children.namedItem("index").style.display = "none";
    }
    menu.style.display = "block";
    menu.style.left = event.clientX + "px";
    menu.style.top = event.clientY + "px";
    menu_shown = true;
}

var selected_drive = null;

function index_drive() {
    selected_drive = selected_drive.replace("/setdir/", "");
    console.log("indexing drive: " + selected_drive);
    // TODO: make POST request to /indexdrive with selected_drive as data
    // let the user know that the drive is being indexed

    selected_drive = null;
}

function hide_menu() {
    var menu = document.getElementById("menu");
    menu.children.namedItem("index").style.display = "none";
    menu.style.display = "none";
    selected_drive = null;
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
            event.stopImmediatePropagation();
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
