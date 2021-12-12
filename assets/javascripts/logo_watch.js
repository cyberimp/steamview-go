let host = window.location.host;
let socket = new WebSocket("ws://" + host + "/socket");
let logo = document.getElementById("logo");
let hero = document.getElementById("hero");

/**
 * hides logo on 404
 * @returns {boolean}
 */
function LogoError() {
    logo.className = "hidden";
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
    logo.className = message.align;
    logo.style.width = message.width + "%"
    logo.style.height = message.height + "%"
    logo.onerror = LogoError;
    logo.src = message.logo;
    hero.src = message.hero;
};