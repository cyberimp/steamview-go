"use strict";

let host = window.location.host;
let socket = new WebSocket("ws://" + host + "/socket");
let logo = document.getElementById("logo");
let hero = document.getElementById("hero");
let name = document.getElementById("name");
let container = document.getElementById("container");


/**
 * hides logo on 404
 * @returns {boolean}
 */
function LogoError() {
    logo.className = "hidden";
    name.className = "";
    container.className = "";
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
     * @property {string} width - logo width in percents
     * @property {string} height - logo height in percents
     * @property {string} name - name of game running
     */
    let message = JSON.parse(msg.data);

    logo.className = message.align;
    if (message.align === "BottomLeft" &&
        (!message.width.includes(".") || parseFloat(message.width) > 90 ))
    {
        container.style.width = "50%";
    }
    else
    {
        container.style.width = "100%";
    }

    logo.style.width = message.width + "%";
    logo.style.height = message.height + "%";

    name.innerText = message.name;
    name.className = (message.align === "hidden")?"":"hidden";

    logo.onerror = LogoError;
    logo.src = message.logo;
    hero.src = message.hero;
};