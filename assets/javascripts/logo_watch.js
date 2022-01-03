/* eslint-env browser */
const { host } = window.location;
const socket = new WebSocket(`ws://${host}/socket`);
const logo = document.getElementById('logo');
const hero = document.getElementById('hero');
const name = document.getElementById('name');
const container = document.getElementById('container');

/**
 * hides logo on 404
 * @returns {boolean}
 */
const LogoError = () => {
  logo.className = 'hidden';
  name.className = '';
  container.className = '';
  logo.onerror = null;
  return true;
};

logo.onerror = LogoError;

const circle = new window.ProgressBar.Circle(container, {
  color: '#fff',
  strokeWidth: 3,
  duration: 100,
  svgStyle: {
    position: 'absolute',
    height: '50%',
    left: '50%',
    top: '50%',
    transform: {
      prefix: true,
      value: 'translate(-50%, -50%)'
    },
  }
});

circle.destroyed = false;

circle.tryDestroy = () => {
  if (!circle.destroyed) {
    circle.destroy();
    circle.destroyed = true;
  }
};

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
  const message = JSON.parse(msg.data);

  logo.className = message.align;
  logo.onerror = LogoError;
  logo.src = message.logo;
  hero.src = message.hero;

  if (message.name === '_VDF_READING') {
    container.style.width = '100%';
    circle.animate(message.width);
    logo.style.width = '50%';
    logo.style.height = '50%';
    name.innerText = 'Loading appinfo.vdf…';
    name.className = '';
    return;
  }
  circle.tryDestroy();

  // some bicycle from Gaben
  if (message.align === 'BottomLeft'
        && (!message.width.includes('.') || parseFloat(message.width) > 90)) {
    container.style.width = '50%';
  } else {
    container.style.width = '100%';
  }

  logo.style.width = `${message.width}%`;
  logo.style.height = `${message.height}%`;
  name.innerText = message.name;
  name.className = (message.align === 'hidden') ? '' : 'hidden';
};

/** *
 * Displays message about server gone down
 */
const setError = () => {
  logo.src = '/images/error.png';
  logo.className = 'CenterCenter';
  logo.style.width = '50%';
  logo.style.height = '50%';
  container.style.width = '100%';
  name.className = 'hidden';
  circle.tryDestroy();
};

/** *
 * Try reloading page on socket disconnect, this includes server crashing, shutdown
 * and going into ACPI sleep state
 */
const tryReload = () => {
  const xhr = new XMLHttpRequest();
  xhr.onload = () => setTimeout(window.location.reload.bind(window.location), 100);
  xhr.onerror = setError;
  // Parameter added to force ignore cache on Chrome
  xhr.open('GET', `http://${host}/?_=${new Date().getTime()}`, true);
  xhr.send();
};

socket.onclose = () => setTimeout(tryReload, 300);
