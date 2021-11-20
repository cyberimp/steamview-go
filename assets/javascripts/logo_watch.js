let host = window.location.host;
let socket = new WebSocket("ws://" + host + "/socket");
let logo = document.getElementById("logo");
let hero = document.getElementById("hero");

let img = new Image();

img.src = "/images/error.png" //cache image for display on shutdown

function LogoError() {
    logo.classList.remove("left", "right", "center", "absolute-center", "left-stretch", "hidden");
    logo.classList.add("hidden");
    logo.onerror = null;
    return true;
}

logo.onerror = LogoError;

socket.onmessage = (msg) => {
    let message = JSON.parse(msg.data);
    logo.classList.remove("left", "right", "center", "absolute-center", "left-stretch", "hidden");
    logo.classList.add(message.align);
    logo.onerror = LogoError;
    logo.src = message.logo;
    hero.src = message.hero;
};