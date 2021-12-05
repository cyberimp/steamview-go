let host = window.location.host;
let socket = new WebSocket("ws://" + host + "/socket");
let logo = document.getElementById("logo");
let hero = document.getElementById("hero");


function LogoError() {
    logo.classList.remove("left", "right", "center", "absolute-center", "left-stretch", "hidden");
    logo.classList.add("hidden");
    logo.onerror = null;
    return true;
}

logo.onerror = LogoError;

socket.onmessage = (msg) => {
    /**
     * message from server should include this fields
     * @type {object} message
     * @property {string} align - align of logo on hero
     * @property {string} hero - image path for background of banner
     * @property {string} logo - image path for game logo
     */
    let message = JSON.parse(msg.data);
    logo.classList.remove("left", "right", "center", "absolute-center", "left-stretch", "hidden");
    logo.classList.add(message.align);
    logo.onerror = LogoError;
    logo.src = message.logo;
    hero.src = message.hero;
};